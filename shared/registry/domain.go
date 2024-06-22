package registry

import "golang.org/x/crypto/sha3"

func HashDomain(input string) [32]byte {
	hash := sha3.NewLegacyKeccak256()
	_, _ = hash.Write([]byte(input))

	// Get the resulting encoded byte slice
	sha3 := hash.Sum(nil)

	var sha32Bytes [32]byte
	copy(sha32Bytes[:], sha3)

	// Convert the encoded byte slice to a string
	return sha32Bytes
}
