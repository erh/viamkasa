package main

import (
	"context"
	"flag"
	"fmt"

	"go.viam.com/rdk/components/switch"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"

	"github.com/erh/viamkasa"
)

func main() {
	err := realMain()
	if err != nil {
		panic(err)
	}
}

func realMain() error {
	ctx := context.Background()
	logger := logging.NewLogger("viamkasacli")

	timeout := flag.Int("timeout", 10, "timeout for discovery")
	probes := flag.Int("probes", 1, "# probes for discovery")
	debug := flag.Bool("debug", false, "debug")
	device := flag.String("device", "", "What device to do")
	setting := flag.Int("set", -1, "What to set the device to")

	flag.Parse()

	if *debug {
		logger.SetLevel(logging.DEBUG)
	}

	d := viamkasa.NewDiscovery(logger)
	all, err := d.DiscoverKasa(ctx, *timeout, *probes)
	if err != nil {
		return err
	}

	var info resource.Config

	for _, c := range all {
		fmt.Printf("%v\n", c)
		if c.Name == *device {
			info = c
		}
	}

	if *device != "" {
		if info.Name == "" {
			return fmt.Errorf("cannot find device [%s]", *device)
		}

		logger.Infof("found device %v", info)

		reg, found := resource.LookupRegistration(info.API, info.Model)
		if !found {
			return fmt.Errorf("cannot find registration")
		}

		info.ConvertedAttributes, err = reg.AttributeMapConverter(info.Attributes)
		if err != nil {
			return err
		}

		thing, err := reg.Constructor(ctx, nil, info, logger)
		if err != nil {
			return err
		}

		realThing, ok := thing.(toggleswitch.Switch)
		if !ok {
			return fmt.Errorf("why aren't you a switch")
		}

		pos, err := realThing.GetPosition(ctx, nil)
		if err != nil {
			return err
		}

		numPositions, err := realThing.GetNumberOfPositions(ctx, nil)
		if err != nil {
			return err
		}

		logger.Infof("starting at position: %v out of %d", pos, numPositions)

		if *setting >= 0 {
			err = realThing.SetPosition(ctx, uint32(*setting), nil)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
