package accountant

import (
	"strings"
	"os"
	"bufio"
	"strconv"
	"log"
	"fmt"
	"database/sql"
	//it's for fun
	_ "github.com/lib/pq"
)

//const for using
const (
	HOURLY = int(0)
	RATE = int(1)
)

type worker struct {
	id int
	name string
	surname string
	payment int
	adress string
	unit uint
	position string
}

func (w worker) String() string {
	var p string
	if w.payment == HOURLY {
		p = "часовая"
	} else {
		p = "оклад"
	}
	return fmt.Sprintf("%v %v - %v. тип оплаты: %v, адрес: %v, номер подразделения: %v",
		w.surname, w.name, w.position, p, w.adress, w.unit)
}

//HandleAccountant handles accountant account
func HandleAccountant(username string, db *sql.DB) {
	for {
		fmt.Println(`Выберите действие:
1. Создать карточку
2. Редактировать карточку
3. Сделать чек
4. Сделать отчет
5. Выйти`)
		var action string
		fmt.Scan(&action)
		switch action{
		case "1":
			createCard(username, db)
		case "2":
			editCard(username, db)
		case "3":
			makeCheque(username, db)
		case "4":
			makeReport(username, db)
		case "5":
			fmt.Println("Good bye!")
			return
		default:
			fmt.Println("Выберете корректное действие!")
		}
	}
}

func createCard(username string, db *sql.DB) {
	var name, surname, adress, position string
	var payment = -1
	var unit = uint(0)
	fmt.Print("Создание новой карточки.\nЧтобы выйти напишите EXIT\nВведите имя работника: ")
	for { 
		for {
			fmt.Scan(&name)
			if name == "EXIT" {
				return
			}
			if len(name) > 20 {
				fmt.Println("\nИмя работника не должно превышать 20 символов. Попробуйте еще раз.")
			} else {
				break
			}
		}
	
		for {
			fmt.Print("Введите фамилию: ")
			fmt.Scan(&surname)
			if surname == "EXIT" {
				return
			}
			if len(surname) > 20 {
				fmt.Println("\nФамилия работника не должна превышать 20 символов. Попробуйте еще раз.")
			} else {
				break
			}
		}

		for payment == -1 {
			fmt.Print("Выберете тип оплаты рабочего:\n1 - оклад\n2 - почасовой\n: ")
			var temp string
			fmt.Scan(&temp)
			switch temp {
			case "EXIT":
				return
			case "1":
				payment = int(RATE)
				break
			case "2":
				payment = int(HOURLY)
				break
			default:
				fmt.Println("\nВведите корректный тип оплаты!")
			}
		}	

		for {
			fmt.Print("Введите адрес рабочего: ")
			fmt.Scan(&adress)
			if adress == "EXIT" {
				return
			}
			if len(adress) > 20 {
				fmt.Println("\nАдрес не должен превышать 20 символов. Попробуйте еще раз.")
			} else {
				break
			}
		}

		for {
			fmt.Print("Введите должность рабочего: ")
			fmt.Scan(&position)
			if position == "EXIT" {
				return
			}
			if len(position) > 20 {
				fmt.Println("\nНазвание должности не должно превышать 20 символов. Попробуйте еще раз.")
			} else {
				break
			}
		}

		unit = getUnit(username, db)
		fmt.Printf(`Проверьте введенные данные:
Имя работника: %v
Фамилия: %v
Тип оплаты: %v
Адрес: %v
Номер подразделения: %v
Должность: %v
Данные корректны? (Y/N)`, name, surname, payment, adress, unit, position)
		var answer string
		fmt.Scan(&answer)
		if answer[0] == 'Y' || answer[0] == 'y' {
			break
		}
	}
	cardToDB(name, surname, adress, payment, unit, position, db)
}

