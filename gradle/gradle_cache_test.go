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
	"github.com/cloudfoundry/libjavabuildpack/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestGradleCache(t *testing.T) {
	spec.Run(t, "GradleCache", testGradleCache, spec.Report(report.Terminal{}))
}

func testGradleCache(t *testing.T, when spec.G, it spec.S) {
	it("contributes .gradle if it doesn't exist", func() {
		f := test.NewBuildFactory(t)

		home := test.ScratchDir(t, "home")
		defer test.ReplaceEnv(t, "HOME", home)()

		g := gradle.NewGradleCache(f.Build)

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

		home := test.ScratchDir(t, "home")
		defer test.ReplaceEnv(t, "HOME", home)()

		g := gradle.NewGradleCache(f.Build)

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
