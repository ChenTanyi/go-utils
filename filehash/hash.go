package filehash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"math"
	"strings"
)

var hashFuncs = map[string]func() hash.Hash{
	"md5":    md5.New,
	"sha1":   sha1.New,
	"sha224": sha256.New224,
	"sha256": sha256.New,
	"sha384": sha512.New384,
	"sha512": sha512.New,
}

type fileReader struct {
	file   io.ReaderAt
	offset int64
	end    int64
}

func (f *fileReader) Read(b []byte) (n int, err error) {
	n, err = f.file.ReadAt(b, f.offset)
	if f.offset+int64(n) >= f.end {
		n = int(f.end - f.offset)
		err = io.EOF
	}
	f.offset += int64(n)
	return
}

// HashAllFile sha1 hash
func HashAllFile(file io.ReaderAt) ([]byte, error) {
	return Hash(sha1.New, file, 0, math.MaxInt64)
}

// HashFile sha1 hash
func HashFile(file io.ReaderAt, begin, end int64) ([]byte, error) {
	return Hash(sha1.New, file, begin, end)
}

// HashAllFileWithFuncName .
func HashAllFileWithFuncName(funcName string, file io.ReaderAt) ([]byte, error) {
	return HashFileWithFuncName(funcName, file, 0, math.MaxInt64)
}

// HashFileWithFuncName .
func HashFileWithFuncName(funcName string, file io.ReaderAt, begin, end int64) ([]byte, error) {
	hashFunc, ok := hashFuncs[funcName]
	if !ok {
		keys := make([]string, 0, len(hashFuncs))
		for key := range hashFuncs {
			keys = append(keys, key)
		}
		return nil, fmt.Errorf("Invalid hash function name: %s, expected (%s)", funcName, strings.Join(keys, ", "))
	}
	return Hash(hashFunc, file, begin, end)
}

// Hash .
func Hash(hashFunc func() hash.Hash, file io.ReaderAt, begin, end int64) ([]byte, error) {
	h := hashFunc()
	_, err := io.Copy(h, &fileReader{
		file: file, offset: begin, end: end,
	})
	return h.Sum(nil), err
}
