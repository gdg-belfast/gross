package implementation

import (
	"crypto/md5"
	"fmt"
	"github.com/gdg-belfast/gross/domain"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// HashString takes a string and returns the MD5 of that string
func HashString(s string) string {
	checksum := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", checksum)
}

// MonitorDirectory takes a local directory and a iterates through it
// placing newly found files into a channel
//
// It is not meant to be returned from
func MonitorDirectory(directory string, additions chan (*domain.MediaFile)) {
	if err := os.Chdir(directory); err != nil {
		log.Fatalln(err)
	}
	fileList := make(map[string]*domain.MediaFile)
	for {
		files, err := ioutil.ReadDir(directory)
		if err != nil {
			log.Fatalln("Error reading")
		}
		for _, file := range files {
			if file.IsDir() {
				// The returned file is a directory. For demo purposes
				// we won't recurse
				continue
			}
			sha := HashString(file.Name())
			if _, ok := fileList[sha]; ok {
				// File exists in the hash. Continue
				continue
			}
			fileList[sha] = &domain.MediaFile{
				file,
				filepath.Join(directory, file.Name()),
				sha,
			}
			log.Println("Adding:", file.Name())
			additions <- fileList[sha]
		}
		time.Sleep(time.Second * 60)
	}
}
