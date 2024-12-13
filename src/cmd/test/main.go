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
	mpi := pflag.Int32("mpi", -1, "")
	local := pflag.Bool("local", false, "")
	pflag.Parse()

	if *seed == 0 {
		*seed = uint64(time.Now().UnixNano() % 100000)
	}
	fmt.Printf("Using seed %d\n", *seed)

	sizes := []int32{1, 5, 10, 50, 100, 500}
	results := make([]float64, 0)
	if *mpi != -1 {
		fmt.Printf("Using %d mpi processes\n", *mpi)
	}

	for _, sizeInMB := range sizes {

		var totalComputeTime float64 = 0
		for i := int32(0); i < *runs; i++ {
			fmt.Printf("Run %d of %d, Size %d; ", i+1, *runs, sizeInMB)
			client.Run(*filename, *seed, sizeInMB, *local, "local")

			var cmd *exec.Cmd
			if *mpi != -1 {
				cmd = exec.Command("mpirun", "-np", strconv.Itoa(int(*mpi)), "--allow-run-as-root", "./compute", *filename, "output")
			} else {
				cmd = exec.Command("./compute", *filename, "output")
			}

			if err := cmd.Run(); err != nil {
				log.Fatal("exec error: ", err)
			}

			computeTime, err := readTimeFromFile("time")
			if err != nil {
				log.Fatal("readTimeFromFile error: ", err)
			}

			if *local {
				if client.DeepCompareFiles("local", "output") {
					fmt.Print("[OK]; ")
				} else {
					fmt.Print("[ERROR]; ")
				}
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
		fmt.Println(strings.Replace(fmt.Sprintf("%.3f", result), ".", ",", 1))
	}
}
