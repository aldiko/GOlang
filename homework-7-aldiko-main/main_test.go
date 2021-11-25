package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)


// тесты запускаю сразу все все, а не какой-то определенный ибо в первом сразу создается тестовый юзер,
// который использвуется в других тестах
func TestSaveOne(t *testing.T){
	req := httptest.NewRequest(http.MethodPost,`/saveOne?login=aldiko&password=orazbek`,nil)
	w :=httptest.NewRecorder()
	saveOne(w,req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err!=nil{
		t.Errorf("Test error saveOne: %v",err)
	}
	if string(data) != "User added!"{
		t.Errorf("User not added, got %v",string(data))
	}
}

func TestGetOne(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost,`/getOne?login=aldiko`,nil)
	w := httptest.NewRecorder()
	getOne(w,req)
	res:=w.Result()
	defer res.Body.Close()
	data, err:= ioutil.ReadAll(res.Body)
	if err!=nil {
		t.Errorf("Test error getOne: %v",err)
	}
	ans := newUser("aldiko","orazbek")
	if string(data)!=fmt.Sprintf("%s",ans) {
		t.Errorf("User not found, get %s",string(data))
	}
}

func TestGetAll(t *testing.T)  {
	req:= httptest.NewRequest(http.MethodGet,`/getAll`,nil)
	w := httptest.NewRecorder()
	getAll(w,req)
	res:= w.Result()
	defer res.Body.Close()
	data, err:= ioutil.ReadAll(res.Body)
	if err!=nil {
		t.Errorf("Test error getAll: %v",err)
	}
	ans := newUser("aldiko","orazbek")
	if string(data)!=fmt.Sprintf("All Users:\nUser: %v\n",ans) {
		t.Errorf("Test error getAll, get %s",string(data))
	}
}
