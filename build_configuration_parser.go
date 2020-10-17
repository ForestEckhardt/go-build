package gobuild

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/buildkite/interpolate"
)

type BuildConfiguration struct {
	Targets    []string
	Flags      []string
	ImportPath string
}

type BuildConfigurationParser struct{}

func NewBuildConfigurationParser() BuildConfigurationParser {
	return BuildConfigurationParser{}
}

func (p BuildConfigurationParser) Parse() (BuildConfiguration, error) {
	var config BuildConfiguration

	config.Targets = []string{"."}

	if targets, ok := os.LookupEnv("BP_GO_BUILD_TARGETS"); ok {
		config.Targets = filepath.SplitList(targets)

		for index, target := range config.Targets {
			if strings.HasPrefix(target, string(filepath.Separator)) {
				return BuildConfiguration{}, fmt.Errorf("failed to determine build targets: %q is an absolute path, targets must be relative to the source directory", target)
			}
			config.Targets[index] = fmt.Sprintf("./%s", filepath.Clean(target))
		}
	}

	if flags, ok := os.LookupEnv("BP_GO_BUILD_FLAGS"); ok {
		interpolatedFlags, err := interpolate.Interpolate(interpolate.NewSliceEnv(os.Environ()), flags)
		if err != nil {
			return BuildConfiguration{}, fmt.Errorf("environment variable expansion failed: %w", err)
		}
		for _, f := range strings.Split(interpolatedFlags, ",") {
			config.Flags = append(config.Flags, splitFlags(f)...)
		}
	}

	config.ImportPath = os.Getenv("BP_GO_BUILD_IMPORT_PATH")

	return config, nil
}

func splitFlags(flag string) []string {
	parts := strings.SplitN(flag, "=", 2)
	if len(parts) == 2 {
		if len(parts[1]) >= 2 {
			if c := parts[1][len(parts[1])-1]; parts[1][0] == c && (c == '"' || c == '\'') {
				parts[1] = parts[1][1 : len(parts[1])-1]
			}
		}
	}

	return parts
}
