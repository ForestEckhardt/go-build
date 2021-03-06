package gobuild

import (
	"path/filepath"

	"github.com/paketo-buildpacks/packit"
)

//go:generate faux --interface ConfigurationParser --output fakes/configuration_parser.go
type ConfigurationParser interface {
	Parse() (BuildConfiguration, error)
}

func Detect(parser ConfigurationParser) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		configuration, err := parser.Parse()
		if err != nil {
			return packit.DetectResult{}, err
		}

		for _, target := range configuration.Targets {
			files, err := filepath.Glob(filepath.Join(target, "*.go"))
			if err != nil {
				return packit.DetectResult{}, err
			}

			if len(files) == 0 {
				return packit.DetectResult{}, packit.Fail
			}
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "go",
						Metadata: map[string]interface{}{
							"build": true,
						},
					},
				},
			},
		}, nil
	}
}
