SHELL := /bin/sh

.PHONY: nav-test nav-run ui-install ui-build vision-check

nav-test:
	cd navigation_core && go test ./...

nav-run:
	cd navigation_core && go run ./cmd/navigator

ui-install:
	cd tactical_map && npm install

ui-build:
	cd tactical_map && npm run build

vision-check:
	python -m compileall vision_provider simulation
