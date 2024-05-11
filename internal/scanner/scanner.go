package scanner

import (
	"bufio"
	"io"
)

type IScanner interface {
	ReadLine() (string, error)
}

type Scanner struct {
	s *bufio.Scanner
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		s: bufio.NewScanner(r),
	}
}

// ReadLine reads line value from reader
func (scanner *Scanner) ReadLine() (string, error) {
	if !scanner.s.Scan() {
		err := scanner.s.Err()
		if err != nil {
			return "", err
		}

		return "", io.EOF
	}
	line := scanner.s.Text()

	return line, nil
}
