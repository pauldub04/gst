#include "shared.h"

#include "mpi.h"
#include <iostream>
#include <vector>
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

    int rows_per_proc = rows / world_size;
    int remainder = rows % world_size;

    std::vector<int> sendcounts(world_size);
    std::vector<int> displs(world_size);
    std::vector<int> local_rows(world_size);

    int offset = 0;
    for (int i = 0; i < world_size; ++i) {
        local_rows[i] = rows_per_proc + (i < remainder);
        sendcounts[i] = local_rows[i] * cols;
        displs[i] = offset;
        offset += sendcounts[i];
    }

    std::vector<int> local_matrix(local_rows[world_rank] * cols);
    MPI_Scatterv(matrix.data(), sendcounts.data(), displs.data(), MPI_INT,
        local_matrix.data(), sendcounts[world_rank], MPI_INT, 0, MPI_COMM_WORLD);

    std::vector<int> local_result(local_rows[world_rank], 0);

    double local_start_time = MPI_Wtime();
    compute_simple(local_rows[world_rank], cols, local_matrix, vector, local_result);
    double local_end_time = MPI_Wtime();

    std::vector<int> displs_result(world_size);
    for (int i = 0; i < world_size; ++i) {
        displs_result[i] = (i == 0) ? 0 : displs_result[i-1] + local_rows[i-1];
    }

    MPI_Gatherv(local_result.data(), local_rows[world_rank], MPI_INT,
        result.data(), local_rows.data(), displs_result.data(), MPI_INT, 0, MPI_COMM_WORLD);

    double global_start_time;
    MPI_Reduce(&local_start_time, &global_start_time, 1, MPI_DOUBLE, MPI_MIN, 0, MPI_COMM_WORLD);
    double global_end_time;
    MPI_Reduce(&local_end_time, &global_end_time, 1, MPI_DOUBLE, MPI_MAX, 0, MPI_COMM_WORLD);

    if (world_rank == 0) {
        double time_taken = global_end_time - global_start_time;
        write_result(output_filename, result, time_taken, true);
    }

    MPI_Finalize();
    return 0;
}