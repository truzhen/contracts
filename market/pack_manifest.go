package market

import "fmt"

// GatewayClass declares which Truzhen gateway owns a provider-backed
// capability. It is a contract value only; it does not grant execution rights.
type GatewayClass string

const (
	GatewayClassExecution     GatewayClass = "execution"
	GatewayClassCommunication GatewayClass = "communication"
	GatewayClassModel         GatewayClass = "model"
	GatewayClassMemory        GatewayClass = "memory"
)

// RiskClass mirrors the product governance risk colors in a stable JSON value.
type RiskClass string

const (
	RiskClassLow      RiskClass = "low"
	RiskClassMedium   RiskClass = "medium"
	RiskClassHigh     RiskClass = "high"
	RiskClassCritical RiskClass = "critical"
)

// SoftwareFallbackPolicy is the fail-closed behavior when a declared software
// dependency cannot be resolved.
type SoftwareFallbackPolicy string

const (
	SoftwareFallbackBlocked         SoftwareFallbackPolicy = "blocked"
	SoftwareFallbackProviderMissing SoftwareFallbackPolicy = "provider_missing"
	SoftwareFallbackManualHandoff   SoftwareFallbackPolicy = "manual_handoff"
	SoftwareFallbackNotReady        SoftwareFallbackPolicy = "not_ready"
)

// SoftwareIsolationPolicy declares what a resolver may do when versions
// conflict. This is a declaration, not permission to install or execute.
type SoftwareIsolationPolicy string

const (
	SoftwareIsolationReusePreferred      SoftwareIsolationPolicy = "reuse_preferred"
	SoftwareIsolationIsolatedInstall     SoftwareIsolationPolicy = "isolated_install"
	SoftwareIsolationCoexistMultiVersion SoftwareIsolationPolicy = "coexist_multi_version"
	SoftwareIsolationBlocked             SoftwareIsolationPolicy = "blocked"
)

// LicensePolicy is a pack author declaration used by review and resolver UI.
type LicensePolicy string

const (
	LicensePolicyAnyOSI           LicensePolicy = "any_osi"
	LicensePolicyCopyleftExcluded LicensePolicy = "copyleft_excluded"
	LicensePolicyReviewRequired   LicensePolicy = "review_required"
)

// SoftwareResolution is the resolver outcome written into a lock.
type SoftwareResolution string

const (
	SoftwareResolutionReused            SoftwareResolution = "reused"
	SoftwareResolutionBound             SoftwareResolution = "bound"
	SoftwareResolutionInstalledIsolated SoftwareResolution = "installed_isolated"
	SoftwareResolutionCoexist           SoftwareResolution = "coexist"
	SoftwareResolutionInstallRequired   SoftwareResolution = "install_required"
	SoftwareResolutionVersionConflict   SoftwareResolution = "version_conflict"
	SoftwareResolutionIsolationRequired SoftwareResolution = "isolation_required"
	SoftwareResolutionBlocked           SoftwareResolution = "blocked"
	SoftwareResolutionNotReady          SoftwareResolution = "not_ready"
	SoftwareResolutionProviderMissing   SoftwareResolution = "provider_missing"
)

// PackSoftwareRequirement is the canonical software dependency declaration for
// scene/capability/role packs. It names the needed software family and version
// range without binding to a user's local provider instance.
type PackSoftwareRequirement struct {
	RequirementID        string                  `json:"requirement_id"`
	SoftwareFamily       string                  `json:"software_family"`
	ProviderFamily       string                  `json:"provider_family,omitempty"`
	VersionRange         string                  `json:"version_range"`
	AdapterRange         string                  `json:"adapter_range,omitempty"`
	RequiredCapabilities []string                `json:"required_capabilities,omitempty"`
	LicensePolicy        LicensePolicy           `json:"license_policy,omitempty"`
	IsolationPolicy      SoftwareIsolationPolicy `json:"isolation_policy"`
	FallbackPolicy       SoftwareFallbackPolicy  `json:"fallback_policy"`
	Optional             bool                    `json:"optional,omitempty"`
	GatewayClass         GatewayClass            `json:"gateway_class"`
	RiskClass            RiskClass               `json:"risk_class"`
}

// ProviderRequirement is the canonical pack-side provider capability shape.
// Runtime readiness and final provider choice remain owned by truzhenos.
type ProviderRequirement struct {
	RequirementID           string                 `json:"requirement_id"`
	ProviderFamily          string                 `json:"provider_family"`
	GatewayClass            GatewayClass           `json:"gateway_class"`
	RequiredCapabilities    []string               `json:"required_capabilities,omitempty"`
	SoftwareRequirementRefs []string               `json:"software_requirement_refs,omitempty"`
	RiskClass               RiskClass              `json:"risk_class"`
	FallbackPolicy          SoftwareFallbackPolicy `json:"fallback_policy"`
	Optional                bool                   `json:"optional,omitempty"`
}

// PackLifecycleStatus is the unified eight-tier pack lifecycle declaration.
// Values follow the governance dictionary verbatim (Chinese tokens are the
// established convention, see truzhen-packs candidate-set usage); this is an
// author-side declaration only — acceptance/release verdicts stay with Owner.
type PackLifecycleStatus string

const (
	PackLifecycleIdea          PackLifecycleStatus = "想法"
	PackLifecycleDesigning     PackLifecycleStatus = "设计中"
	PackLifecycleContractFixed PackLifecycleStatus = "契约已定"
	PackLifecycleImplemented   PackLifecycleStatus = "已实现"
	PackLifecycleWired         PackLifecycleStatus = "已接线"
	PackLifecycleAccepted      PackLifecycleStatus = "已验收"
	PackLifecycleReleased      PackLifecycleStatus = "已发布"
	PackLifecycleDeprecated    PackLifecycleStatus = "已弃用"
)

