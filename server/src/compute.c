#include "compute.h"

void compute(int rows, int cols, int** matrix, int* vector, int* result) {
    #pragma omp parallel for num_threads(THREAD_NUM)
    for (int i = 0; i < rows; ++i) {
        int sum = 0;
        for (int j = 0; j < cols; ++j) {
            sum += matrix[i][j] * vector[j];
        }
        result[i] = sum;
    }
}
