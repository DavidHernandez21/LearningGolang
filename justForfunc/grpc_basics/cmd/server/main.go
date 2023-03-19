package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"

	pb "src/github.com/DavidHernandez21/justForfunc/grpc_basics/todo"

	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func main() {
	srv := grpc.NewServer()
	var tasks taskServer
	pb.RegisterTasksServer(srv, tasks)
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("could not listen to :8888: %v", err)
	}
	log.Println("listening on port 8888")
	log.Fatal(srv.Serve(l))
}

type taskServer struct {
	pb.UnimplementedTasksServer
}

type length int64

const (
	sizeOfLength = 8
	dbPath       = "mydb.pb"
)

var endianness = binary.LittleEndian

func (taskServer) Add(ctx context.Context, text *pb.Text) (task *pb.Task, err error) {

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %v", dbPath, err)
	}

	defer func() {

		if errCloseFile := f.Close(); errCloseFile != nil {
			err = fmt.Errorf("could not close file %s: %v", dbPath, errCloseFile)
			task = nil
		}

	}()

	task = &pb.Task{
		Text: text.Text,
		Done: false,
	}

	b, err := proto.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("could not encode task: %v", err)
	}

	if err := binary.Write(f, endianness, length(len(b))); err != nil {
		return nil, fmt.Errorf("could not encode length of message: %v", err)
	}

	writer := bufio.NewWriter(f)

	defer func() {

		if errWriterFlush := writer.Flush(); errWriterFlush != nil {
			err = fmt.Errorf("could not flush the buffer: %v", errWriterFlush)
			task = nil
		}
	}()

	_, err = writer.Write(b)
	if err != nil {
		return nil, fmt.Errorf("could not write task to file: %v", err)
	}

	return task, err
}

func (taskServer) List(ctx context.Context, void *pb.Void) (*pb.TaskList, error) {
	b, err := os.ReadFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %v", dbPath, err)
	}

	var tasks pb.TaskList
	for {
		if len(b) == 0 {
			tasks.Tasklist[len(tasks.Tasklist)-1].Done = true
			return &tasks, nil
		}
		if len(b) < sizeOfLength {
			return nil, fmt.Errorf("remaining odd %d bytes, what to do?", len(b))
		}

		var l length

		if err := binary.Read(bytes.NewReader(b[:sizeOfLength]), endianness, &l); err != nil {
			return nil, fmt.Errorf("could not decode message length: %v", err)
		}

		b = b[sizeOfLength:]

		var task pb.Task
		if err := proto.Unmarshal(b[:l], &task); err != nil {
			return nil, fmt.Errorf("could not read task: %v", err)
		}
		b = b[l:]
		tasks.Tasklist = append(tasks.Tasklist, &task)
	}
}
