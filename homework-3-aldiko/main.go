package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
)

type Target struct {
	InputFilePath  string // Путь до файла с паролями
	OutputFilePath string // Путь до файла, куда должны записываться результаты
	*Connection
}

func main() {
	requestChan := make(chan *Request)
	responseChan := make(chan *Response)
	defer close(requestChan)
	defer close(responseChan)

	connection := &Connection{
		RequestConn:  requestChan,
		ResponseConn: responseChan,
	}

	target := &Target{
		InputFilePath:  "darkweb2017-top10000.txt",
		OutputFilePath: "output.txt",
		Connection:     connection,
	}

	server := NewVulnerableServer("MaprCheM56458", connection)

	go server.Run()

	// Пробовать запускать с разными контекстами
	ctx := context.Background()

	HackServer(ctx, target)
}

// все работает хорошо, все горутины точно успевают пробегаться
// и выполнить запрос и т.д (коментарий на 86 строке для проверки)
// но есть некоторый факт рандома, те если брать прям совсем последние пароли из списка
// то горутина найдет нужнный пароль, однако алгоритм записи может не успеть запистаь его 
// в файл. Это происходит только с последними паролями и то через раз. Происходит потому что 
// после пароли ИНОГДА не успевают записаться, из-за того сервер заканчивает программу через
// os exit. А так вроде все остальное работает
func HackServer(ctx context.Context, target *Target) {
	passList, err := readTxt(target.InputFilePath)

	if err != nil {
		fmt.Printf("Error HackServer: %s \n",err.Error())
		return	
	}

	messages := make(chan string)

	LOOP:
		for _, v := range passList {
			select {
			case <- ctx.Done():
				break LOOP
			default:
				go hackGoroutine(ctx, v, messages, target)
			}
		}

	err=writeLines(messages, target.OutputFilePath)
	if err!=nil {
		fmt.Printf("Write txt Error:%s",err.Error())
	}
}


func hackGoroutine(ctx context.Context, pass string, ch chan string, target *Target){
	req := Request{
		ctx,
		pass,
	}

	SendRequest(target.Connection, &req)

	res := <-target.Connection.ResponseConn

	if res.Pass {
		ch <- res.Password + ":true"
		//fmt.Println("Found")//причина почему закоментил выше
		ctx.Done()
	} else {
		ch <- res.Password + ":false"
	}
}

func readTxt(path string)([]string,error){
	file, err := os.Open(path)
    if err != nil{
    	return nil, fmt.Errorf("readTXT error open: %w",err)     
    }
    defer file.Close() 
    
	var res []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        res = append(res, scanner.Text())
    }
    return res, scanner.Err()
}


// writeLines writes the lines to the given file.
func writeLines(lines chan string, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    w := bufio.NewWriter(file)
    for line := range lines {
        fmt.Fprintln(w, line)
    }
    return w.Flush()
}
