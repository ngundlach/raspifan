package fan

type FanState string

const (
	On  FanState = "on"
	Off FanState = "off"
)

type FanController interface {
	SwitchOn()
	SwitchOff()
	FanState() FanState
	CleanUp()
}
