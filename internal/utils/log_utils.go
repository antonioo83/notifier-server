package utils

import (
	"io"
	"log"
)

// LogErr write error to the log.
func LogErr(n int, err error) error {
	if err != nil {
		log.Printf("Write failed %d byte: %v", n, err)
	}

	return err
}

// ResourceClose close resource.
func ResourceClose(body io.ReadCloser) error {
	if err := body.Close(); err != nil {
		log.Printf("Can't close resource: %v", err)
		return err
	}

	return nil
}
