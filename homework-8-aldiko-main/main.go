package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type User struct {
	ID       int
	Username string
	Password string
}

type CryptoWallet struct {
	Name   string
	Amount int64
	sync.RWMutex
}

func (c *CryptoWallet) Mine() {
	time.Sleep(10 * time.Second)
	c.Lock()
	c.Amount++
	c.Unlock()
}

type UserCryptoWallet struct {
	UserCrypto *User
	CryptoWallets []*CryptoWallet
}

type MyServer struct {
	ServerUsers []*UserCryptoWallet
}

// не делал проверку на Nil, ибо мидлвейр при отсутствии данных не даст запустить функцию
func (serv *MyServer) authFunc(login string, pass string) bool{
	for _,v := range serv.ServerUsers {
		if v.UserCrypto.Username == login  && v.UserCrypto.Password == pass{
			return true
		}
	}
	return false
}

func (serv *MyServer) GetWalletById(writer http.ResponseWriter, request *http.Request) {

	// провека на аунтефикацию
	userData:= strings.Fields(request.Context().Value("user").(string))
	ok := serv.authFunc(userData[0],userData[1])
	if !ok {
		writer.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(writer,"login or password is not correct (login fail)")
		return
	}

	params := mux.Vars(request)

	// логика ввыода по id криптокошельков
	id,err := strconv.Atoi(params["id"])
	if err!=nil {
		log.Printf("GetCryptoWalletById: convert to Int: %v",err)
		fmt.Fprintf(writer,"Invalid value of user Id\n")
		writer.WriteHeader(500)
		return
	}
	existUser := false
	for _, user := range serv.ServerUsers {
		if user.UserCrypto.ID == id {
			existUser = true
			fmt.Fprintf(writer, "CryptoWallets for user id %d:\n",id)
			for i, wallet := range user.CryptoWallets {
				if wallet== nil {
					continue
				}
				fmt.Fprintf(writer, "%d. %s, amount %d\n",i+1,wallet.Name,wallet.Amount)
			}
			break
		}
	}
	// для проверки есть пользователь с таким Id
	if !existUser {
		fmt.Fprintf(writer, "User with that %d ID not found!\n",id)
		return
	}
}

func (serv *MyServer) RegUserById(writer http.ResponseWriter, request *http.Request) {
	// провека на аунтефикацию
	userData:= strings.Fields(request.Context().Value("user").(string))
	ok := serv.authFunc(userData[0],userData[1])
	if !ok {
		writer.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(writer,"login or password is not correct (login fail)")
		return
	}

	params := mux.Vars(request)

	// логика регистрации юзера
	username := request.FormValue("username")
	password := request.FormValue("password")
	if username == "" || password == "" {
		fmt.Fprintln(writer,"Incorrect data for registration")
		return
	}

	id,err := strconv.Atoi(params["id"])
	if err!=nil {
		log.Printf("GetCryptoWalletById: convert to Int: %v",err)
		fmt.Fprintf(writer,"Invalid value of user Id\n")
		writer.WriteHeader(500)
		return
	}

	//проверка есть ли юзер с таким id
	for _,v := range serv.ServerUsers {
		if v.UserCrypto.ID == id {
			writer.WriteHeader(400)
			fmt.Fprintf(writer, "User with %d already exist!\n", id)
			return
		}
	}

	sliceCrypto := make([]*CryptoWallet,3)
	sliceCrypto[0] = &CryptoWallet{}
	userCryptoWallet := &UserCryptoWallet{
		UserCrypto: &User{
			ID: id,
			Username: username,
			Password: password,
		},
		CryptoWallets: sliceCrypto,
	}
	serv.ServerUsers = append(serv.ServerUsers,userCryptoWallet)
	fmt.Fprintln(writer,"User registration success!")
}

