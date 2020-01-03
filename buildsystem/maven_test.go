/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package buildsystem_test

import (
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/build-system-cnb/buildsystem"
	"github.com/cloudfoundry/jvm-application-cnb/jvmapplication"
	"github.com/cloudfoundry/libcfbuildpack/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/cloudfoundry/openjdk-cnb/jdk"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestMaven(t *testing.T) {
	spec.Run(t, "Maven", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)
		})

		it("contains maven, jvm-application, and openjdk-jdk in build plan", func() {
			g.Expect(buildsystem.MavenPlan()).To(gomega.Equal(buildplan.Plan{
				Provides: []buildplan.Provided{
					{Name: buildsystem.MavenDependency},
					{Name: jvmapplication.Dependency},
				},
				Requires: []buildplan.Required{
					{Name: buildsystem.MavenDependency},
					{Name: jdk.Dependency},
				},
			}))
		})

		when("Contribute", func() {

			it("contributes maven if mvnw does not exist", func() {
				f.AddDependency(buildsystem.MavenDependency, filepath.Join("testdata", "stub-maven.tar.gz"))
				f.AddPlan(buildpackplan.Plan{Name: buildsystem.MavenDependency})

				b, _, err := buildsystem.NewMavenBuildSystem(f.Build)
				g.Expect(err).NotTo(gomega.HaveOccurred())

				g.Expect(b.Contribute()).To(gomega.Succeed())

				layer := f.Build.Layers.Layer("maven")
				g.Expect(layer).To(test.HaveLayerMetadata(false, true, false))
				g.Expect(filepath.Join(layer.Root, "fixture-marker")).To(gomega.BeARegularFile())
			})

			it("does not contribute maven if mvnw does exist", func() {
				f.AddDependency(buildsystem.MavenDependency, filepath.Join("testdata", "stub-maven.tar.gz"))
				f.AddPlan(buildpackplan.Plan{Name: buildsystem.MavenDependency})

				test.TouchFile(t, f.Build.Application.Root, "mvnw")

				b, _, err := buildsystem.NewMavenBuildSystem(f.Build)
				g.Expect(err).NotTo(gomega.HaveOccurred())

				g.Expect(b.Contribute()).To(gomega.Succeed())

				layer := f.Build.Layers.Layer("maven")
				g.Expect(filepath.Join(layer.Root, "fixture-marker")).NotTo(gomega.BeAnExistingFile())
			})
		})

		when("IsMaven", func() {

			it("returns false if pom.xml does not exist", func() {
				g.Expect(buildsystem.IsMaven(f.Build.Application)).To(gomega.BeFalse())
			})

			it("returns true if pom.xml does exist", func() {
				test.TouchFile(t, f.Build.Application.Root, "pom.xml")

				g.Expect(buildsystem.IsMaven(f.Build.Application)).To(gomega.BeTrue())
			})
		})

		when("NewMavenBuildSystem", func() {

			it("returns true if build plan exists", func() {
				f.AddDependency(buildsystem.MavenDependency, filepath.Join("testdata", "stub-maven.tar.gz"))
				f.AddPlan(buildpackplan.Plan{Name: buildsystem.MavenDependency})

				_, ok, err := buildsystem.NewMavenBuildSystem(f.Build)
				g.Expect(ok).To(gomega.BeTrue())
				g.Expect(err).NotTo(gomega.HaveOccurred())
			})

			it("returns false if build plan does not exist", func() {
				_, ok, err := buildsystem.NewMavenBuildSystem(f.Build)
				g.Expect(ok).To(gomega.BeFalse())
				g.Expect(err).NotTo(gomega.HaveOccurred())
			})
		})
	}, spec.Report(report.Terminal{}))
}
