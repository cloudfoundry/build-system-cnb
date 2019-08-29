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

package runner

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpack/libbuildpack/application"
)

// BuiltArtifactProvider returns the artifact built as part of running a build system.
type BuiltArtifactProvider struct {
	target string
}

// Get returns the built artifact if exactly one exists.  If less than or more than one exists, returns an error.
func (b BuiltArtifactProvider) Get(application application.Application) (string, error) {
	candidates, err := filepath.Glob(filepath.Join(application.Root, b.target))
	if err != nil {
		return "", err
	}

	if len(candidates) != 1 {
		return "", fmt.Errorf("unable to find built artifact in %s, candidates: %s", b.target, candidates)
	}

	return candidates[0], nil
}

// NewBuiltArtifactProvider creates a new instance using the default target if not otherwise configured.
func NewBuiltArtifactProvider(defaultTarget ...string) BuiltArtifactProvider {
	if target, ok := os.LookupEnv("BP_BUILT_ARTIFACT"); ok {
		return BuiltArtifactProvider{target}
	}

	target := filepath.Join(defaultTarget...)

	if module, ok := os.LookupEnv("BP_BUILT_MODULE"); ok {
		target = filepath.Join(module, target)
	}

	return BuiltArtifactProvider{target}
}
