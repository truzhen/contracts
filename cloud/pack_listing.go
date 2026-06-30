package cloud

import "time"

type PackPublicationStatus string

const (
	PackPublicationStatusCandidate PackPublicationStatus = "candidate"
	PackPublicationStatusReviewing PackPublicationStatus = "reviewing"
	PackPublicationStatusPublished PackPublicationStatus = "published"
	PackPublicationStatusSuspended PackPublicationStatus = "suspended"
	PackPublicationStatusWithdrawn PackPublicationStatus = "withdrawn"
)

type PackArtifactDigest struct {
	ArtifactRef string `json:"artifact_ref,omitempty"`
	Algorithm   string `json:"algorithm"`
	Digest      string `json:"digest"`
	SizeBytes   int64  `json:"size_bytes,omitempty"`
}

type CloudPackListing struct {
	ListingRef       string                `json:"listing_ref"`
	PackRef          string                `json:"pack_ref"`
	VersionRef       string                `json:"version_ref"`
	AuthorAccountRef CloudAccountRef       `json:"author_account_ref,omitempty"`
	Status           PackPublicationStatus `json:"status"`
	Artifact         PackArtifactDigest    `json:"artifact"`
	DisplayName      string                `json:"display_name,omitempty"`
	Summary          string                `json:"summary,omitempty"`
	PriceCents       int64                 `json:"price_cents,omitempty"`
	Currency         string                `json:"currency,omitempty"`
	PublishedAt      *time.Time            `json:"published_at,omitempty"`
	UpdatedAt        *time.Time            `json:"updated_at,omitempty"`
}
