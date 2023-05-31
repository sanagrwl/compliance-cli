## Compliance CLI

CLI to check if code, docker images etc are compliant based on policies.

### Setup
Install
- [go lang](https://go.dev/)
- [cobra cli](https://github.com/spf13/cobra-cli/blob/main/README.md)

```bash
# copy conftest cli under pkg/cli and name it as `conftest`.
https://github.com/open-policy-agent/conftest/releases

# run a command without building
go run main.go test

# build cli
go build -o dist/compliance-cli

# run cli
dist/compliance-cli

# Add command
cobra-cli add <command-name>
```