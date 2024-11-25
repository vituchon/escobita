package util

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

type IntegerSequence struct {
	mu           sync.Mutex
	filename     string
	initialValue int
	increment    int
}

func NewIntegerSequence(filename string, initialValue int, increment int) *IntegerSequence {
	return &IntegerSequence{filename: filename, initialValue: initialValue, increment: increment}
}

func (seq *IntegerSequence) GetNext() (int, error) {
	seq.mu.Lock()
	defer seq.mu.Unlock()

	file, err := openOrCreate(seq.filename, seq.initialValue)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	current, err := parseContent(file)
	if err != nil {
		return 0, err
	}
	next := *current + seq.increment
	err = setContent(file, next)
	if err != nil {
		return 0, err
	}

	return next, nil
}

func openOrCreate(filename string, initialValue int) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return create(filename, initialValue)
		}
		return nil, err
	}
	return file, nil
}

func create(filename string, initialValue int) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		fmt.Println(err)
		if err != nil && file != nil {
			file.Close()
		}
	}()
	err = setContent(file, initialValue)
	if err != nil {
		return nil, err
	}
	_, err = file.Seek(0, os.SEEK_SET)
	return file, err
}

func setContent(file *os.File, value int) error {
	_, err := file.WriteString(strconv.Itoa(value))
	if err != nil {
		return err
	}
	return nil
}

func parseContent(file *os.File) (*int, error) {
	var content string
	_, err := fmt.Fscan(file, &content)
	if err != nil {
		return nil, err
	}
	number, err := strconv.Atoi(content)
	if err != nil {
		return nil, err
	}
	return &number, nil
}
