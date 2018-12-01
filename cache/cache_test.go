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

package cache_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/build-system-buildpack/cache"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestCache(t *testing.T) {
	spec.Run(t, "Cache", testCache, spec.Report(report.Terminal{}))
}

func testCache(t *testing.T, when spec.G, it spec.S) {

	it("contributes destination if it does not exist", func() {
		f := test.NewBuildFactory(t)

		destination := filepath.Join(f.Home, "destination")

		c, err := cache.NewCache(f.Build, destination)
		if err != nil {
			t.Fatal(err)
		}

		if err := c.Contribute(); err != nil {
			t.Fatal(err)
		}

		fi, err := os.Lstat(destination)
		if err != nil {
			t.Fatal(err)
		}

		if fi.Mode()&os.ModeSymlink != os.ModeSymlink {
			t.Errorf("destination.Mode() = %s, expected symlink", fi.Mode())
		}

		test.BeLayerLike(t, f.Build.Layers.Layer("build-system-cache"), false, true, false)
	})
}
