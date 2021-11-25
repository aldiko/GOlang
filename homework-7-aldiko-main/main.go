package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sync"
)

type User struct {
	login string
	pass string
}

func newUser(l string, p string) *User {
	return &User{
		login: l,
		pass: p,
	}
}

var users = sync.Map{}

func saveOne(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost{
		writer.WriteHeader(405)
		return
	}

	login := request.FormValue("login")
	pass := request.FormValue("password")

	user := newUser(login,pass)

	_, exist := users.LoadOrStore(login,user)
	if exist {
		_,_ = writer.Write([]byte("User already exist"))
		return
	}

	_,_ = writer.Write([]byte("User added!"))
}

func getOne(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost{
		writer.WriteHeader(405)
		return
	}

	l :=request.FormValue("login")
	user, ok := users.Load(l)
	if !ok {
		_, _ = writer.Write([]byte("User with that login not found "+l))
		return
	}

	_,_ = fmt.Fprintf(writer,"%v",user)
}

func getAll(writer http.ResponseWriter, request *http.Request) {
	_, _ = fmt.Fprintf(writer,"All Users:\n")
	users.Range(func(key, value interface{}) bool {
		_,_ = fmt.Fprintf(writer,"User: %v\n",value)
		return true
	})
}

func main() {

	// тестовые юзеры или же как имитация уже существующих пользователей в БД/системе
	users.Store("test",newUser("test","test"))
	users.Store("test2",newUser("test2","test2"))
	users.Store("test3",newUser("test3","test3"))

	http.HandleFunc("/saveOne", saveOne)

	http.HandleFunc("/getOne", getOne)

	http.HandleFunc("/getAll", getAll)


	http.HandleFunc("/",
		func(writer http.ResponseWriter, request *http.Request) {
			_ , _ = writer.Write([]byte("My server"))
		})

	err := http.ListenAndServe(":8080",nil)
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println("end")
}
