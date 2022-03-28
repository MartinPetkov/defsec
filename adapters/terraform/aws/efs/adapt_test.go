package efs

import (
	"testing"

	"github.com/aquasecurity/defsec/adapters/terraform/testutil"
	"github.com/aquasecurity/defsec/parsers/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aquasecurity/defsec/providers/aws/efs"
)

func Test_adaptFileSystem(t *testing.T) {
	tests := []struct {
		name      string
		terraform string
		expected  efs.FileSystem
	}{
		{
			name: "configured",
			terraform: `
			resource "aws_efs_file_system" "example" {
				name       = "bar"
				encrypted  = true
				kms_key_id = "my_kms_key"
			  }
`,
			expected: efs.FileSystem{
				Metadata:  types.NewTestMetadata(),
				Encrypted: types.Bool(true, types.NewTestMetadata()),
			},
		},
		{
			name: "defaults",
			terraform: `
			resource "aws_efs_file_system" "example" {
			  }
`,
			expected: efs.FileSystem{
				Metadata:  types.NewTestMetadata(),
				Encrypted: types.Bool(false, types.NewTestMetadata()),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modules := testutil.CreateModulesFromSource(test.terraform, ".tf", t)
			adapted := adaptFileSystem(modules.GetBlocks()[0])
			testutil.AssertDefsecEqual(t, test.expected, adapted)
		})
	}
}

func TestLines(t *testing.T) {
	src := `
	resource "aws_efs_file_system" "example" {
		name       = "bar"
		encrypted  = true
		kms_key_id = "my_kms_key"
	  }
	`
	modules := testutil.CreateModulesFromSource(src, ".tf", t)
	adapted := Adapt(modules)

	require.Len(t, adapted.FileSystems, 1)
	fileSystem := adapted.FileSystems[0]

	assert.Equal(t, 2, fileSystem.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 6, fileSystem.GetMetadata().Range().GetEndLine())

	assert.Equal(t, 4, fileSystem.Encrypted.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 4, fileSystem.Encrypted.GetMetadata().Range().GetEndLine())
}
