package base

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// ─────────────────── Artifact 留痕 / 过闸契约（命根：文件进柜不审批、出柜办事才审批） ───────────────────
//
// 红线（V3_GOVERNANCE「文件纳管区分留痕 vs 过闸」）：本地文件入项目 / 入柜 /
// 入工作区 / 入证据候选库**只自动留痕**——生成 ArtifactBinding + ArtifactTraceRecord，
// trace-only、不走 Owner / Base Gate、不要求 OwnerDecision。只有当文件被**用于对外动作**
// （外发附件 / 外部系统提交 / 对外发布 / 真实执行 / 转正式记忆）时，那个 ArtifactUseIntent
// 才经 ArtifactUseGateRequirement 判定过对应 Gate。类型层把"留痕"与"过闸"彻底分开，
// 杜绝把本地纳管误判成需要审批的副作用，也杜绝把留痕当正式对象绕过主权链。

// ArtifactBinding 表达一次本地文件纳管（trace-only）。它**故意不含** OwnerDecisionRef
// 等主权裁定字段：本地纳管在类型层就不携带、也不需要任何 Owner / Base 裁定。
type ArtifactBinding struct {
	BindingID        string    `json:"binding_id"`
	OwnerID          string    `json:"owner_id"`
	SourceRef        string    `json:"source_ref"` // 本地来源（项目 / 柜子 / 工作区 / 证据候选库）
	ContentHash      string    `json:"content_hash"`
	ContentSize      int64     `json:"content_size,omitempty"`
	Filename         string    `json:"filename,omitempty"`
	MediaType        string    `json:"media_type,omitempty"`
	TransactionRef   string    `json:"transaction_ref,omitempty"` // 可选挂事务主线
	Scope            string    `json:"scope,omitempty"`
	ArtifactTraceRef string    `json:"artifact_trace_ref"` // 指向本次留痕记录
	CreatedAt        time.Time `json:"created_at,omitempty"`
}

// ArtifactTraceRecord 是一条轻量留痕记录（append-only trace）。它本身就是 trace-only，
// 永远不要求过闸，可派生出 evidence_ref 供后续证据链引用。
type ArtifactTraceRecord struct {
	TraceRef    string    `json:"trace_ref"`
	BindingID   string    `json:"binding_id"`
	ContentHash string    `json:"content_hash"`
	SourceRef   string    `json:"source_ref"`
	EvidenceRef string    `json:"evidence_ref,omitempty"`
	RecordedAt  time.Time `json:"recorded_at,omitempty"`
}

// ArtifactBindingInput is the local-file metadata used to build a trace-only
// binding. No OwnerDecision is part of the input by design.
type ArtifactBindingInput struct {
	OwnerID        string
	SourceRef      string
	ContentHash    string
	ContentSize    int64
	Filename       string
	MediaType      string
	TransactionRef string
	Scope          string
}

// NewArtifactBindingTrace builds a trace-only ArtifactBinding plus its
// ArtifactTraceRecord from local file metadata. It never requires an
// OwnerDecision: local intake is trace-only by construction. Refs are
// content-addressed (deterministic), timestamps are left to the caller / store.
func NewArtifactBindingTrace(in ArtifactBindingInput) (ArtifactBinding, ArtifactTraceRecord, error) {
	if in.OwnerID == "" {
		return ArtifactBinding{}, ArtifactTraceRecord{}, errors.New("artifact binding owner_id is required")
	}
	if in.SourceRef == "" {
		return ArtifactBinding{}, ArtifactTraceRecord{}, errors.New("artifact binding source_ref is required")
	}
	if in.ContentHash == "" {
		return ArtifactBinding{}, ArtifactTraceRecord{}, errors.New("artifact binding content_hash is required")
	}
	if in.ContentSize < 0 {
		return ArtifactBinding{}, ArtifactTraceRecord{}, errors.New("artifact binding content_size must not be negative")
	}
	bindingID := artifactRef("artifact_binding", in.OwnerID, in.ContentHash, in.SourceRef)
	traceRef := artifactRef("artifact_trace", bindingID, in.ContentHash)
	binding := ArtifactBinding{
		BindingID:        bindingID,
		OwnerID:          in.OwnerID,
		SourceRef:        in.SourceRef,
		ContentHash:      in.ContentHash,
		ContentSize:      in.ContentSize,
		Filename:         in.Filename,
		MediaType:        in.MediaType,
		TransactionRef:   in.TransactionRef,
		Scope:            in.Scope,
		ArtifactTraceRef: traceRef,
	}
	trace := ArtifactTraceRecord{
		TraceRef:    traceRef,
		BindingID:   bindingID,
		ContentHash: in.ContentHash,
		SourceRef:   in.SourceRef,
		EvidenceRef: artifactRef("evidence", traceRef, in.ContentHash),
	}
	return binding, trace, nil
}

func artifactRef(prefix string, parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	return prefix + "_" + hex.EncodeToString(h.Sum(nil))[:16]
}