func cardToDB(name, surname, adress string, payment int, unit uint, position string, db *sql.DB) {
	rows, _ := db.Query("SELECT MAX(id) FROM employees")
	var id int
	for rows.Next() {
		rows.Scan(&id)
	}
	_, err := db.Exec("INSERT INTO employees (id, name, surname, payment, adress, unit, position) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		id + 1, name, surname, payment, adress, unit, position)
	if err != nil {
		log.Fatal(err)
	}
}


func editCard(username string, db *sql.DB) {
	unit := getUnit(username, db)
	
	workers := getWorkers(unit, db)

	for {
		fmt.Println("Чтобы выйти, напишите EXIT.")
		fmt.Println("Выберете номер карточки работника для изменения: ")
		for i, v := range (workers) {
			fmt.Println(i + 1, v)
		}
		var query string
		fmt.Scan(&query)
		fmt.Println()
		if query == "EXIT" {
			return
		}
		

		number, err := strconv.Atoi(query)
		fmt.Println(number)
		if err != nil || number <= 0 || number > len(workers){
			fmt.Println("Введите корректный номер карточки!")
		} else {
			fmt.Println("Чтобы сохранить изменения, введите SAVE, чтобы выйти, не сохраняя, введите EXIT.")
			fmt.Println("Выберете поле, которое хотите изменить:")
			for {
				fmt.Printf("(НЕ ИЗМЕНЯЕМО) Имя сотрудника: %v\n", workers[number - 1].name)
				fmt.Printf("(НЕ ИЗМЕНЯЕМО) Фамилия сотрудника: %v\n", workers[number - 1].surname)
				fmt.Printf("(НЕ ИЗМЕНЯЕМО) Должность: %v\n", workers[number - 1].position)
			
				var p string
				if workers[number - 1].payment == HOURLY {
					p = "почасовой"
				} else {
					p = "оклад"
				}

				fmt.Printf("1. Тип оплаты сотрудника: %v\n", p)
				fmt.Printf("2. Адрес сотрудника: %v\n", workers[number - 1].adress)

				fmt.Scan(&query)

				if query == "EXIT" {
					return
				}
				if query == "SAVE" {
					editCardDB(workers[number - 1], db)
					break
				}
				
				num, err2 := strconv.Atoi(query)

				if err2 != nil || num < 1 || num > 2 {
					fmt.Println("Введите корректный запрос!")
				} else {
					if num == 1 {
						var newP = -1
						for newP == -1 {
							fmt.Print("Выберете тип оплаты рабочего:\n1 - оклад\n2 - почасовой\n: ")
							var temp string
							fmt.Scan(&temp)
							switch temp {
							case "EXIT":
								return
							case "1":
								newP = int(RATE)
								workers[number - 1].payment = newP
								break
							case "2":
								newP = int(HOURLY)
								workers[number - 1].payment = newP
								break
							default:
								fmt.Println("\nВведите корректный тип оплаты!")
							}
						}	
					}

					if num == 2 {
						fmt.Print("Введите новый адресс: ")
						scanner := bufio.NewReader(os.Stdin)
						newAdr, _ := scanner.ReadString('\n')
						workers[number - 1].adress = string(newAdr[:len(newAdr) - 1])
					}
				}

			}
		}
	}
}

func makeCheque(username string, db *sql.DB) {
	fmt.Println("Выберете сотрудника, для которого сформировать чек: ")
	workers := getWorkers(getUnit(username, db), db)
	for i, v := range (workers) {
		fmt.Println(i + 1, v)
	}
	
	var number int

	fmt.Scan(&number)

	hours := getHours(workers[number - 1], db)
	
	var salary float32
	for _, v := range (hours) {
		if v > 8 {
			salary += v * 10 + (8 - v) * 20
		}
	}
	fmt.Printf("Сотруднику выписан чек на %v рублей.\n", salary)
}

func makeReport(username string, db *sql.DB) {
	fmt.Println("Выберете тип отчета:")
	fmt.Println("1. Количесвто отработанных часов сотрудниками подразделения.")
	fmt.Println("2. Отчет по работникам на ставке.")

	var number int
	var report string
	fmt.Scan(&number)
	workers := getWorkers(getUnit(username, db), db)

	switch number{
	case 1:
		fmt.Println("Выберете сотрудников для отчета")
		for i, v := range(workers) {
			fmt.Println(i + 1, v)
		}
		s := bufio.NewReader(os.Stdin)
		n, _ := s.ReadString('\n')
		n = string(n[:len(n) - 1])
		nums := strings.Fields(n)
		var vv []int
		for _, v := range(nums) {
			f, _ := strconv.Atoi(v)
			vv = append(vv, f)
		}
		for _, v := range(vv) {
			hours := getHours(workers[v - 1], db)
			var h float32
			for i := range(hours) {
				h += hours[i]
			}
			report += fmt.Sprintf("%v. Отработано %v часов.\n", workers[v - 1], h)
		}
	case 2:
		for _, v := range(workers) {
			if v.payment == RATE {
				report += fmt.Sprintf("%v. Отработано %v часов.\n", v, getHours(v, db))
			}
		}
		if report == "" {
			report = "На данный момент отсутствуют работники на ставке."
		}
	}
	sendReport(username, report, db)
}

func getUnit(username string, db *sql.DB) uint {
	var unit uint
	row := db.QueryRow("SELECT unit FROM users WHERE username=$1", username)
	row.Scan(&unit)
	return unit
}

func getWorkers(unit uint, db *sql.DB) []worker {
	var w []worker

	rows, _ := db.Query("SELECT * FROM employees WHERE unit=$1", unit)

	for rows.Next() {
		var tw worker
		rows.Scan(&tw.id, &tw.name, &tw.surname, &tw.payment, &tw.adress, &tw.unit, &tw.position)
		w = append(w, tw)
	}

	return w
}

func editCardDB(w worker, db *sql.DB) {
	_, err := db.Exec("UPDATE employees SET adress=$1 WHERE id=$2", w.adress, w.id)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("UPDATE employees SET payment=$1 WHERE id=$2", w.payment, w.id)
}

func getHours(w worker, db *sql.DB) []float32 {
	rows, err := db.Query("SELECT hours FROM workedtime WHERE id=$1", w.id)
	if err != nil {
		log.Fatal(err)
	}
	var hours []float32
	for rows.Next() {
		var hour float32
		rows.Scan(&hour)
		hours = append(hours, hour)
	}
	return hours
}

func sendReport(username, report string, db *sql.DB) {
	_, err := db.Exec("INSERT INTO reports (type, username, report) VALUES ('accountant', $1, $2)", username, report)
	if err != nil {
		log.Fatal(err)
	}
}