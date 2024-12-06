#pragma once
#include <iostream>
#include <vector>
#include <fstream>
#include <cassert>
#include <omp.h>


void read_data(const std::string &filename, int &rows, int &cols, std::vector<int> &matrix, std::vector<int> &vector) {
    std::ifstream fin(filename);
    if (!fin) {
        std::cerr << "Error reading file " << filename << std::endl;
        exit(1);
    }

    fin >> rows >> cols;
    matrix.resize(rows * cols);
    vector.resize(cols);
    for (int i = 0; i < rows * cols; ++i) {
        fin >> matrix[i];
    }
    for (int i = 0; i < cols; ++i) {
        fin >> vector[i];
    }
}

void write_result(const std::string &filename, const std::vector<int> &result, double time_taken, bool save = true) {
    std::ofstream ftime("time");
    if (!ftime) {
        std::cerr << "Error writing to file time" << std::endl;
        exit(1);
    }
    ftime << time_taken << std::endl;

    if (save) {
        std::ofstream fout(filename);
        if (!fout) {
            std::cerr << "Error writing to file " << filename << std::endl;
            exit(1);
        }

        for (const auto &val : result) {
            fout << val << " ";
        }
        fout << std::endl;
    }
}

void compute(int rows, int cols, const std::vector<int>& flat_matrix, const std::vector<int>& vector, std::vector<int>& local_result) {
    assert(int(flat_matrix.size()) == rows * cols);
    assert(int(local_result.size()) == rows);
    assert(int(vector.size()) == cols);

    // #pragma omp parallel for num_threads(4)
    for (int i = 0; i < rows; ++i) {
        for (int j = 0; j < cols; ++j) {
            local_result[i] += flat_matrix[i * cols + j] * vector[j];
        }
    }
}
