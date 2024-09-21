default: test gen build


build:
	go build -o sshtrust .

test:
	go test -fullpath ./...

gen:
	swag init
