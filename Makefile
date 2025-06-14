node:
	go build -o ./bin/node ./cmd/Node/main.go

clean-logs:
	find . -type f -name log.txt -delete