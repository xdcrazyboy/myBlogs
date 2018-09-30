// 变量定义  var a,b,c bool
// :=   只能在函数内使用， 简约，可以省略var

//变量类型写在变量名之后
//编译器可以推测变量类型

package main

import (
	"fmt"
)
//只能用var，只是包内部变量，不是全局
var aa = 3
var(
	ss = "def"
	vv = 3
)

func variableTypeDeduction(){
	var a,c,b,s = 2,3,true,"def"
	fmt.Println(a,c,b,s)
}

func variableShorter(){
	a,c,b,s := 3,4,true,"def"
	c = 5
	fmt.Println(a,c,b,s)
}

func main(){
	fmt.Println("Hello World!")
	variableTypeDeduction()
	variableShorter()

}