package database

import (
	"context"
)

type Executor interface {
	LookupByID(ctx context.Context, tableName string, id uint64) (interface{}, error)
	LookupAll(ctx context.Context, tableName string) ([]interface{}, error)
	Insert(ctx context.Context, tableName string, value interface{}) error
	CreateTable(ctx context.Context, tableName string) error
}

type ExecutorReader interface {
	LookupByID(ctx context.Context, tableName string, id uint64) (interface{}, error)
	LookupAll(ctx context.Context, tableName string) ([]interface{}, error)
}

type ExecutorWriter interface {
	Insert(ctx context.Context, tableName string, value interface{}) error
	CreateTable(ctx context.Context, tableName string) error
}

// Sample executor implementation
type constError string
type contextKeyType int

const (
	ErrNotFoundRecord                     = constError("could not find record")
	ErrNotFoundTable                      = constError("could not find table")
	ErrTableFull                          = constError("table is full")
	ErrMaxTablesReach                     = constError("Max number of tables reached")
	ErrTablesAlreadyexists                = constError("Table already exists")
	ExecutorContextKey     contextKeyType = iota
)

type (
	inMemory struct {
		numTables  uint64
		numObjects uint64
		data       map[string]*object
	}

	object struct {
		data map[uint64]interface{}
	}
)

func (e constError) Error() string {
	return string(e)
}

func (i *inMemory) LookupByID(ctx context.Context, tableName string, id uint64) (interface{}, error) {

	table, exists := i.data[tableName]
	if !exists {
		return nil, ErrNotFoundTable
	}

	val, exists := table.data[id]

	if !exists {
		return nil, ErrNotFoundRecord
	}

	return val, nil
}

func (i *inMemory) LookupAll(ctx context.Context, tableName string) ([]interface{}, error) {

	table, exists := i.data[tableName]
	if !exists {
		return nil, ErrNotFoundTable
	}

	values := make([]interface{}, 0, len(table.data))
	for _, value := range table.data {
		values = append(values, value)
	}

	return values, nil

}

func (i *inMemory) Insert(ctx context.Context, tableName string, value interface{}) error {

	table, exists := i.data[tableName]
	if !exists {
		return ErrNotFoundTable
	}

	tableLen := len(table.data)

	if tableLen >= int(i.numObjects) {
		return ErrTableFull
	}

	// This could overflow - don't use it in production code
	table.data[uint64(tableLen)] = value

	return nil
}

func (i *inMemory) CreateTable(ctx context.Context, tableName string) error {

	_, exists := i.data[tableName]

	if exists {
		return ErrTablesAlreadyexists
		// return fmt.Errorf("table: '%v' already exists", tableName)
	}

	if len(i.data) >= int(i.numTables) {
		return ErrMaxTablesReach
	}

	o := &object{
		data: make(map[uint64]interface{}, i.numObjects),
	}

	i.data[tableName] = o

	return nil
}

func NewInMemoryDB(numTables, numObjects uint64) *inMemory {

	memory := &inMemory{
		numTables:  numTables,
		numObjects: numObjects,
	}

	memory.data = make(map[string]*object, memory.numTables)

	return memory
}
