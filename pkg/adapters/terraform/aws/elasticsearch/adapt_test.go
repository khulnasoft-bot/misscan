package elasticsearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft-lab/misscan/internal/testutil"
	"github.com/khulnasoft-lab/misscan/pkg/adapters/terraform/tftestutil"
	"github.com/khulnasoft-lab/misscan/pkg/providers/aws/elasticsearch"
	misscanTypes "github.com/khulnasoft-lab/misscan/pkg/types"
)

func Test_adaptDomain(t *testing.T) {
	tests := []struct {
		name      string
		terraform string
		expected  elasticsearch.Elasticsearch
	}{
		{
			name: "configured",
			terraform: `
			resource "aws_elasticsearch_domain" "example" {
				domain_name = "domain-foo"
			  
				node_to_node_encryption {
					enabled = true
				}
	 
				encrypt_at_rest {
					enabled = true
				}

				domain_endpoint_options {
				  enforce_https = true
				  tls_security_policy = "Policy-Min-TLS-1-2-2019-07"
				}

				log_publishing_options {
					cloudwatch_log_group_arn = aws_cloudwatch_log_group.example.arn
					log_type                 = "AUDIT_LOGS"
					enabled                  = true  
				}
			  }
`,
			expected: elasticsearch.Elasticsearch{
				Domains: []elasticsearch.Domain{
					{
						Metadata:   misscanTypes.NewTestMetadata(),
						DomainName: misscanTypes.String("domain-foo", misscanTypes.NewTestMetadata()),
						LogPublishing: elasticsearch.LogPublishing{
							Metadata:     misscanTypes.NewTestMetadata(),
							AuditEnabled: misscanTypes.Bool(true, misscanTypes.NewTestMetadata()),
						},
						TransitEncryption: elasticsearch.TransitEncryption{
							Metadata: misscanTypes.NewTestMetadata(),
							Enabled:  misscanTypes.Bool(true, misscanTypes.NewTestMetadata()),
						},
						AtRestEncryption: elasticsearch.AtRestEncryption{
							Metadata: misscanTypes.NewTestMetadata(),
							Enabled:  misscanTypes.Bool(true, misscanTypes.NewTestMetadata()),
						},
						Endpoint: elasticsearch.Endpoint{
							Metadata:     misscanTypes.NewTestMetadata(),
							EnforceHTTPS: misscanTypes.Bool(true, misscanTypes.NewTestMetadata()),
							TLSPolicy:    misscanTypes.String("Policy-Min-TLS-1-2-2019-07", misscanTypes.NewTestMetadata()),
						},
					},
				},
			},
		},
		{
			name: "defaults",
			terraform: `
			resource "aws_elasticsearch_domain" "example" {
			  }
`,
			expected: elasticsearch.Elasticsearch{
				Domains: []elasticsearch.Domain{
					{
						Metadata:   misscanTypes.NewTestMetadata(),
						DomainName: misscanTypes.String("", misscanTypes.NewTestMetadata()),
						LogPublishing: elasticsearch.LogPublishing{
							Metadata:     misscanTypes.NewTestMetadata(),
							AuditEnabled: misscanTypes.Bool(false, misscanTypes.NewTestMetadata()),
						},
						TransitEncryption: elasticsearch.TransitEncryption{
							Metadata: misscanTypes.NewTestMetadata(),
							Enabled:  misscanTypes.Bool(false, misscanTypes.NewTestMetadata()),
						},
						AtRestEncryption: elasticsearch.AtRestEncryption{
							Metadata: misscanTypes.NewTestMetadata(),
							Enabled:  misscanTypes.Bool(false, misscanTypes.NewTestMetadata()),
						},
						Endpoint: elasticsearch.Endpoint{
							Metadata:     misscanTypes.NewTestMetadata(),
							EnforceHTTPS: misscanTypes.Bool(false, misscanTypes.NewTestMetadata()),
							TLSPolicy:    misscanTypes.String("", misscanTypes.NewTestMetadata()),
						},
					},
				},
			},
		},
		{
			name: "opensearch",
			terraform: `resource "aws_opensearch_domain" "example" {
  domain_name    = "example"

  node_to_node_encryption {
    enabled = true
  }

  encrypt_at_rest {
    enabled = true
  }

  domain_endpoint_options {
	enforce_https = true
	tls_security_policy = "Policy-Min-TLS-1-2-2019-07"
  }

  log_publishing_options {
    cloudwatch_log_group_arn = aws_cloudwatch_log_group.example.arn
    log_type                 = "AUDIT_LOGS"
  }
}
`,
			expected: elasticsearch.Elasticsearch{
				Domains: []elasticsearch.Domain{
					{
						Metadata:   misscanTypes.NewTestMetadata(),
						DomainName: misscanTypes.String("example", misscanTypes.NewTestMetadata()),
						LogPublishing: elasticsearch.LogPublishing{
							Metadata:     misscanTypes.NewTestMetadata(),
							AuditEnabled: misscanTypes.Bool(true, misscanTypes.NewTestMetadata()),
						},
						TransitEncryption: elasticsearch.TransitEncryption{
							Metadata: misscanTypes.NewTestMetadata(),
							Enabled:  misscanTypes.Bool(true, misscanTypes.NewTestMetadata()),
						},
						AtRestEncryption: elasticsearch.AtRestEncryption{
							Metadata: misscanTypes.NewTestMetadata(),
							Enabled:  misscanTypes.Bool(true, misscanTypes.NewTestMetadata()),
						},
						Endpoint: elasticsearch.Endpoint{
							Metadata:     misscanTypes.NewTestMetadata(),
							EnforceHTTPS: misscanTypes.Bool(true, misscanTypes.NewTestMetadata()),
							TLSPolicy:    misscanTypes.String("Policy-Min-TLS-1-2-2019-07", misscanTypes.NewTestMetadata()),
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
			testutil.AssertDefsecEqual(t, test.expected, adapted)
		})
	}
}

func TestLines(t *testing.T) {
	src := `
	resource "aws_elasticsearch_domain" "example" {
		domain_name = "domain-foo"
	  
		node_to_node_encryption {
			enabled = true
		}

		encrypt_at_rest {
			enabled = true
		}

		domain_endpoint_options {
		  enforce_https = true
		  tls_security_policy = "Policy-Min-TLS-1-2-2019-07"
		}

		log_publishing_options {
			cloudwatch_log_group_arn = aws_cloudwatch_log_group.example.arn
			log_type                 = "AUDIT_LOGS"
			enabled                  = true  
		}
	  }`

	modules := tftestutil.CreateModulesFromSource(t, src, ".tf")
	adapted := Adapt(modules)

	require.Len(t, adapted.Domains, 1)
	domain := adapted.Domains[0]

	assert.Equal(t, 2, domain.Metadata.Range().GetStartLine())
	assert.Equal(t, 23, domain.Metadata.Range().GetEndLine())

	assert.Equal(t, 3, domain.DomainName.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 3, domain.DomainName.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 5, domain.TransitEncryption.Metadata.Range().GetStartLine())
	assert.Equal(t, 7, domain.TransitEncryption.Metadata.Range().GetEndLine())

	assert.Equal(t, 6, domain.TransitEncryption.Enabled.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 6, domain.TransitEncryption.Enabled.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 9, domain.AtRestEncryption.Metadata.Range().GetStartLine())
	assert.Equal(t, 11, domain.AtRestEncryption.Metadata.Range().GetEndLine())

	assert.Equal(t, 10, domain.AtRestEncryption.Enabled.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 10, domain.AtRestEncryption.Enabled.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 13, domain.Endpoint.Metadata.Range().GetStartLine())
	assert.Equal(t, 16, domain.Endpoint.Metadata.Range().GetEndLine())

	assert.Equal(t, 14, domain.Endpoint.EnforceHTTPS.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 14, domain.Endpoint.EnforceHTTPS.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 15, domain.Endpoint.TLSPolicy.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 15, domain.Endpoint.TLSPolicy.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 18, domain.LogPublishing.Metadata.Range().GetStartLine())
	assert.Equal(t, 22, domain.LogPublishing.Metadata.Range().GetEndLine())

	assert.Equal(t, 21, domain.LogPublishing.AuditEnabled.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 21, domain.LogPublishing.AuditEnabled.GetMetadata().Range().GetEndLine())
}
