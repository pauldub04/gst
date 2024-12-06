package main

import (
	"bufio"
	"client/internal/client"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

func readTimeFromFile(filename string) (float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return strconv.ParseFloat(scanner.Text(), 64)
	}
	return 0, scanner.Err()
}

func main() {
	filename := pflag.String("filename", "input", "")
	seed := pflag.Uint64("seed", 0, "")
	runs := pflag.Int32("runs", 1, "")
	np := pflag.Int32("np", 4, "")
	pflag.Parse()

	if *seed == 0 {
		*seed = uint64(time.Now().UnixNano() % 100000)
	}
	fmt.Printf("Using seed %d\n", *seed)

	sizes := []int32{1, 5, 10, 50, 100, 500}
	// sizes := []int32{1, 5, 10, 50, 100, 500, 1000}
	results := make([]float64, 0)
	fmt.Printf("Using %d processes\n", *np)

	for _, sizeInMB := range sizes {

		var totalComputeTime float64 = 0
		for i := int32(0); i < *runs; i++ {
			fmt.Printf("Run %d of %d, Size %d; ", i+1, *runs, sizeInMB)
			client.Run(*filename, *seed, sizeInMB)

			cmd := exec.Command("mpirun", "-np", strconv.Itoa(int(*np)), "./compute", *filename, "output")
			// cmd := exec.Command("./compute", *filename, "output")
			err := cmd.Run()
			if err != nil {
				log.Fatal("exec error: ", err)
			}

			computeTime, err := readTimeFromFile("time")
			if err != nil {
				log.Fatal("readTimeFromFile error: ", err)
			}

			computeTime *= 1000.0
			fmt.Printf("Compute time: %.3f ms\n", computeTime)
			totalComputeTime += computeTime
		}
		results = append(results, totalComputeTime/float64(*runs))
		fmt.Println("------------------------------")
	}

	fmt.Println("RESULTS")
	for _, result := range results {
		fmt.Printf(strings.Replace(fmt.Sprintf("%.3f\n", result), ".", ",", 1))
	}
}
