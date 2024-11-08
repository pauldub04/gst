CLIENT := client
TEST := test

all: compile

run: compile
	./$(CLIENT)

compile:
	go build ./cmd/$(CLIENT)
	go build ./cmd/$(TEST)

clean:
	rm -rf $(CLIENT) $(TEST) result.txt statistics.txt

deps:
	go mod download

update-deps:
	go get -u ./...
	go mod tidy -v

tidy:
	go mod tidy -v

.PHONY: all run compile clean deps update-deps tidy
