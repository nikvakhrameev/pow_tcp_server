package pow

import "crypto/sha256"

type Sha256Hasher struct{}

func NewSha256Hasher() Sha256Hasher {
	return Sha256Hasher{}
}

func (s Sha256Hasher) HashData(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}
