package main

import (
	"encoding/json"
	"fmt"
)

type UserRequest struct {
	User string `json:"login"`
}

type GroupRequest struct {
	Group string `json:"login"`
}

type Request interface {
	UserRequest | GroupRequest
}

type Ptr[T Request] interface {
	*T
}

func UnmarshalJSON[R Request, T Ptr[R]](data []byte, v T) error {
	return json.Unmarshal(data, v)
}

func main() {
	data := []byte(`{"login": "elliot"}`)
	var r UserRequest
	fmt.Println(UnmarshalJSON(data, &r))
	fmt.Println(r)
	// fmt.Println(UnmarshalJSON(data, r)) doesn't compile
}
