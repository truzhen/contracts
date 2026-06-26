package gates

// AccessDecision 代表系统网关自动阻挡或放行的决策
// 轻量布尔决策(仅 allowed+reason),非 Base 完整治理裁定;禁用于正式动作主权裁定(须用 base.GateDecision)。
type AccessDecision struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
}

// OwnerVerdict 代表主人的最终人工裁定
// 轻量人工表态(仅 approved+comment),非 Base 完整 OwnerDecision;正式动作 Owner 裁定须用 base.OwnerDecision。
type OwnerVerdict struct {
	Approved bool   `json:"approved"`
	Comment  string `json:"comment"`
}
