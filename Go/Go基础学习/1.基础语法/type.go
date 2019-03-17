// bool,string
// (u)int,(u)int8、16/32/64，uintptr
// byte,rune(相当于char)
// float32，float64，complex64,complex64(复数)

package main
import (
	"math/cmplx"
	"fmt"
	"math"
)

//欧拉公式，使用复数
func euler(){
	// c := 3 + 4i
	// fmt.Println(cmplx.Abs(c))
	fmt.Println(cmplx.Exp(1i * math.Pi) + 1)
	
}

//类型转换都是强制的，没有隐形

func triangle(){
	var a,b int = 3,4
	var c int 
	c = int(math.Sqrt(float64(a * a + b * b)))
	fmt.Println(c)
}

//常量定义

func consts(){
	const filename = "abc.txt"
	const a,b = 3,4
	var x int
	x  = int(math.Sqrt(a*a + b*b))
	fmt.Println(x)
}

//枚举类型
func enums(){
	const(
		// cpp = 0
		// java = 1
		// python = 2
		// golang = 3
		cpp = iota
		_
		python
		golang
		javascript
	)
	fmt.Println(cpp,javascript,python,golang)

	const(
		B = 1 << (10 * iota)
		KB
		MB
		GB
		TB
		PB
	)
	fmt.Println(B,KB,MB,GB,TB,PB)
}

func main(){
	triangle()

	euler()

	consts()

	enums()

}

