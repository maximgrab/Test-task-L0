all:
	go build -v ./pub.go
	./pub &
	go build -v ./cmd/main.go
	./main
