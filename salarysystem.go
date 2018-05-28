package main

import (
	"fmt"
	"log"
	"golang.org/x/crypto/ssh/terminal"
	"./accountant"
	"./administrator"
	"syscall"
	"database/sql"
	_ "github.com/lib/pq"
)

type user struct {
	username string
	password string
}

func main() {
	var username string
	for success := false; success != true; {
		fmt.Print("Введите логин: (или введите \"Ctrl + C\", чтобы выйти) ")
		_, err1 := fmt.Scan(&username)
		
		db, _ := sql.Open("postgres", "user=dmitry password=dmitry")


		if err1 != nil {
			log.Fatal("Неудалось прочитать логин: %v", err1)
		}

		fmt.Print("Введите пароль: ")
		password, err2 := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err2 != nil {
			log.Fatal("Неудалось прочитать пароль: %v", err2)
		}
		switch checkUser(username, string(password), db) {
		case "admin":
			administrator.HandleAdministrator(username, db)
			return
		case "accountant":
			accountant.HandleAccountant(username, db)
			return
		default:
			fmt.Println("Неправильное имя или пароль. \nПопробуйте еще.")
		}
	}
}

func checkUser(username, password string, db *sql.DB) string {
	rows, err := db.Query("SELECT type FROM users WHERE username=$1 and password=$2", username, password)

	if err != nil {
		log.Fatal(err)
	}
	var t string
	for i := 0; rows.Next(); i++ {

		if err2 := rows.Scan(&t); err2 != nil {
			log.Fatal(err2)
		}

		if i != 0 {
			return "Incorrect"
		}
	}
	return t
}