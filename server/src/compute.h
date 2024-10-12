#ifndef COMPUTE_H
#define COMPUTE_H

#include <omp.h>

#define THREAD_NUM 2

void compute(int rows, int cols, int** matrix, int* vector, int* result);

#endif /* COMPUTE_H */
