package fakes

import (
	"sync"

	gobuild "github.com/paketo-buildpacks/go-build"
)

type ConfigurationParser struct {
	ParseCall struct {
		sync.Mutex
		CallCount int
		Returns   struct {
			BuildConfiguration gobuild.BuildConfiguration
			Error              error
		}
		Stub func() (gobuild.BuildConfiguration, error)
	}
}

func (f *ConfigurationParser) Parse() (gobuild.BuildConfiguration, error) {
	f.ParseCall.Lock()
	defer f.ParseCall.Unlock()
	f.ParseCall.CallCount++
	if f.ParseCall.Stub != nil {
		return f.ParseCall.Stub()
	}
	return f.ParseCall.Returns.BuildConfiguration, f.ParseCall.Returns.Error
}
