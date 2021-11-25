package main

import "fmt"

func main() {
	sl :=[]int{10,2,3,7,19,11,8,4}
	fmt.Printf("INITIAL slice:\t %v\n", sl)
	SortSlice(sl)
	IncrementOdd(sl)
	PrintSlice(sl)
	ReverseSlice(sl)
	fmt.Println("\n\tAppendFunc:")
	a := AppendFunc(SortSlice,IncrementOdd,ReverseSlice)
	a(sl)
}


// SortSlice - функция сортировки слайса. Запрещается использовать пакет sort.
func SortSlice(slice []int) {
	//buble
	for i:=0; i< len(slice)-1; i++ {
		for j:=0; j < len(slice)-i-1; j++ {
		   if (slice[j] > slice[j+1]) {
			  slice[j], slice[j+1] = slice[j+1], slice[j]
		   }
		}
	 }
	 fmt.Printf("SortSlice:\t %v\n", slice)
}

// IncrementOdd - функция, которая проходится по нечетным позициям и увеличивает число на 1.
func IncrementOdd(slice []int) {
	for i := range slice {
		if(i%2!=0){
			slice[i]+=1
		}
	}
	fmt.Printf("IncrementOdd:\t %v\n", slice)
}

// PrintSlice - функция, печатающая слайс.
func PrintSlice(slice []int) {
	fmt.Printf("PrintSlice:\n")
	for i, v := range slice {
		fmt.Printf("\t\tindex: %d, value: %d\n",i,v)
	}
}

// ReverseSlice - функция, переворачивающая слайс
func ReverseSlice(slice []int) {
	for i,j := len(slice)-1, 0; i > j; i,j=i-1,j+1 {
		slice[i],slice[j]=slice[j],slice[i]
	}
	fmt.Printf("ReverseSlice:\t %v\n", slice)
}



// AppendFunc - данная функция принимает в качестве аргумента некоторую функцию по обработке слайсов dst
// и неограниченное число других функций по обработке, которые нужно "присоединить" к функции dst и вернуть уже новую функцию
func AppendFunc(dst func([]int), src ...func([]int)) func([]int) {
	return func (sl []int)  {
		dst(sl)
		for _, v := range src {
			v(sl)
		}
	}
}
