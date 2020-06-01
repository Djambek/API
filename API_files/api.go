package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var tok_secret string

func registration(rw http.ResponseWriter, re *http.Request) {
	//fmt.Println("Registration")
	//в get-запросе ищем параметр login и присваем его
	login := re.URL.Query().Get("login")

	// открываем файл где есть логины
	file2, _ := os.Open("log.txt")
	b2, _ := ioutil.ReadAll(file2)

	// делаем список из логинов
	list_login := strings.Split(string(b2), " ")
	counter := 0

	//цикл для проверки есть ли такой логин уже
	for i := 0; i < len(list_login)-1; i++ {
		if list_login[i] == login {
			rw.Write([]byte("Your login used another human. Please choose other."))
			counter = 1
		}
	}
	//fmt.Println(len(list_login))
	//Смотрим есть ли у нас такой логин, если нет то пишем ваш логин
	if counter == 0 {
		check_token()
		//fmt.Println("Put login into file")
		message := []byte(string(b2) + string(login) + " ")
		ioutil.WriteFile("log.txt", message, 0644)
		//file2.WriteString(login)
		rw.Write([]byte("Your login, " + login + "\n"))
		//string_token := fmt.Sprintf("%x", sha256.Sum256([]byte(tok_secret)))
		rw.Write([]byte("Your token " + tok_secret))

	}
}
func check_token() {
	try := 0
	for try < 1 {
		//открываем файл с токенами
		file_token, _ := os.Open("tok.txt")
		tokens, _ := ioutil.ReadAll(file_token)
		//делаем список токенов

		list_tokens := strings.Split(string(tokens), " ")
		//смотрим есть ли у нас такой токен
		coun := 0
		generate_token()
		// проверяем на повторения
		string_token := fmt.Sprintf("%x", sha256.Sum256([]byte(tok_secret)))
		for i := 0; i < len(list_tokens)-1; i++ {
			if list_tokens[i] == string_token {
				coun += 1
			}
		}
		// если их нет то записываем токен в файл
		if coun == 0 {
			message := []byte(string(tokens) + string_token + " ")
			ioutil.WriteFile("tok.txt", message, 0644)
			try = 1
		}
	}
}

func generate_token() {
	// алфавит для генерации токена
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	token := ""
	for i := 0; i < 20; i++ {
		rand.Seed(time.Now().UnixNano())
		some_number := rand.Intn(len(alphabet) - 1)
		token += string(alphabet[some_number])
	}
	//sha256.Sum256([]byte(token))
	tok_secret = token
	//fmt.Println(tok_secret)
}

func main() {
	http.HandleFunc("/reg", registration)
	http.ListenAndServe(":8000", nil)
}
