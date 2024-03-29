.DEFAULT_GOAL := help

export GO111MODULE=on

VERSION := $(shell git rev-parse HEAD)
ENV     := development
HOST    := localhost:8080

.PHONY: help deps run build fmt vet clean test

help:
	@cat Makefile

deps: export GO111MODULE=off
deps: env/env.go dev.db
	which sql-migrate || go get -u -v github.com/rubenv/sql-migrate/...

run:
	go run server.go

build: fmt vet
	go build -ldflags "-X=main.version=$(VERSION)" server.go

fmt:
	go fmt $$(go list ./...)

vet:
	go vet $$(go list ./...)

clean:
	rm -rf vendor
	rm dev.db

test: fmt vet
	@rm -f test.db
	@cp -i _etc/seed.db test.db
	GIN_MODE=test go test -v

env/env.go:
	cp env/env.go.tmpl env/env.go

dev.db:
	cp -i _etc/seed.db dev.db

.PHONY: migrate_*
## Migrate db schema
migrate_up:
	sql-migrate up -env=$(ENV)

## Migrate db schema(dryrun)
migrate_dryrun:
	sql-migrate up -env=$(ENV) -dryrun

## Show migration status
migrate_status:
	sql-migrate status -env=$(ENV)

.PHONY: curl_*
curl_ping:
	curl -i $(HOST)/api/ping

curl_messages_get_all:
	curl -i $(HOST)/api/messages

ID :=
curl_messages_get:
	curl -i $(HOST)/api/messages/$(ID)

BODY :=
curl_message_post:
	curl -i -X POST $(HOST)/api/messages -d '{"BODY": "$(BODY)"}'



