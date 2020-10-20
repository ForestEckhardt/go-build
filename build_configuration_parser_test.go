package gobuild_test

import (
	"os"
	"testing"

	gobuild "github.com/paketo-buildpacks/go-build"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBuildConfigurationParser(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		parser gobuild.BuildConfigurationParser
	)

	it.Before(func() {
		Expect(os.Setenv("BP_GO_TARGETS", "first:./second")).To(Succeed())
		Expect(os.Setenv("BP_GO_BUILD_FLAGS", `-first=value,-second=value,-third="value",-fourth='value'`)).To(Succeed())
		Expect(os.Setenv("BP_GO_BUILD_IMPORT_PATH", "some-import-path")).To(Succeed())

		parser = gobuild.NewBuildConfigurationParser()
	})

	it.After(func() {
		Expect(os.Unsetenv("BP_GO_TARGETS")).To(Succeed())
		Expect(os.Unsetenv("BP_GO_BUILD_FLAGS")).To(Succeed())
		Expect(os.Unsetenv("BP_GO_BUILD_IMPORT_PATH")).To(Succeed())
	})

	it("parses the targets and flags from a env vars", func() {
		configuration, err := parser.Parse()
		Expect(err).NotTo(HaveOccurred())
		Expect(configuration).To(Equal(gobuild.BuildConfiguration{
			Targets: []string{"./first", "./second"},
			Flags: []string{
				"-first", "value",
				"-second", "value",
				"-third", "value",
				"-fourth", "value",
			},
			ImportPath: "some-import-path",
		}))
	})

	context("when the targets list is empty", func() {
		it.Before(func() {
			Expect(os.Unsetenv("BP_GO_TARGETS")).To(Succeed())
		})

		it("returns a list of targets with . as the only target", func() {
			configuration, err := parser.Parse()
			Expect(err).NotTo(HaveOccurred())
			Expect(configuration.Targets).To(Equal([]string{"."}))
		})
	})

	context("when the build flags reference an env var", func() {
		it.Before(func() {
			Expect(os.Setenv("BP_GO_BUILD_FLAGS", `-first=${SOME_VALUE},-second=$SOME_OTHER_VALUE`)).To(Succeed())

			os.Setenv("SOME_VALUE", "some-value")
			os.Setenv("SOME_OTHER_VALUE", "some-other-value")
		})

		it.After(func() {
			os.Unsetenv("SOME_VALUE")
			os.Unsetenv("SOME_OTHER_VALUE")
		})

		it("replaces the targets list with the values in the env var", func() {
			configuration, err := parser.Parse()
			Expect(err).NotTo(HaveOccurred())
			Expect(configuration.Flags).To(Equal([]string{
				"-first", "some-value",
				"-second", "some-other-value",
			}))
		})
	})

	context("failure cases", func() {

		context("when a the env var expansion fails", func() {
			it.Before(func() {
				Expect(os.Setenv("BP_GO_BUILD_FLAGS", `-first=$& `)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := parser.Parse()
				Expect(err).To(MatchError(ContainSubstring("environment variable expansion failed:")))
			})
		})

		context("when a target is an absolute path", func() {
			it.Before(func() {
				Expect(os.Setenv("BP_GO_TARGETS", "/some-target")).To(Succeed())
			})

			it("returns an error", func() {
				_, err := parser.Parse()
				Expect(err).To(MatchError(ContainSubstring("failed to determine build targets: \"/some-target\" is an absolute path, targets must be relative to the source directory")))
			})
		})
	})
}
