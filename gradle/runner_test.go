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
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bouk/monkey"
	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/build-system-buildpack/gradle"
	"github.com/cloudfoundry/libjavabuildpack"
	"github.com/cloudfoundry/libjavabuildpack/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestRunner(t *testing.T) {
	spec.Run(t, "Runner", testRunner, spec.Report(report.Terminal{}))
}

func testRunner(t *testing.T, when spec.G, it spec.S) {

	it("builds application", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, gradle.GradleDependency, "stub-gradle.zip")
		f.AddBuildPlan(t, gradle.GradleDependency, libbuildpack.BuildPlanDependency{})

		if err := libjavabuildpack.WriteToFile(strings.NewReader(""), filepath.Join(f.Build.Application.Root, "gradlew"), 0644); err != nil {
			t.Fatal(err)
		}

		g, _, err := gradle.NewGradle(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		r := gradle.NewRunner(f.Build, g)

		pg := monkey.Patch((*exec.Cmd).Run, func(c *exec.Cmd) error {
			expected := []string{filepath.Join(f.Build.Application.Root, "gradlew"), "-x", "test", "build"}
			if !reflect.DeepEqual(c.Args, expected) {
				t.Errorf("Cmd.Args = %s, expected %s", c.Args, expected)
			}

			return nil
		})
		defer pg.Unpatch()

		source := filepath.Join(test.FindRoot(t), "fixtures", "stub-application.jar")
		destination := filepath.Join(f.Build.Application.Root, "build", "libs", "stub-application.jar")
		if err := libjavabuildpack.CopyFile(source, destination); err != nil {
			t.Fatal(err)
		}

		if err := r.Contribute(); err != nil {
			t.Fatal(err)
		}
	})

	it("removes source code", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, gradle.GradleDependency, "stub-gradle.zip")
		f.AddBuildPlan(t, gradle.GradleDependency, libbuildpack.BuildPlanDependency{})

		if err := libjavabuildpack.WriteToFile(strings.NewReader(""), filepath.Join(f.Build.Application.Root, "gradlew"), 0644); err != nil {
			t.Fatal(err)
		}

		g, _, err := gradle.NewGradle(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		r := gradle.NewRunner(f.Build, g)

		pg := monkey.Patch((*exec.Cmd).Run, func(c *exec.Cmd) error {
			return nil
		})
		defer pg.Unpatch()

		source := filepath.Join(test.FindRoot(t), "fixtures", "stub-application.jar")
		destination := filepath.Join(f.Build.Application.Root, "build", "libs", "stub-application.jar")
		if err := libjavabuildpack.CopyFile(source, destination); err != nil {
			t.Fatal(err)
		}

		if err := r.Contribute(); err != nil {
			t.Fatal(err)
		}

		exists, err := libjavabuildpack.FileExists(filepath.Join(f.Build.Application.Root, "gradlew"))
		if err != nil {
			t.Fatal(err)
		}

		if exists {
			t.Errorf("Expected source code to be removed, but was not")
		}
	})

	it("explodes built application", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, gradle.GradleDependency, "stub-gradle.zip")
		f.AddBuildPlan(t, gradle.GradleDependency, libbuildpack.BuildPlanDependency{})

		if err := libjavabuildpack.WriteToFile(strings.NewReader(""), filepath.Join(f.Build.Application.Root, "gradlew"), 0644); err != nil {
			t.Fatal(err)
		}

		g, _, err := gradle.NewGradle(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		r := gradle.NewRunner(f.Build, g)

		pg := monkey.Patch((*exec.Cmd).Run, func(c *exec.Cmd) error {
			return nil
		})
		defer pg.Unpatch()

		source := filepath.Join(test.FindRoot(t), "fixtures", "stub-application.jar")
		destination := filepath.Join(f.Build.Application.Root, "build", "libs", "stub-application.jar")
		if err := libjavabuildpack.CopyFile(source, destination); err != nil {
			t.Fatal(err)
		}

		if err := r.Contribute(); err != nil {
			t.Fatal(err)
		}

		exists, err := libjavabuildpack.FileExists(filepath.Join(f.Build.Application.Root, "fixture-marker"))
		if err != nil {
			t.Fatal(err)
		}

		if !exists {
			t.Errorf("Expected application to be expanded, but was not")
		}
	})
}
