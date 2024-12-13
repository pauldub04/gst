package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

type TReq struct {
	Rows   int32
	Cols   int32
	Matrix [][]int32
	Vector []int32
}

type TRsp struct {
	Result      []int32
	ComputeTime float64
	DataSize    int32
}

func SendData(conn net.Conn, req *TReq) error {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, req.Rows)
	binary.Write(&buf, binary.LittleEndian, req.Cols)
	for _, row := range req.Matrix {
		for _, item := range row {
			binary.Write(&buf, binary.LittleEndian, item)
		}
	}
	for _, item := range req.Vector {
		binary.Write(&buf, binary.LittleEndian, item)
	}

	if _, err := conn.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func RecvData(conn net.Conn, rows int32) (*TRsp, error) {
	var rsp TRsp
	rsp.Result = make([]int32, rows)
	if err := binary.Read(conn, binary.LittleEndian, &rsp.Result); err != nil {
		return nil, err
	}

	if err := binary.Read(conn, binary.LittleEndian, &rsp.ComputeTime); err != nil {
		return nil, err
	}

	if err := binary.Read(conn, binary.LittleEndian, &rsp.DataSize); err != nil {
		return nil, err
	}

	return &rsp, nil
}

func SaveDataToFile(filename string, req *TReq) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "%d %d\n", req.Rows, req.Cols)

	for _, row := range req.Matrix {
		for _, value := range row {
			fmt.Fprintf(file, "%d ", value)
		}
		fmt.Fprintln(file)
	}

	for _, value := range req.Vector {
		fmt.Fprintf(file, "%d ", value)
	}
	fmt.Fprintln(file)
	return nil
}

func SaveVectorToFile(filename string, vector []int32) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, value := range vector {
		fmt.Fprintf(file, "%d ", value)
	}
	fmt.Fprintln(file)
	return nil
}
