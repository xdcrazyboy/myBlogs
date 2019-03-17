package main
import(
	"fmt"
	"strconv"
	"os"
	"bufio"

)

func converToBin(n int) string {
	res := ""
	if n == 0 {
		return "0"
	}
	for ; n > 0; n /= 2 {
		lsb := n % 2
		res = strconv.Itoa(lsb) + res
	}
	return res
}

func printFile(filename string){
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	//省略其他
	for scanner.Scan(){
		fmt.Println(scanner.Text())
	}
}

//死循环，没有while
func forever(){
	for {
		fmt.Println("abc")
	}
}

func main(){
	fmt.Println(
		converToBin(5),
		converToBin(13),
		converToBin(325576135),
		converToBin(0),
	)

	printFile("abc.txt")

	// forever()
}