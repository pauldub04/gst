package main

import (
	"bytes"
	"client/internal/client"
	"client/internal/config"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"reflect"
	"time"

	"github.com/spf13/pflag"
)

func generateRequest(seed uint64, sizeInMB int32) *client.TReq {
	totalBytes := sizeInMB * 1024 * 1024
	elements := totalBytes / config.ElementByteSize
	rows := int32(math.Sqrt(float64(elements)))
	cols := rows

	matrix := client.GenerateMatrix(seed, rows, cols, config.ElementMin, config.ElementMax)
	vector := client.GenerateVector(seed, cols, config.ElementMin, config.ElementMax)
	dataHash := client.CalculateHash(matrix, vector)

	fmt.Println("Generated SHA-256 hash:")
	fmt.Printf("%x\n", dataHash)

	return &client.TReq{
		Rows:   rows,
		Cols:   cols,
		Matrix: matrix,
		Vector: vector,
		Hash:   dataHash,
	}
}

func saveResult(clientTime float64, rsp *client.TRsp) {
	if err := client.SaveResults(config.ResultFile, rsp.Result); err != nil {
		log.Fatal("Error saving results:", err)
	}
	if err := client.SaveStatistics(config.StatisticsFile, clientTime, rsp.ComputeTime, rsp.DataSize); err != nil {
		log.Fatal("Error saving statistics:", err)
	}
}

func printResult() {
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

func run(host string, port string, seed uint64, sizeInMB int32, check bool, save bool) float64 {
	clientStartTime := time.Now()

	if seed == 0 {
		seed = uint64(time.Now().UnixNano() % 100000)
	}
	fmt.Printf("Using seed %d\n", seed)

	req := generateRequest(seed, sizeInMB)

	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	if err := client.SendData(conn, req); err != nil {
		log.Fatal("Error sending data:", err)
	}

	rsp, err := client.RecvData(conn, req.Rows)
	if err != nil {
		log.Fatal("Error receiving data:", err)
	}
	clientTime := time.Since(clientStartTime).Seconds()

	fmt.Printf("Client time: %.3f seconds\n", clientTime)
	fmt.Printf("Compute time: %.3f ms\n", rsp.ComputeTime*1000.0)
	fmt.Printf("Processed data size: %.3f mb\n", float64(rsp.DataSize)/(1024*1024))
	if save {
		saveResult(clientTime, rsp)
		printResult()
	}

	if check {
		fmt.Println("Checking result using local compute")
		localResut := client.MultiplyMatrixVector(req.Matrix, req.Vector)
		if reflect.DeepEqual(localResut, rsp.Result) {
			fmt.Println("Client and server results match!")
		} else {
			fmt.Println("Error! Client and server results do NOT match!")
		}
	}

	return rsp.ComputeTime * 1000.0
}

func main() {
	sizeInMB := pflag.Int32("size", 0, "The size of data in MB to process")
	host := pflag.String("host", "", "The server host to connect to")
	port := pflag.String("port", "", "The port to connect to on the server")
	seed := pflag.Uint64("seed", 0, "Seed for random number generation (optional)")
	check := pflag.Bool("check", false, "Check result using local compute (optional)")
	save := pflag.Bool("save", true, "Save result to file (optional)")
	runs := pflag.Int32("runs", 1, "Number of runs")
	pflag.Parse()

	if *sizeInMB == 0 || *host == "" || *port == "" {
		fmt.Printf("Usage: %s --size <size_in_mb> --host <host> --port <port>\n", os.Args[0])
		return
	}

	var totalComputeTime float64 = 0
	for i := int32(0); i < *runs; i++ {
		fmt.Printf("RUN №%d\n", i+1)
		totalComputeTime += run(*host, *port, *seed, *sizeInMB, *check, *save)
		fmt.Println("----------------------------")
	}
	fmt.Printf("Average compute time: %.3f ms\n", totalComputeTime/float64(*runs))
}
