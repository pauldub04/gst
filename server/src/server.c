#include "compute.h"
#include "utils.h"

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <arpa/inet.h>

#define MAX_CONNECTIONS 128

int recv_data(int client_fd, int* rows, int* cols, int*** matrix, int** vector, unsigned char* client_hash) {
    if (
        recv_full(client_fd, rows, sizeof(int)) != sizeof(int) ||
        recv_full(client_fd, cols, sizeof(int)) != sizeof(int)
    ) {
        error_return("reading data");
    }

    *matrix = (int**) malloc(*rows * sizeof(int*));
    *vector = (int*) malloc(*cols * sizeof(int));
    for (int i = 0; i < *rows; i++) {
        (*matrix)[i] = (int*) malloc(*cols * sizeof(int));

        if (recv_full(client_fd, (*matrix)[i], *cols * sizeof(int)) != (ssize_t)(*cols * sizeof(int))) {
            error_return("reading data");
        }
    }

    if (
        recv_full(client_fd, *vector, *cols * sizeof(int)) != (ssize_t)(*cols * sizeof(int)) ||
        recv_full(client_fd, client_hash, SHA256_DIGEST_LENGTH) != SHA256_DIGEST_LENGTH
    ) {
        error_return("reading data");
    }
    return 0;
}

int send_data(int client_fd, int* result, int rows, int cols, double compute_time) {
    int data_size = rows * cols * sizeof(int) + cols * sizeof(int);
    if (
        send_full(client_fd, result, rows * sizeof(int)) != (ssize_t)(rows * sizeof(int)) ||
        send_full(client_fd, &compute_time, sizeof(double)) != sizeof(double) ||
        send_full(client_fd, &data_size, sizeof(int)) != sizeof(int)
    ) {
        error_return("sending data");
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

int check_hash(int rows, int cols, int** matrix, int* vector, unsigned char* client_hash) {
    unsigned char server_hash[EVP_MAX_MD_SIZE];
    unsigned int hash_len = calculate_sha256(rows, cols, matrix, vector, server_hash);
    print_sha256(server_hash, hash_len);

    if (memcmp(client_hash, server_hash, SHA256_DIGEST_LENGTH) != 0) {
        error_return("hash check failed");
    }
    info("hash check passed")
    return 0;
}

int process_client(int client_fd) {
    int rows = 0;
    int cols = 0;
    int** matrix = NULL;
    int* vector  = NULL;
    int* result  = NULL;
    unsigned char client_hash[EVP_MAX_MD_SIZE];
    int err = 0;

    if ((err = recv_data(client_fd, &rows, &cols, &matrix, &vector, client_hash)) < 0) {
        free_resources(rows, cols, &matrix, &vector, &result);
        return err;
    }

    if ((err = check_hash(rows, cols, matrix, vector, client_hash)) < 0) {
        free_resources(rows, cols, &matrix, &vector, &result);
        return err;
    }

    result = (int*) malloc(rows * sizeof(int));
    clock_t start = clock();
    compute(rows, cols, matrix, vector, result);
    double compute_time = (double)(clock() - start) / CLOCKS_PER_SEC;
    
    if ((err = send_data(client_fd, result, rows, cols, compute_time)) < 0) {
        free_resources(rows, cols, &matrix, &vector, &result);
        return err;
    }

    free_resources(rows, cols, &matrix, &vector, &result);
    return err;
}

int main(int argc, char *argv[]) {
    if (argc != 2) {
        printf("Usage: %s <port>\n", argv[0]);
        return 1;
    }
    int port = atoi(argv[1]);
    int server_fd = init_server(port, MAX_CONNECTIONS);

    while (1) {
        struct sockaddr_in client;
        socklen_t client_len = sizeof(client);
        int client_fd = accept(server_fd, (struct sockaddr *) &client, &client_len);
        if (client_fd < 0) {
            error_continue("accept");
        }

        info("accepting socket %d", client_fd);
        if (process_client(client_fd) == 0) {
            info("success");
        }
        info("closing socket %d", client_fd);

        close(client_fd);
        printf("\n");
    }
    
    return 0;
}
