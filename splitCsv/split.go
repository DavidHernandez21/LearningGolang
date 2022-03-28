package splitCsv

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

//minFileChunkSize min file chunk size in bytes
const minFileChunkSize = 100

var (
	ErrSmallFileChunkSize = errors.New("file chunk size is too small")
	ErrBigFileChunkSize   = errors.New("file chunk size is bigger than input file")
	// ErrFileAlreadyClosed     = errors.New("file already closed")
	fileSize int64
)

//Splitter struct which contains options for splitting
//FileChunkSize - a size of chunk in bytes, should be set by client
//WithHeader - whether split csv with header (true by default)
type Splitter struct {
	FileChunkSize        int //in bytes
	WithHeader           bool
	bufferSize           int   //in bytes
	cumulativeBufferSize int64 // in bytes
}

//New initializes Splitter struct
func New() Splitter {
	return Splitter{
		WithHeader:           true,
		bufferSize:           os.Getpagesize() * 128,
		cumulativeBufferSize: 0,
	}
}

//Split splits file in smaller chunks
func (s Splitter) Split(inputFilePath, outputDirPath, jsonKeyPath string) ([]string, error) {

	client, err := getClient(jsonKeyPath)

	if err != nil {
		return nil, err
	}

	defer client.Close()

	if s.FileChunkSize < minFileChunkSize {
		return nil, ErrSmallFileChunkSize
	}
	file, err := os.Open(inputFilePath)
	if err != nil {
		msg := fmt.Sprintf("Couldn't open file %s : %v", inputFilePath, err)
		return nil, errors.New(msg)
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		msg := fmt.Sprintf("Couldn't get file stat %s : %v", inputFilePath, err)
		return nil, errors.New(msg)
	}
	fileSize = stat.Size()
	if fileSize <= int64(s.FileChunkSize) {
		return nil, ErrBigFileChunkSize
	}

	bufBulk := make([]byte, s.bufferSize)
	fileName, ext := getFileName(inputFilePath)
	st := state{
		s:             s,
		inputFilePath: inputFilePath,
		fileName:      fileName,
		ext:           ext,
		resultDirPath: prepareResultDirPath(outputDirPath),
		inputFile:     file,
		firstLine:     true,
		chunk:         1,
		bulkBuffer:    bytes.NewBuffer(make([]byte, 0, s.bufferSize)),
	}
	fileNameChanIn := make(chan string, fileSize/int64(s.FileChunkSize)+3)
	fileNameChanOut := make(chan string, fileSize/int64(s.FileChunkSize)+3)
	var numWorkers int = runtime.NumCPU()
	errChanUpload := make(chan error, 1)
	errChanDelete := make(chan error, 1)

	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	wg1.Add(numWorkers)
	wg2.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		uploadFileToBucket("gpjegld01-1812-convintel-ivr-csv-tr-gcf", fileNameChanIn, fileNameChanOut, errChanUpload, &wg1, client)
		deleteFiles(fileNameChanOut, &wg2, ctx, errChanDelete)
	}

	logFileName := "errors.log"
	writeErrorsToFile(errChanUpload, logFileName)

	for {
		//Read bulk from file
		size, errEOF := file.Read(bufBulk)
		if errEOF == io.EOF {
			_, err = st.chunkFile.Write(st.brokenLine)
			if err != nil {
				msg := fmt.Sprintf("Couldn't write chunk file %s : %v", st.chunkFilePath, err)
				return nil, errors.New(msg)
			}
			break
		}
		if err != nil {
			msg := fmt.Sprintf("Couldn't read file bulk %s : %v", inputFilePath, err)
			return nil, errors.New(msg)
		}
		st.fileBuffer = bytes.NewBuffer(bufBulk[:size])
		if len(st.brokenLine) > 0 {
			st.bulkBuffer.Write(st.brokenLine)
			st.brokenLine = []byte{}
		}

		err = readLinesFromBulk(&st, fileNameChanIn)
		if err != nil {
			return nil, err
		}

		err = saveBulkToFile(&st, fileNameChanIn)
		if err != nil {
			return nil, err
		}

	}

	if err := st.chunkFile.Close(); err != nil {
		return nil, fmt.Errorf("could not close file '%v': %v", st.chunkFilePath, err)
	}
	fileNameChanIn <- st.chunkFilePath
	close(fileNameChanIn)

	wg1.Wait()

	close(fileNameChanOut)
	close(errChanUpload)

	wg2.Wait()
	close(errChanDelete)

	return st.result, nil
}

