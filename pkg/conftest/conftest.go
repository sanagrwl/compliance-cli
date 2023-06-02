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
	fmt.Println("Downloading policies\n")
	downloadPolicies()
}

func printMsgs(msgType string, msgs []ConftestMsg) {
	if msgs == nil && len(msgs) == 0 {
		return
	}

	fmt.Println(msgType)
	for i := 0; i < len(msgs); i++ {
		fmt.Printf("%v: %s", i+1, msgs[0].Msg)
	}
}

func printResults(results []ConftestResult) {
	for i := 0; i < len(results); i++ {
		result := results[i]
		if len(result.Warnings) > 0 || len(result.Failures) > 0 {
			fmt.Printf("File: %s\n", results[0].Filename)
			fmt.Printf("Namespace: %s\n", result.Namespace)
			printMsgs("Failures:", result.Failures)
			printMsgs("Warnings:", result.Warnings)
			fmt.Println("\n")
		}
	}
}

func testAllRepoFiles(root string) error {
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() && (f.Name() == ".git" || f.Name() == "downloads" || f.Name() == "dist") {
			return filepath.SkipDir
		}

		if !f.IsDir() {
			results := testCompliance(path)
			if results != nil {
				printResults(results)
			}

		}
		return nil
	})
	return err
}

func CheckCompliance() {
	testAllRepoFiles(workingDir())
}