// ValidateArtifactBinding enforces the trace-only binding shape.
func ValidateArtifactBinding(b *ArtifactBinding) error {
	if b == nil {
		return errors.New("artifact binding is required")
	}
	if b.BindingID == "" {
		return errors.New("artifact binding binding_id is required")
	}
	if b.OwnerID == "" {
		return errors.New("artifact binding owner_id is required")
	}
	if b.SourceRef == "" {
		return errors.New("artifact binding source_ref is required")
	}
	if b.ContentHash == "" {
		return errors.New("artifact binding content_hash is required")
	}
	if b.ContentSize < 0 {
		return errors.New("artifact binding content_size must not be negative")
	}
	if b.ArtifactTraceRef == "" {
		return errors.New("artifact binding artifact_trace_ref is required")
	}
	return nil
}

// ValidateArtifactTraceRecord enforces the trace record shape.
func ValidateArtifactTraceRecord(r *ArtifactTraceRecord) error {
	if r == nil {
		return errors.New("artifact trace record is required")
	}
	if r.TraceRef == "" {
		return errors.New("artifact trace record trace_ref is required")
	}
	if r.BindingID == "" {
		return errors.New("artifact trace record binding_id is required")
	}
	if r.ContentHash == "" {
		return errors.New("artifact trace record content_hash is required")
	}
	return nil
}

// ArtifactUseTarget describes what a (already traced) file is being used for.
// It is the single discriminator that decides trace-only vs gated.
type ArtifactUseTarget string

const (
	// trace-only / local uses — never gated.
	ArtifactUseLocalIntake    ArtifactUseTarget = "local_intake"
	ArtifactUseLocalReference ArtifactUseTarget = "local_reference"
	// gated uses — the file is leaving the cabinet to do work.
	ArtifactUseExternalSendAttachment   ArtifactUseTarget = "external_send_attachment"
	ArtifactUseExternalSystemSubmission ArtifactUseTarget = "external_system_submission"
	ArtifactUseExternalPublish          ArtifactUseTarget = "external_publish"
	ArtifactUseRealExecution            ArtifactUseTarget = "real_execution"
	ArtifactUseFormalMemoryWrite        ArtifactUseTarget = "formal_memory_write"
	ArtifactUseFormalKnowledgeWrite     ArtifactUseTarget = "formal_knowledge_write"
	ArtifactUseFormalTaskWrite          ArtifactUseTarget = "formal_task_write"
)

// ArtifactGateKind names which Base gate an artifact use must pass through.
type ArtifactGateKind string

const (
	ArtifactGateNone               ArtifactGateKind = "none"
	ArtifactGateSend               ArtifactGateKind = "send_gate"
	ArtifactGateExternalSubmission ArtifactGateKind = "external_submission_gate"
	ArtifactGateExecution          ArtifactGateKind = "execution_gate"
	ArtifactGateFormalization      ArtifactGateKind = "formalization_gate"
)

// ArtifactUseIntent is a request to use an already-traced artifact for some
// target. It is built by the module owning the candidate, never carries a
// self-minted decision; gating is derived from the use target.
type ArtifactUseIntent struct {
	UseIntentRef   string            `json:"use_intent_ref,omitempty"`
	BindingID      string            `json:"binding_id"`
	ContentHash    string            `json:"content_hash"`
	UseTarget      ArtifactUseTarget `json:"use_target"`
	TransactionRef string            `json:"transaction_ref,omitempty"`
	EvidenceRef    string            `json:"evidence_ref,omitempty"`
	TargetRef      string            `json:"target_ref,omitempty"`
	RiskClass      RiskClass         `json:"risk_class,omitempty"`
}

// ArtifactUseGateRequirement is the derived gating verdict for a use target.
type ArtifactUseGateRequirement struct {
	UseTarget              ArtifactUseTarget  `json:"use_target"`
	RequiredGate           bool               `json:"required_gate"`
	GateKinds              []ArtifactGateKind `json:"gate_kinds"`
	SideEffectClass        SideEffectClass    `json:"side_effect_class"`
	RequiresTransactionRef bool               `json:"requires_transaction_ref"`
	RequiresEvidenceRef    bool               `json:"requires_evidence_ref"`
	// RequiresBaseIssuer means the action may only be authorized by a
	// Base-issued triple (decision_ref / run_id / nonce); a caller-supplied
	// owner decision is never trusted.
	RequiresBaseIssuer bool   `json:"requires_base_issuer"`
	Reason             string `json:"reason"`
}

// ValidArtifactUseTarget reports whether a use target is known.
func ValidArtifactUseTarget(t ArtifactUseTarget) bool {
	switch t {
	case ArtifactUseLocalIntake, ArtifactUseLocalReference,
		ArtifactUseExternalSendAttachment, ArtifactUseExternalSystemSubmission,
		ArtifactUseExternalPublish, ArtifactUseRealExecution,
		ArtifactUseFormalMemoryWrite, ArtifactUseFormalKnowledgeWrite, ArtifactUseFormalTaskWrite:
		return true
	}
	return false
}

