// check_go_schema_consistency verifies the intentionally small, explicit
// Go-to-JSON-Schema mapping declared in scripts/go-schema-map.json.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type mapping struct {
	GoDir  string `json:"go_dir"`
	GoType string `json:"go_type"`
	Schema string `json:"schema"`
}

type mappingFile struct {
	Pairs []mapping `json:"pairs"`
}

type shape struct {
	Kind string
	Elem *shape
}

type structField struct {
	Name      string
	Shape     shape
	OmitEmpty bool
}

type checker struct {
	root       string
	drifts     []string
	toolIssues []string
}

func main() {
	rootFlag := flag.String("repo-root", ".", "repository root")
	mapFlag := flag.String("map", "scripts/go-schema-map.json", "mapping JSON file")
	flag.Parse()

	root, err := filepath.Abs(*rootFlag)
	if err != nil {
		fmt.Printf("[TOOL_ERROR] repo root: %v\n", err)
		fmt.Println("RESULT: TOOL_ERROR")
		os.Exit(2)
	}
	c := checker{root: root}
	mappings := c.readMappings(*mapFlag)
	schemaFiles := c.schemaFiles()
	mappedSchemaFiles := map[string]bool{}
	passed := 0
	seenTypes := map[string]bool{}

	for _, pair := range mappings {
		seenTypes[pair.GoDir+"."+pair.GoType] = true
		if relative, ok := relativeToRoot(root, pair.Schema); ok {
			mappedSchemaFiles[filepath.ToSlash(relative)] = true
		}
		beforeDrifts, beforeTools := len(c.drifts), len(c.toolIssues)
		c.checkPair(pair)
		if beforeDrifts == len(c.drifts) && beforeTools == len(c.toolIssues) {
			passed++
		}
	}

	mappedExisting := 0
	for file := range mappedSchemaFiles {
		if schemaFiles[file] {
			mappedExisting++
		}
	}
	fmt.Printf("mapped_pairs=%d\n", len(mappings))
	fmt.Printf("passed_pairs=%d\n", passed)
	fmt.Printf("schema_files_total=%d\n", len(schemaFiles))
	fmt.Printf("mapped_schema_files=%d\n", mappedExisting)
	fmt.Printf("unmapped_schema_files=%d\n", len(schemaFiles)-mappedExisting)
	fmt.Printf("go_types_in_scope=%d\n", len(seenTypes))
	c.printIssues()

	if len(c.toolIssues) > 0 {
		fmt.Println("RESULT: TOOL_ERROR")
		os.Exit(2)
	}
	if len(c.drifts) > 0 {
		fmt.Println("RESULT: DRIFT")
		os.Exit(1)
	}
	fmt.Println("RESULT: PASS")
}

func (c *checker) readMappings(path string) []mapping {
	path = resolvePath(c.root, path)
	content, err := os.ReadFile(path)
	if err != nil {
		c.toolIssues = append(c.toolIssues, fmt.Sprintf("mapping file %s: %v", path, err))
		return nil
	}
	var config mappingFile
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&config); err != nil {
		c.toolIssues = append(c.toolIssues, fmt.Sprintf("mapping JSON %s: %v", path, err))
		return nil
	}
	if len(config.Pairs) == 0 {
		c.toolIssues = append(c.toolIssues, "mapping must declare at least one pair")
	}
	for index, pair := range config.Pairs {
		if pair.GoDir == "" || pair.GoType == "" || pair.Schema == "" {
			c.toolIssues = append(c.toolIssues, fmt.Sprintf("mapping pair %d lacks go_dir, go_type, or schema", index))
		}
	}
	return config.Pairs
}

func (c *checker) schemaFiles() map[string]bool {
	files := map[string]bool{}
	err := filepath.WalkDir(c.root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			relative, err := filepath.Rel(c.root, path)
			if err != nil {
				return err
			}
			relative = filepath.ToSlash(relative)
			if relative == ".git" || strings.HasPrefix(relative, ".git/") || relative == "scripts/tests" || strings.HasPrefix(relative, "scripts/tests/") {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(entry.Name(), ".schema.json") {
			relative, err := filepath.Rel(c.root, path)
			if err != nil {
				return err
			}
			files[filepath.ToSlash(relative)] = true
		}
		return nil
	})
	if err != nil {
		c.toolIssues = append(c.toolIssues, fmt.Sprintf("scan schema files: %v", err))
	}
	return files
}