// PackRiskType is the structured skeleton of one scene-pack judgment-policy
// risk type (统一决策表 #11 五件套：定义/触发/证据要求/升级路径/回退). A
// declaration never grants authority: escalation_path only *requests* an
// escalation; the ruling stays with Owner + Base Gate, and packs without
// declarations keep today's behavior unchanged (conservative tier).
// escalation_path values mirror spines.RiskEscalationPath ("none"|"owner_gate").
type PackRiskType struct {
	RiskTypeID          string   `json:"risk_type_id"`
	Definition          string   `json:"definition"`
	TriggerActionTypes  []string `json:"trigger_action_types"`
	EvidenceRequirement string   `json:"evidence_requirement,omitempty"`
	EscalationPath      string   `json:"escalation_path"`
	Fallback            string   `json:"fallback,omitempty"`
}

// PackManifest is the canonical cloud-facing manifest shape. It intentionally
// stays descriptive: pack runtime state, local provider resolution and product
// listing state live in their owning repositories.
type PackManifest struct {
	PackID               string                    `json:"pack_id"`
	Name                 string                    `json:"name"`
	Version              string                    `json:"version"`
	Kind                 ProductKind               `json:"kind"`
	MinTruzhenVersion    string                    `json:"min_truzhen_version"`
	Description          string                    `json:"description,omitempty"`
	LifecycleStatus      PackLifecycleStatus       `json:"lifecycle_status,omitempty"`
	RegionCode           string                    `json:"region_code,omitempty"`
	ArchTags             []string                  `json:"arch_tags,omitempty"`
	SoftwareRequirements []PackSoftwareRequirement `json:"software_requirements,omitempty"`
	ProviderRequirements []ProviderRequirement     `json:"provider_requirements,omitempty"`
	RiskTypes            []PackRiskType            `json:"risk_types,omitempty"`
	ExternalSoftwareRefs []string                  `json:"external_software_refs,omitempty"`
}

// SoftwareResolutionLock is produced by truzhenos after resolving a pack
// requirement against the local software registry. It is not authored by packs
// or clients.
type SoftwareResolutionLock struct {
	LockID                   string             `json:"lock_id"`
	PackRef                  string             `json:"pack_ref"`
	RequirementID            string             `json:"requirement_id"`
	SoftwareFamily           string             `json:"software_family"`
	ProviderRequirementRef   string             `json:"provider_requirement_ref,omitempty"`
	ProviderFamily           string             `json:"provider_family,omitempty"`
	ResolvedSoftwareRef      string             `json:"resolved_software_ref,omitempty"`
	ResolvedVersion          string             `json:"resolved_version,omitempty"`
	AdapterVersion           string             `json:"adapter_version,omitempty"`
	ProviderResourceRef      string             `json:"provider_resource_ref,omitempty"`
	ResolvedControlMethodRef string             `json:"resolved_control_method_ref,omitempty"`
	ResolvedExecutionMode    string             `json:"resolved_execution_mode,omitempty"`
	Resolution               SoftwareResolution `json:"resolution"`
	ConflictNote             string             `json:"conflict_note,omitempty"`
	DecisionRef              string             `json:"decision_ref,omitempty"`
	ReceiptRef               string             `json:"receipt_ref,omitempty"`
	ResolvedAt               string             `json:"resolved_at"`
	SourceRegistryRef        string             `json:"source_registry_ref,omitempty"`
	SourceLockFileRef        string             `json:"source_lock_file_ref,omitempty"`
	ResolverVersion          string             `json:"resolver_version,omitempty"`
	IdempotencyKey           string             `json:"idempotency_key,omitempty"`
	TransactionRef           string             `json:"transaction_ref,omitempty"`
	PackVersionRef           string             `json:"pack_version_ref,omitempty"`
}

// ValidateProviderRequirementSoftwareRefs checks that a provider requirement
// references only software requirements from the same fixed Pack version.
// It is a pure contract-boundary check: it does not resolve local software,
// choose a provider instance, install anything, or grant execution authority.
func ValidateProviderRequirementSoftwareRefs(provider ProviderRequirement, softwareRequirements []PackSoftwareRequirement) error {
	softwareByID := make(map[string]PackSoftwareRequirement, len(softwareRequirements))
	for _, software := range softwareRequirements {
		if software.RequirementID == "" {
			return fmt.Errorf("software requirement id must not be empty")
		}
		if _, exists := softwareByID[software.RequirementID]; exists {
			return fmt.Errorf("duplicate software requirement id %q", software.RequirementID)
		}
		softwareByID[software.RequirementID] = software
	}

	seen := make(map[string]struct{}, len(provider.SoftwareRequirementRefs))
	for _, ref := range provider.SoftwareRequirementRefs {
		if ref == "" {
			return fmt.Errorf("provider requirement %q contains an empty software requirement ref", provider.RequirementID)
		}
		if _, exists := seen[ref]; exists {
			return fmt.Errorf("provider requirement %q contains duplicate software requirement ref %q", provider.RequirementID, ref)
		}
		seen[ref] = struct{}{}

		software, exists := softwareByID[ref]
		if !exists {
			return fmt.Errorf("provider requirement %q references software requirement %q outside the fixed Pack version", provider.RequirementID, ref)
		}
		if provider.ProviderFamily != "" && software.ProviderFamily != "" && provider.ProviderFamily != software.ProviderFamily {
			return fmt.Errorf("provider requirement %q family %q conflicts with software requirement %q family %q", provider.RequirementID, provider.ProviderFamily, ref, software.ProviderFamily)
		}
	}
	return nil
}
