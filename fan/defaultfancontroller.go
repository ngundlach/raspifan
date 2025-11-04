package fan

import (
	"fmt"
	"log"

	"github.com/stianeikeland/go-rpio/v4"
)

type defaultFanController struct {
	pin rpio.Pin
}

func NewDefaultFc(fanpin int) (*defaultFanController, error) {
	f := new(defaultFanController)
	if err := f.setup(fanpin); err != nil {
		return nil, err
	}

	return f, nil
}

func (fc *defaultFanController) setup(fanpin int) error {
	fmt.Println("opening pi")
	if err := rpio.Open(); err != nil {
		return err
	}
	fmt.Println("setting fan")
	fc.pin = rpio.Pin(fanpin)
	fmt.Println("setting output")
	fc.pin.Output()
	return nil
}

func (fc *defaultFanController) SwitchOn() {
	fc.pin.High()
}

func (fc *defaultFanController) SwitchOff() {
	fc.pin.Low()
}

func (fc *defaultFanController) CleanUp() {
	rpio.Close()
	log.Println("cleaning up...")
}

func (fc *defaultFanController) FanState() FanState {
	fan := fc.pin.Read()
	switch fan {
	case rpio.High:
		return On
	default:
		return Off
	}
	// if fan == rpio.High {
	// 	return true
	// } else {
	// 	return false
	// }
}

var _ FanController = (*defaultFanController)(nil)
