#include "shared.h"

#include <mpi.h>
#include <iostream>
#include <vector>
#include <chrono>
#include <cassert>

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
    std::vector<int> matrix, vector, local_matrix;
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

    int rows_per_proc = rows / world_size;
    int extra_rows = rows % world_size;

    int local_rows = rows_per_proc;
    if (world_rank < extra_rows) {
        ++local_rows;
    }
    local_matrix.resize(local_rows * cols);

    if (world_rank == 0) {
        int offset = 0;
        for (int i = 0; i < world_size; ++i) {
            int count = (i < extra_rows ? rows_per_proc + 1 : rows_per_proc) * cols;
            if (i == 0) {
                std::copy(matrix.begin(), matrix.begin() + count, local_matrix.begin());
            } else {
                MPI_Send(matrix.data() + offset, count, MPI_INT, i, 0, MPI_COMM_WORLD);
            }
            offset += count;
        }
    } else {
        MPI_Recv(local_matrix.data(), local_matrix.size(), MPI_INT, 0, 0, MPI_COMM_WORLD, MPI_STATUS_IGNORE);
    }

    auto start = std::chrono::high_resolution_clock::now();

    std::vector<int> local_result(local_rows, 0);
    compute(local_rows, cols, local_matrix, vector, local_result);

    if (world_rank == 0) {
        int offset = local_rows;
        std::copy(local_result.begin(), local_result.end(), result.begin());
        for (int i = 1; i < world_size; ++i) {
            int recv_count = (i < extra_rows ? rows_per_proc + 1 : rows_per_proc);
            MPI_Recv(result.data() + offset, recv_count, MPI_INT, i, 0, MPI_COMM_WORLD, MPI_STATUS_IGNORE);
            offset += recv_count;
        }
    } else {
        MPI_Send(local_result.data(), local_rows, MPI_INT, 0, 0, MPI_COMM_WORLD);
    }

    auto end = std::chrono::high_resolution_clock::now();
    double time_taken = std::chrono::duration<double>(end - start).count();

    if (world_rank == 0) {
        write_result(output_filename, result, time_taken, false);
    }

    MPI_Finalize();
    return 0;
}
