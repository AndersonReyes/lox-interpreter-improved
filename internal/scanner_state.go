package scanner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

type ScannerState struct {
	reader *bufio.Reader
	head   string
	next   string
}

func NewScannerState(reader *bufio.Reader) (*ScannerState, error) {

	state := &ScannerState{
		reader: reader,
	}
	// Load head and next charsjk
	_, err := state.GetNext()
	if err != nil {
		return nil, err
	}

	_, err = state.GetNext()
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (s *ScannerState) Head() string {
	return s.head
}

func (s *ScannerState) Peek() string {
	return s.next
}

func (s *ScannerState) GetNext() (string, error) {

	r, _, err := s.reader.ReadRune()
	if err != nil && err != io.EOF {
		return "", errors.Join(err, fmt.Errorf("failed to read next character from reader"))
	}

	s.head = s.next
	s.next = string(r)

	return s.head, err // err is only io.EOF here
}

func (s *ScannerState) Skip(n int) {

	for range n {
		_, _ = s.GetNext()
	}
}
