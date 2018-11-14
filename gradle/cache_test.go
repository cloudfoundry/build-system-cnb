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
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/build-system-buildpack/gradle"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestCache(t *testing.T) {
	spec.Run(t, "Cache", testCache, spec.Report(report.Terminal{}))
}

func testCache(t *testing.T, when spec.G, it spec.S) {

	it("contributes .gradle if it doesn't exist", func() {
		f := test.NewBuildFactory(t)

		home := filepath.Join(f.Build.Application.Root, "home")

		g, err := gradle.NewCache(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		g.Gradle = filepath.Join(home, ".gradle")

		if err := g.Contribute(); err != nil {
			t.Fatal(err)
		}

		fi, err := os.Lstat(filepath.Join(home, ".gradle"))
		if err != nil {
			t.Fatal(err)
		}

		if fi.Mode()&os.ModeSymlink != os.ModeSymlink {
			t.Errorf("$HOME/.gradle.Mode() = %s, expected symlink", fi.Mode())
		}
	})

	it("does not contribute .gradle if it does exist", func() {
		f := test.NewBuildFactory(t)

		home := filepath.Join(f.Build.Application.Root, "home")

		g, err := gradle.NewCache(f.Build);
		if err != nil {
			t.Fatal(err)
		}

		if err := os.MkdirAll(filepath.Join(home, ".gradle"), 0755); err != nil {
			t.Fatal(err)
		}

		if err := g.Contribute(); err != nil {
			t.Fatal(err)
		}

		fi, err := os.Lstat(filepath.Join(home, ".gradle"))
		if err != nil {
			t.Fatal(err)
		}

		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			t.Errorf("$HOME/.gradle.Mode() = %s, expected not symlink", fi.Mode())
		}
	})
}
