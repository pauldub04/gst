CXX = mpic++
CXXFLAGS := -O3 -Wall -Werror -Wextra
LDFLAGS :=

# CXXFLAGS += -fopenmp
# LDFLAGS += -fopenmp

# CFLAGS += -I/usr/local/opt/openssl/include
# LDFLAGS += -L/usr/local/opt/openssl/lib -lssl -lcrypto

TARGET := server
SRC_DIR := src
# SRC_DIR := src
# BUILD_DIR := build

all: $(TARGET)

$(TARGET):
	$(CXX) $(CXXFLAGS) $(LDFLAGS) $(SRC_DIR)/server.cpp -o $(TARGET)

clean:
	rm $(TARGET)
