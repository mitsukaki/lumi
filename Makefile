# Compiler options
GO := go
GOFLAGS := -v

# Install dependencies target
deps:
	$(GO) mod download

# App build target
app: deps lumi.go
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -o build/lumi lumi.go

# web build target
webui: web/*
	cd web && npm install && npx vite build . --outDir ../build/public --emptyOutDir

# Build target
build: app webui

# Run target
run: build
	./build/lumi

# Docker build target
docker: clean
	docker build -t lumi .

# Test target
test:
	$(GO) test ./...

# Clean target
clean:
	rm -rf build
	