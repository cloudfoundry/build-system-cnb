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

package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/buildpack/libbuildpack/application"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/fatih/color"
)

// Runner represents the behavior of running the build system command to build an application.
type Runner struct {
	// Executor is the function that isolates execution
	Executor Executor

	application           application.Application
	builtArtifactProvider BuiltArtifactProvider
	command               *exec.Cmd
	layer                 layers.Layer
	logger                logger.Logger
}

// Contributes builds the application from source code, expands the built artifact, and symlinks the expanded artifact
// to $APPLICATION_ROOT.
func (r Runner) Contribute() error {
	r.logger.FirstLine("%s application", color.YellowString("Building"))

	if err := r.Executor(r.application, r.command, r.logger); err != nil {
		return err
	}

	artifact, err := r.builtArtifactProvider(r.application)
	if err != nil {
		return err
	}
	r.logger.Debug("Built artifact: %s", artifact)

	if err := os.RemoveAll(r.layer.Root); err != nil {
		return err
	}

	r.logger.Debug("Expanding %s to %s", artifact, r.layer.Root)
	if err := layers.ExtractZip(artifact, r.layer.Root, 0); err != nil {
		return err
	}

	r.logger.SubsequentLine("Removing source code")
	if err := os.RemoveAll(r.application.Root); err != nil {
		return err
	}

	r.logger.Debug("Linking %s => %s", r.layer.Root, r.application.Root)
	if err := os.Symlink(r.layer.Root, r.application.Root); err != nil {
		return err
	}

	return r.layer.WriteMetadata(nil, layers.Build, layers.Launch)
}

// String makes Runner satisfy the Stringer interface.
func (r Runner) String() string {
	return fmt.Sprintf("Runner{ Executor: %v, application: %s, builtArtifactProvider: %v, command: %v, layer:%s, logger: %s }",
		r.Executor, r.application, r.builtArtifactProvider, r.command, r.layer, r.logger)
}

type BuiltArtifactProvider func(application application.Application) (string, error)

type Executor func(application application.Application, cmd *exec.Cmd, logger logger.Logger) error

var defaultExecutor = func(application application.Application, cmd *exec.Cmd, logger logger.Logger) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = application.Root

	logger.SubsequentLine("Running %s", strings.Join(cmd.Args, " "))
	return cmd.Run()
}

func NewRunner(build build.Build, builtArtifactProvider BuiltArtifactProvider, cmd *exec.Cmd) Runner {
	return Runner{
		defaultExecutor,
		build.Application,
		builtArtifactProvider,
		cmd,
		build.Layers.Layer("build-system-application"),
		build.Logger,
	}
}