// DeriveArtifactUseGateRequirement maps a use target to its gate requirement.
// Fail-closed: an unknown target is an error, never a silent "no gate".
func DeriveArtifactUseGateRequirement(target ArtifactUseTarget) (ArtifactUseGateRequirement, error) {
	req := ArtifactUseGateRequirement{UseTarget: target}
	switch target {
	case ArtifactUseLocalIntake:
		req.RequiredGate = false
		req.GateKinds = []ArtifactGateKind{ArtifactGateNone}
		req.SideEffectClass = SideEffectLocalFileWrite
		req.Reason = "本地文件纳管是 trace-only，自动留痕不过闸"
	case ArtifactUseLocalReference:
		req.RequiredGate = false
		req.GateKinds = []ArtifactGateKind{ArtifactGateNone}
		req.SideEffectClass = SideEffectReadOnly
		req.Reason = "本地引用文件作上下文是 trace-only，不过闸"
	case ArtifactUseExternalSendAttachment:
		req.RequiredGate = true
		req.GateKinds = []ArtifactGateKind{ArtifactGateSend}
		req.SideEffectClass = SideEffectExternalSend
		req.RequiresEvidenceRef = true
		req.Reason = "文件作为附件外发是对外动作，必须过 SendGate"
	case ArtifactUseExternalPublish:
		req.RequiredGate = true
		req.GateKinds = []ArtifactGateKind{ArtifactGateSend}
		req.SideEffectClass = SideEffectExternalSend
		req.RequiresEvidenceRef = true
		req.Reason = "文件对外发布是对外动作，必须过 SendGate"
	case ArtifactUseExternalSystemSubmission:
		req.RequiredGate = true
		req.GateKinds = []ArtifactGateKind{ArtifactGateExternalSubmission}
		req.SideEffectClass = SideEffectExternalSend
		req.RequiresTransactionRef = true
		req.RequiresEvidenceRef = true
		req.Reason = "文件提交外部系统（执法 / OA / CRM / ERP / 市场 / 云端）必须过闸并绑 transaction_ref + evidence_ref"
	case ArtifactUseRealExecution:
		req.RequiredGate = true
		req.GateKinds = []ArtifactGateKind{ArtifactGateExecution}
		req.SideEffectClass = SideEffectRealExecution
		req.RequiresTransactionRef = true
		req.RequiresEvidenceRef = true
		req.Reason = "文件触发真实执行必须过 ExecutionGate"
	case ArtifactUseFormalMemoryWrite, ArtifactUseFormalKnowledgeWrite, ArtifactUseFormalTaskWrite:
		req.RequiredGate = true
		req.GateKinds = []ArtifactGateKind{ArtifactGateFormalization}
		req.SideEffectClass = SideEffectFormalWrite
		req.RequiresTransactionRef = true
		req.RequiresEvidenceRef = true
		req.RequiresBaseIssuer = true
		req.Reason = "文件转正式对象必须经 Base Formalization（Base 签发 decision_ref / run_id / nonce）"
	default:
		return ArtifactUseGateRequirement{}, fmt.Errorf("unknown artifact use target %q", target)
	}
	return req, nil
}

// ValidateArtifactUseIntent validates the use intent shape and enforces the
// transaction_ref / evidence_ref bindings the derived requirement demands.
func ValidateArtifactUseIntent(intent *ArtifactUseIntent) error {
	if intent == nil {
		return errors.New("artifact use intent is required")
	}
	if intent.BindingID == "" {
		return errors.New("artifact use intent binding_id is required")
	}
	if intent.ContentHash == "" {
		return errors.New("artifact use intent content_hash is required")
	}
	req, err := DeriveArtifactUseGateRequirement(intent.UseTarget)
	if err != nil {
		return err
	}
	if req.RequiresTransactionRef && intent.TransactionRef == "" {
		return fmt.Errorf("artifact use %q requires a transaction_ref", intent.UseTarget)
	}
	if req.RequiresEvidenceRef && intent.EvidenceRef == "" {
		return fmt.Errorf("artifact use %q requires an evidence_ref", intent.UseTarget)
	}
	return nil
}

// ArtifactTraceOnlyTargets are the use targets that stay trace-only (no gate).
func ArtifactTraceOnlyTargets() []ArtifactUseTarget {
	return []ArtifactUseTarget{ArtifactUseLocalIntake, ArtifactUseLocalReference}
}

// ArtifactGatedTargets are the use targets that must pass a Base gate (the file
// is leaving the cabinet to do work).
func ArtifactGatedTargets() []ArtifactUseTarget {
	return []ArtifactUseTarget{
		ArtifactUseExternalSendAttachment, ArtifactUseExternalSystemSubmission,
		ArtifactUseExternalPublish, ArtifactUseRealExecution,
		ArtifactUseFormalMemoryWrite, ArtifactUseFormalKnowledgeWrite, ArtifactUseFormalTaskWrite,
	}
}
