package viamkasa

import (
	"context"
	"encoding/json"
	"fmt"

	"go.viam.com/rdk/components/switch"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"

	"github.com/cloudkucooland/go-kasa"
)

var (
	KasaLight = family.WithModel("kasa-light")
)

func init() {
	resource.RegisterComponent(toggleswitch.API, KasaLight,
		resource.Registration[toggleswitch.Switch, *LightConfig]{
			Constructor: newViamkasaKasaLight,
		},
	)
}

type LightConfig struct {
	Host string
}

func (cfg *LightConfig) Validate(path string) ([]string, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("need a host")
	}
	return nil, nil
}

type viamkasaKasaLight struct {
	resource.AlwaysRebuild
	resource.TriviallyCloseable

	name resource.Name

	logger logging.Logger
	cfg    *LightConfig

	dev      *kasa.Device
	settings *kasa.Sysinfo

	lastPosition uint32
}

func newViamkasaKasaLight(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (toggleswitch.Switch, error) {
	conf, err := resource.NativeConfig[*LightConfig](rawConf)
	if err != nil {
		return nil, err
	}

	s := &viamkasaKasaLight{
		name:   rawConf.ResourceName(),
		logger: logger,
		cfg:    conf,
	}

	s.dev, err = kasa.NewDevice(conf.Host)
	if err != nil {
		return nil, fmt.Errorf("can't connect to kasa swtich @ (%s): %w", conf.Host, err)
	}

	s.settings, err = s.dev.GetSettings()
	if err != nil {
		return nil, fmt.Errorf("couldn't get settings: %w", err)
	}

	if len(s.settings.PreferredState) == 0 {
		return nil, fmt.Errorf("no preferred states, need for light for now")
	}

	return s, nil
}

func (s *viamkasaKasaLight) Name() resource.Name {
	return s.name
}

func (s *viamkasaKasaLight) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (s *viamkasaKasaLight) SetPosition(ctx context.Context, position uint32, extra map[string]interface{}) error {
	sub := kasa.Preset{}

	if position == 0 {
		sub.OnOff = 0
	} else if position == 1 {
		sub.OnOff = 1
	} else {
		pref := int(position - 2)
		if pref >= len(s.settings.PreferredState) {
			return fmt.Errorf("invalid state")
		}
		sub = s.settings.PreferredState[pref]
		sub.OnOff = 1
	}

	command := map[string]interface{}{
		"smartlife.iot.smartbulb.lightingservice": map[string]interface{}{
			"transition_light_state": sub,
		},
	}

	jsonData, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf("Failed to marshal command: %w", err)
	}

	s.logger.Debugf("sending [%s]", string(jsonData))

	err = s.dev.SendRawCommand(string(jsonData))
	if err != nil {
		return fmt.Errorf("error from command: %v", err)
	}

	s.lastPosition = position

	return nil
}

func (s *viamkasaKasaLight) GetPosition(ctx context.Context, extra map[string]interface{}) (uint32, error) {
	return s.lastPosition, nil
}

func (s *viamkasaKasaLight) GetNumberOfPositions(ctx context.Context, extra map[string]interface{}) (uint32, error) {
	return uint32(2 + len(s.settings.PreferredState)), nil
}
