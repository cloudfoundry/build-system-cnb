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

package main

import (
	"fmt"
	"os"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/build-system-buildpack/gradle"
	"github.com/cloudfoundry/build-system-buildpack/maven"
	buildPkg "github.com/cloudfoundry/libcfbuildpack/build"
)

func main() {
	build, err := buildPkg.DefaultBuild()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize Build: %s\n", err)
		os.Exit(101)
	}

	if code, err := b(build); err != nil {
		build.Logger.Info(err.Error())
		os.Exit(code)
	} else {
		os.Exit(code)
	}
}

func b(build buildPkg.Build) (int, error) {
	build.Logger.FirstLine(build.Logger.PrettyIdentity(build.Buildpack))

	if g, ok, err := gradle.NewGradle(build); err != nil {
		return build.Failure(102), err
	} else if ok {
		if err = g.Contribute(); err != nil {
			return build.Failure(103), err
		}

		if cache, err := gradle.NewCache(build); err != nil {
			return build.Failure(102), err
		} else {
			if err = cache.Contribute(); err != nil {
				return build.Failure(103), err
			}
		}

		if err = gradle.NewRunner(build, g).Contribute(); err != nil {
			return build.Failure(103), err
		}
	}

	if m, ok, err := maven.NewMaven(build); err != nil {
		return build.Failure(102), err
	} else if ok {
		if err = m.Contribute(); err != nil {
			return build.Failure(103), err
		}

		if cache, err := maven.NewCache(build); err != nil {
			return build.Failure(101), err
		} else {
			if err = cache.Contribute(); err != nil {
				return build.Failure(103), err
			}
		}

		if err = maven.NewRunner(build, m).Contribute(); err != nil {
			return build.Failure(103), err
		}
	}

	return build.Success(buildplan.BuildPlan{})
}
