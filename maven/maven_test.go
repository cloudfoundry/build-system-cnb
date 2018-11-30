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

package maven_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/build-system-buildpack/maven"
	"github.com/cloudfoundry/jvm-application-buildpack/jvmapplication"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/cloudfoundry/openjdk-buildpack/jdk"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestMaven(t *testing.T) {
	spec.Run(t, "Maven", testMaven, spec.Report(report.Terminal{}))
}

func testMaven(t *testing.T, when spec.G, it spec.S) {

	when("BuildPlan Contribution", func() {

		it("contains maven", func() {
			_, ok := maven.BuildPlanContribution()[maven.Dependency]

			if !ok {
				t.Errorf("BuildPlan[\"maven\"] = %t, expected to exist", ok)
			}
		})

		it("contains jvm-application", func() {
			_, ok := maven.BuildPlanContribution()[jvmapplication.Dependency]

			if !ok {
				t.Errorf("BuildPlan[\"jvm-application\"] = %t, expected to exist", ok)
			}
		})

		it("contains openjdk-jdk", func() {
			_, ok := maven.BuildPlanContribution()[jdk.Dependency]

			if !ok {
				t.Errorf("BuildPlan[\"openjdk-jdk\"] = %t, expected to exist", ok)
			}
		})
	})

	when("Contribute", func() {

		it("contributes maven if mvnw does not exist", func() {
			f := test.NewBuildFactory(t)
			f.AddDependency(t, maven.Dependency, "stub-maven.tar.gz")
			f.AddBuildPlan(t, maven.Dependency, buildplan.Dependency{})

			m, _, err := maven.NewMaven(f.Build)
			if err != nil {
				t.Fatal(err)
			}

			if err := m.Contribute(); err != nil {
				t.Fatal(err)
			}

			layerRoot := filepath.Join(f.Build.Layers.Root, "maven")
			test.BeFileLike(t, filepath.Join(layerRoot, "fixture-marker"), 0644, "")
		})

		it("does not contribute maven if mvnw does exist", func() {
			f := test.NewBuildFactory(t)
			f.AddDependency(t, maven.Dependency, "stub-maven.tar.gz")
			f.AddBuildPlan(t, maven.Dependency, buildplan.Dependency{})

			if err := layers.WriteToFile(strings.NewReader(""), filepath.Join(f.Build.Application.Root, "mvnw"), 0755); err != nil {
				t.Fatal(err)
			}

			m, _, err := maven.NewMaven(f.Build)
			if err != nil {
				t.Fatal(err)
			}

			if err := m.Contribute(); err != nil {
				t.Fatal(err)
			}

			exist, err := layers.FileExists(filepath.Join(f.Build.Layers.Root, "maven", "fixture-marker"))
			if err != nil {
				t.Fatal(err)
			}

			if exist {
				t.Errorf("Expected mvn not to be contributed, but was")
			}
		})
	})

	when("IsMaven", func() {

		it("returns false if pom.xml does not exist", func() {
			f := test.NewBuildFactory(t)

			actual := maven.IsMaven(f.Build.Application)
			if actual {
				t.Errorf("IsMaven = %t, expected false", actual)
			}
		})

		it("returns true if pom.xml does exist", func() {
			f := test.NewBuildFactory(t)

			if err := layers.WriteToFile(strings.NewReader(""), filepath.Join(f.Build.Application.Root, "pom.xml"), 0644); err != nil {
				t.Fatal(err)
			}

			actual := maven.IsMaven(f.Build.Application)
			if !actual {
				t.Errorf("IsMaven = %t, expected true", actual)
			}
		})
	})

	when("NewMaven", func() {

		it("returns true if build plan exists", func() {
			f := test.NewBuildFactory(t)
			f.AddDependency(t, maven.Dependency, "stub-maven.tar.gz")
			f.AddBuildPlan(t, maven.Dependency, buildplan.Dependency{})

			_, ok, err := maven.NewMaven(f.Build)
			if err != nil {
				t.Fatal(err)
			}
			if !ok {
				t.Errorf("NewMaven = %t, expected true", ok)
			}
		})

		it("returns false if build plan does not exist", func() {
			f := test.NewBuildFactory(t)

			_, ok, err := maven.NewMaven(f.Build)
			if err != nil {
				t.Fatal(err)
			}
			if ok {
				t.Errorf("NewMaven = %t, expected false", ok)
			}
		})
	})
}
