[![Open in Visual Studio Code](https://classroom.github.com/assets/open-in-vscode-f059dc9a6f8d3a56e377f745f24479a46679e63a5d9fe6f495e02850cd0d8118.svg)](https://classroom.github.com/online_ide?assignment_repo_id=5734755&assignment_repo_type=AssignmentRepo)
# homework-1-starter
homework-1-starter

1) Необходимо зарегистрироваться на github.com
2) Необходимо создать приватный репозиторий в github для домашнего задания
3) Необходимо написать программу и реализовать в ней следующие функции:
func SortSlice(slice []int) - функция сортировки слайса. Запрещается использовать пакет sort
func IncrementOdd(slice []int) - функция, которая проходится по нечетным позициям и увеличивает число на 1
func PrintSlice(slice []int) - функция, печатающая слайс
func ReverseSlice(slice []int) - функция, переворачивающая слайс
4)Необходимо реализовать функцию следующего вида:
 func appendFunc(dst func([]int), src ...func([]int)) func([]int) 
Данная функция принимает в качестве аргумента некоторую функцию по обработке слайсов dst и неограниченное число других функций по обработке, которые нужно "присоединить" к функции dst и вернуть уже новую функцию
