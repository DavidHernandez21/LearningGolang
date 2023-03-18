package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var getChannel = make(chan string)
var setChannel = make(chan string)
var deleteChannel = make(chan string)
var incrChannel = make(chan string)

func TestHandleGet(t *testing.T) {
	t.Parallel()
	const key1 = "key1"
	tests := map[string]struct {
		commandSlice []string
		response     string
	}{
		"key present":               {commandSlice: []string{"GET", key1}, response: "value1"},
		"key not present":           {commandSlice: []string{"GET", "key2"}, response: ""},
		"wrong number of arguments": {commandSlice: []string{"GET"}, response: "ERR wrong number of arguments for 'get' command"},
	}

	commandChannel := make(chan commandMessage)
	// getChannel := make(chan string)

	go handleDB(commandChannel)

	// set key1
	const value1 = "value1"
	setCommand := commandMessage{
		commandName:     Set,
		key:             "key1",
		value:           value1,
		responseChannel: setChannel}

	commandChannel <- setCommand
	ok := <-setCommand.responseChannel
	if ok != "OK" {
		t.Fatalf("Expected OK, got %s", ok)
	}

	for name, tt := range tests {
		commandSlice := tt.commandSlice
		response := tt.response
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, want := handleGet(commandChannel, getChannel, commandSlice), response; got != want {
				t.Errorf("input %q; got %q; want %q", commandSlice, got, want)

			}
		})
	}

	t.Cleanup(func() {
		close(commandChannel)
	})
	// close(commandChannel)
}

func TestHandleSet(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		commandSlice []string
		response     string
	}{
		"exact number of arguments":                  {commandSlice: []string{"SET", "key1", "value1"}, response: "OK"},
		"greater number of arguments than necessary": {commandSlice: []string{"SET", "key2", "value2", "value3"}, response: "OK"},
		"wrong number of arguments":                  {commandSlice: []string{"SET"}, response: "ERR wrong number of arguments for 'set' command"},
	}

	commandChannel := make(chan commandMessage)
	// setChannel := make(chan string)

	go handleDB(commandChannel)

	for name, tt := range tests {
		commandSlice := tt.commandSlice
		response := tt.response
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, want := handleSet(commandChannel, setChannel, commandSlice), response; got != want {
				t.Errorf("input %q; got %q; want %q", commandSlice, got, want)

			}

			if len(commandSlice) < 3 {
				return
			}

			getCommand := commandMessage{
				commandName:     Get,
				key:             commandSlice[1],
				responseChannel: getChannel}

			commandChannel <- getCommand
			value := <-getCommand.responseChannel
			if value != commandSlice[2] {
				t.Errorf("input %q; got %q; want %q", commandSlice, value, commandSlice[2])
			}

		})
	}

	t.Cleanup(func() {
		close(commandChannel)
	})
	// close(commandChannel)
}

func TestHandleIncrement(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		commandSlice []string
		response     string
	}{
		"key first time": {commandSlice: []string{"INCR", "key1"}, response: "1"},
		// "key second time":           {commandSlice: []string{"INCR", "key1"}, response: "2"},
		"wrong number of arguments": {commandSlice: []string{"INCR"}, response: "ERR wrong number of arguments for 'incr' command"},
		"value not a number":        {commandSlice: []string{"INCR", "key2"}, response: "ERR value is not an integer or out of range"},
	}

	commandChannel := make(chan commandMessage)
	// incrChannel := make(chan string)

	go handleDB(commandChannel)

	// set key2 to a non-integer value
	setCommand := commandMessage{
		commandName:     Set,
		key:             "key2",
		value:           "value2",
		responseChannel: setChannel}

	commandChannel <- setCommand
	ok := <-setCommand.responseChannel
	if ok != "OK" {
		t.Fatalf("Expected OK, got %s", ok)
	}

	for name, tt := range tests {
		commandSlice := tt.commandSlice
		response := tt.response
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, want := handleIncrement(commandChannel, incrChannel, commandSlice), response; got != want {
				t.Errorf("input %q; got %q; want %q", commandSlice, got, want)

			}
		})
	}
}

func TestHandleIncrementAfterFirstTime(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		commandSlice []string
		response     string
	}{
		"key second time": {commandSlice: []string{"INCR", "key1"}, response: "2"},
		"key third time":  {commandSlice: []string{"INCR", "key1"}, response: "3"},
	}

	commandChannel := make(chan commandMessage)
	// incrChannel := make(chan string)

	go handleDB(commandChannel)

	// increment key1 once
	incrCommand := commandMessage{
		commandName:     Incr,
		key:             "key1",
		value:           "1",
		responseChannel: incrChannel}

	commandChannel <- incrCommand
	one := <-incrCommand.responseChannel
	if one != "1" {
		t.Fatalf("Expected OK, got %s", one)
	}

	for name, tt := range tests {
		commandSlice := tt.commandSlice
		response := tt.response
		t.Run(name, func(t *testing.T) {
			// t.Parallel()
			if got, want := handleIncrement(commandChannel, incrChannel, commandSlice), response; got != want {
				t.Errorf("input %q; got %q; want %q", commandSlice, got, want)

			}
		})
	}
}

func TestHandleDelete(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		commandSlice []string
		response     string
	}{
		"key present":               {commandSlice: []string{"DEL", "key1"}, response: "key1"},
		"key not present":           {commandSlice: []string{"DEL", "key2"}, response: "0"},
		"wrong number of arguments": {commandSlice: []string{"DEL"}, response: "ERR wrong number of arguments for 'del' command"},
	}

	commandChannel := make(chan commandMessage)
	// deleteChannel := make(chan string)

	go handleDB(commandChannel)

	// set key1
	setCommand := commandMessage{
		commandName:     Set,
		key:             "key1",
		value:           "value1",
		responseChannel: setChannel}

	commandChannel <- setCommand
	ok := <-setCommand.responseChannel
	if ok != "OK" {
		t.Fatalf("Expected OK, got %s", ok)
	}

	for name, tt := range tests {
		commandSlice := tt.commandSlice
		response := tt.response
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, want := handleDelete(commandChannel, deleteChannel, commandSlice), response; got != want {
				t.Errorf("input %q; got %q; want %q", commandSlice, got, want)

			}

			if len(commandSlice) < 2 {
				return
			}

			if response != "1" {
				return
			}

			getCommand := commandMessage{
				commandName:     Get,
				key:             commandSlice[1],
				responseChannel: getChannel}

			commandChannel <- getCommand
			value := <-getCommand.responseChannel
			if value != "" {
				t.Errorf("key %q; got %q; want %q", commandSlice[1], value, "")
			}

		})
	}
}

func TestPreProcessInput(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		input        string
		commandSlice []string
	}{
		"empty input":        {input: "", commandSlice: []string{}},
		"SET without spaces": {input: "SET key1 value1", commandSlice: []string{"SET", "key1", "value1"}},
		"SET with spaces":    {input: "  SET   key1   value1  ", commandSlice: []string{"SET", "key1", "value1"}},
	}

	for name, tt := range tests {
		commandSlice := tt.commandSlice
		input := tt.input
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, want := preProcessInput(input), commandSlice
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("input %q; got %q; want %q; diff %q", input, got, want, diff)
			}
		})
	}
}

func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
