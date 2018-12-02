/*
 * Copyright 2018, Pivotal Software, Inc. All Rights Reserved.
 * Proprietary and Confidential.
 * Unauthorized use, copying or distribution of this source code via any medium is
 * strictly prohibited without the express written consent of Pivotal Software,
 * Inc.
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
