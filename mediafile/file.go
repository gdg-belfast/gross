package mediafile

import (
	"os"
)

// MediaFile is the base filetype for files being transformed into RssItems
type File struct {
	os.FileInfo
	// Filepath is the full path on the local filestem to the file
	Filepath string
	// Hash is the MD5 of the full path name
	Hash string
}
