package util

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type IntegerSequence interface {
	GetNext() (int, error)
}

type FsIntegerSequence struct {
	mu           sync.Mutex
	filename     string
	initialValue int
	increment    int
}

func NewFsIntegerSequence(filename string, initialValue int, increment int) *FsIntegerSequence {
	return &FsIntegerSequence{filename: filename, initialValue: initialValue, increment: increment}
}

func (seq *FsIntegerSequence) GetNext() (int, error) {
	seq.mu.Lock()
	defer seq.mu.Unlock()

	file, err := openOrCreate(seq.filename, seq.initialValue)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	current, err := getContent(file)
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
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return os.OpenFile(filename, os.O_RDWR, 0644)
		} else {
			return nil, err
		}
	} else {
		err = setContent(file, initialValue)
		if err != nil {
			return nil, err
		}
		return file, nil
	}

}

func setContent(file *os.File, value int) error {
	_, err := file.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}
	_, err = file.WriteString(strconv.Itoa(value))
	if err != nil {
		return err
	}
	return nil
}

func getContent(file *os.File) (*int, error) {
	_, err := file.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	var content string
	_, err = fmt.Fscan(file, &content)
	if err != nil {
		return nil, err
	}
	number, err := strconv.Atoi(content)
	if err != nil {
		return nil, err
	}
	return &number, nil
}
