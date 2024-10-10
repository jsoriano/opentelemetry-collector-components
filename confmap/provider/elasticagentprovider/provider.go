package elasticagentprovider

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/elastic/elastic-agent-client/v7/pkg/client"
	"go.opentelemetry.io/collector/confmap"
	"go.uber.org/zap"
)

const schemeName = "elasticagent"

type provider struct {
	mutex    sync.Mutex
	client   client.V2
	services []client.Service

	logger *zap.Logger
}

// NewFactory returns a factory for a confmap.Provider that gets the configuration through the
// Elastic Agent control protocol.
func NewFactory() confmap.ProviderFactory {
	return confmap.NewProviderFactory(newProvider)
}

func newProvider(settings confmap.ProviderSettings) confmap.Provider {
	return &provider{
		logger: settings.Logger.Named(schemeName),
	}
}

func (p *provider) Retrieve(ctx context.Context, _ string, watcher confmap.WatcherFunc) (*confmap.Retrieved, error) {
	// TODO: Use the uri to setup the reader, defaulting to stdin.
	err := p.ensureInitialized(ctx, os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("could not initialize Elastic Agent provider: %w", err)
	}

	return nil, nil
}

func (p *provider) ensureInitialized(ctx context.Context, r io.Reader) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.client != nil {
		return nil
	}

	versionInfo := client.VersionInfo{
		Name:      "otel-collector-client",
		BuildHash: "unknown",
		Meta: map[string]string{
			"commit":     "unknown",
			"build_time": "unknown",
		},
	}

	client, services, err := client.NewV2FromReader(r, versionInfo)
	if err != nil {
		return fmt.Errorf("failed to create agent client: %w", err)
	}

	err = client.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start agent client: %w", err)
	}

	p.client = client
	p.services = services

	config, more, err := p.updateConfig()
	if err != nil {
		return fmt.Errorf("failed to get initial configuration: %w", err)
	}
	if !more {
		return errors.New("client closed before getting initial configuration")
	}

	go clientLoop()

	return nil
}

func (p *provider) clientLoop() {
	for {
		config, more, err := p.updateConfig()
		if err != nil {
		}
		if !more {
			return
		}
	}
}

func (p *provider) updateConfig() (*confmap.Retrieved, bool, error) {
	select {
	case changes, more := <-p.client.UnitChanges():
		if !more {
			return nil, false, nil
		}
		changes.Unit.Expected
	case err, more := <-p.client.Errors() {
		if !more {
			return nil, false, nil
		}
	}
}

// Scheme returns the location scheme used by Retrieve.
func (p *provider) Scheme() string {
	return schemeName
}

// Shutdown signals that the configuration for which this Provider was used to
// retrieve values is no longer in use and the Provider should close and release
// any resources that it may have created.
func (p *provider) Shutdown(_ context.Context) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.client == nil {
		return nil
	}

	p.client.Stop()
	p.client = nil
	p.services = nil
	return nil
}
