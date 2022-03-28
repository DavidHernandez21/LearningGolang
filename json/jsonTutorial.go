package main

import (
	"encoding/json"
	"log"
)

// func main() {
// 	unstructuredJson := `{"os": {"Windows": "Windows OS","Mac": "OSX","Linux": {"Ubuntu": 18}},"compilers": "gcc"}`

// 	var result map[string]interface{}

// 	json.Unmarshal([]byte(unstructuredJson), &result)

// 	// fmt.Printf("%T\n", result["os"]) // map[Linux:Ubuntu Mac:OSX Windows:Windows OS]

// 	for _, v := range result {

// 		// fmt.Printf("key: %v\tvalue: %v\ttype: %T\n", k, v, v)
// 		switch v.(type) {
// 		case map[string]interface{}:
// 			for k, v1 := range v.(map[string]interface{}) {
// 				fmt.Printf("key: %v\tvalue: %v\ttype: %T\n", k, v1, v1)
// 			}
// 			// case map[string]map[string]int:
// 			// 	for k, v1 := range v.(map[string]map[string]int) {
// 			// 		fmt.Printf("key: %v\tvalue: %v\n", k, v1)
// 		// }
// 		case string:
// 			fmt.Println(v.(string))

// 		default:
// 			fmt.Printf("type: %T\n", v)

// 		}

// 	}
// }

func main() {

	// "bucket": "bucketName"
	unstructuredJson := `{"filename": "fileName"}`

	var d struct {
		FileNAme string `json:"filename,omitempty"`
		Bucket   string `json:"bucket,omitempty"`
	}

	if err := json.Unmarshal([]byte(unstructuredJson), &d); err != nil {
		log.Fatalf("error unmarshalling data: %v", err)
	}

	if d.Bucket == "" {
		d.Bucket = "Daje"
	}

	log.Printf("data recieved: %v\n", d)

	// if err := json.NewDecoder([]byte{unstructuredJson}).Decode(&d); err != nil {

	// }

}
