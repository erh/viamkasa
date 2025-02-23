package main

import (
	"context"
	"fmt"

	"go.viam.com/rdk/logging"

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
	logger.SetLevel(logging.DEBUG)

	d := viamkasa.NewDiscovery(logger)
	all, err := d.DiscoverResources(ctx, nil)
	if err != nil {
		return err
	}
	for _, c := range all {
		fmt.Printf("%v\n", c)
	}
	return nil
}
