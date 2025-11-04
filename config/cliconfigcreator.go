package config

import (
	"flag"

	"github.com/ngundlach/raspifan/temps"
)

type CliConfigCreator struct{}

func (cc *CliConfigCreator) CreateConfig() *Config {
	cfg := new(Config)

	flag.StringVar(&cfg.TempProvider, "tempprovider", string(temps.SysFsTemp), "Provider for temperature readout")
	flag.StringVar(&cfg.SysFsPath, "sysfspath", "/sys/class/thermal/thermal_zone0/temp", "Sysfs-path for temperature readout")
	flag.IntVar(&cfg.FanPin, "fanpin", 18, "GPIO control pin")
	flag.IntVar(&cfg.Poll, "poll", 5, "Temperature pollingtime in seconds")
	flag.Float64Var(&cfg.TurnOnTemp, "turnontemp", 48, "Temperature threshold that turns on the fan")
	flag.Float64Var(&cfg.TurnOffTemp, "turnofftemp", 40, "Temperature threshold that turns off the fan")
	flag.BoolVar(&cfg.WebServer, "webserver", false, "Start the webserver")
	flag.IntVar(&cfg.Port, "port", 3333, "Port the webserver is listening on")
	flag.Parse()
	return cfg
}

var _ configCreator = (*CliConfigCreator)(nil)
