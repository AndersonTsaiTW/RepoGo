# ===== repogo Makefile =====
APP := repogo
PKG := .
# Use git tag or manual VERSION, will be injected into main.Version
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo 0.1.0)

# Go build flags
LDFLAGS := -s -w -X 'main.Version=$(VERSION)'

# Target platforms (can be customized)
TARGETS := \
	darwin/amd64 \
	darwin/arm64 \
	linux/amd64  \
	linux/arm64  \
	windows/amd64 \
	windows/arm64

# Default: build for host platform
.PHONY: build
build:
	@echo "Building $(APP) $(VERSION) for host..."
	@go build -trimpath -ldflags="$(LDFLAGS)" -o bin/$(APP) $(PKG)

# Cross-platform build
.PHONY: cross
cross: clean
	@mkdir -p dist
	@set -e; \
	for tgt in $(TARGETS); do \
		GOOS=$${tgt%/*}; GOARCH=$${tgt#*/}; \
		OUT="dist/$(APP)-$${GOOS}-$${GOARCH}"; \
		EXT=""; \
		if [ "$$GOOS" = "windows" ]; then EXT=".exe"; fi; \
		echo "Building $$OUT$$EXT ..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o "$$OUT$$EXT" $(PKG); \
	done

# Package compression (zip/tar.gz)
.PHONY: package
package: cross
	@set -e; \
	for f in dist/$(APP)-*; do \
		if echo $$f | grep -q windows; then \
			zip -j "$$f.zip" "$$f.exe" >/dev/null; \
			rm -f "$$f.exe"; \
		else \
			tar -C dist -czf "$$f.tar.gz" "$$(basename $$f)"; \
			rm -f "$$f"; \
		fi; \
	done
	@echo "Artifacts in ./dist"

# Generate SHA256 checksum file
.PHONY: checksum
checksum:
	@cd dist && \
	for a in *.{zip,tar.gz}; do \
		[ -f "$$a" ] || continue; \
		if command -v shasum >/dev/null 2>&1; then \
			shasum -a 256 "$$a" >> SHA256SUMS; \
		else \
			sha256sum "$$a" >> SHA256SUMS; \
		fi \
	done && \
	echo "Wrote dist/SHA256SUMS"

# All-in-one: build -> package -> checksum
.PHONY: release
release: package checksum
	@echo "Release artifacts ready in ./dist"

.PHONY: clean
clean:
	@rm -rf bin dist
