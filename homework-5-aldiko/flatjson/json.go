package flatjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Unmarshal , не смог написать функцию рекурсии для любой другой структуры (про метод construct)
// поэтому Unmarshal рабоатет для примера приведенной в main
func Unmarshal(data []byte, dst interface{})  error {
	ref:=reflect.ValueOf(dst)

	switch ref.Kind() {
	// провека на указатель
	case reflect.Ptr:
		// проверка на валидность json
		if json.Valid(data){
			jsonMap := make(map[string]interface{})// для создания мапа с json
			err := json.Unmarshal([]byte(data), &jsonMap)
			if err!=nil {
				return fmt.Errorf("unmarshal error: %w",err)
			}

			obj := OutStruct{}
			path := []string{}

			construct(jsonMap, &obj, path)// рекурсия для парсинга мапа и распределние по филдам объекта

			val :=ref.Elem()
			for i:=0;i<val.NumField();i++ {
				field := val.Field(i)
				if field.CanSet() && val.Type().Field(i).Name == "Name" {
					field.SetString(obj.Name)
				}
				if field.CanSet() && val.Type().Field(i).Name == "Age" {
					field.SetInt(obj.Age)
				}
				if field.CanSet() && val.Type().Field(i).Name == "EducationDegree" {
					field.SetString(obj.EducationDegree)
				}
				if field.CanSet() && val.Type().Field(i).Name== "EducationUniversity" {
					field.SetString(obj.EducationUniversity)
				}
				if field.CanSet() && val.Type().Field(i).Name == "EducationFacultyName" {
					field.SetString(obj.EducationFacultyName)
				}
				if field.CanSet() && val.Type().Field(i).Name == "EducationFacultyDepartment" {
					field.SetString(obj.EducationFacultyDepartment)
				}
				if field.CanSet() &&val.Type().Field(i).Name == "EducationFacultyAdviserFirstName" {
					field.SetString(obj.EducationFacultyAdviserFirstName)
				}
				if field.CanSet() && val.Type().Field(i).Name == "EducationFacultyAdviserLastName" {
					field.SetString(obj.EducationFacultyAdviserLastName)
				}
				if field.CanSet() && val.Type().Field(i).Name == "EducationFacultyAdviserDegree" {
					field.SetString(obj.EducationFacultyAdviserDegree)
				}
				if field.CanSet() && val.Type().Field(i).Name == "EducationFacultyAdviserArticleCount" {
					field.SetInt(int64(obj.EducationFacultyAdviserArticleCount))
				}
			}
		}else {
			return fmt.Errorf("unmarshal error: %w",errors.New("not valid data for json"))
		}
	default:
		return fmt.Errorf("unmarshal error: %w",errors.New("not pointer"))
	}
	return nil
}

func Marshal(src interface{}) ([]byte, error) {
	// TODO: Write code here
	return nil, nil
}


type OutStruct struct {
	Name                                string
	Age                                 int64
	EducationDegree                     string
	educationAverageGrade               float32 // Приватное поле должно остаться пустым
	EducationUniversity                 string
	EducationFacultyName                string
	EducationFacultyDepartment          string
	EducationFacultyAdviserFirstName    string
	EducationFacultyAdviserLastName     string
	EducationFacultyAdviserDegree       string
	EducationFacultyAdviserArticleCount int32
}


func construct(json map[string]interface{}, obj *OutStruct, path []string) {
	val := reflect.ValueOf(obj)

	for key := range json {
		value := reflect.ValueOf(json[key])

		switch value.Kind() {
		case reflect.Map:
			construct(value.Interface().(map[string]interface{}), obj, append(path, key))
		case reflect.Float64:
			field := val.Elem().FieldByName(strings.Join(path[:], "") + key)
			if field.IsValid() {
				field.SetInt(int64(value.Float()))
			}
		case reflect.String:
			field := val.Elem().FieldByName(strings.Join(path[:], "") + key)
			field.SetString(value.String())
		}
	}
}



