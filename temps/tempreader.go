package temps

const (
	VcTemp    tempProvider = "vc"
	SysFsTemp tempProvider = "sysfs"
)

type tempProvider string

type TempReader interface {
	ReadTemp() (float64, error)
}
