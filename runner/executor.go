/*
 * Copyright 2018-2019 the original author or authors.
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

package runner

import (
	"os"
	"os/exec"
	"strings"

	"github.com/buildpack/libbuildpack/application"
	"github.com/cloudfoundry/libcfbuildpack/logger"
)

// Executor is an interface to mock out actual execution.
type Executor interface {
	// Execute configures a command and executes it, sending output to stdout and stderr.
	Execute(application application.Application, cmd *exec.Cmd, logger logger.Logger) error

	// Execute configures a command and executes it, collecting output and returning it
	ExecuteWithOutput(application application.Application, cmd *exec.Cmd, logger logger.Logger) ([]byte, error)
}

type defaultExecutor struct{}

func (defaultExecutor) Execute(application application.Application, cmd *exec.Cmd, logger logger.Logger) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = application.Root

	logger.SubsequentLine("Running %s", strings.Join(cmd.Args, " "))
	return cmd.Run()
}

func (defaultExecutor) ExecuteWithOutput(application application.Application, cmd *exec.Cmd, logger logger.Logger) ([]byte, error) {
	cmd.Dir = application.Root

	return cmd.CombinedOutput()
}
