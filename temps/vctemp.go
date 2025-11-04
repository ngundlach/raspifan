package temps

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type vcTemp struct {
	command string
	args    string
}

func NewVcTemp() *vcTemp {
	tr := new(vcTemp)
	tr.command = "/usr/bin/vcgencmd"
	tr.args = "measure_temp"
	fmt.Println("using vcgencmd to retrieve temperature")
	return tr
}

func (tr *vcTemp) ReadTemp() (float64, error) {
	cmd := exec.Command(tr.command, tr.args)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	strtemp, _ := strings.CutPrefix(string(out), "temp=") // removing unnecessary output
	strtemp, _ = strings.CutSuffix(strtemp, "'C\n")
	temp, err := strconv.ParseFloat(strtemp, 64)
	if err != nil {
		return 0, err
	}
	return temp, nil
}
