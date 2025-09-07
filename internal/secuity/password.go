// Package security
package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/rtmelsov/adv-keeper/internal/models"

	"golang.org/x/crypto/argon2"
)

var DefaultParams = models.Argon2idParams{
	Time:    3,
	Memory:  64 * 1024,
	Threads: 1,
	SaltLen: 16,
	KeyLen:  32,
}

func HashPasswordPHC(password string, p models.Argon2idParams) (string, error) {

	salt := make([]byte, p.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Threads, p.KeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	phc := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s%s", p.Memory, p.Time, p.Threads, b64Salt, b64Hash)
	return phc, nil

}
