package viamkasa

import (
	"context"
	"fmt"

	"go.viam.com/rdk/components/switch"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/discovery"
	"go.viam.com/rdk/utils"

	"github.com/cloudkucooland/go-kasa"
)

var KasaDiscovery = family.WithModel("kasa-discovery")

func init() {
	resource.RegisterService(discovery.API, KasaDiscovery,
		resource.Registration[discovery.Service, resource.NoNativeConfig]{
			Constructor: newViamkasaKasaDiscover,
		},
	)
}

type viamkasaKasaDiscover struct {
	resource.AlwaysRebuild
	resource.TriviallyCloseable

	name resource.Name

	logger logging.Logger
}

func newViamkasaKasaDiscover(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (discovery.Service, error) {
	s := &viamkasaKasaDiscover{
		name:   rawConf.ResourceName(),
		logger: logger,
	}

	return s, nil
}

func (s *viamkasaKasaDiscover) Name() resource.Name {
	return s.name
}

func (s *viamkasaKasaDiscover) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (s *viamkasaKasaDiscover) DiscoverResources(ctx context.Context, extra map[string]any) ([]resource.Config, error) {
	all, err := kasa.BroadcastDiscovery(10, 1)
	if err != nil {
		return nil, fmt.Errorf("cannot do kasa discovery: %w", err)
	}

	configs := []resource.Config{}
	for host, info := range all {
		configs = append(configs, resource.Config{
			Name:  info.Alias,
			API:   toggleswitch.API,
			Model: KasaSwitch,
			Attributes: utils.AttributeMap{
				"Host": host,
			},
		})
	}
	return configs, nil
}
