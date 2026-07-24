package readmodels

import (
	"fmt"

	"github.com/truzhen/contracts/market"
)

// PackHandsRequirementReadModel is the server-produced, client-safe
// explanation of one Pack Hands requirement. It is a projection only: it
// cannot authorize execution, install a provider, or force a provider
// instance selection.
type PackHandsRequirementReadModel struct {
	PackRef                      string   `json:"pack_ref"`
	PackVersionRef               string   `json:"pack_version_ref"`
	RequirementRef               string   `json:"requirement_ref"`
	RequiredCapabilities         []string `json:"required_capabilities"`
	ProviderFamily               string   `json:"provider_family"`
	SoftwareRequirementRefs      []string `json:"software_requirement_refs"`
	SoftwareResolutionLockRefs   []string `json:"software_resolution_lock_refs"`
	ResolutionStatus             string   `json:"resolution_status"`
	BlockedReason                string   `json:"blocked_reason,omitempty"`
	ResolvedProviderResourceRef  string   `json:"resolved_provider_resource_ref,omitempty"`
	ResolvedProviderDisplayLabel string   `json:"resolved_provider_display_label,omitempty"`
	ResolvedControlMethodRef     string   `json:"resolved_control_method_ref,omitempty"`
	ResolvedExecutionMode        string   `json:"resolved_execution_mode,omitempty"`
	Installable                  bool     `json:"installable"`
	InstallCandidateRef          string   `json:"install_candidate_ref,omitempty"`
	RiskClass                    string   `json:"risk_class"`
	Optional                     bool     `json:"optional"`
	LastProbeReceiptRef          string   `json:"last_probe_receipt_ref,omitempty"`
	LastResolutionReceiptRef     string   `json:"last_resolution_receipt_ref,omitempty"`
	Stale                        bool     `json:"stale"`
	ExpiresAt                    string   `json:"expires_at,omitempty"`
}

// ValidatePackHandsRequirementReadModel checks projection integrity without
// treating readiness, installability, or binding as authorization.
func ValidatePackHandsRequirementReadModel(model PackHandsRequirementReadModel) error {
	for field, value := range map[string]string{
		"pack_ref":          model.PackRef,
		"pack_version_ref":  model.PackVersionRef,
		"requirement_ref":   model.RequirementRef,
		"provider_family":   model.ProviderFamily,
		"resolution_status": model.ResolutionStatus,
		"risk_class":        model.RiskClass,
	} {
		if value == "" {
			return fmt.Errorf("pack hands read model %s must not be empty", field)
		}
	}
	if err := validateUniqueRefs("software_requirement_refs", model.SoftwareRequirementRefs); err != nil {
		return err
	}
	if err := validateUniqueRefs("software_resolution_lock_refs", model.SoftwareResolutionLockRefs); err != nil {
		return err
	}
	for index, capability := range model.RequiredCapabilities {
		if capability == "" {
			return fmt.Errorf("required_capabilities[%d] must not be empty", index)
		}
	}
	if !validPackHandsResolutionStatus(model.ResolutionStatus) {
		return fmt.Errorf("pack hands read model has unknown resolution status %q", model.ResolutionStatus)
	}
	if !validPackHandsRiskClass(model.RiskClass) {
		return fmt.Errorf("pack hands read model has unknown risk class %q", model.RiskClass)
	}
	if model.ResolvedExecutionMode != "" && !validPackHandsExecutionMode(model.ResolvedExecutionMode) {
		return fmt.Errorf("pack hands read model has unknown execution mode %q", model.ResolvedExecutionMode)
	}
	if model.Installable && model.InstallCandidateRef == "" {
		return fmt.Errorf("installable pack hands requirement must carry an install candidate ref")
	}
	return nil
}

func validateUniqueRefs(field string, refs []string) error {
	seen := make(map[string]struct{}, len(refs))
	for _, ref := range refs {
		if ref == "" {
			return fmt.Errorf("%s contains an empty ref", field)
		}
		if _, exists := seen[ref]; exists {
			return fmt.Errorf("%s contains duplicate ref %q", field, ref)
		}
		seen[ref] = struct{}{}
	}
	return nil
}

func validPackHandsResolutionStatus(value string) bool {
	switch market.SoftwareResolution(value) {
	case market.SoftwareResolutionReused,
		market.SoftwareResolutionBound,
		market.SoftwareResolutionInstalledIsolated,
		market.SoftwareResolutionCoexist,
		market.SoftwareResolutionInstallRequired,
		market.SoftwareResolutionVersionConflict,
		market.SoftwareResolutionIsolationRequired,
		market.SoftwareResolutionBlocked,
		market.SoftwareResolutionNotReady,
		market.SoftwareResolutionProviderMissing:
		return true
	default:
		return false
	}
}

func validPackHandsRiskClass(value string) bool {
	switch market.RiskClass(value) {
	case market.RiskClassLow, market.RiskClassMedium, market.RiskClassHigh, market.RiskClassCritical:
		return true
	default:
		return false
	}
}

func validPackHandsExecutionMode(value string) bool {
	switch value {
	case "mcp", "cli", "gui":
		return true
	default:
		return false
	}
}
