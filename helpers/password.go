package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(pass []byte) ([]byte, error) {
	h, err := bcrypt.GenerateFromPassword(pass, 5)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func ComparePassword(h, raw []byte) error {
	return bcrypt.CompareHashAndPassword(h, raw)
}
