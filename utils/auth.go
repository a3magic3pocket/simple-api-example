package utils

import "golang.org/x/crypto/bcrypt"

// HashAndSalt : 입력 받은 password를 암호화
func HashAndSalt(password string) (hash string, err error) {
	cost := 10
	byteHash, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	return string(byteHash), err
}

// ComparePasswords : 암호호된 패스워드와 평문 패스워드 비교
func ComparePasswords(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
