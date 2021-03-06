package data

import (
	"encoding/json"
	"errors"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Person struct {
		ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
		Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	}

	People []*Person
)

var ErrNotFound = errors.New("not found")

func (p *Person) ToJSON(w io.Writer) error {

	return json.NewEncoder(w).Encode(&p)
}

func (p *Person) FromJSON(r io.Reader) error {

	return json.NewDecoder(r).Decode(&p)
}

func (p *People) ToJSON(w io.Writer) error {

	if len(*p) == 0 {
		return ErrNotFound
	}

	return json.NewEncoder(w).Encode(&p)
}

func (p *People) FromJSON(r io.Reader) error {

	return json.NewDecoder(r).Decode(&p)

}
