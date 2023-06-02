package conftest

import (
	"fmt"
	"os"
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

func testAllRepoFiles(root string) error {
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() && (f.Name() == ".git" || f.Name() == "downloads" || f.Name() == "dist") {
			return filepath.SkipDir
		}

		if !f.IsDir() {
			result := testCompliance(path)
			if result != nil {
				fmt.Println("File: %s", result[0]["filename"])
				failures := result[0]["failures"]
				if failures != nil {
					fmt.Println("Failure: %s", failures)
				}
				warnings := result[0]["warnings"]
				if warnings != nil {
					fmt.Println("Warnings: %s", warnings)
				}

			}

		}
		return nil
	})
	return err
}

func CheckCompliance() {
	testAllRepoFiles(workingDir())
}
