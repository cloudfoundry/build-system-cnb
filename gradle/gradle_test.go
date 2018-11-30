/*
 * Copyright 2018 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gradle_test

import (
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/build-system-buildpack/gradle"
	"github.com/cloudfoundry/jvm-application-buildpack/jvmapplication"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/cloudfoundry/openjdk-buildpack/jdk"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestGradle(t *testing.T) {
	spec.Run(t, "Gradle", testGradle, spec.Report(report.Terminal{}))
}

func testGradle(t *testing.T, when spec.G, it spec.S) {

	it("contains gradle, jvm-application, and openjdk-jdk in build plan", func() {
		test.BeBuildPlanLike(t, gradle.BuildPlanContribution(), buildplan.BuildPlan{
			gradle.Dependency:         buildplan.Dependency{},
			jvmapplication.Dependency: buildplan.Dependency{},
			jdk.Dependency: buildplan.Dependency{
				Version: "1.*",
			},
		})
	})

	when("Contribute", func() {

		it("contributes gradle if gradlew does not exist", func() {
			f := test.NewBuildFactory(t)
			f.AddDependency(t, gradle.Dependency, "stub-gradle.zip")
			f.AddBuildPlan(t, gradle.Dependency, buildplan.Dependency{})

			g, _, err := gradle.NewGradle(f.Build)
			if err != nil {
				t.Fatal(err)
			}

			if err := g.Contribute(); err != nil {
				t.Fatal(err)
			}

			layer := f.Build.Layers.Layer("gradle")
			test.BeLayerLike(t, layer, true, true, false)
			test.BeFileLike(t, filepath.Join(layer.Root, "fixture-marker"), 0644, "")
		})

		it("does not contribute gradle if gradlew does exist", func() {
			f := test.NewBuildFactory(t)
			f.AddDependency(t, gradle.Dependency, "stub-gradle.zip")
			f.AddBuildPlan(t, gradle.Dependency, buildplan.Dependency{})

			test.TouchFile(t, f.Build.Application.Root, "gradlew")

			g, _, err := gradle.NewGradle(f.Build)
			if err != nil {
				t.Fatal(err)
			}

			if err := g.Contribute(); err != nil {
				t.Fatal(err)
			}

			exist, err := layers.FileExists(filepath.Join(f.Build.Layers.Root, "gradle", "fixture-marker"))
			if err != nil {
				t.Fatal(err)
			}

			if exist {
				t.Errorf("Expected gradle not to be contributed, but was")
			}
		})
	})

	when("IsGradle", func() {

		it("returns false if build.gradle does not exist", func() {
			f := test.NewBuildFactory(t)

			actual := gradle.IsGradle(f.Build.Application)
			if actual {
				t.Errorf("Gradle = %t, expected false", actual)
			}
		})

		it("returns true if build.gradle does exist", func() {
			f := test.NewBuildFactory(t)

			test.TouchFile(t, f.Build.Application.Root, "build.gradle")

			actual := gradle.IsGradle(f.Build.Application)
			if !actual {
				t.Errorf("IsGradle = %t, expected true", actual)
			}
		})
	})

	when("NewGradle", func() {

		it("returns true if build plan exists", func() {
			f := test.NewBuildFactory(t)
			f.AddDependency(t, gradle.Dependency, "stub-gradle.zip")
			f.AddBuildPlan(t, gradle.Dependency, buildplan.Dependency{})

			_, ok, err := gradle.NewGradle(f.Build)
			if err != nil {
				t.Fatal(err)
			}
			if !ok {
				t.Errorf("NewGradle = %t, expected true", ok)
			}
		})

		it("returns false if build plan does not exist", func() {
			f := test.NewBuildFactory(t)

			_, ok, err := gradle.NewGradle(f.Build)
			if err != nil {
				t.Fatal(err)
			}
			if ok {
				t.Errorf("NewGradle = %t, expected false", ok)
			}
		})
	})
}
