build:
	GOOS=linux go build cmd/topic/topic.go

clean:
	rm -f topic
