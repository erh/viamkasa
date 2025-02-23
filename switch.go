package viamkasa

import (
	"context"
	"fmt"

	"go.viam.com/rdk/components/switch"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"

	"github.com/cloudkucooland/go-kasa"
)

var (
	KasaSwitch = resource.NewModel("erh", "viamkasa", "kasa-switch")
)

func init() {
	resource.RegisterService(toggleswitch.API, KasaSwitch,
		resource.Registration[toggleswitch.Switch, *Config]{
			Constructor: newViamkasaKasaSwitch,
		},
	)
}

type Config struct {
	Host string
}

func (cfg *Config) Validate(path string) ([]string, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("need a host")
	}
	return nil, nil
}

type viamkasaKasaSwitch struct {
	resource.AlwaysRebuild
	resource.TriviallyCloseable

	name resource.Name

	logger logging.Logger
	cfg    *Config

	dev *kasa.Device
}

func newViamkasaKasaSwitch(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (toggleswitch.Switch, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	s := &viamkasaKasaSwitch{
		name:   rawConf.ResourceName(),
		logger: logger,
		cfg:    conf,
	}

	s.dev, err = kasa.NewDevice(conf.Host)
	if err != nil {
		return nil, fmt.Errorf("can't connect to kasa swtich @ (%s): %w", conf.Host, err)
	}

	return s, nil
}

func (s *viamkasaKasaSwitch) Name() resource.Name {
	return s.name
}

func (s *viamkasaKasaSwitch) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (s *viamkasaKasaSwitch) SetPosition(ctx context.Context, position uint32, extra map[string]interface{}) error {
	state := false
	if position > 0 {
		state = true
	}
	return s.dev.SetRelayState(state)
}

func (s *viamkasaKasaSwitch) GetPosition(ctx context.Context, extra map[string]interface{}) (uint32, error) {
	settings, err := s.dev.GetSettings()
	if err != nil {
		return 0, err
	}
	state := uint32(0)
	if settings.RelayState > 0 {
		state = 1
	}
	return state, nil
}

func (s *viamkasaKasaSwitch) GetNumberOfPositions(ctx context.Context, extra map[string]interface{}) (uint32, error) {
	return 2, nil
}
