CLIENT := client
TEST := test
TARGET := compute
SRC_DIR := server

CXXFLAGS := 
LDFLAGS :=
CXXFLAGS += -O3 -Wall -Werror -Wextra
# CXXFLAGS += -Wno-deprecated-gpu-targets


all: run

run: client_go server_mpi
	cd build && ./test --seed 1 --local --runs 3 --mpi 4
# cd build && ./test --seed 1 --local

server_simple:
	@mkdir -p build
	g++ $(CXXFLAGS) $(LDFLAGS) $(SRC_DIR)/$@.cpp -o build/$(TARGET)

server_mpi:
	@mkdir -p build
	mpic++ $(CXXFLAGS) $(LDFLAGS) $(SRC_DIR)/$@.cpp -o build/$(TARGET)

server_mpi_send:
	@mkdir -p build
	mpic++ $(CXXFLAGS) $(LDFLAGS) $(SRC_DIR)/$@.cpp -o build/$(TARGET)

server_cuda:
	@mkdir -p build
	nvcc -arch=compute_35 -code=sm_35 $(CXXFLAGS) $(LDFLAGS) $(SRC_DIR)/$@.cu -o build/$(TARGET)

client_go:
	@mkdir -p build
	go build -o build/$(CLIENT) ./cmd/$(CLIENT)
	go build -o build/$(TEST) ./cmd/$(TEST)

clean:
	rm -rf build

deps:
	go mod download

update-deps:
	go get -u ./...
	go mod tidy -v

tidy:
	go mod tidy -v

.PHONY: all client_go clean deps update-deps tidy
