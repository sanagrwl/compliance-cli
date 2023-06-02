package conftest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/walle/targz"
)

const (
	cliInfoURL                = "https://api.github.com/repos/open-policy-agent/conftest/releases/latest"
	cliZipDownloadUrlTemplate = "https://github.com/open-policy-agent/conftest/releases/download/v%s/conftest_%s_Darwin_arm64.tar.gz"
	zipFilename               = "conftest.tar.gz"
)

func workingDir() string {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	return workingDir
}

func downloadsDir() string {
	return filepath.Join(workingDir(), "downloads")
}

func cliDir() string {
	return filepath.Join(downloadsDir(), "conftest-cli")
}

func cliPath() string {
	return filepath.Join(cliDir(), "conftest")
}

func zipDownloadPath() string {
	return filepath.Join(downloadsDir(), filepath.Base(zipFilename))
}

func policiesDir() string {
	return filepath.Join(downloadsDir(), "policies")
}

func cliLatestVersion() string {
	resp, err := http.Get(cliInfoURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	x := map[string]string{}

	json.Unmarshal(body, &x)
	tagName := x["tag_name"]
	latestVersion := strings.ReplaceAll(tagName, "v", "")
	return latestVersion
}

func createDownloadsDir() {
	path := downloadsDir()
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func downloadCLI() {
	latestVersion := cliLatestVersion()
	downloadUrl := fmt.Sprintf(cliZipDownloadUrlTemplate, latestVersion, latestVersion)
	resp, err := http.Get(downloadUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(zipDownloadPath())
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
}

func extractCLI() {
	err := targz.Extract(zipDownloadPath(), cliDir())
	if err != nil {
		log.Fatalln(err)
	}
	_, err = os.Stat(cliPath())
	if err != nil {
		log.Fatal(err)
	}
}

func makeCLIExecutable() {
	err := os.Chmod(cliPath(), 0700)
	if err != nil {
		log.Println("what error")
		log.Fatalln(err)
	}
}

func latestCLIExists() bool {
	_, err := os.Stat(cliPath())
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	makeCLIExecutable()
	latestVersion := cliLatestVersion()
	result, err := execCLI("-v")
	if err != nil {
		log.Fatalln(result)
	}

	return strings.Contains(result, latestVersion)
}

func downloadPolicies() {
	policyDirRelativePath := strings.ReplaceAll(policiesDir(), workingDir(), "")
	policyUrl := "localhost:8080/policies"
	result, err := execCLI("pull", policyUrl, "--policy", policyDirRelativePath)
	if err != nil {
		log.Fatalln(result)
	}
}

func execCLI(cmdArgs ...string) (string, error) {
	out, err := exec.Command(cliPath(), cmdArgs...).CombinedOutput()
	return string(out), err
}

type ConftestMsg struct {
	Msg string `json:"msg"`
}

type ConftestResult struct {
	Filename  string        `json:"filename"`
	Namespace string        `json:"namespace"`
	Successes int           `json:"successes"`
	Warnings  []ConftestMsg `json:"warnings"`
	Failures  []ConftestMsg `json:"failures"`
}

func testCompliance(filepath string) []ConftestResult {
	out, err := exec.Command(cliPath(), "test", filepath, "--policy", policiesDir(), "--all-namespaces", "--no-fail", "--output", "json").CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "unknown parser") {
			return nil
		}
		log.Fatalln(string(out))
	}
	// log.Println(string(out))

	result := []ConftestResult{}

	json.Unmarshal(out, &result)
	return result
}