//readLinesFromBulk reads bulk line by line
func readLinesFromBulk(st *state, fileNameChan chan string) error {
	for {
		bytesLine, err := st.fileBuffer.ReadBytes('\n')
		if err == io.EOF {
			st.brokenLine = bytesLine
			break
		}
		if err != nil {
			msg := fmt.Sprintf("Couldn't read bytes from buffer of file %s : %v", st.inputFilePath, err)
			return errors.New(msg)
		}
		if st.firstLine && st.s.WithHeader {
			st.firstLine = false
			st.header = bytesLine
			continue
		}
		st.bulkBuffer.Write(bytesLine)
		if st.s.FileChunkSize < st.s.bufferSize && st.bulkBuffer.Len() >= (st.s.FileChunkSize-len(st.header)) {
			err = saveBulkToFile(st, fileNameChan)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// saveBulkToFile saves lines from bulk to a new file
func saveBulkToFile(st *state, fileNameChan chan string) error {
	st.chunkFilePath = st.resultDirPath + st.fileName + "_" + strconv.Itoa(st.chunk) + "." + st.ext

	_, err := os.Stat(st.chunkFilePath)

	if os.IsNotExist(err) {
		chunkFile, err := os.Create(st.chunkFilePath)
		if err != nil {
			msg := fmt.Sprintf("Couldn't create file %s : %v", st.chunkFilePath, err)
			return errors.New(msg)
		}
		st.setChunkFile(chunkFile)
		_, err = st.chunkFile.Write(st.header)
		if err != nil {
			msg := fmt.Sprintf("Couldn't write header of chunk file %s : %v", st.chunkFilePath, err)
			return errors.New(msg)
		}
		st.result = append(st.result, st.chunkFilePath)
	}
	_, err = st.chunkFile.Write(st.bulkBuffer.Bytes())
	if err != nil {
		msg := fmt.Sprintf("Couldn't write chunk file %s : %v", st.chunkFilePath, err)
		return errors.New(msg)
	}
	stat, _ := os.Stat(st.chunkFilePath)
	if st.chunk == 1 {
		st.s.cumulativeBufferSize += int64(len(st.header))
	}

	if stat.Size() > int64(st.s.FileChunkSize-st.s.bufferSize) {

		st.chunk++
		st.s.cumulativeBufferSize += stat.Size()
		// fmt.Fprintf(os.Stdout, "cummulative buffer size: %v, file size: %v\n", st.s.cumulativeBufferSize, fileSize)
		if st.s.cumulativeBufferSize < int64(fileSize) {
			if err := st.chunkFile.Close(); err != nil {
				return fmt.Errorf("could not close file '%v': %v", st.chunkFile.Name(), err)
			}
			fmt.Fprintf(os.Stdout, "file send to upload in gcs: %v, size: %v\n", st.chunkFile.Name(), stat.Size())
			fileNameChan <- st.chunkFile.Name()

		}

	}

	st.bulkBuffer.Reset()

	return nil
}

//getFileName extracts name and extension from path
func getFileName(path string) (string, string) {
	split := strings.Split(path, "/")
	name := split[len(split)-1]
	split = strings.Split(name, ".")
	ext := split[len(split)-1]
	nSplit := split[:len(split)-1]
	name = strings.Join(nSplit, "")

	return name, ext
}

//prepareResultDirPath adds '/' to the end of path if needed
func prepareResultDirPath(path string) string {
	if path == "" {
		return ""
	}
	p := []byte(path)
	if p[len(p)-1] != '/' {
		p = append(p, '/')
	}

	return string(p)
}

//deletes files from the system, names of the files to delete come from a channel
func deleteFiles(ch chan string, wg *sync.WaitGroup, ctx context.Context, errChan chan<- error) {

	go func() {
		defer wg.Done()

		for filename := range ch {

			select {
			case <-ctx.Done():
				fmt.Fprintf(os.Stderr, "context error: %v\n", ctx.Err())
				if err := deleteFile(filename); err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				errChan <- ctx.Err()

			default:

				if err := deleteFile(filename); err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}

			}

		}

	}()

}

func deleteFile(filename string) error {

	err := os.Remove(filename)
	if err != nil {
		return fmt.Errorf("error deleting file '%v': %v", filename, err)

	}
	fmt.Fprintf(os.Stdout, "file '%v', deleted\n", filename)
	return nil

}

// writes errors coming from a channel to a file
func writeErrorsToFile(ch <-chan error, logFileName string) {

	var wg sync.WaitGroup
	wg.Add(1)

	f, err := os.Create(logFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating file: %v", err)
	}

	// create new buffer
	// buffer := bufio.NewWriter(f)
	//set output of logs to f
	log.SetOutput(f)

	go func() {
		defer wg.Done()
		for errFromChannel := range ch {
			log.Println(errFromChannel)
			// _, err := buffer.WriteString(errFromChannel.Error() + "\n")
			// if err != nil {
			// 	fmt.Fprintf(os.Stderr, "error writing line to the file: %v", err)
			// }
		}
	}()

	go func() {
		// remember to close the file
		defer f.Close()
		wg.Wait()
		// flush buffered data to the file
		// if err := buffer.Flush(); err != nil {
		// 	fmt.Fprintf(os.Stderr, "error flushing the buffer: %v\n", err)
		// }

	}()

}
