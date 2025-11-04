package temps

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type sysFsTemp struct {
	path string
}

func NewSysFsTemp(path string) (*sysFsTemp, error) {
	tr := new(sysFsTemp)
	tr.path = path
	fmt.Println("using sysfs to retrieve temperature on path: ", path)
	if err := tr.testTempRetrieval(); err != nil {
		return nil, err
	}
	return tr, nil
}

func (tr *sysFsTemp) testTempRetrieval() error {
	_, err := tr.ReadTemp()
	return err
}

func (tr *sysFsTemp) ReadTemp() (float64, error) {
	content, err := os.ReadFile(tr.path)
	if err != nil {
		return 0, err
	}
	strtemp, _ := strings.CutSuffix(string(content), "\n") // remove trailing newline
	temp, err := strconv.ParseFloat(strtemp, 64)
	if err != nil {
		return 0, err
	}
	return temp / 1000, nil
}
