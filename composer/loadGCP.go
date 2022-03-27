// package main

// import (
// 	"io/ioutil"
// 	"log"
// 	"os/exec"
// 	"regexp"
// 	"sync"
// )

// func main() {

// 	files, err := ioutil.ReadDir("./")
// 	if err != nil {
// 		log.Fatalf("error reading the directory ./: %v", err)
// 	}

// 	pattern := regexp.MustCompile(`DATA.+_[0-9]{1,2}.csv$`)
// 	var wg sync.WaitGroup
// 	for _, f := range files {
// 		find := pattern.FindAllString(f.Name(), 100000)
// 		if len(find) > 0 {
// 			wg.Add(1)
// 			go sendFileToBucket(find[0], &wg)
// 		}

// 		// log.Println(pattern.FindAllString(f.Name(), 100000))
// 		// log.Println(f.Name())
// 	}

// 	wg.Wait()
// }

// func sendFileToBucket(file string, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	log.Printf("copying file %v, to bucket gs://datastage-gia-csv-repo-prod", file)
// 	cmd, err := exec.Command("gsutil", "-m", "cp", "-n", file, "gs://datastage-gia-csv-repo-prod").CombinedOutput()

// 	if err != nil {
// 		log.Printf("error copying file %v", file)
// 	}

// 	log.Println(string(cmd[:]))
// }
