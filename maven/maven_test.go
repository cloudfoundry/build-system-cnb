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
	"reflect"
	"strings"
	"testing"

	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/build-system-buildpack/maven"
	"github.com/cloudfoundry/jvm-application-buildpack"
	"github.com/cloudfoundry/libjavabuildpack"
	"github.com/cloudfoundry/libjavabuildpack/test"
	"github.com/cloudfoundry/openjdk-buildpack"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestMaven(t *testing.T) {
	spec.Run(t, "Maven", testMaven, spec.Report(report.Terminal{}))
}

func testMaven(t *testing.T, when spec.G, it spec.S) {

	it("contains maven", func() {
		bp := maven.BuildPlanContribution()

		actual := bp[maven.MavenDependency]

		expected := libbuildpack.BuildPlanDependency{}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("BuildPlan[\"maven\"] = %s, expected = %s", actual, expected)
		}
	})

	it("contains jvm-application", func() {
		bp := maven.BuildPlanContribution()

		actual := bp[jvm_application_buildpack.JVMApplication]

		expected := libbuildpack.BuildPlanDependency{}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("BuildPlan[\"jvm-application\"] = %s, expected = %s", actual, expected)
		}
	})

	it("contains openjdk-jdk", func() {
		bp := maven.BuildPlanContribution()

		actual := bp[openjdk_buildpack.JDKDependency]

		expected := libbuildpack.BuildPlanDependency{
			Version: "1.*",
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("BuildPlan[\"openjdk-jdk\"] = %s, expected = %s", actual, expected)
		}
	})

	it("returns true if build plan exists", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, maven.MavenDependency, "stub-maven.tar.gz")
		f.AddBuildPlan(t, maven.MavenDependency, libbuildpack.BuildPlanDependency{})

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

	it("contributes maven if mvnw does not exist", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, maven.MavenDependency, "stub-maven.tar.gz")
		f.AddBuildPlan(t, maven.MavenDependency, libbuildpack.BuildPlanDependency{})

		m, _, err := maven.NewMaven(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		if err := m.Contribute(); err != nil {
			t.Fatal(err)
		}

		layerRoot := filepath.Join(f.Build.Cache.Root, "maven")
		test.BeFileLike(t, filepath.Join(layerRoot, "fixture-marker"), 0644, "")
	})

	it("does not contribute maven if mvnw does exist", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, maven.MavenDependency, "stub-maven.tar.gz")
		f.AddBuildPlan(t, maven.MavenDependency, libbuildpack.BuildPlanDependency{})

		if err := libjavabuildpack.WriteToFile(strings.NewReader(""), filepath.Join(f.Build.Application.Root, "mvnw"), 0755); err != nil {
			t.Fatal(err)
		}

		m, _, err := maven.NewMaven(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		if err := m.Contribute(); err != nil {
			t.Fatal(err)
		}

		exist, err := libjavabuildpack.FileExists(filepath.Join(f.Build.Cache.Root, "maven", "fixture-marker"))
		if err != nil {
			t.Fatal(err)
		}

		if exist {
			t.Errorf("Expected mvn not to be contributed, but was")
		}
	})

}
