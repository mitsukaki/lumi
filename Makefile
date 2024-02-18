# Compiler options
GO := go
GOFLAGS := -v

# App build target
app: lumi.go
	$(GO) build $(GOFLAGS) -o build/lumi lumi.go

# web build target
web:
	./scripts/web_build.sh

# Build target
build: web app

# Run target
run: build
	./build/lumi

# Clean target
clean:
	rm -f build