func (c *checker) checkPair(pair mapping) {
	label := pair.GoDir + "." + pair.GoType + " <-> " + pair.Schema
	fields, ok := c.goFields(pair, label)
	if !ok {
		return
	}
	properties, required, ok := c.schemaFields(pair, label)
	if !ok {
		return
	}

	for name, field := range fields {
		property, exists := properties[name]
		if !exists {
			c.drifts = append(c.drifts, fmt.Sprintf("%s: Go JSON field %q lacks schema property", label, name))
			continue
		}
		if field.OmitEmpty == required[name] {
			c.drifts = append(c.drifts, fmt.Sprintf("%s: %q Go omitempty and schema required disagree", label, name))
		}
		if schemaShape, known := c.propertyShape(property, label, name); known {
			c.compareShape(label, name, field.Shape, schemaShape)
		}
	}
	for name := range properties {
		if _, exists := fields[name]; !exists {
			c.drifts = append(c.drifts, fmt.Sprintf("%s: schema property %q lacks Go JSON field", label, name))
		}
	}
}

func (c *checker) goFields(pair mapping, label string) (map[string]structField, bool) {
	directory := resolvePath(c.root, pair.GoDir)
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, directory, func(info fs.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: parse Go dir: %v", label, err))
		return nil, false
	}
	typeExpressions := map[string]ast.Expr{}
	var target *ast.StructType
	for _, pkg := range packages {
		for _, file := range pkg.Files {
			for _, declaration := range file.Decls {
				general, ok := declaration.(*ast.GenDecl)
				if !ok || general.Tok != token.TYPE {
					continue
				}
				for _, spec := range general.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					typeExpressions[typeSpec.Name.Name] = typeSpec.Type
					if typeSpec.Name.Name == pair.GoType {
						if structType, ok := typeSpec.Type.(*ast.StructType); ok {
							target = structType
						} else {
							c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: Go type is not a struct", label))
						}
					}
				}
			}
		}
	}
	if target == nil {
		if !containsPrefix(c.toolIssues, label+": Go type is not a struct") {
			c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: Go type not found", label))
		}
		return nil, false
	}

	fields := map[string]structField{}
	for _, field := range target.Fields.List {
		if len(field.Names) == 0 {
			c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: anonymous embedded field is unsupported", label))
			continue
		}
		if field.Tag == nil {
			c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: field %s lacks JSON tag", label, field.Names[0].Name))
			continue
		}
		rawTag, err := strconv.Unquote(field.Tag.Value)
		if err != nil {
			c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: invalid struct tag on %s", label, field.Names[0].Name))
			continue
		}
		jsonTag, present := reflect.StructTag(rawTag).Lookup("json")
		if !present {
			c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: field %s lacks JSON tag", label, field.Names[0].Name))
			continue
		}
		parts := strings.Split(jsonTag, ",")
		if parts[0] == "-" {
			continue
		}
		if parts[0] == "" {
			c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: field %s has incomplete JSON tag", label, field.Names[0].Name))
			continue
		}
		fieldShape, known := goShape(field.Type, typeExpressions, map[string]bool{})
		if !known {
			c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: field %s has unsupported Go type", label, parts[0]))
			continue
		}
		if _, exists := fields[parts[0]]; exists {
			c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: duplicate Go JSON field %q", label, parts[0]))
			continue
		}
		fields[parts[0]] = structField{Name: parts[0], Shape: fieldShape, OmitEmpty: contains(parts[1:], "omitempty")}
	}
	return fields, true
}

func (c *checker) schemaFields(pair mapping, label string) (map[string]any, map[string]bool, bool) {
	path := resolvePath(c.root, pair.Schema)
	content, err := os.ReadFile(path)
	if err != nil {
		c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: read schema: %v", label, err))
		return nil, nil, false
	}
	var root map[string]any
	if err := json.Unmarshal(content, &root); err != nil {
		c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: parse schema JSON: %v", label, err))
		return nil, nil, false
	}
	properties, ok := root["properties"].(map[string]any)
	if !ok {
		c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: schema root lacks object properties", label))
		return nil, nil, false
	}
	required := map[string]bool{}
	if rawRequired, exists := root["required"]; exists {
		values, ok := rawRequired.([]any)
		if !ok {
			c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: schema required is not an array", label))
			return nil, nil, false
		}
		for _, value := range values {
			name, ok := value.(string)
			if !ok {
				c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: schema required includes non-string", label))
				return nil, nil, false
			}
			required[name] = true
		}
	}
	return properties, required, true
}

