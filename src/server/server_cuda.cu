#include <iostream>
#include <vector>
#include <chrono>
#include <cuda_runtime.h>
#include "shared.h"

std::chrono::high_resolution_clock::time_point start{};
std::chrono::high_resolution_clock::time_point end{};

__global__ void compute_cuda(int* d_matrix, int* d_vector, int* d_result, int rows, int cols) {
    int row = blockIdx.x * blockDim.x + threadIdx.x;
    if (row < rows) {
        int sum = 0;
        for (int j = 0; j < cols; ++j) {
            sum += d_matrix[row * cols + j] * d_vector[j];
        }
        d_result[row] = sum;
    }
}

void run_gpu(int rows, int cols, const std::vector<int>& flat_matrix, const std::vector<int>& vector, std::vector<int>& result) {
    int* d_matrix;
    int* d_vector;
    int* d_result;

    cudaMalloc((void**)&d_matrix, flat_matrix.size() * sizeof(int));
    cudaMalloc((void**)&d_vector, vector.size() * sizeof(int));
    cudaMalloc((void**)&d_result, rows * sizeof(int));

    cudaMemcpy(d_matrix, flat_matrix.data(), flat_matrix.size() * sizeof(int), cudaMemcpyHostToDevice);
    cudaMemcpy(d_vector, vector.data(), vector.size() * sizeof(int), cudaMemcpyHostToDevice);

    int threadsPerBlock = 128;
    int blocksPerGrid = (rows + threadsPerBlock - 1) / threadsPerBlock;

    start = std::chrono::high_resolution_clock::now();
    compute_cuda<<<blocksPerGrid, threadsPerBlock>>>(d_matrix, d_vector, d_result, rows, cols);
    cudaDeviceSynchronize();
    end = std::chrono::high_resolution_clock::now();

    cudaMemcpy(result.data(), d_result, rows * sizeof(int), cudaMemcpyDeviceToHost);
    cudaFree(d_matrix);
    cudaFree(d_vector);
    cudaFree(d_result);
    assert(cudaGetLastError() == cudaSuccess);
}

int main(int argc, char** argv) {
    if (argc < 3) {
        std::cerr << "Usage: " << argv[0] << " <input_file> <output_file>" << std::endl;
        return 1;
    }

    std::string input_filename = argv[1];
    std::string output_filename = argv[2];

    int rows, cols;
    std::vector<int> matrix, vector;
    read_data(input_filename, rows, cols, matrix, vector);
    std::vector<int> result(rows, 0);

    run_gpu(rows, cols, matrix, vector, result);

    double time_taken = std::chrono::duration<double>(end - start).count();
    write_result(output_filename, result, time_taken, true);
    return 0;
}
