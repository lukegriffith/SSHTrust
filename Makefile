
test:
	go test -fullpath ./...

gen:
	swag init -g cmd/server/main.go
