CXXFLAGS := -Wno-deprecated-gpu-targets
# CXXFLAGS := -O3 -Wall -Werror -Wextra
LDFLAGS :=

# CXXFLAGS += -fopenmp
# LDFLAGS += -fopenmp

# CFLAGS += -I/usr/local/opt/openssl/include
# LDFLAGS += -L/usr/local/opt/openssl/lib -lssl -lcrypto

TARGET := compute
SRC_DIR := src
# SRC_DIR := src
# BUILD_DIR := build

all: server_mpi

server_simple:
	g++ $(CXXFLAGS) $(LDFLAGS) $(SRC_DIR)/$@.cpp -o $(TARGET)

server_mpi:
	mpic++ $(CXXFLAGS) $(LDFLAGS) $(SRC_DIR)/$@.cpp -o $(TARGET)

server_cuda:
	nvcc -arch=compute_35 -code=sm_35 $(CXXFLAGS) $(LDFLAGS) $(SRC_DIR)/$@.cu -o $(TARGET)

clean:
	rm $(TARGET)
