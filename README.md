# Setup

## Install

Download the archive for your system from the [releases page](https://github.com/tech-dipesh/drop/releases), unpack it and put `drop` (`drop.exe` on Windows) somewhere on your PATH.

Or with Go 1.22 or newer:

```bash
go install github.com/tech-dipesh/drop@latest
```

## First run

Set the name other devices will see and where incoming files go. Both are saved, so you only do this once:

```bash
drop -name "Dipesh laptop" -dir ~/Downloads/drop
```
## Build from source

```bash
git clone https://github.com/tech-dipesh/drop
cd drop
go test ./...
go build
```

Cross-compile for another system, for example Android:

```bash
GOOS=linux GOARCH=arm64 go build -o drop-arm64 .
```

Push a version tag and GitHub Actions builds and publishes every binary:

```bash
git tag v0.1.0
git push origin v0.1.0
```