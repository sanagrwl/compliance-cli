## Compliance CLI

CLI to check if code, docker images etc are compliant based on policies.

### Setup
Install
- [go lang](https://go.dev/)
- [cobra cli](https://github.com/spf13/cobra-cli/blob/main/README.md)
- [conftest cli](https://github.com/open-policy-agent/conftest/releases)

### CLI Commands:

#### test
- Downloads conftest cli under `<repo>/downloads` folder
- Downloads policies from registry
- Runs the tests and displays results


```bash
# Run docker registry for conftest policies
docker run -d --rm -p 8080:5000 --name registry registry:latest

# policies
conftest push --policy policies localhost:8080/policies:latest

# test policies in registry
conftest pull localhost:8080/policies:latest --policy foo
ls foo/policies

# run a command without building
go run main.go test

# build cli
go build -o dist/compliance-cli

# run cli
dist/compliance-cli

# Add command
cobra-cli add <command-name>
```