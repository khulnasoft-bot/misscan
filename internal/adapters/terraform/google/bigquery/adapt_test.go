package bigquery

import (
	"testing"

	misscanTypes "github.com/khulnasoft-lab/misscan/pkg/types"

	"github.com/khulnasoft-lab/misscan/pkg/providers/google/bigquery"

	"github.com/khulnasoft-lab/misscan/internal/adapters/terraform/tftestutil"

	"github.com/khulnasoft-lab/misscan/test/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Adapt(t *testing.T) {
	tests := []struct {
		name      string
		terraform string
		expected  bigquery.BigQuery
	}{
		{
			name: "basic",
			terraform: `
			resource "google_bigquery_dataset" "my_dataset" {
				access {
				  role          = "OWNER"
				  special_group = "allAuthenticatedUsers"
				}
			  
				access {
				  role   = "READER"
				  domain = "hashicorp.com"
				}
			  }
`,
			expected: bigquery.BigQuery{
				Datasets: []bigquery.Dataset{
					{
						Metadata: misscanTypes.NewTestMetadata(),
						ID:       misscanTypes.String("", misscanTypes.NewTestMetadata()),
						AccessGrants: []bigquery.AccessGrant{
							{
								Metadata:     misscanTypes.NewTestMetadata(),
								Role:         misscanTypes.String("OWNER", misscanTypes.NewTestMetadata()),
								Domain:       misscanTypes.String("", misscanTypes.NewTestMetadata()),
								SpecialGroup: misscanTypes.String(bigquery.SpecialGroupAllAuthenticatedUsers, misscanTypes.NewTestMetadata()),
							},
							{
								Metadata:     misscanTypes.NewTestMetadata(),
								Role:         misscanTypes.String("READER", misscanTypes.NewTestMetadata()),
								Domain:       misscanTypes.String("hashicorp.com", misscanTypes.NewTestMetadata()),
								SpecialGroup: misscanTypes.String("", misscanTypes.NewTestMetadata()),
							},
						},
					},
				},
			},
		},
		{
			name: "no access blocks",
			terraform: `
			resource "google_bigquery_dataset" "my_dataset" {
				dataset_id                  = "example_dataset"
			  }
`,
			expected: bigquery.BigQuery{
				Datasets: []bigquery.Dataset{
					{
						Metadata: misscanTypes.NewTestMetadata(),
						ID:       misscanTypes.String("example_dataset", misscanTypes.NewTestMetadata()),
					},
				},
			},
		},
		{
			name: "access block without fields",
			terraform: `
			resource "google_bigquery_dataset" "my_dataset" {
				access {
				}
			  }
`,
			expected: bigquery.BigQuery{
				Datasets: []bigquery.Dataset{
					{
						Metadata: misscanTypes.NewTestMetadata(),
						ID:       misscanTypes.String("", misscanTypes.NewTestMetadata()),
						AccessGrants: []bigquery.AccessGrant{
							{
								Metadata:     misscanTypes.NewTestMetadata(),
								Role:         misscanTypes.String("", misscanTypes.NewTestMetadata()),
								Domain:       misscanTypes.String("", misscanTypes.NewTestMetadata()),
								SpecialGroup: misscanTypes.String("", misscanTypes.NewTestMetadata()),
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modules := tftestutil.CreateModulesFromSource(t, test.terraform, ".tf")
			adapted := Adapt(modules)
			testutil.AssertMisscanEqual(t, test.expected, adapted)
		})
	}
}

func TestLines(t *testing.T) {
	src := `
	resource "google_bigquery_dataset" "my_dataset" {
		dataset_id                  = "example_dataset"
		friendly_name               = "test"
		description                 = "This is a test description"
		location                    = "EU"
		default_table_expiration_ms = 3600000
	  
		labels = {
		  env = "default"
		}
	  
		access {
		  role          = "OWNER"
		  special_group = "allAuthenticatedUsers"
		}
	  
		access {
		  role   = "READER"
		  domain = "hashicorp.com"
		}
	}`

	modules := tftestutil.CreateModulesFromSource(t, src, ".tf")
	adapted := Adapt(modules)

	require.Len(t, adapted.Datasets, 1)
	dataset := adapted.Datasets[0]
	require.Len(t, dataset.AccessGrants, 2)

	assert.Equal(t, 14, dataset.AccessGrants[0].Role.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 14, dataset.AccessGrants[0].Role.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 15, dataset.AccessGrants[0].SpecialGroup.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 15, dataset.AccessGrants[0].SpecialGroup.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 19, dataset.AccessGrants[1].Role.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 19, dataset.AccessGrants[1].Role.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 20, dataset.AccessGrants[1].Domain.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 20, dataset.AccessGrants[1].Domain.GetMetadata().Range().GetEndLine())
}
