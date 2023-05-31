package conftest

import (
	"fmt"
	"path/filepath"
)

func init() {
	createDownloadsDir()
	if !latestCLIExists() {
		fmt.Println("Downloading latest conftest CLI")
		downloadCLI()
		extractCLI()
		makeCLIExecutable()
		fmt.Println(fmt.Sprintf("Latest conftest CLI installed at %s", cliDir()))
	}
	fmt.Println("Downloading policies")
	downloadPolicies()
}

func ExecConftest() {
	fmt.Println("Testing dockerfile")
	result := testCompliance(filepath.Join(workingDir(), "Dockerfile.test"))

	fmt.Println(result[0]["filename"])
	fmt.Println(result[0]["failures"])

}