func (c *checker) propertyShape(raw any, label, field string) (shape, bool) {
	property, ok := raw.(map[string]any)
	if !ok {
		c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: property %q is not an object", label, field))
		return shape{}, false
	}
	if _, exists := property["$ref"]; exists {
		c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: property %q uses $ref (v1 unsupported)", label, field))
		return shape{}, false
	}
	for _, keyword := range []string{"anyOf", "oneOf", "allOf", "not", "if", "then", "else"} {
		if _, exists := property[keyword]; exists {
			c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: property %q uses %s composition (v1 unsupported)", label, field, keyword))
			return shape{}, false
		}
	}
	rawType, exists := property["type"]
	if !exists {
		return shape{Kind: "any"}, true
	}
	kind, ok := rawType.(string)
	if !ok {
		c.toolIssues = append(c.toolIssues, fmt.Sprintf("%s: property %q has non-string type", label, field))
		return shape{}, false
	}
	result := shape{Kind: kind}
	if kind == "array" {
		if rawItems, exists := property["items"]; exists {
			item, known := c.propertyShape(rawItems, label, field+"[]")
			if !known {
				return shape{}, false
			}
			result.Elem = &item
		}
	}
	return result, true
}

func (c *checker) compareShape(label, field string, goValue, schemaValue shape) {
	if goValue.Kind == "any" || schemaValue.Kind == "any" {
		return
	}
	if goValue.Kind != schemaValue.Kind {
		c.drifts = append(c.drifts, fmt.Sprintf("%s: %q kind differs: Go=%s schema=%s", label, field, goValue.Kind, schemaValue.Kind))
		return
	}
	if goValue.Kind == "array" && goValue.Elem != nil && schemaValue.Elem != nil {
		c.compareShape(label, field+"[]", *goValue.Elem, *schemaValue.Elem)
	}
}

func goShape(expression ast.Expr, named map[string]ast.Expr, resolving map[string]bool) (shape, bool) {
	switch value := expression.(type) {
	case *ast.Ident:
		switch value.Name {
		case "string":
			return shape{Kind: "string"}, true
		case "bool":
			return shape{Kind: "boolean"}, true
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
			return shape{Kind: "integer"}, true
		case "float32", "float64":
			return shape{Kind: "number"}, true
		}
		if resolving[value.Name] {
			return shape{}, false
		}
		underlying, exists := named[value.Name]
		if !exists {
			return shape{}, false
		}
		resolving[value.Name] = true
		result, ok := goShape(underlying, named, resolving)
		delete(resolving, value.Name)
		return result, ok
	case *ast.ArrayType:
		element, ok := goShape(value.Elt, named, resolving)
		if !ok {
			return shape{}, false
		}
		return shape{Kind: "array", Elem: &element}, true
	case *ast.MapType:
		return shape{Kind: "object"}, true
	case *ast.InterfaceType:
		return shape{Kind: "any"}, true
	case *ast.StarExpr:
		return goShape(value.X, named, resolving)
	case *ast.StructType:
		return shape{Kind: "object"}, true
	default:
		return shape{}, false
	}
}

func relativeToRoot(root, path string) (string, bool) {
	resolved := resolvePath(root, path)
	relative, err := filepath.Rel(root, resolved)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", false
	}
	return relative, true
}

func resolvePath(root, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func containsPrefix(values []string, target string) bool {
	for _, value := range values {
		if strings.HasPrefix(value, target) {
			return true
		}
	}
	return false
}

func (c *checker) printIssues() {
	sort.Strings(c.toolIssues)
	sort.Strings(c.drifts)
	for _, issue := range c.drifts {
		fmt.Println("[DRIFT] " + issue)
	}
	for _, issue := range c.toolIssues {
		fmt.Println("[TOOL_ERROR] " + issue)
	}
}
