package config_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/thazelart/terraform-validator/internal/config"
	"github.com/thazelart/terraform-validator/internal/fs"
	"github.com/thazelart/terraform-validator/internal/utils"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func TestParseArgs(t *testing.T) {
	expectedResult := "/tmp"
	os.Args = []string{"terraform-validator", expectedResult}

	testResult := config.ParseArgs("version")

	if diff := cmp.Diff(expectedResult, testResult); diff != "" {
		t.Errorf("ParseArgs() mismatch (-want +got):\n%s", diff)
	}
}

func TestUnmarshalYAML(t *testing.T) {
	// first test: using custom config from example
	expectedCustomResult := config.DefaultTerraformConfig
	expectedCustomResult.Files = map[string]config.FileConfig{
		"default": {
			Mandatory: false,
			AuthorizedBlocks: []string{
				"variable",
				"output",
				"provider",
				"terraform",
				"resource",
				"module",
				"data",
				"locals",
			},
		},
	}
	expectedCustomResult.EnsureProvidersVersion = false
	expectedCustomResult.EnsureReadmeUpdated = false

	customConfigFile := fs.NewFile("../../examples/custom_config/terraform-validator.yaml")
	var testCustomResult config.TerraformConfig
	err := yaml.Unmarshal(customConfigFile.Content, &testCustomResult)
	utils.EnsureOrFatal(err)

	if diff := cmp.Diff(expectedCustomResult, testCustomResult); diff != "" {
		t.Errorf("GetCustomConfig(custom) mismatch (-want +got):\n%s", diff)
	}

	// second test with the others possibility of custmization
	expectedCustomResult2 := config.DefaultTerraformConfig
	expectedCustomResult2.EnsureTerraformVersion = false
	expectedCustomResult2.BlockPatternName = "foo"

	customConfigContent := []byte("ensure_terraform_version: false\nblock_pattern_name: 'foo'")
	var testCustomResult2 config.TerraformConfig
	err = yaml.Unmarshal(customConfigContent, &testCustomResult2)
	utils.EnsureOrFatal(err)

	if diff := cmp.Diff(expectedCustomResult2, testCustomResult2); diff != "" {
		t.Errorf("GetCustomConfig(custom) mismatch (-want +got):\n%s", diff)
	}
}

func TestGetTerraformConfig(t *testing.T) {
	// First test case: no custom config
	defaultFolder := fs.NewTerraformFolder("../../examples/default_config/")
	expectedDefaultResult := config.DefaultTerraformConfig
	testDefaultResult := config.GetTerraformConfig(defaultFolder)

	if diff := cmp.Diff(expectedDefaultResult, testDefaultResult); diff != "" {
		t.Errorf("GetCustomConfig(default) mismatch (-want +got):\n%s", diff)
	}

	// Second test case: with custom config
	expectedCustomResult := config.DefaultTerraformConfig
	expectedCustomResult.Files = map[string]config.FileConfig{
		"default": {
			Mandatory: false,
			AuthorizedBlocks: []string{
				"variable",
				"output",
				"provider",
				"terraform",
				"resource",
				"module",
				"data",
				"locals",
			},
		},
	}
	expectedCustomResult.EnsureProvidersVersion = false
	expectedCustomResult.EnsureReadmeUpdated = false

	customFolder := fs.NewTerraformFolder("../../examples/custom_config/")
	testCustomResult := config.GetTerraformConfig(customFolder)

	if diff := cmp.Diff(expectedCustomResult, testCustomResult); diff != "" {
		t.Errorf("GetCustomConfig(custom) mismatch (-want +got):\n%s", diff)
	}
}

func TestGenerateGlobalConfig(t *testing.T) {
	// First test case: no custom config
	os.Args = []string{"terraform-validator", "../../examples/default_config/"}

	defaultConfig := config.DefaultTerraformConfig
	defaultFolder := fs.NewTerraformFolder("../../examples/default_config/")
	expectedDefaultConfig := config.GlobalConfig{WorkDir: defaultFolder, TerraformConfig: defaultConfig}
	testDefaultResult := config.GenerateGlobalConfig("dev")

	if diff := cmp.Diff(expectedDefaultConfig, testDefaultResult); diff != "" {
		t.Errorf("GetCustomConfig(default) mismatch (-want +got):\n%s", diff)
	}

	// Second test case: with custom config
	os.Args = []string{"terraform-validator", "../../examples/custom_config/"}

	customConfig := config.DefaultTerraformConfig
	customConfig.Files = map[string]config.FileConfig{
		"default": {
			AuthorizedBlocks: []string{
				"variable",
				"output",
				"provider",
				"terraform",
				"resource",
				"module",
				"data",
				"locals",
			},
		},
	}
	customConfig.EnsureProvidersVersion = false
	customConfig.EnsureReadmeUpdated = false
	customFolder := fs.NewTerraformFolder("../../examples/custom_config/")
	expectedCustomConfig := config.GlobalConfig{WorkDir: customFolder, TerraformConfig: customConfig}
	testCustomResult := config.GenerateGlobalConfig("dev")

	if diff := cmp.Diff(expectedCustomConfig, testCustomResult); diff != "" {
		t.Errorf("GetCustomConfig(custom) mismatch (-want +got):\n%s", diff)
	}
}