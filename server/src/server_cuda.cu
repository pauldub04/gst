#include <iostream>
#include <vector>
#include <chrono>
#include <cuda_runtime.h>
#include "shared.h"

__global__ void matrix_vector_multiply(const int* d_matrix, const int* d_vector, int* d_result, int rows, int cols) {
    int row = blockIdx.x * blockDim.x + threadIdx.x;
    if (row < rows) {
        int sum = 0;
        for (int j = 0; j < cols; ++j) {
            sum += d_matrix[row * cols + j] * d_vector[j];
        }
        d_result[row] = sum;
    }
}

void compute_cuda(int rows, int cols, const std::vector<int>& flat_matrix, const std::vector<int>& vector, std::vector<int>& result) {
    int *d_matrix, *d_vector, *d_result;
    cudaMalloc((void**)&d_matrix, rows * cols * sizeof(int));
    cudaMalloc((void**)&d_vector, cols * sizeof(int));
    cudaMalloc((void**)&d_result, rows * sizeof(int));

    cudaMemcpy(d_matrix, flat_matrix.data(), rows * cols * sizeof(int), cudaMemcpyHostToDevice);
    cudaMemcpy(d_vector, vector.data(), cols * sizeof(int), cudaMemcpyHostToDevice);

    int blockSize = 256;
    int numBlocks = (rows + blockSize - 1) / blockSize;
    matrix_vector_multiply<<<numBlocks, blockSize>>>(d_matrix, d_vector, d_result, rows, cols);
    cudaMemcpy(result.data(), d_result, rows * sizeof(int), cudaMemcpyDeviceToHost);

    cudaFree(d_matrix);
    cudaFree(d_vector);
    cudaFree(d_result);
}

int main(int argc, char** argv) {
    if (argc < 3) {
        std::cerr << "Usage: " << argv[0] << " <input_file> <output_file>" << std::endl;
        return 1;
    }

    std::string input_filename = argv[1];
    std::string output_filename = argv[2];

    int rows, cols;
    std::vector<int> matrix, vector, local_matrix;
    std::vector<int> result;

    read_data(input_filename, rows, cols, matrix, vector);
    result.assign(rows, 0);
    vector.assign(cols, 0);

    auto start = std::chrono::high_resolution_clock::now();

    std::vector<int> local_result(rows, 0);
    compute_cuda(rows, cols, local_matrix, vector, local_result);

    auto end = std::chrono::high_resolution_clock::now();
    double time_taken = std::chrono::duration<double>(end - start).count();

    write_result(output_filename, result, time_taken, false);
    return 0;
}
