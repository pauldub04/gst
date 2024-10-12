
TARGET := client

all: compile

run: compile
	./$(TARGET)

compile:
	go build ./cmd/$(TARGET)

clean:
	rm -rf $(TARGET) result.txt statistics.txt

deps:
	go mod download

update-deps:
	go get -u ./...
	go mod tidy -v

tidy:
	go mod tidy -v

.PHONY: all run compile clean deps update-deps tidy
