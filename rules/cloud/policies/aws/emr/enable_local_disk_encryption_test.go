package emr

import (
	"testing"

	misscanTypes "github.com/khulnasoft-lab/misscan/pkg/types"

	"github.com/khulnasoft-lab/misscan/pkg/providers/aws/emr"
	"github.com/khulnasoft-lab/misscan/pkg/scan"
	"github.com/khulnasoft-lab/misscan/pkg/state"
	"github.com/stretchr/testify/assert"
)

func TestEnableLocalDiskEncryption(t *testing.T) {
	tests := []struct {
		name     string
		input    emr.EMR
		expected bool
	}{
		{
			name: "EMR cluster with local-disk encryption disabled",
			input: emr.EMR{
				SecurityConfiguration: []emr.SecurityConfiguration{
					{
						Name: misscanTypes.String("test", misscanTypes.NewTestMetadata()),
						Configuration: misscanTypes.String(`{
							"EncryptionConfiguration": {
							  "AtRestEncryptionConfiguration": {
								"S3EncryptionConfiguration": {
								  "EncryptionMode": "SSE-S3"
								},
								"LocalDiskEncryptionConfiguration": {
								  "EncryptionKeyProviderType": "",
								  "AwsKmsKey": ""
								}
							  },
							  "EnableInTransitEncryption": true,
							  "EnableAtRestEncryption": true
							}
						  }`, misscanTypes.NewTestMetadata()),
					},
				},
			},
			expected: true,
		},
		{
			name: "EMR cluster with local-disk encryption enabled",
			input: emr.EMR{
				SecurityConfiguration: []emr.SecurityConfiguration{
					{
						Name: misscanTypes.String("test", misscanTypes.NewTestMetadata()),
						Configuration: misscanTypes.String(`{
							"EncryptionConfiguration": {
							  "AtRestEncryptionConfiguration": {
								"S3EncryptionConfiguration": {
								  "EncryptionMode": "SSE-S3"
								},
								"LocalDiskEncryptionConfiguration": {
								  "EncryptionKeyProviderType": "AwsKms",
								  "AwsKmsKey": "arn:aws:kms:us-west-2:187416307283:alias/tf_emr_test_key"
								}
							  },
							  "EnableInTransitEncryption": true,
							  "EnableAtRestEncryption": true
							}
						  }`, misscanTypes.NewTestMetadata()),
					},
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testState state.State
			testState.AWS.EMR = test.input
			results := CheckEnableLocalDiskEncryption.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == scan.StatusFailed && result.Rule().LongID() == CheckEnableLocalDiskEncryption.Rule().LongID() {
					found = true
				}
			}
			if test.expected {
				assert.True(t, found, "Rule should have been found")
			} else {
				assert.False(t, found, "Rule should not have been found")
			}
		})
	}
}
