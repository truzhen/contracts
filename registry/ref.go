package registry

type RegistryRef struct {
	Namespace string `json:"namespace"`
	ID        string `json:"id"`
}

type SkillRef struct {
	RegistryRef
}