func (serv *MyServer) GetCryptoWalletByName(writer http.ResponseWriter, request *http.Request) {
	// провека есть такой юзер в базе
	userData:= strings.Fields(request.Context().Value("user").(string))
	ok := serv.authFunc(userData[0],userData[1])
	if !ok {
		writer.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(writer,"login or password is not correct (login fail)")
		return
	}

	params := mux.Vars(request)

	// основная логика

	nameOfWallet := strings.ToLower(params["name"])
	if nameOfWallet=="" {
		fmt.Fprintf(writer,"Incorrect parameter of wallet name")
		return
	}

	// чекер для прокери кошелька
	exist := false

Loop:
	for _,v := range serv.ServerUsers {
		for _,v2 := range v.CryptoWallets {
			// по сути флаг того что больше каошельков нет
			if v2 == nil {
				break Loop
			}
			if strings.ToLower(v2.Name) == nameOfWallet {
				exist = true
				// проверка владельца кошелька с юзером сделавший запрос
				if v.UserCrypto.Username != userData[0] && v.UserCrypto.Password != userData[1] {
					writer.WriteHeader(404)
					fmt.Fprintf(writer,"You do not have access for this wallet!\n")
					return
				}
				fmt.Fprintf(writer,"%s :\n\tAmount: %d\n",v2.Name,v2.Amount)
				break Loop
			}
		}
	}

	if !exist {
		fmt.Fprintf(writer,"Wallet with %s name not found!\n",nameOfWallet)
		return
	}
}

func (serv *MyServer) CreateWalletByName(writer http.ResponseWriter, request *http.Request) {
	// провека есть такой юзер в базе
	userData:= strings.Fields(request.Context().Value("user").(string))
	ok := serv.authFunc(userData[0],userData[1])
	if !ok {
		writer.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(writer,"login or password is not correct (login fail)")
		return
	}

	params := mux.Vars(request)

	// основная логика

	nameOfWallet := strings.ToLower(params["name"])
	if nameOfWallet=="" {
		fmt.Fprintf(writer,"Incorrect parameter of wallet name")
		return
	}

	// проверка есть ли у этого юзеря кошелек с таким именем
	exists := false
	for _,v := range serv.ServerUsers {
		if v.UserCrypto.Username == userData[0] && v.UserCrypto.Password == userData[1] {

			for _,wallet := range v.CryptoWallets {
				if wallet==nil {
					continue
				}
				if strings.ToLower(wallet.Name) == strings.ToLower(nameOfWallet) {
					exists = true
				}
			}

			if exists {
				writer.WriteHeader(400)
				fmt.Fprintf(writer,"Wallet with %s name already exists!\n",nameOfWallet)
				return
			}

			userWallet := &CryptoWallet{
				Name: nameOfWallet,
				Amount: 0,
			}

			v.CryptoWallets = append(v.CryptoWallets,userWallet)
			fmt.Fprintln(writer,"Adding new Wallet is success!")

			break
		}
	}
}


var mineStop = make(chan bool) // не смог придумать реализацию без глобальной переменной

func (serv *MyServer) StartMine(writer http.ResponseWriter, request *http.Request) {
	// провека есть такой юзер в базе
	userData:= strings.Fields(request.Context().Value("user").(string))
	ok := serv.authFunc(userData[0],userData[1])
	if !ok {
		writer.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(writer,"login or password is not correct (login fail)")
		return
	}

	params := mux.Vars(request)

	nameOfWallet := strings.ToLower(params["name"])
	if nameOfWallet=="" {
		fmt.Fprintf(writer,"Incorrect parameter of wallet name")
		return
	}

	ok = false

	for _,v := range serv.ServerUsers {
		if v.UserCrypto.Username == userData[0] && v.UserCrypto.Password == userData[1] {

			for _,wallet := range v.CryptoWallets {
				if wallet==nil {
					continue
				}
				if strings.ToLower(wallet.Name) == strings.ToLower(nameOfWallet) {
					ok = true
					fmt.Fprintf(writer,"Mine Started...\n")
					go func() {
					Loop:
						for  {
							select {
							case <- mineStop:
								break Loop
							default:
								wallet.Mine()
							}
						}
					}()
				}
			}

			if !ok {
				fmt.Fprintf(writer,"you dont have aceess or wallet '%s' not exists\n",nameOfWallet)
				return
			}

			break
		}
	}
}

