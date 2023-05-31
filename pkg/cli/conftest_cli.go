package cli

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
)

//go:embed *conftest
var conftestCLIData []byte //assign the variable conftest to embeded file

var conftestCliName string = "conftest-cli"

func init() {

	_ = os.WriteFile(conftestCliName, conftestCLIData, 0755)
}

func ExecConftest() {
	out, _ := exec.Command("./" + conftestCliName).Output()
	fmt.Printf("Output: %s\n", out)

}
