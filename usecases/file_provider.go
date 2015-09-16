package usecases

import (
	"os"
)

type ErrNotFound struct {
	error
}

type FileProvider interface {
	Get(string, string) (*os.File, error)
}
