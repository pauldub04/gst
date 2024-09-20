#include "compute.h"
#include "utils.h"

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <arpa/inet.h>

#define MAX_CONNECTIONS 128

int read_data(int client_fd, int* rows, int* cols, int*** matrix, int** vector, unsigned char* client_hash) {
    if (read_full(client_fd, rows, sizeof(int)) != sizeof(int)) {
        error_return("reading rows");
    }
    if (read_full(client_fd, cols, sizeof(int)) != sizeof(int)) {
        error_return("reading cols");
    }

    *matrix = (int**) malloc(*rows * sizeof(int*));
    for (int i = 0; i < *rows; i++) {
        (*matrix)[i] = (int*) malloc(*cols * sizeof(int));

        if (read_full(client_fd, (*matrix)[i], *cols * sizeof(int)) != (ssize_t)(*cols * sizeof(int))) {
            error_return("reading matrix row");
        }
    }

    *vector = (int*) malloc(*cols * sizeof(int));
    if (read_full(client_fd, *vector, *cols * sizeof(int)) != (ssize_t)(*cols * sizeof(int))) {
        error_return("reading vector");
    }
    if (read_full(client_fd, client_hash, SHA256_DIGEST_LENGTH) != SHA256_DIGEST_LENGTH) {
        error_return("reading hash");
    }
    return 0;
}

void free_resources(int rows, int cols, int*** matrix, int** vector, int** result) {
    if (matrix != NULL && *matrix != NULL) {
        for (int i = 0; i < rows; ++i) {
            if ((*matrix)[i] != NULL) {
                free((*matrix)[i]);
            }
        }
        free(*matrix);
    }
    
    if (vector != NULL && *vector != NULL) {
        free(*vector);
    }
    if (result != NULL && *result != NULL) {
        free(*result);
    }
}

int process_client(int client_fd) {
    int rows = 0, cols = 0;
    int** matrix = NULL;
    int* vector  = NULL;
    unsigned char client_hash[EVP_MAX_MD_SIZE];

    int err = read_data(client_fd, &rows, &cols, &matrix, &vector, client_hash);
    if (err < 0) {
        return err;
    }
    int* result = (int*) malloc(rows * sizeof(int));

    unsigned char server_hash[EVP_MAX_MD_SIZE];
    unsigned int hash_len = calculate_sha256(rows, cols, matrix, vector, server_hash);
    print_sha256(server_hash, hash_len);
    
    if (memcmp(client_hash, server_hash, SHA256_DIGEST_LENGTH) != 0) {
        free_resources(rows, cols, &matrix, &vector, &result);
        error_return("hash check failed");
    } else {
        info("hash check passed")
    }

    clock_t start = clock();
    compute(rows, cols, matrix, vector, result);
    clock_t end = clock();

    double compute_time = (double)(end - start) / CLOCKS_PER_SEC;
    info("compute time: %f s", compute_time);

    // send(client_fd, result, rows * sizeof(int), 0);
    // send(client_fd, &exec_time, sizeof(double), 0);

    free_resources(rows, cols, &matrix, &vector, &result);
    return 0;
}

int main(int argc, char *argv[]) {
    if (argc != 3) {
        printf("Usage: %s <port> <return_port>\n", argv[0]);
        return 1;
    }
    int port = atoi(argv[1]);
    int return_port = atoi(argv[2]);

    int server_fd = init_server(port, MAX_CONNECTIONS);

    while (1) {
        struct sockaddr_in client;
        socklen_t client_len = sizeof(client);

        int client_fd = accept(server_fd, (struct sockaddr *) &client, &client_len);
        if (client_fd < 0) {
            error_continue("accept");
        }

        info("accepted socket %d", client_fd);
        int err = process_client(client_fd);
        if (err == 0) {
            info("success");
        }

        info("closing socket %d", client_fd);
        close(client_fd);
        printf("\n");
    }
    
    return 0;
}
