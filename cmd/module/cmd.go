package main

import (
	"go.viam.com/rdk/components/switch"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"

	"github.com/erh/viamkasa"
)

func main() {
	module.ModularMain(resource.APIModel{toggleswitch.API, viamkasa.KasaSwitch})
}
