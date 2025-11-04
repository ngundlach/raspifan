package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/ngundlach/raspifan/config"
	"github.com/ngundlach/raspifan/fan"
	"github.com/ngundlach/raspifan/sensor"
	"github.com/ngundlach/raspifan/temps"
	"github.com/ngundlach/raspifan/web"
)

type app struct {
	fc fan.FanController
	tr temps.TempReader
	sd *sensor.SensorData
}

func main() {
	cfg := createConfig()
	app, err := createAppContext(cfg)
	if err != nil {
		log.Println(err.Error())
		os.Exit(3)
	}

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	var wg sync.WaitGroup
	// log.Println(fc.pin)

	// starting program loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		Run(app, cfg, ctx)
	}()

	// waiting for signal to stop
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-c
		cancel()
	}()

	if cfg.WebServer {
		server := web.NewRestService(app.sd, ":"+strconv.Itoa(cfg.Port))
		go func() {
			log.Println("Starting webserver")
			err := server.Srv.ListenAndServe()
			if err != nil {
				log.Println("error starting webserver: ", err)
				return
			}
		}()
	}

	wg.Wait()

	app.fc.CleanUp()
}

func Run(app *app, cfg *config.Config, ctx context.Context) {
	adjustFan(app, cfg)
	tkr := time.NewTicker(time.Duration(cfg.Poll) * time.Second)
	defer log.Println("stopping...")
	defer tkr.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tkr.C:
			adjustFan(app, cfg)
		}
	}
}

func adjustFan(app *app, cfg *config.Config) {
	cTemp, err := app.tr.ReadTemp()
	// log.Println(cTemp)
	if err != nil {
		log.Println(err.Error())
		os.Exit(2)
	}
	switch {
	case cTemp >= cfg.TurnOnTemp:
		app.fc.SwitchOn()
	case cTemp <= cfg.TurnOffTemp:
		app.fc.SwitchOff()
	}
	fan := app.fc.FanState()
	app.sd.Store(cTemp, fan)
}

func createAppContext(cfg *config.Config) (*app, error) {
	app := new(app)
	switch cfg.TempProvider {
	case string(temps.SysFsTemp):
		tr, err := temps.NewSysFsTemp(cfg.SysFsPath)
		if err != nil {
			return nil, err
		}
		app.tr = tr
	case string(temps.VcTemp):
		app.tr = temps.NewVcTemp()
	default:
		return nil, errors.New("unsupported temperature provider")
	}
	fc, err := fan.NewDefaultFc(cfg.FanPin)
	if err != nil {
		log.Println("Error while initializing FanController: ", err)
		os.Exit(3)
	}
	app.fc = fc
	app.sd = new(sensor.SensorData)
	return app, nil
}

func createConfig() *config.Config {
	cc := new(config.CliConfigCreator)
	return cc.CreateConfig()
}
