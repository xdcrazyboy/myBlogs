package main

import (
	"reflect"
	"fmt"
	"runtime"
	"math"
)

func eval(a int, b int, op string) (int, error) {
	switch op{
	case "+":
		return a + b,nil
	case "-":
		return a - b,nil
	case "*":
		return a * b,nil
	case "/":
		// return a / b
		q, _ := div(a,b)
		return q, nil
	default:
		return 0, fmt.Errorf("unsupported operation: %s", op)

	}
}

func div(a,b int) (q, r int){
	return a / b, a % b
}

func apply(op func(int,int) int, a,b int) int {
	//反射获取函数名
	p := reflect.ValueOf(op).Pointer()
	opName := runtime.FuncForPC(p).Name()
	fmt.Println("Calling function %s with args " + "(%d, %d)", opName, a, b)
	return op(a,b)
}

func pow(a,b int) int {
	return int(math.Pow(float64(a),float64(b)))
}

//可变参数列表
func sum(numbers ...int) int {
	s := 0
	for i := range numbers {
		s += numbers[i]
	}
	return s
}


func main(){
	// fmt.Println(eval(3,4,"/"))
	if res, err := eval(3,4,"x"); err != nil{
		fmt.Println("Error:", err)
	}else{
		fmt.Println(res)
	}
	q,r := div(13,3)
	fmt.Println(q,r)

	//调用
	// fmt.Println(apply(pow,3,4))
	//匿名函数
	fmt.Println(apply(func(a,b int) int{
		return int(math.Pow(float64(a), float64(b)))
	},3,4))

	fmt.Println(sum(1,2,3,4,5,6))
}

/** 要点总结
1. 返回值类型写在最后面
2. 可返回多个值
3. 函数作为参数，还可以协成匿名函数
4. 没有默认参数，可选参数。 有参数列表

*/