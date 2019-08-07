BINARY            = feegiver
GITHUB_USERNAME   = terra-project
VERSION           = v0.1.0
GOARCH            = amd64
ARTIFACT_DIR      = build

COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
FLAG_PATH=github.com/${GITHUB_USERNAME}/${BINARY}/cmd

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X ${FLAG_PATH}.Version=${VERSION} -X ${FLAG_PATH}.Commit=${COMMIT} -X ${FLAG_PATH}.Branch=${BRANCH}"

# Build the project
all: clean linux darwin windows

# Build and Install project into GOPATH using current OS setup
install:
	go install ${LDFLAGS} ./...

test:
	go test -v ./utils/...

# Build binary for Linux
linux: clean
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${ARTIFACT_DIR}/${BINARY}-linux-${GOARCH} . ;

# Build binary for MacOS
darwin:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${ARTIFACT_DIR}/${BINARY}-darwin-${GOARCH} . ;

# Build binary for Windows
windows:
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${ARTIFACT_DIR}/${BINARY}-windows-${GOARCH}.exe . ;

# Remove all the built binaries
clean:
	rm -rf ${ARTIFACT_DIR}/*

.PHONY: all install test linux darwin windows clean
