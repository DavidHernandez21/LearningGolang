package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// get speculative number of files/directories in a directory
// assuming two files per directory
func getNumFiles(dir string) (int, error) {
	numFiles := 0
	d, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	for k := range d {
		numFiles++
		if d[k].IsDir() {
			numFiles++
		}
	}

	return numFiles, nil
}

// walkFiles starts a goroutine to walk the directory tree at root and send the
// path of each regular file on the string channel.  It sends the result of the
// walk on the error channel.  If done is closed, walkFiles abandons its work.
func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)
	go func() { // HL
		// Close the paths channel after Walk returns.
		defer close(paths) // HL
		// No select needed for this send, since errc is buffered.
		errc <- filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error { // HL
			if err != nil {
				return err
			}
			if !d.Type().IsRegular() {
				return nil
			}
			select {
			case paths <- path: // HL
			case <-done: // HL
				return errors.New("walk canceled")
			}
			return nil
		})
	}()
	return paths, errc
}

// A result is the product of reading and summing a file using MD5.
type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

// digester reads path names from paths and sends digests of the corresponding
// files on c until either paths or done is closed.
func digester(done <-chan struct{}, paths <-chan string, c chan<- result) {
	for path := range paths { // HLpaths
		data, err := os.ReadFile(path)
		select {
		case c <- result{path, md5.Sum(data), err}:
		case <-done:
			return
		}
	}
}

// MD5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, MD5All returns an error.  In that case,
// MD5All does not wait for inflight read operations to complete.
func MD5All(root string, numDigesters int) (map[string][md5.Size]byte, error) {
	// MD5All closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc.
	done := make(chan struct{})
	defer close(done)

	paths, errc := walkFiles(done, root)

	// Start a fixed number of goroutines to read and digest files.
	digesterChannel := make(chan result) // HLc
	var wg sync.WaitGroup

	wg.Add(numDigesters)
	for i := 0; i < numDigesters; i++ {
		go func() {
			digester(done, paths, digesterChannel) // HLc
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(digesterChannel) // HLc
	}()
	// End of pipeline. OMIT

	//speculative length of map for the number of files
	// assuming two files per directory
	mapLen, err := getNumFiles(root)
	if err != nil {
		return nil, err
	}
	m := make(map[string][md5.Size]byte, mapLen*2)
	for r := range digesterChannel {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	// Check whether the Walk failed.
	if err := <-errc; err != nil { // HLerrc
		return nil, err
	}
	return m, nil
}

func main() {
	// Calculate the MD5 sum of all files under the specified directory,
	// then print the results sorted by path name.
	const numDigesters = 20
	m, err := MD5All(os.Args[1], numDigesters)
	if err != nil {
		fmt.Println(err)
		return
	}
	// var paths []string
	paths := make([]string, len(m))
	k := 0
	for path := range m {
		paths[k] = path
		k++
		// paths = append(paths, path)
	}

	sort.Strings(paths)

	for k := range paths {
		fmt.Printf("%x  %s\n", m[paths[k]], paths[k])
	}
}
