package cmd_test

import (
	"os"
	"testing"

	"github.com/nanovms/ops/cmd"
	"github.com/nanovms/ops/config"
	"github.com/nanovms/ops/lepton"
	"github.com/stretchr/testify/assert"
)

func TestMergeMultipleFlags(t *testing.T) {

	buildImageFlagSet := newBuildImageFlagSet()
	buildImageFlagSet.Set("imagename", "build-image")
	buildImageFlagSet.Set("args", "test")
	buildImageFlags := cmd.NewBuildImageCommandFlags(buildImageFlagSet)

	configFileName := "test-" + String(5) + ".json"
	configFile := &config.Config{
		RunConfig: config.RunConfig{
			Imagename: "config-image-name",
		},
	}
	writeConfigToFile(configFile, configFileName)
	defer os.Remove(configFileName)
	configFlagSet := newConfigFlagSet()
	configFlagSet.Set("config", configFileName)
	configFlags := cmd.NewConfigCommandFlags(configFlagSet)

	imagesPath := lepton.GetOpsHome() + "/images"

	t.Run("if config flags are placed before the build image flags imagename overrides the value", func(t *testing.T) {
		container := cmd.NewMergeConfigContainer(configFlags, buildImageFlags)

		config := &config.Config{}

		err := container.Merge(config)

		assert.Nil(t, err)
		assert.Equal(t, config.RunConfig.Imagename, imagesPath+"/build-image.img")
	})

	t.Run("if build image flags are placed before the config image flags imagename overrides the value", func(t *testing.T) {
		container := cmd.NewMergeConfigContainer(buildImageFlags, configFlags)

		config := &config.Config{}

		err := container.Merge(config)

		assert.Nil(t, err)

		assert.Equal(t, config.RunConfig.Imagename, "config-image-name")
	})
}