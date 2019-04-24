package main
import(
	"fmt"
)
func main(){
	var year int32
	for{
		fmt.Print("import year:")
		_,err:=fmt.Scanln(&year)
		if err!=nil{
			fmt.Println("err:",err)
			continue
		}
		fmt.Println("year:",year)
	}
}