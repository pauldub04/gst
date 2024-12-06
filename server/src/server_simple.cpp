#include "shared.h"

#include <iostream>
#include <vector>
#include <chrono>

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
    compute(rows, cols, local_matrix, vector, local_result);

    auto end = std::chrono::high_resolution_clock::now();
    double time_taken = std::chrono::duration<double>(end - start).count();

    write_result(output_filename, result, time_taken, false);
    return 0;
}
