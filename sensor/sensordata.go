package sensor

import (
	"sync"

	"github.com/ngundlach/raspifan/fan"
)

type SensorData struct {
	mu sync.RWMutex
	D  data
}

type data struct {
	Temp float64
	Fan  fan.FanState
}

func (sd *SensorData) Store(temp float64, fan fan.FanState) {
	sd.mu.Lock()
	// log.Println("writing data")
	defer sd.mu.Unlock()

	sd.D = data{temp, fan}
}

func (sd *SensorData) Get() data {
	sd.mu.RLock()
	// log.Println("reading data")
	defer sd.mu.RUnlock()

	return sd.D
}
