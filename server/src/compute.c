#include "compute.h"

void compute(int rows, int cols, int** matrix, int* vector, int* result) {
    for (int i = 0; i < rows; ++i) {
        result[i] = 0;
        for (int j = 0; j < cols; ++j) {
            result[i] += matrix[i][j] * vector[j];
        }
    }
}
