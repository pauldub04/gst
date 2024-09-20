#ifndef UTILS_H
#define UTILS_H

#include <unistd.h>
#include <stdio.h>
#include <openssl/evp.h>
#include <openssl/ssl.h>
#include <arpa/inet.h>
#include <netdb.h>


#define fatal(...)          { fprintf(stderr, "[FATAL]: "__VA_ARGS__); fprintf(stderr, "\n"); fflush(stderr); exit(1); }
#define error_continue(...) { fprintf(stderr, "[ERROR]: "__VA_ARGS__); fprintf(stderr, "\n"); fflush(stderr); continue; }
#define error_return(...)   { fprintf(stderr, "[ERROR]: "__VA_ARGS__); fprintf(stderr, "\n"); fflush(stderr); return -1; }
#define info(...)           { fprintf(stdout, "[INFO]: " __VA_ARGS__); fprintf(stdout, "\n"); fflush(stdout); }

int init_server(int port, int connections);
ssize_t read_full(int sockfd, void* buffer, size_t length);

unsigned int calculate_sha256(int rows, int cols, int** matrix, int* vector, unsigned char* output_hash);
void print_sha256(unsigned char* hash, size_t len);

#endif /* UTILS_H */
