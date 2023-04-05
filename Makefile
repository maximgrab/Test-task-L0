all:
	go build -v ./cmd/pub/pub.go
	./pub &
	go build -v ./cmd/main/main.go
	./main
