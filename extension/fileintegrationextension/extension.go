// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package configintegrationextension // import "github.com/elastic/opentelemetry-collector-components/extension/configintegrationextension"

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"

	"github.com/elastic/opentelemetry-collector-components/internal/integrations"
)

type fileTemplateExtension struct {
	config *Config
}

var _ integrations.TemplateFinder = &fileTemplateExtension{}

func newFileTemplateExtension(config *Config) *fileTemplateExtension {
	return &fileTemplateExtension{
		config: config,
	}
}

func (e *fileTemplateExtension) FindTemplate(ctx context.Context, name, version string) (integrations.Template, error) {
	path := filepath.Join(e.config.Path, name+".yml")
	_, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, integrations.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &integrationFile{
		path: path,
	}, nil
}

func (*fileTemplateExtension) Start(context.Context, component.Host) error {
	return nil
}

func (*fileTemplateExtension) Shutdown(context.Context) error {
	return nil
}

type integrationFile struct {
	path string
}

func (t *integrationFile) URI() string {
	return "file:" + t.path
}

func (t *integrationFile) ProviderFactory() confmap.ProviderFactory {
	return fileprovider.NewFactory()
}
