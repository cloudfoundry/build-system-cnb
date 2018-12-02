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
	c, err := r.compiledCode()
	if err != nil {
		return err
	}

	if err := r.layer.Contribute(c, func(layer layers.Layer) error {
		if err := r.Executor.Execute(r.application, r.command, r.logger); err != nil {
			return err
		}

		artifact, err := r.builtArtifactProvider(r.application)
		if err != nil {
			return err
		}
		r.logger.Debug("Built artifact: %s", artifact)

		r.logger.Debug("Expanding %s to %s", artifact, r.layer.Root)
		return layers.ExtractZip(artifact, r.layer.Root, 0)
	}, layers.Build, layers.Launch); err != nil {
		return nil
	}

	r.logger.SubsequentLine("Removing source code")
	if err := os.RemoveAll(r.application.Root); err != nil {
		return err
	}

	r.logger.Debug("Linking %s => %s", r.layer.Root, r.application.Root)
	return os.Symlink(r.layer.Root, r.application.Root)
}

// String makes Runner satisfy the Stringer interface.
func (r Runner) String() string {
	return fmt.Sprintf("Runner{ Executor: %v, application: %s, builtArtifactProvider: %v, command: %v, layer:%s, logger: %s }",
		r.Executor, r.application, r.builtArtifactProvider, r.command, r.layer, r.logger)
}

func (r Runner) compiledCode() (compiledCode, error) {
	v, err := r.Executor.ExecuteWithOutput(r.application, exec.Command("javac", "-version"), r.logger)
	if err != nil {
		return compiledCode{}, err
	}

	return compiledCode{strings.TrimSpace(string(v))}, nil
}

type compiledCode struct {
	JavaVersion string `toml:"java-version"`
}

func (c compiledCode) Identity() (string, string) {
	return "Compiled Code", ""
}

// BuildArtifactProvider returns the location of the build artifact.
type BuiltArtifactProvider func(application application.Application) (string, error)

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

func NewRunner(build build.Build, builtArtifactProvider BuiltArtifactProvider, cmd *exec.Cmd) Runner {
	return Runner{
		defaultExecutor{},
		build.Application,
		builtArtifactProvider,
		cmd,
		build.Layers.Layer("build-system-application"),
		build.Logger,
	}
}
