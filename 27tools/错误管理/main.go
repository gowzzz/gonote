package main

import (
	"github.com/pkg/errors"
	"fmt"
	"os"
)
func openfile()error{
	f,err:=os.Open("./1.txt")
	// err = errors.Errorf("whoops: %s", "foo")
	// err = errors.New("whoops")
	err = errors.WithStack(err)
	if err!=nil{
		// fmt.Println("err:",err)
		// fmt.Println(errors.WithMessage(err, "err"))
		// fmt.Println(errors.WithStack(err))
		// fmt.Println(errors.Wrap(err, "err"))
		return err
	}
	defer f.Close()
	return nil
}
type stackTracer interface {
    StackTrace() errors.StackTrace
}
func main(){
	
	err:=openfile()
	fmt.Printf("%+v", err)
	return 
	if err!=nil{
		fmt.Printf("err:%+v\n",err)
		fmt.Println("------------------")

		fmt.Println(errors.Cause(err))
	}
	
}