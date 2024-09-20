#include "utils.h"

int init_server(int port, int connections) {
    int server_fd;
    struct sockaddr_in server;

    server_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (server_fd < 0) {
        fatal("socket\n");
    }

    server.sin_family = AF_INET;
    server.sin_port = htons(port);
    server.sin_addr.s_addr = htonl(INADDR_ANY);

    int opt_val = 1;
    setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR, &opt_val, sizeof(opt_val));
    if (bind(server_fd, (struct sockaddr*)&server, sizeof(server)) < 0) {
        fatal("bind\n");
    }
    if (listen(server_fd, connections) < 0) {
        fatal("listen\n");
    }

    printf("Server is listening on port %d\n", port);
    return server_fd;
}

ssize_t read_full(int sockfd, void* buffer, size_t length) {
    size_t total_read = 0;
    ssize_t bytes_read;

    while (total_read < length) {
        bytes_read = recv(sockfd, (char*) buffer+total_read, length-total_read, 0);
        if (bytes_read < 0) {
            return -1;
        } else if (bytes_read == 0) {
            break;
        }
        total_read += bytes_read;
    }
    return total_read;
}

unsigned int calculate_sha256(int rows, int cols, int** matrix, int* vector, unsigned char* output_hash) {
    EVP_MD_CTX* mdctx = EVP_MD_CTX_new();
    const EVP_MD* md = EVP_sha256();

    EVP_DigestInit_ex(mdctx, md, NULL);
    for (int i = 0; i < rows; i++) {
        EVP_DigestUpdate(mdctx, matrix[i], cols * sizeof(int));
    }
    EVP_DigestUpdate(mdctx, vector, cols * sizeof(int));

    unsigned int hash_len;
    EVP_DigestFinal_ex(mdctx, output_hash, &hash_len);
    EVP_MD_CTX_free(mdctx);
    return hash_len;
}

void print_sha256(unsigned char* hash, size_t len) {
    info("calculated SHA-256 hash:");
    for (size_t i = 0; i < len; i++) {
        printf("%02x", hash[i]);
    }
    printf("\n");
}
