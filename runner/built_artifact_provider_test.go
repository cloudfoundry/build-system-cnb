/*
 * Copyright 2018-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package runner_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/build-system-cnb/runner"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuiltArtifactProvider(t *testing.T) {
	spec.Run(t, "BuiltArtifactProvider", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)
		})

		it("fails with no files", func() {
			_, err := runner.NewBuiltArtifactProvider("*.[jw]ar").Get(f.Build.Application)

			g.Expect(err).To(gomega.MatchError("unable to find built artifact (executable JAR or WAR) in *.[jw]ar, candidates: []"))
		})

		it("fails with multiple candidates", func() {
			test.CopyFile(t, filepath.Join("testdata", "stub-application.jar"),
				filepath.Join(f.Build.Application.Root, "stub-application.jar"))
			test.CopyFile(t, filepath.Join("testdata", "stub-application.war"),
				filepath.Join(f.Build.Application.Root, "stub-application.war"))
			test.CopyFile(t, filepath.Join("testdata", "stub-executable.jar"),
				filepath.Join(f.Build.Application.Root, "stub-executable.jar"))

			_, err := runner.NewBuiltArtifactProvider("*.[jw]ar").Get(f.Build.Application)

			g.Expect(err).To(gomega.MatchError(
				fmt.Sprintf("unable to find built artifact (executable JAR or WAR) in *.[jw]ar, candidates: [%s %s %s]",
					filepath.Join(f.Build.Application.Root, "stub-application.jar"),
					filepath.Join(f.Build.Application.Root, "stub-application.war"),
					filepath.Join(f.Build.Application.Root, "stub-executable.jar"))))
		})

		it("passes with a single candidate", func() {
			test.CopyFile(t, filepath.Join("testdata", "stub-application.jar"),
				filepath.Join(f.Build.Application.Root, "stub-application.jar"))

			g.Expect(runner.NewBuiltArtifactProvider("*.[jw]ar").Get(f.Build.Application)).
				To(gomega.Equal(filepath.Join(f.Build.Application.Root, "stub-application.jar")))
		})

		it("passes with single executable JAR", func() {
			test.CopyFile(t, filepath.Join("testdata", "stub-application.jar"),
				filepath.Join(f.Build.Application.Root, "stub-application.jar"))
			test.CopyFile(t, filepath.Join("testdata", "stub-executable.jar"),
				filepath.Join(f.Build.Application.Root, "stub-executable.jar"))

			g.Expect(runner.NewBuiltArtifactProvider("*.[jw]ar").Get(f.Build.Application)).
				To(gomega.Equal(filepath.Join(f.Build.Application.Root, "stub-executable.jar")))
		})

		it("passes with single WAR", func() {
			test.CopyFile(t, filepath.Join("testdata", "stub-application.jar"),
				filepath.Join(f.Build.Application.Root, "stub-application.jar"))
			test.CopyFile(t, filepath.Join("testdata", "stub-application.war"),
				filepath.Join(f.Build.Application.Root, "stub-application.war"))

			g.Expect(runner.NewBuiltArtifactProvider("*.[jw]ar").Get(f.Build.Application)).
				To(gomega.Equal(filepath.Join(f.Build.Application.Root, "stub-application.war")))
		})
	}, spec.Report(report.Terminal{}))
}
