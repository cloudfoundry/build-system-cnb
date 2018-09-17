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

package gradle

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/libjavabuildpack"
)

// GradleCache represents the location that Gradle caches its downloaded artifacts for reuse.
type GradleCache struct {
	layer  libbuildpack.CacheLayer
	logger libjavabuildpack.Logger
}

// Contribute links the cache layer to $HOME/.gradle.
func (g GradleCache) Contribute() error {
	gradle := g.gradle()

	exists, err := libjavabuildpack.FileExists(gradle)
	if err != nil {
		return err
	}

	if exists {
		g.logger.Debug("Gradle cache already exists")
		return nil
	}

	g.logger.SubsequentLine("Linking Gradle Cache to %s", gradle)

	g.logger.Debug("Creating cache directory %s", g.layer.Root)
	if err := os.MkdirAll(g.layer.Root, 0755); err != nil {
		return err
	}

	g.logger.Debug("Linking %s => %s", g.layer.Root, gradle)
	return os.Symlink(g.layer.Root, gradle)
}

// String makes GradleCache satisfy the Stringer interface.
func (g GradleCache) String() string {
	return fmt.Sprintf("GradleCache{ layer :%s , logger: %s}", g.layer, g.logger)
}

func (g GradleCache) gradle() string {
	return filepath.Join(os.Getenv("HOME"), ".gradle")
}

// NewGradleCache creates a new GradleCache instance.
func NewGradleCache(build libjavabuildpack.Build) GradleCache {
	return GradleCache{
		build.Cache.Layer("gradle-cache"),
		build.Logger,
	}
}
