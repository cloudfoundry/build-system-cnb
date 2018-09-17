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
	"reflect"
	"strings"
	"testing"

	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/build-system-buildpack/gradle"
	"github.com/cloudfoundry/libjavabuildpack"
	"github.com/cloudfoundry/libjavabuildpack/test"
	"github.com/cloudfoundry/openjdk-buildpack"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestGradle(t *testing.T) {
	spec.Run(t, "Gradle", testGradle, spec.Report(report.Terminal{}))
}

func testGradle(t *testing.T, when spec.G, it spec.S) {

	it("contains gradle", func() {
		bp := gradle.BuildPlanContribution()

		actual := bp[gradle.GradleDependency]

		expected := libbuildpack.BuildPlanDependency{}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("BuildPlan[\"maven\"] = %s, expected = %s", actual, expected)
		}
	})

	it("contains jvm-application", func() {
		bp := gradle.BuildPlanContribution()

		actual := bp["jvm-application"] // TODO use constants for jvm-application

		expected := libbuildpack.BuildPlanDependency{}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("BuildPlan[\"jvm-application\"] = %s, expected = %s", actual, expected)
		}
	})

	it("contains openjdk-jdk", func() {
		bp := gradle.BuildPlanContribution()

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
		f.AddDependency(t, gradle.GradleDependency, "stub-gradle.zip")
		f.AddBuildPlan(t, gradle.GradleDependency, libbuildpack.BuildPlanDependency{})

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

	it("contributes maven if gradlew does not exist", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, gradle.GradleDependency, "stub-gradle.zip")
		f.AddBuildPlan(t, gradle.GradleDependency, libbuildpack.BuildPlanDependency{})

		g, _, err := gradle.NewGradle(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		if err := g.Contribute(); err != nil {
			t.Fatal(err)
		}

		layerRoot := filepath.Join(f.Build.Cache.Root, "gradle")
		test.BeFileLike(t, filepath.Join(layerRoot, "fixture-marker"), 0644, "")
	})

	it("does not contribute maven if gradlew does exist", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, gradle.GradleDependency, "stub-gradle.zip")
		f.AddBuildPlan(t, gradle.GradleDependency, libbuildpack.BuildPlanDependency{})

		if err := libjavabuildpack.WriteToFile(strings.NewReader(""), filepath.Join(f.Build.Application.Root, "gradlew"), 0755); err != nil {
			t.Fatal(err)
		}

		g, _, err := gradle.NewGradle(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		if err := g.Contribute(); err != nil {
			t.Fatal(err)
		}

		exist, err := libjavabuildpack.FileExists(filepath.Join(f.Build.Cache.Root, "gradle", "fixture-marker"))
		if err != nil {
			t.Fatal(err)
		}

		if exist {
			t.Errorf("Expected gradle not to be contributed, but was")
		}
	})

}