// у меня выходит ошибка при этом запросе, не знаю как исправить
func  (serv *MyServer) StopMine(writer http.ResponseWriter, request *http.Request) {
	// провека есть такой юзер в базе
	userData:= strings.Fields(request.Context().Value("user").(string))
	ok := serv.authFunc(userData[0],userData[1])
	if !ok {
		writer.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(writer,"login or password is not correct (login fail)")
		return
	}

	params := mux.Vars(request)

	nameOfWallet := strings.ToLower(params["name"])
	if nameOfWallet=="" {
		fmt.Fprintf(writer,"Incorrect parameter of wallet name")
		return
	}

	ok = false

	for _,v := range serv.ServerUsers {
		if v.UserCrypto.Username == userData[0] && v.UserCrypto.Password == userData[1] {

			for _,wallet := range v.CryptoWallets {
				if wallet==nil {
					continue
				}
				if strings.ToLower(wallet.Name) == strings.ToLower(nameOfWallet) {
					ok = true
					mineStop <-true
					fmt.Fprintf(writer,"Mine stopped!\n")
				}
			}

			if !ok {
				fmt.Fprintf(writer,"you dont have aceess or wallet '%s' not exists\n",nameOfWallet)
				return
			}

			break
		}
	}
}


// middlewares

func TimerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		writer.Header().Set("content-type", "application/json")
		next.ServeHTTP(writer,request)
		// запись в header Тут не сработает, судя по обсужденням в дискорде, сказали пока не париться
		// сказали что может покажут другой вариант
		writer.Header().Add("execution",time.Now().Sub(start).String())
	})
}

// UserValidationCheckMiddleware не смог полностью реализовать проверку аунтефикации в мидлвейре, так как не могу получить доступ к сущности сервер
// поэтому тут реализовал только валидацию, а саму проверку по "базе" сделаю в другой фунцкии
func UserValidationCheckMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		login := request.Header.Get("login")
		password := request.Header.Get("password")
		if login == "" || password == "" {
			fmt.Fprintln(writer,"user data not found in Headers")
			writer.WriteHeader(http.StatusForbidden)
			return
		}
		user:= fmt.Sprintf("%s %s",login,password)
		ctx := context.WithValue(request.Context(), "user", user)
		request = request.WithContext(ctx)
		next.ServeHTTP(writer,request)
	})
}

func main() {
	// создание подобие админ-юзера
	// log: admin, pass: admin
	// считывание идет с хедера запроса
	sliceCrypto := make([]*CryptoWallet,3)
	sliceCrypto[0] = &CryptoWallet{
		Name: "SteamBit",
		Amount: 35,
	}
	sliceCrypto[1] = &CryptoWallet{
		Name: "BadBat",
		Amount: 142,
	}
	userCryptoWallet := &UserCryptoWallet{
		UserCrypto: &User{
			ID: 1,
			Username: "admin",
			Password: "admin",
		},
		CryptoWallets: sliceCrypto,
	}


	server:= &MyServer{
		ServerUsers: []*UserCryptoWallet{
			userCryptoWallet,
		},
	}

	r := mux.NewRouter()

	r.HandleFunc("/app/user/{id:[0-9]+}",server.GetWalletById).Methods("GET")
	r.HandleFunc("/app/user/{id:[0-9]+}",server.RegUserById).Methods("POST")
	r.HandleFunc("/app/wallet/{name:[a-zA-Z]+}",server.GetCryptoWalletByName).Methods("GET")
	r.HandleFunc("/app/wallet/{name:[a-zA-Z]+}",server.CreateWalletByName).Methods("POST")
	r.HandleFunc("/app/wallet/{name:[a-zA-Z]+}/start",server.StartMine).Methods(http.MethodOptions)
	r.HandleFunc("/app/wallet/{name:[a-zA-Z]+}/stop",server.StopMine).Methods(http.MethodOptions)

	r.Use(TimerMiddleWare,UserValidationCheckMiddleware)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("Server Error:%v",err)
		return
	}
}
