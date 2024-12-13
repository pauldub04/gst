package client

import (
	"client/internal/config"
	"log"
	"math"
)

func GenerateRequest(seed uint64, sizeInMB int32) *TReq {
	totalBytes := sizeInMB * 1024 * 1024
	elements := totalBytes / config.ElementByteSize
	rows := int32(math.Sqrt(float64(elements)))
	cols := rows

	matrix := GenerateMatrix(seed, rows, cols, config.ElementMin, config.ElementMax)
	vector := GenerateVector(seed, cols, config.ElementMin, config.ElementMax)

	return &TReq{
		Rows:   rows,
		Cols:   cols,
		Matrix: matrix,
		Vector: vector,
	}
}

func Run(filename string, seed uint64, sizeInMB int32, runLocal bool, localFilename string) {
	req := GenerateRequest(seed, sizeInMB)

	if err := SaveDataToFile(filename, req); err != nil {
		log.Fatal("Error saving data to file:", err)
	}

	if runLocal {
		result := MultiplyMatrixVector(req.Matrix, req.Vector)
		SaveVectorToFile(localFilename, result)
	}
}
