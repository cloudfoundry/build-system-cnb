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
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/build-system-buildpack/maven"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestRunner(t *testing.T) {
	spec.Run(t, "Runner", testRunner, spec.Report(report.Terminal{}))
}

func testRunner(t *testing.T, when spec.G, it spec.S) {

	it("builds application", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, maven.Dependency, "stub-maven.tar.gz")
		f.AddBuildPlan(t, maven.Dependency, buildplan.Dependency{})

		test.TouchFile(t, f.Build.Application.Root, "mvnw")

		m, _, err := maven.NewMaven(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		r := maven.NewRunner(f.Build, m)
		r.Exec = func(cmd *exec.Cmd) error {
			expected := []string{filepath.Join(f.Build.Application.Root, "mvnw"), "-Dmaven.test.skip=true", "package"}

			if !reflect.DeepEqual(cmd.Args, expected) {
				t.Errorf("Cmd.Args = %s, expected %s", cmd.Args, expected)
			}

			return nil
		}

		source := test.FixturePath(t, "stub-application.jar")
		destination := filepath.Join(f.Build.Application.Root, "target", "stub-application.jar")
		if err := layers.CopyFile(source, destination); err != nil {
			t.Fatal(err)
		}

		if err := r.Contribute(); err != nil {
			t.Fatal(err)
		}
	})

	it("removes source code", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, maven.Dependency, "stub-maven.tar.gz")
		f.AddBuildPlan(t, maven.Dependency, buildplan.Dependency{})

		test.TouchFile(t, f.Build.Application.Root, "mvnw")

		m, _, err := maven.NewMaven(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		r := maven.NewRunner(f.Build, m)
		r.Exec = func(cmd *exec.Cmd) error {
			return nil
		}

		source := test.FixturePath(t, "stub-application.jar")
		destination := filepath.Join(f.Build.Application.Root, "target", "stub-application.jar")
		if err := layers.CopyFile(source, destination); err != nil {
			t.Fatal(err)
		}

		if err := r.Contribute(); err != nil {
			t.Fatal(err)
		}

		exists, err := layers.FileExists(filepath.Join(f.Build.Application.Root, "mvnw"))
		if err != nil {
			t.Fatal(err)
		}

		if exists {
			t.Errorf("Expected source code to be removed, but was not")
		}
	})

	it("explodes built application", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, maven.Dependency, "stub-maven.tar.gz")
		f.AddBuildPlan(t, maven.Dependency, buildplan.Dependency{})

		test.TouchFile(t, f.Build.Application.Root, "mvnw")

		m, _, err := maven.NewMaven(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		r := maven.NewRunner(f.Build, m)
		r.Exec = func(cmd *exec.Cmd) error {
			return nil
		}

		source := test.FixturePath(t, "stub-application.jar")
		destination := filepath.Join(f.Build.Application.Root, "target", "stub-application.jar")
		if err := layers.CopyFile(source, destination); err != nil {
			t.Fatal(err)
		}

		if err := r.Contribute(); err != nil {
			t.Fatal(err)
		}

		exists, err := layers.FileExists(filepath.Join(f.Build.Application.Root, "fixture-marker"))
		if err != nil {
			t.Fatal(err)
		}

		if !exists {
			t.Errorf("Expected application to be expanded, but was not")
		}
	})
}
