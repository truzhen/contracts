package cloud

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCloudContractJSONShapes(t *testing.T) {
	now := time.Date(2026, 6, 30, 10, 0, 0, 0, time.UTC)
	samples := []struct {
		name string
		v    any
		want []string
	}{
		{
			name: "entitlement",
			v: CloudEntitlement{
				EntitlementRef: "cloud_entitlement://owner/pack",
				OwnerRef:       "owner://demo",
				AccountRef:     CloudAccountRef{Provider: "truzhen-cloud", Subject: "acct_demo"},
				PackRef:        "pack://demo",
				Scope:          EntitlementScopePack,
				Status:         EntitlementStatusActive,
				IssuedAt:       now,
			},
			want: []string{"entitlement_ref", "account_ref", "pack_ref", "status"},
		},
		{
			name: "license",
			v: LicenseToken{
				LicenseRef:       "license://demo",
				EntitlementRef:   "cloud_entitlement://owner/pack",
				DeviceBindingRef: "device://local",
				Status:           LicenseStatusActive,
				IssuedAt:         now,
			},
			want: []string{"license_ref", "entitlement_ref", "device_binding_ref", "status"},
		},
		{
			name: "payment",
			v: PaymentWebhook{
				WebhookRef: "payment_webhook://demo",
				Provider:   PaymentProviderWechat,
				OrderRef:   "payment_order://demo",
				EventID:    "evt_demo",
				ReceivedAt: now,
			},
			want: []string{"webhook_ref", "provider", "order_ref", "event_id"},
		},
		{
			name: "pack listing",
			v: CloudPackListing{
				ListingRef:  "cloud_pack_listing://demo",
				PackRef:     "pack://demo",
				VersionRef:  "pack_version://demo/0.1.0",
				Status:      PackPublicationStatusPublished,
				Artifact:    PackArtifactDigest{Algorithm: "sha256", Digest: strings.Repeat("a", 64)},
				PublishedAt: &now,
			},
			want: []string{"listing_ref", "pack_ref", "version_ref", "artifact", "published_at"},
		},
		{
			name: "session",
			v: CloudSession{
				SessionRef: "cloud_session://demo",
				AccountRef: CloudAccountRef{
					Provider: "truzhen-cloud",
					Subject:  "acct_demo",
				},
				Roles:     []CloudRole{CloudRoleAuthor},
				IssuedAt:  now,
				ExpiresAt: now.Add(time.Hour),
			},
			want: []string{"session_ref", "account_ref", "roles", "expires_at"},
		},
		{
			name: "release",
			v: CloudReleaseReceipt{
				ReleaseRef:     "cloud_release://demo",
				CandidateRef:   "cloud_release_candidate://demo",
				DeployTarget:   CloudDeployTargetStaging,
				ArtifactDigest: "sha256:" + strings.Repeat("b", 64),
				SmokeState:     CloudReleaseSmokePassed,
				CreatedAt:      now,
			},
			want: []string{"release_ref", "candidate_ref", "deploy_target", "artifact_digest", "smoke_state"},
		},
		{
			name: "web surface",
			v: CloudWebSurface{
				SurfaceRef:   "cloud_web_surface://market",
				Route:        "/market.html",
				Entry:        "static/legacy-consoles/pack-portal/market.html",
				OwnerModule:  "05-cloud-web-surfaces",
				PublishState: CloudWebPublishReleaseCandidate,
				AssetDigests: []CloudWebAssetDigest{{AssetRef: "cloud-web-asset-digest://market", Algorithm: "sha256", Digest: strings.Repeat("c", 64)}},
			},
			want: []string{"surface_ref", "route", "entry", "owner_module", "publish_state", "asset_digests"},
		},
	}

	for _, sample := range samples {
		t.Run(sample.name, func(t *testing.T) {
			b, err := json.Marshal(sample.v)
			if err != nil {
				t.Fatal(err)
			}
			got := string(b)
			for _, field := range sample.want {
				if !strings.Contains(got, `"`+field+`"`) {
					t.Fatalf("missing json field %q in %s", field, got)
				}
			}
		})
	}
}

func TestCloudContractSchemasExistAndParse(t *testing.T) {
	names := []string{
		"entitlement.schema.json",
		"license.schema.json",
		"payment.schema.json",
		"pack_listing.schema.json",
		"session.schema.json",
		"release.schema.json",
		"web_surface.schema.json",
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			b, err := os.ReadFile(filepath.Join(".", name))
			if err != nil {
				t.Fatal(err)
			}
			var v map[string]any
			if err := json.Unmarshal(b, &v); err != nil {
				t.Fatal(err)
			}
			if v["$schema"] == "" || v["title"] == "" || v["type"] != "object" {
				t.Fatalf("schema missing required metadata: %#v", v)
			}
		})
	}
}
