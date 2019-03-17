package main
import (
	"fmt"
	"io/ioutil"
)

//switch ,不需要写break，会自动break

func eval(a,b int, op string) int {
	var result int
	switch op{
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		result = a / b
	default:
		panic("unsupported operator:" + op) 
	}
	return result
}

//switch 后可以没有表达式
func grade(score int) string {
	result := ""
	switch {
	case score < 0 || score > 100:
		panic(fmt.Sprintf("Wrong score : %d", score)) 
	case score < 60:
		result = "F"
	case score < 80:
		result = "C"
	case score < 90:
		result = "B"
	case score <= 100:
		result = "A"
		
	}
	return result
}

func main(){
	const filename = "abc.txt"
	// contents, err := ioutil.ReadFile(filename)
	// if err != nil{
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Printf("%s\n", contents)
	// }
	if contents, err := ioutil.ReadFile(filename); err == nil {
		fmt.Println(string(contents))
	}else{
		fmt.Println("cannot print file contents:",err)
	}

	fmt.Println(eval(4,5,"*"))

	fmt.Println(
		grade(90),
		grade(102),
	)

}