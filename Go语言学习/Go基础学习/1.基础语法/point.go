//指针 1. 指针不能运算  2.参数传递？ go只有值传递  3.cache类型

package main

import(
	"fmt"
)

func swapByVal(a, b int){
	b, a = a, b
}

func swapByRef(a, b *int){
	*b, *a = *a, *b
}


func swap(a, b int) (int,int){
	return b,a
}

func main(){
	a, b := 3, 4
	swapByVal(a,b)
	fmt.Println(a,b)
	
	swapByRef(&a,&b)
	fmt.Println(a,b)

	a,b = swap(a,b)
	fmt.Println(a,b)

}

