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
	return filepath.Join(cliDir(), "policies")
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
	out := execCLI("-v")

	return strings.Contains(out, latestVersion)
}

func downloadPolicies() {
	policyDirRelativePath := strings.ReplaceAll(policiesDir(), workingDir(), "")
	policyUrl := "https://raw.githubusercontent.com/open-policy-agent/conftest/master/examples/docker/policy/images.rego"
	execCLI("pull", policyUrl, "--policy", policyDirRelativePath)

	policyUrl = "https://raw.githubusercontent.com/open-policy-agent/conftest/master/examples/kustomize/policy/base.rego"
	execCLI("pull", policyUrl, "--policy", policyDirRelativePath)
}

func execCLI(cmdArgs ...string) string {
	out, err := exec.Command(cliPath(), cmdArgs...).Output()
	if err != nil {
		log.Fatalln(err)
	}
	return string(out)
}

func testCompliance(filepath string) []map[string]interface{} {
	out, err := exec.Command(cliPath(), "test", filepath, "--policy", policiesDir(), "--no-fail", "--output", "json").Output()
	if err != nil {
		log.Println(string(out))
		log.Fatalln(err)
	}
	var result []map[string]interface{}

	json.Unmarshal(out, &result)
	return result
}
