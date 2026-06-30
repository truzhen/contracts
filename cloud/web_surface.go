package cloud

type CloudWebPublishState string

const (
	CloudWebPublishDraft            CloudWebPublishState = "draft"
	CloudWebPublishReleaseCandidate CloudWebPublishState = "release_candidate"
	CloudWebPublishPublished        CloudWebPublishState = "published"
	CloudWebPublishNotReady         CloudWebPublishState = "not_ready"
)

type CloudWebAssetDigest struct {
	AssetRef  string `json:"asset_ref"`
	Algorithm string `json:"algorithm"`
	Digest    string `json:"digest"`
	Path      string `json:"path,omitempty"`
}

type CloudWebRoute struct {
	Route             string   `json:"route"`
	Entry             string   `json:"entry,omitempty"`
	DataSources       []string `json:"data_sources,omitempty"`
	RequiresLogin     bool     `json:"requires_login"`
	RuntimeConfigKeys []string `json:"runtime_config_keys,omitempty"`
	ForbiddenInlines  []string `json:"forbidden_inlines,omitempty"`
	BlockedReason     string   `json:"blocked_reason,omitempty"`
}

type CloudWebSurface struct {
	SurfaceRef   string                `json:"surface_ref"`
	Route        string                `json:"route"`
	Entry        string                `json:"entry,omitempty"`
	OwnerModule  string                `json:"owner_module"`
	PublishState CloudWebPublishState  `json:"publish_state"`
	AssetDigests []CloudWebAssetDigest `json:"asset_digests,omitempty"`
	Routes       []CloudWebRoute       `json:"routes,omitempty"`
}
