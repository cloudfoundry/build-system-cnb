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
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/bouk/monkey"
	"github.com/cloudfoundry/build-system-buildpack/maven"
	"github.com/cloudfoundry/libjavabuildpack/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestCache(t *testing.T) {
	spec.Run(t, "Cache", testCache, spec.Report(report.Terminal{}))
}

func testCache(t *testing.T, when spec.G, it spec.S) {

	it("contributes .m2 if it doesn't exist", func() {
		f := test.NewBuildFactory(t)

		home := test.ScratchDir(t, "home")

		pg := monkey.Patch(user.Current, func() (*user.User, error) {
			return &user.User{HomeDir: home}, nil
		})
		defer pg.Unpatch()

		m := maven.NewCache(f.Build)

		if err := m.Contribute(); err != nil {
			t.Fatal(err)
		}

		fi, err := os.Lstat(filepath.Join(home, ".m2"))
		if err != nil {
			t.Fatal(err)
		}

		if fi.Mode()&os.ModeSymlink != os.ModeSymlink {
			t.Errorf("$HOME/.m2.Mode() = %s, expected symlink", fi.Mode())
		}
	})

	it("does not contribute .m2 if it does exist", func() {
		f := test.NewBuildFactory(t)

		home := test.ScratchDir(t, "home")

		pg := monkey.Patch(user.Current, func() (*user.User, error) {
			return &user.User{HomeDir: home}, nil
		})
		defer pg.Unpatch()

		m := maven.NewCache(f.Build)

		if err := os.MkdirAll(filepath.Join(home, ".m2"), 0755); err != nil {
			t.Fatal(err)
		}

		if err := m.Contribute(); err != nil {
			t.Fatal(err)
		}

		fi, err := os.Lstat(filepath.Join(home, ".m2"))
		if err != nil {
			t.Fatal(err)
		}

		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			t.Errorf("$HOME/.m2.Mode() = %s, expected not symlink", fi.Mode())
		}
	})
}
