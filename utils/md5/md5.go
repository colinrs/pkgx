package md5

import (
	cryptoMD5 "crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/spaolacci/murmur3"
)

var _ MD5 = (*md5)(nil)

type MD5 interface {
	i()
	// Encrypt 加密
	Encrypt(encryptStr string) string
}

type md5 struct{}

func New() MD5 {
	return &md5{}
}

func (m *md5) i() {}

func (m *md5) Encrypt(encryptStr string) string {
	s := cryptoMD5.New()
	s.Write([]byte(encryptStr))
	return hex.EncodeToString(s.Sum(nil))
}

// Hash returns the hash value of data.
func Hash(data []byte) uint64 {
	return murmur3.Sum64(data)
}

// Md5 returns the md5 bytes of data.
func Md5(data []byte) []byte {
	digest := cryptoMD5.New()
	digest.Write(data)
	return digest.Sum(nil)
}

// Md5Hex returns the md5 hex string of data.
func Md5Hex(data []byte) string {
	return fmt.Sprintf("%x", Md5(data))
}
