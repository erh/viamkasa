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

func NewDiscovery(logger logging.Logger) *ViamkasaKasaDiscover {
	return &ViamkasaKasaDiscover{logger: logger}
}

type ViamkasaKasaDiscover struct {
	resource.AlwaysRebuild
	resource.TriviallyCloseable

	name resource.Name

	logger logging.Logger
}

func newViamkasaKasaDiscover(ctx context.Context, _ resource.Dependencies, rawConf resource.Config, logger logging.Logger) (discovery.Service, error) {
	s := &ViamkasaKasaDiscover{
		name:   rawConf.ResourceName(),
		logger: logger,
	}

	return s, nil
}

func (s *ViamkasaKasaDiscover) Name() resource.Name {
	return s.name
}

func (s *ViamkasaKasaDiscover) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (s *ViamkasaKasaDiscover) DiscoverResources(ctx context.Context, extra map[string]any) ([]resource.Config, error) {
	return s.DiscoverKasa(ctx, 10, 1)
}

func (s *ViamkasaKasaDiscover) DiscoverKasa(ctx context.Context, timeout, probes int) ([]resource.Config, error) {
	all, err := kasa.BroadcastDiscovery(timeout, probes)
	if err != nil {
		return nil, fmt.Errorf("cannot do kasa discovery: %w", err)
	}

	configs := []resource.Config{}
	for host, info := range all {
		s.logger.Debugf("discovery result host: %v\n\t%#v", host, info)

		c := resource.Config{
			Name: info.Alias,
			Attributes: utils.AttributeMap{
				"Host": host,
			},
		}

		switch info.MIC {
		case "IOT.SMARTBULB":
			c.API = toggleswitch.API
			c.Model = KasaLight
		case "IOT.SMARTPLUGSWITCH":
			c.API = toggleswitch.API
			c.Model = KasaSwitch
		default:
			s.logger.Warnf("unknown MIC [%s] using basic switch", info.MIC)
			c.API = toggleswitch.API
			c.Model = KasaSwitch
		}

		configs = append(configs, c)
	}
	return configs, nil
}
