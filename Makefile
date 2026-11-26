# Let's keep all our changelog commands the same across all our packages:
update-changelog:
	conventional-changelog -p angular -i CHANGELOG.md -s

# Rewrite all Go source files in gofmt's simplified style.
fmt:
	gofmt -s -w .

# Check formatting without rewriting anything, then run go vet and
# staticcheck (go install honnef.co/go/tools/cmd/staticcheck@latest).
# Stops at the first tool that reports a problem.
lint:
	@unformatted="$$(gofmt -s -l .)" || exit 1; \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt: these files need 'make fmt':"; \
		echo "$$unformatted"; \
		exit 1; \
	fi
	go vet ./...
	staticcheck ./...

test:
	go test ./...

.PHONY: update-changelog fmt lint test
