.PHONY: install_wire
install_wire:
	go install github.com/google/wire/cmd/wire@latest \
	&& [ -n "$$(command -v asdf 2>/dev/null)" ] && asdf reshim golang || true

.PHONY: install_package
install_package:
	go mod tidy

.PHONY: tidy
tidy: install_package install_wire

.PHONY: remove_local
remove_local:
	git remote update --prune
	git switch --detach origin/main
	@git for-each-ref --format '%(refname:short)' refs/heads | xargs -r -t git branch -D

.PHONY: wire
wire:
	wire ./internal/ui/restapi

.PHONY: test
test:
	go test -cover ./...