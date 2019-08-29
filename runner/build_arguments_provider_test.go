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
	"testing"

	"github.com/cloudfoundry/build-system-cnb/runner"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuildArgumentsProvider(t *testing.T) {
	spec.Run(t, "BuildArgumentsProvider", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("parses value from $BP_BUILD_ARGUMENTS", func() {
			defer test.ReplaceEnv(t, "BP_BUILD_ARGUMENTS", "test configured arguments")()

			p, err := runner.NewBuildArgumentsProvider("test", "default", "arguments")

			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(p.Arguments).To(gomega.Equal([]string{"test", "configured", "arguments"}))

		})

		it("uses default arguments", func() {
			p, err := runner.NewBuildArgumentsProvider("test", "default", "arguments")

			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(p.Arguments).To(gomega.Equal([]string{"test", "default", "arguments"}))
		})
	}, spec.Report(report.Terminal{}))
}
