package eks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft-lab/misscan/internal/testutil"
	"github.com/khulnasoft-lab/misscan/pkg/adapters/terraform/tftestutil"
	"github.com/khulnasoft-lab/misscan/pkg/providers/aws/eks"
	misscanTypes "github.com/khulnasoft-lab/misscan/pkg/types"
)

func Test_adaptCluster(t *testing.T) {
	tests := []struct {
		name      string
		terraform string
		expected  eks.Cluster
	}{
		{
			name: "configured",
			terraform: `
			resource "aws_eks_cluster" "example" {
				encryption_config {
					resources = [ "secrets" ]
					provider {
						key_arn = "key-arn"
					}
				}
			
				enabled_cluster_log_types = ["api", "authenticator", "audit", "scheduler", "controllerManager"]
			
				name = "good_example_cluster"
				role_arn = var.cluster_arn
				vpc_config {
					endpoint_public_access = false
					public_access_cidrs = ["10.2.0.0/8"]
				}
			}
`,
			expected: eks.Cluster{
				Metadata: misscanTypes.NewTestMetadata(),
				Logging: eks.Logging{
					Metadata:          misscanTypes.NewTestMetadata(),
					API:               misscanTypes.Bool(true, misscanTypes.NewTestMetadata()),
					Authenticator:     misscanTypes.Bool(true, misscanTypes.NewTestMetadata()),
					Audit:             misscanTypes.Bool(true, misscanTypes.NewTestMetadata()),
					Scheduler:         misscanTypes.Bool(true, misscanTypes.NewTestMetadata()),
					ControllerManager: misscanTypes.Bool(true, misscanTypes.NewTestMetadata()),
				},
				Encryption: eks.Encryption{
					Metadata: misscanTypes.NewTestMetadata(),
					Secrets:  misscanTypes.Bool(true, misscanTypes.NewTestMetadata()),
					KMSKeyID: misscanTypes.String("key-arn", misscanTypes.NewTestMetadata()),
				},
				PublicAccessEnabled: misscanTypes.Bool(false, misscanTypes.NewTestMetadata()),
				PublicAccessCIDRs: []misscanTypes.StringValue{
					misscanTypes.String("10.2.0.0/8", misscanTypes.NewTestMetadata()),
				},
			},
		},
		{
			name: "defaults",
			terraform: `
			resource "aws_eks_cluster" "example" {
			}
`,
			expected: eks.Cluster{
				Metadata: misscanTypes.NewTestMetadata(),
				Logging: eks.Logging{
					Metadata:          misscanTypes.NewTestMetadata(),
					API:               misscanTypes.Bool(false, misscanTypes.NewTestMetadata()),
					Authenticator:     misscanTypes.Bool(false, misscanTypes.NewTestMetadata()),
					Audit:             misscanTypes.Bool(false, misscanTypes.NewTestMetadata()),
					Scheduler:         misscanTypes.Bool(false, misscanTypes.NewTestMetadata()),
					ControllerManager: misscanTypes.Bool(false, misscanTypes.NewTestMetadata()),
				},
				Encryption: eks.Encryption{
					Metadata: misscanTypes.NewTestMetadata(),
					Secrets:  misscanTypes.Bool(false, misscanTypes.NewTestMetadata()),
					KMSKeyID: misscanTypes.String("", misscanTypes.NewTestMetadata()),
				},
				PublicAccessEnabled: misscanTypes.Bool(true, misscanTypes.NewTestMetadata()),
				PublicAccessCIDRs:   nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modules := tftestutil.CreateModulesFromSource(t, test.terraform, ".tf")
			adapted := adaptCluster(modules.GetBlocks()[0])
			testutil.AssertDefsecEqual(t, test.expected, adapted)
		})
	}
}

func TestLines(t *testing.T) {
	src := `
	resource "aws_eks_cluster" "example" {
		encryption_config {
			resources = [ "secrets" ]
			provider {
				key_arn = "key-arn"
			}
		}
	
		enabled_cluster_log_types = ["api", "authenticator", "audit", "scheduler", "controllerManager"]
	
		name = "good_example_cluster"
		role_arn = var.cluster_arn
		vpc_config {
			endpoint_public_access = false
			public_access_cidrs = ["10.2.0.0/8"]
		}
	}`

	modules := tftestutil.CreateModulesFromSource(t, src, ".tf")
	adapted := Adapt(modules)

	require.Len(t, adapted.Clusters, 1)
	cluster := adapted.Clusters[0]

	assert.Equal(t, 2, cluster.Metadata.Range().GetStartLine())
	assert.Equal(t, 18, cluster.Metadata.Range().GetEndLine())

	assert.Equal(t, 3, cluster.Encryption.Metadata.Range().GetStartLine())
	assert.Equal(t, 8, cluster.Encryption.Metadata.Range().GetEndLine())

	assert.Equal(t, 4, cluster.Encryption.Secrets.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 4, cluster.Encryption.Secrets.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 6, cluster.Encryption.KMSKeyID.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 6, cluster.Encryption.KMSKeyID.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 10, cluster.Logging.Metadata.Range().GetStartLine())
	assert.Equal(t, 10, cluster.Logging.Metadata.Range().GetEndLine())

	assert.Equal(t, 10, cluster.Logging.API.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 10, cluster.Logging.API.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 10, cluster.Logging.Audit.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 10, cluster.Logging.Audit.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 10, cluster.Logging.Authenticator.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 10, cluster.Logging.Authenticator.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 10, cluster.Logging.Scheduler.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 10, cluster.Logging.Scheduler.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 10, cluster.Logging.ControllerManager.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 10, cluster.Logging.ControllerManager.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 15, cluster.PublicAccessEnabled.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 15, cluster.PublicAccessEnabled.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 16, cluster.PublicAccessCIDRs[0].GetMetadata().Range().GetStartLine())
	assert.Equal(t, 16, cluster.PublicAccessCIDRs[0].GetMetadata().Range().GetEndLine())

}
