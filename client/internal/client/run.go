package client

import (
	"bytes"
	"client/internal/config"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"reflect"
	"time"
)

func GenerateRequest(seed uint64, sizeInMB int32) *TReq {
	totalBytes := sizeInMB * 1024 * 1024
	elements := totalBytes / config.ElementByteSize
	rows := int32(math.Sqrt(float64(elements)))
	cols := rows

	matrix := GenerateMatrix(seed, rows, cols, config.ElementMin, config.ElementMax)
	vector := GenerateVector(seed, cols, config.ElementMin, config.ElementMax)
	dataHash := CalculateHash(matrix, vector)

	fmt.Println("Generated SHA-256 hash:")
	fmt.Printf("%x\n", dataHash)

	return &TReq{
		Rows:   rows,
		Cols:   cols,
		Matrix: matrix,
		Vector: vector,
		Hash:   dataHash,
	}
}

func SaveResult(clientTime float64, rsp *TRsp) {
	if err := SaveResults(config.ResultFile, rsp.Result); err != nil {
		log.Fatal("Error saving results:", err)
	}
	if err := SaveStatistics(config.StatisticsFile, clientTime, rsp.ComputeTime, rsp.DataSize); err != nil {
		log.Fatal("Error saving statistics:", err)
	}
}

func PrintResult() {
	fmt.Println()
	fmt.Printf("%s:\n", config.ResultFile)
	resultData, err := os.ReadFile(config.ResultFile)
	if err != nil {
		log.Fatal("Error reading file", err)
	}
	fmt.Printf("%d elements saved\n", len(bytes.Split(resultData, []byte(" "))))

	fmt.Println()
	fmt.Printf("%s:\n", config.StatisticsFile)
	statisticsData, err := os.ReadFile(config.StatisticsFile)
	if err != nil {
		log.Fatal("Error reading file", err)
	}
	fmt.Println(string(statisticsData))
}

func Run(host string, port string, seed uint64, sizeInMB int32, check bool, save bool) float64 {
	clientStartTime := time.Now()

	if seed == 0 {
		seed = uint64(time.Now().UnixNano() % 100000)
	}
	fmt.Printf("Using seed %d\n", seed)

	req := GenerateRequest(seed, sizeInMB)

	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	if err := SendData(conn, req); err != nil {
		log.Fatal("Error sending data:", err)
	}

	rsp, err := RecvData(conn, req.Rows)
	if err != nil {
		log.Fatal("Error receiving data:", err)
	}
	clientTime := time.Since(clientStartTime).Seconds()

	fmt.Printf("Client time: %.3f seconds\n", clientTime)
	fmt.Printf("Compute time: %.3f ms\n", rsp.ComputeTime*1000.0)
	fmt.Printf("Processed data size: %.3f mb\n", float64(rsp.DataSize)/(1024*1024))
	if save {
		SaveResult(clientTime, rsp)
		PrintResult()
	}

	if check {
		fmt.Println("Checking result using local compute")
		localResut := MultiplyMatrixVector(req.Matrix, req.Vector)
		if reflect.DeepEqual(localResut, rsp.Result) {
			fmt.Println("Client and server results match!")
		} else {
			fmt.Println("Error! Client and server results do NOT match!")
		}
	}

	return rsp.ComputeTime * 1000.0
}
