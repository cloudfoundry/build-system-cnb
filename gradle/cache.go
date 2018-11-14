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
	"os/user"
	"path/filepath"

	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/logger"
)

// Cache represents the location that Gradle caches its downloaded artifacts for reuse.
type Cache struct {
	// Gradle is the location of the .gradle directory.
	Gradle string

	layer  layers.Layer
	logger logger.Logger
}

// Contribute links the cache layer to $HOME/.gradle.
func (c Cache) Contribute() error {
	exists, err := layers.FileExists(c.Gradle)
	if err != nil {
		return err
	}

	if exists {
		c.logger.Debug("Gradle cache already exists")
		return nil
	}

	c.logger.SubsequentLine("Linking Gradle Cache to %s", c.Gradle)

	c.logger.Debug("Creating cache directory %s", c.layer.Root)
	if err := os.MkdirAll(c.layer.Root, 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(c.Gradle), 0755); err != nil {
		return err
	}

	c.logger.Debug("Linking %s => %s", c.layer.Root, c.Gradle)
	if err := os.Symlink(c.layer.Root, c.Gradle); err != nil {
		return err
	}

	return c.layer.WriteMetadata(nil, layers.Build, layers.Cache)
}

// String makes Cache satisfy the Stringer interface.
func (c Cache) String() string {
	return fmt.Sprintf("Cache{ Gradle:%s, layer: %s , logger: %s}", c.Gradle, c.layer, c.logger)
}

// NewCache creates a new Cache instance.
func NewCache(build build.Build) (Cache, error) {
	gradle, err := gradle()
	if err != nil {
		return Cache{}, err
	}

	return Cache{
		gradle,
		build.Layers.Layer("gradle-cache"),
		build.Logger,
	}, nil
}

func gradle() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(u.HomeDir, ".gradle"), nil
}
