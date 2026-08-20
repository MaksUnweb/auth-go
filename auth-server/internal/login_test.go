package internal

import "testing"

// Тестируем функцию валидации пароля:
func TestValidatePass(t *testing.T) {
	varyLongPass := "2131312312312312312314zdkbfadjfhffudhjfjhdfjlsdlkjfsdkjlfjkldsfjifjsdjfdjsfkdsfjksdfdsfsdfdsfhdshfhjdsfhjdshfjsdhfjdshfjhdsfhdjsfhjdsfjdsfjdshfjdsfhjsdf"
	noValideSymbols := "!@#!$DSF$^&*"
	successPass := "test"

	// Тестируем условие слишком длинного пароля:
	if validatePass(varyLongPass) {
		t.Fatal("Error! This password vary long!")
	}

	if validatePass(noValideSymbols) {
		t.Fatal("Error! No valide symbols!")
	}

	if !validatePass(successPass) {
		t.Fatal("Error! This password valide!")
	}
}

// Тестируем функцию валидации логина:
func TestValidateLogin(t *testing.T) {
	varyLongLogin := "2131312312312312312314zdkbfadjfhffudhjfjhdfjlsdlkjfsdkjlfjkldsfjifjsdjfdjsfkdsfjksdfdsfsdfdsfhdshfhjdsfhjdshfjsdhfjdshfjhdsfhdjsfhjdsfjdsfjdshfjdsfhjsdf"
	noValideSymbols := "!@#!$DSF$^&*"
	successLogin := "Test"

	// Тестируем условие слишком длинного пароля:
	if validateLogin(varyLongLogin) {
		t.Fatal("Error! This login vary long!")
	}

	if validateLogin(noValideSymbols) {
		t.Fatal("Error! No valide symbols!")
	}

	if !validateLogin(successLogin) {
		t.Fatal("Error! This login valide!")
	}
}
