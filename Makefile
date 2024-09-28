export SHELL:=/bin/bash
export SHELLOPTS:=$(if $(SHELLOPTS),$(SHELLOPTS):)pipefail:errexit

.ONESHELL:

default: test gen build launchtest


build:
	go build -o sshtrust .

test:
	go test -fullpath ./...

gen:
	swag init |sed 's/^[0-9]\{4\}\/[0-9]\{2\}\/[0-9]\{2\} [0-9]\{2\}:[0-9]\{2\}:[0-9]\{2\} //'

launchtest:
	bash launch-server.sh
