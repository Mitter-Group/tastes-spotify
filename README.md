# Go template

Go template - Managed via CloudFormation

```
== Go API ==
+ Fiber
+ Gitflow
+ Full pipeline (dev-stg-prd, security-orbs, gitflow, lint, go vet, private github packages)
+ Linter (golangci-lint mandatory [https://golangci-lint.run/usage/install/] + revive optional [https://revive.run/docs#installation])
+ Config
+ Swagger
+ Hooks(auto swagger + linters + tests/coverage)
+ Auth0
+ New Relic (with zap)
+ Sonarqube
+ Readiness/Liveness
+ Utilities
```

### Install required dependencies and tools

```
go mod tidy
sh ./tools/install-hooks.sh
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/mgechev/revive
brew install golangci-lint
```

### Run app:
You can start vscode debugger or by terminal:
```
go run ./cmd/api/main.go
```

### Run linters:

```
golangci-lint run --max-issues-per-linter=0 --max-same-issues=0 --config=./.golangci.yml

revive -config revive.toml -formatter friendly ./...
```

### Generate swagger documentation:

```
swag init -g routes.go --dir internal/handlers --parseDependency
```

### Check go vulnerabilities
```
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```
