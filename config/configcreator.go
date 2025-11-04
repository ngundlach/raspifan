package config

type Config struct {
	TempProvider string
	SysFsPath    string
	TurnOnTemp   float64
	TurnOffTemp  float64
	FanPin       int
	Poll         int
	Port         int
	WebServer    bool
}

type configCreator interface {
	CreateConfig() *Config
}
