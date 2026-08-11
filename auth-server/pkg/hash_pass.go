package pkg

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Функция для хэширования пароля
// Принимает обычный пароль в формате строки, возвращает строку с паролем в уже захэшированном формате или ошибку,
// если что-то пошло не так
func HashPass(pass string) (string, error) {
	if pass == "" {
		return "", fmt.Errorf("Pass can't be empty!")
	}
	hashPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hashPass), err
}



// Функция для верификации пароля.
// Принимает 2 строки: пароль в первозданном виде, а также в захэшированном
// если err == nil, значит пароль верифицирован, если нет, значит пароль не совпадает
func VerifyPass(pass, hashPass string) error {
	if pass == "" || hashPass == "" {
		return fmt.Errorf("Pass or hashPass can't be empty!")
	}
	return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(pass))

}
