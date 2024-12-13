#include "shared.h"

#include <mpi.h>
#include <iostream>
#include <vector>
#include <cassert>
#include <cuda_runtime.h>

const int gpu_rank = 0;
const int gpu_weight = 10;
const int cpu_weight = 1;

double start = 0;
double end = 0;

void compute_cpu(int rows, int cols, const std::vector<int>& flat_matrix, const std::vector<int>& vector, std::vector<int>& local_result) {
    assert(int(flat_matrix.size()) == rows * cols);
    assert(int(local_result.size()) == rows);
    assert(int(vector.size()) == cols);

    start = MPI_Wtime();
    for (int i = 0; i < rows; ++i) {
        for (int j = 0; j < cols; ++j) {
            local_result[i] += flat_matrix[i * cols + j] * vector[j];
        }
    }
    end = MPI_Wtime();
}

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

    int threadsPerBlock = 256;
    int blocksPerGrid = (rows + threadsPerBlock - 1) / threadsPerBlock;

    start = MPI_Wtime();
    compute_cuda<<<blocksPerGrid, threadsPerBlock>>>(d_matrix, d_vector, d_result, rows, cols);
    cudaDeviceSynchronize();
    end = MPI_Wtime();

    cudaMemcpy(result.data(), d_result, rows * sizeof(int), cudaMemcpyDeviceToHost);
    cudaFree(d_matrix);
    cudaFree(d_vector);
    cudaFree(d_result);
    assert(cudaGetLastError() == cudaSuccess);
}

int main(int argc, char** argv) {
    MPI_Init(&argc, &argv);

    int world_rank, world_size;
    MPI_Comm_rank(MPI_COMM_WORLD, &world_rank);
    MPI_Comm_size(MPI_COMM_WORLD, &world_size);

    if (argc < 3) {
        if (world_rank == 0) {
            std::cerr << "Usage: " << argv[0] << " <input_file> <output_file>" << std::endl;
        }
        MPI_Finalize();
        return 1;
    }

    std::string input_filename = argv[1];
    std::string output_filename = argv[2];

    int rows, cols;
    std::vector<int> matrix, vector;
    std::vector<int> result;

    if (world_rank == 0) {
        read_data(input_filename, rows, cols, matrix, vector);
        result.assign(rows, 0);
    }

    MPI_Bcast(&rows, 1, MPI_INT, 0, MPI_COMM_WORLD);
    MPI_Bcast(&cols, 1, MPI_INT, 0, MPI_COMM_WORLD);

    if (world_rank != 0) {
        vector.assign(cols, 0);
    }
    MPI_Bcast(vector.data(), cols, MPI_INT, 0, MPI_COMM_WORLD);

    int total_weight = gpu_weight + (world_size - 1) * cpu_weight;
    std::vector<int> local_rows(world_size);
    int offset = 0;
    for (int i = 0; i < world_size; ++i) {
        int weight = (i == gpu_rank) ? gpu_weight : cpu_weight;
        local_rows[i] = rows * weight / total_weight;
        if (i < rows % total_weight) {
            ++local_rows[i];
        }
        offset += local_rows[i];
    }

    std::vector<int> sendcounts(world_size);
    std::vector<int> displs(world_size);
    offset = 0;
    for (int i = 0; i < world_size; ++i) {
        sendcounts[i] = local_rows[i] * cols;
        displs[i] = offset;
        offset += sendcounts[i];
    }

    std::vector<int> local_matrix(local_rows[world_rank] * cols);
    MPI_Scatterv(matrix.data(), sendcounts.data(), displs.data(), MPI_INT,
        local_matrix.data(), sendcounts[world_rank], MPI_INT, 0, MPI_COMM_WORLD);

    std::vector<int> local_result(local_rows[world_rank], 0);

    if (world_rank == gpu_rank) {
        run_gpu(local_rows[world_rank], cols, local_matrix, vector, local_result);
    } else {
        compute_simple(local_rows[world_rank], cols, local_matrix, vector, local_result);
    }

    std::vector<int> displs_result(world_size);
    offset = 0;
    for (int i = 0; i < world_size; ++i) {
        displs_result[i] = offset;
        offset += local_rows[i];
    }

    MPI_Gatherv(local_result.data(), local_rows[world_rank], MPI_INT,
        result.data(), local_rows.data(), displs_result.data(), MPI_INT, 0, MPI_COMM_WORLD);

    double global_start_time;
    MPI_Reduce(&start, &global_start_time, 1, MPI_DOUBLE, MPI_MIN, 0, MPI_COMM_WORLD);
    double global_end_time;
    MPI_Reduce(&end, &global_end_time, 1, MPI_DOUBLE, MPI_MAX, 0, MPI_COMM_WORLD);

    if (world_rank == 0) {
        double time_taken = global_end_time - global_start_time;
        write_result(output_filename, result, time_taken, true);
    }

    MPI_Finalize();
    return 0;
}