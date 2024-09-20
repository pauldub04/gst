CC = cc
CFLAGS := -std=c99 -g
#-O3 -pedantic -Wall -Werror -Wextra
LDFLAGS :=

CFLAGS += -I/usr/local/opt/openssl/include
LDFLAGS += -L/usr/local/opt/openssl/lib -lssl -lcrypto

TARGET := server
SRC_DIR := src
BUILD_DIR := build

SRCS := $(wildcard $(SRC_DIR)/*.c)
OBJS := $(SRCS:%=$(BUILD_DIR)/%.o)
DEPS = $(OBJS:%.o=%.d)


all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) $(CFLAGS) $(LDFLAGS) $(OBJS) -o $@

-include $(DEPS)

$(BUILD_DIR)/%.c.o: %.c
	@mkdir -p $(dir $@)
	$(CC) $(CFLAGS) -MMD -c $< -o $@

clean:
	rm -rf $(BUILD_DIR) $(TARGET)

.PHONY: all clean
