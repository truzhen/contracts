package base

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

// ─────────────────── 安全核心三档授权模式（薄语义层，映射到 OwnerDelegationGrant） ───────────────────
//
// 三档授权只改变 Owner 的确认摩擦，不改变 Base 硬地板。它是一层产品语义，
// 映射到已有的 delegation 真机制，绝不新增平行权限系统：
//   - request_approval（请求批准）：无 grant，所有 gate-required 动作停下回 Owner。
//   - prompt_owner（提我审批）：无 grant，自动整理候选 + 审批卡集中提审；不真实执行 /
//     不真实发送 / 不正式写入。
//   - scoped_full_access（完全访问）：创建 / 存在有边界、可撤销、会过期的
//     OwnerDelegationGrant；high / critical / 付款 / 删除 / credential / 跨边界动作永不可委托，
//     EmergencyStop 启用后全部 active grant suspended。

// AuthorizationMode is the product-facing authorization friction setting.
type AuthorizationMode string

const (
	AuthorizationModeRequestApproval  AuthorizationMode = "request_approval"
	AuthorizationModePromptOwner      AuthorizationMode = "prompt_owner"
	AuthorizationModeScopedFullAccess AuthorizationMode = "scoped_full_access"
)

// ValidAuthorizationMode reports whether a mode value is known.
func ValidAuthorizationMode(m AuthorizationMode) bool {
	switch m {
	case AuthorizationModeRequestApproval, AuthorizationModePromptOwner, AuthorizationModeScopedFullAccess:
		return true
	}
	return false
}

// AuthorizationModeSemantics is the behavioral contract of an authorization
// mode. ExecutesSideEffect is always false: no mode lets a candidate skip the
// gate and run a side effect directly; modes only change confirmation friction.
type AuthorizationModeSemantics struct {
	Mode                    AuthorizationMode `json:"mode"`
	CreatesDelegationGrant  bool              `json:"creates_delegation_grant"`
	AutoStagesCandidate     bool              `json:"auto_stages_candidate"`
	ExecutesSideEffect      bool              `json:"executes_side_effect"`
	GateRequiredGoesToOwner bool              `json:"gate_required_goes_to_owner"`
	Explanation             string            `json:"explanation"`
	// HardDenies lists side effects no mode can ever delegate or auto-approve.
	HardDenies []SideEffectClass `json:"hard_denies"`
}

// AuthorizationHardDenies are the side effects that always fall to Owner + Base
// regardless of authorization mode or delegation grant (in addition to the
// high/critical risk hard floor enforced by DelegationRiskWithinHardFloor).
func AuthorizationHardDenies() []SideEffectClass {
	return []SideEffectClass{
		SideEffectPayment,
		SideEffectDelete,
		SideEffectCredentialAccess,
	}
}

// DeriveAuthorizationModeSemantics returns the behavioral contract of a mode.
// Fail-closed on unknown modes.
func DeriveAuthorizationModeSemantics(mode AuthorizationMode) (AuthorizationModeSemantics, error) {
	sem := AuthorizationModeSemantics{
		Mode:               mode,
		ExecutesSideEffect: false,
		HardDenies:         AuthorizationHardDenies(),
	}
	switch mode {
	case AuthorizationModeRequestApproval:
		sem.CreatesDelegationGrant = false
		sem.AutoStagesCandidate = false
		sem.GateRequiredGoesToOwner = true
		sem.Explanation = "过闸动作停下请你确认；本地文件纳管不打扰"
	case AuthorizationModePromptOwner:
		sem.CreatesDelegationGrant = false
		sem.AutoStagesCandidate = true
		sem.GateRequiredGoesToOwner = true
		sem.Explanation = "先整理候选和证据、生成审批卡集中提审；不真实执行、不真实发送、不正式写入"
	case AuthorizationModeScopedFullAccess:
		sem.CreatesDelegationGrant = true
		sem.AutoStagesCandidate = true
		// 在 grant 边界内可代签以减少低中风险确认；边界外（含 high/critical/
		// payment/delete/credential）仍回 Owner，故此项不为 true。
		sem.GateRequiredGoesToOwner = false
		sem.Explanation = "仅在授权范围内代签，减少低中风险频繁确认；high / critical / 付款 / 删除 / credential 不可委托"
	default:
		return AuthorizationModeSemantics{}, fmt.Errorf("unknown authorization mode %q", mode)
	}
	return sem, nil
}

// DeriveAuthorizationMode reports the effective mode from current delegation
// state: any active grant means scoped_full_access is in effect; otherwise the
// owner's non-grant preference applies. A scoped_full_access preference with no
// active grant degrades fail-closed to prompt_owner (no grant → no delegated
// signing), and any unknown preference falls back to prompt_owner.
func DeriveAuthorizationMode(activeGrantCount int, fallbackPreference AuthorizationMode) (AuthorizationMode, error) {
	if activeGrantCount < 0 {
		return "", errors.New("active grant count must not be negative")
	}
	if activeGrantCount > 0 {
		return AuthorizationModeScopedFullAccess, nil
	}
	if fallbackPreference == AuthorizationModeScopedFullAccess {
		return AuthorizationModePromptOwner, nil
	}
	if !ValidAuthorizationMode(fallbackPreference) {
		return AuthorizationModePromptOwner, nil
	}
	return fallbackPreference, nil
}

// AuthorizationModeChangeCandidateRef is the backend-owned candidate-ref formula
// an authorization-mode change must target. The OwnerDecision on the confirm
// path must point at this ref; a frontend-minted decision_ref will not match and
// is rejected (the frontend never derives decision targets itself).
func AuthorizationModeChangeCandidateRef(targetMode AuthorizationMode) string {
	sum := sha256.Sum256([]byte(string(targetMode) + "\x00authorization_mode_change"))
	return "authorization_mode_change_candidate_" + hex.EncodeToString(sum[:])[:16]
}
