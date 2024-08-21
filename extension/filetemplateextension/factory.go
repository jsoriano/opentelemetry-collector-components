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

package filetemplateextension // import "github.com/elastic/opentelemetry-collector-components/extension/filetemplateextension"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

var (
	extensionType = component.MustNewType("file_templates")
)

const (
	extensionStability = component.StabilityLevelDevelopment
)

// NewFactory creates a factory for ack extension.
func NewFactory() extension.Factory {
	return extension.NewFactory(
		extensionType,
		createDefaultConfig,
		createExtension,
		extensionStability,
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		Path: "templates",
	}
}

func createExtension(_ context.Context, _ extension.Settings, cfg component.Config) (extension.Extension, error) {
	return newFileTemplateExtension(cfg.(*Config)), nil
}