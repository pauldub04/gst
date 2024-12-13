package client

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

func CalculateHash(matrix [][]int32, vector []int32) []byte {
	hash := sha256.New()
	for _, row := range matrix {
		for _, item := range row {
			binary.Write(hash, binary.LittleEndian, item)
		}
	}
	for _, item := range vector {
		binary.Write(hash, binary.LittleEndian, item)
	}
	return hash.Sum(nil)
}

func int32SliceToIntSlice(arr []int32) []int {
	res := make([]int, len(arr))
	for i, v := range arr {
		res[i] = int(v)
	}
	return res
}

func SaveResults(file string, result []int32) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, value := range result {
		fmt.Fprintf(f, "%d ", value)
	}
	return nil
}

func SaveStatistics(file string, clientTime, computeTime float64, dataSize int32) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "Client time: %.3f seconds\n", clientTime)
	fmt.Fprintf(f, "Compute time: %.3f ms\n", computeTime*1000.0)
	fmt.Fprintf(f, "Processed data size: %.3f mb\n", float64(dataSize)/(1024*1024))
	return nil
}

const chunkSize = 64000

func DeepCompareFiles(file1, file2 string) bool {
	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}
