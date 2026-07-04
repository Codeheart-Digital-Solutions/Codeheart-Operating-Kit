package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// ReaderSHA256 streams data from reader and returns the lowercase SHA-256 digest.
func ReaderSHA256(reader io.Reader) (string, error) {
	digest := sha256.New()
	if _, err := io.Copy(digest, reader); err != nil {
		return "", err
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}

// FileSHA256 streams a file and returns the lowercase SHA-256 digest.
func FileSHA256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return ReaderSHA256(file)
}
