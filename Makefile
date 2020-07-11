cli:
	go build -mod vendor -o bin/emit cmd/emit/main.go
	go build -mod vendor -o bin/images cmd/images/main.go
