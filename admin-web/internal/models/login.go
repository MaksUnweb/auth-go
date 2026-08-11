package models 

//Основная структура для десериализации POST-запроса в структуру:
type Login struct {
	Login string `json:"login"`
	Password string `json:"password"`
}
