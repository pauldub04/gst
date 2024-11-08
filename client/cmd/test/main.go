package main

import (
	"client/internal/client"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

func main() {
	host := pflag.String("host", "", "The server host to connect to")
	port := pflag.String("port", "", "The port to connect to on the server")
	seed := pflag.Uint64("seed", 0, "Seed for random number generation (optional)")
	runs := pflag.Int32("runs", 1, "Number of runs")
	pflag.Parse()

	if *host == "" || *port == "" {
		fmt.Printf("Usage: %s --host <host> --port <port>\n", os.Args[0])
		return
	}

	sizes := []int32{1, 5, 10, 50, 100, 500, 1000}
	results := make([]float64, 0)

	for _, sizeInMB := range sizes {
		var totalComputeTime float64 = 0
		for i := int32(0); i < *runs; i++ {
			totalComputeTime += client.Run(*host, *port, *seed, sizeInMB, false, false)
		}
		results = append(results, totalComputeTime/float64(*runs))
	}

	for _, result := range results {
		fmt.Printf(strings.Replace(fmt.Sprintf("%.3f\n", result), ".", ",", 1))
	}
}
