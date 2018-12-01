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

package maven

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/logger"
)

// Cache represents the location that Maven caches its downloaded artifacts for reuse.
type Cache struct {
	// Maven is the location of the .m2 directory
	Maven string

	layer  layers.Layer
	logger logger.Logger
}

// Contribute links the cache layer to $HOME/.m2.
func (c Cache) Contribute() error {
	exists, err := layers.FileExists(c.Maven)
	if err != nil {
		return err
	}

	if exists {
		c.logger.Debug("Maven cache already exists")
		return nil
	}

	c.logger.SubsequentLine("Linking Maven Cache to %s", c.Maven)

	c.logger.Debug("Creating cache directory %s", c.layer.Root)
	if err := os.MkdirAll(c.layer.Root, 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(c.Maven), 0755); err != nil {
		return err
	}

	c.logger.Debug("Linking %s => %s", c.layer.Root, c.Maven)
	if err := os.Symlink(c.layer.Root, c.Maven); err != nil {
		return err
	}

	return c.layer.WriteMetadata(nil, layers.Cache)
}

// String makes Cache satisfy the Stringer interface.
func (c Cache) String() string {
	return fmt.Sprintf("Cache{ Maven: %s, layer :%s , logger: %s}", c.Maven, c.layer, c.logger)
}

// NewCache creates a new Cache instance.
func NewCache(build build.Build) (Cache, error) {
	m2, err := m2(build.Logger)
	if err != nil {
		return Cache{}, err
	}

	return Cache{
		m2,
		build.Layers.Layer("maven-cache"),
		build.Logger,
	}, nil
}

func m2(logger logger.Logger) (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	m2 := filepath.Join(u.HomeDir, ".m2")
	logger.Debug(".m2 directory: %s", m2)
	return m2, nil
}
