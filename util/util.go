package util

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
)

type Sha1Stream struct {
	sha1 hash.Hash
}

func (sha1Stream *Sha1Stream) Update(data []byte) {
	if sha1Stream.sha1 == nil {
		sha1Stream.sha1 = sha1.New()
	}
	sha1Stream.sha1.Write(data)
}

func (sha1Stream *Sha1Stream) Sum() string {
	sha1 := sha1Stream.sha1
	return hex.EncodeToString(sha1.Sum([]byte("")))
}
func Sha1(data []byte) string {
	sha1 := crypto.SHA1.New()
	sha1.Write(data)
	return hex.EncodeToString(sha1.Sum([]byte("")))
}

func FileSha1(file *os.File) string {
	sha1 := sha1.New()
	io.Copy(sha1, file)
	return hex.EncodeToString(sha1.Sum(nil))
}

func MD5(data []byte) string {
	md5 := md5.New()
	md5.Write(data)
	return hex.EncodeToString(md5.Sum([]byte("")))
}

func FileMD5(file *os.File) string {
	md5 := md5.New()
	io.Copy(md5, file)
	return hex.EncodeToString(md5.Sum(nil))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileSize(fileName string) uint64 {
	var res uint64
	filepath.Walk(fileName, func(path string, fileInfo os.FileInfo, err error) error {
		res = uint64(fileInfo.Size())
		return nil
	})
	return res
}
