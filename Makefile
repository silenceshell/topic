build:
	GOOS=linux go build cmd/topic/topic.go

image: build
	docker build . -t silenceshell/topic

clean:
	rm -f topic
