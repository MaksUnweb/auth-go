package repositories

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


//Структура сериализации данных из базы данных
type Admin struct {
	AdminID string
	Login string 
	Password string	
}


//Функция для проверки и получения администратора в бд:
func CheckInDb(login string, pool *pgxpool.Pool, ctx context.Context) (*Admin, error) {
	var adminStruct Admin 
	err := pool.QueryRow(ctx, "SELECT admin_id, login, password FROM admins WHERE login = $1", login).Scan(&adminStruct.AdminID, &adminStruct.Login, &adminStruct.Password)
	return &adminStruct, err
}


//Функция для проверки по admin_id полю
func CheckByAdminID(ctx context.Context, adminID string, pool *pgxpool.Pool) (string, error) {
	var login string
	err := pool.QueryRow(ctx, "SELECT login FROM admins WHERE admin_id = $1", adminID).Scan(&login)
	if err != nil{
		if err != pgx.ErrNoRows {
			log.Printf("Error request from Postgres: %v", err)
		}
	}
	
	return login, err
}
