package main

import (
	"client/internal/client"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

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
		totalComputeTime += client.Run(*host, *port, *seed, *sizeInMB, *check, *save)
		fmt.Println("----------------------------")
	}
	fmt.Printf("Average compute time: %.3f ms\n", totalComputeTime/float64(*runs))
}
