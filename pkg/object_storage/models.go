package object_storage

import "io"

type File struct {
	Name       string
	Extension  string
	FileReader io.Reader
	Size       int64
}
