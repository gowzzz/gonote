package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	str := `GPU 0: GeForce GTX 1080 Ti (UUID: GPU-10e4f209-2fab-93fe-c6d9-e0b120462b5e)
	GPU 1: GeForce GTX 1080 Ti (UUID: GPU-2391a0e8-b6bc-3094-c382-bfd218dbe57d)
	GPU 2: GeForce GTX 1080 Ti (UUID: GPU-ccf2e9ff-afbc-a6f9-0d62-dfa962549d14)
	GPU 3: GeForce GTX 1080 Ti (UUID: GPU-f88ffeec-ebe8-2d20-9cba-6a5666a798a8)`
	buf := bufio.NewReader(bytes.NewReader([]byte(str)))
	var lines []string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil && err != io.EOF {
			fmt.Println("err2:", err)
			break
		}
		res := strings.SplitN(line, ":", 2)
		fmt.Printf("res:%q\n", res)
		fmt.Println("line:", line)
		lines = append(lines, line)
		if err == io.EOF { //读取结束，会报EOF
			fmt.Println("eof")
			break
		}
	}
	fmt.Println("lines:", lines)
}
func main1() {
	str := "        Gpu                         : 222 %   "
	if strings.Contains(str, "Gpu") {
		resp := regexp.MustCompile("[0-9]+").FindString(str)
		if len(resp) > 0 {
			num, err := strconv.Atoi(resp)
			if err != nil {
				fmt.Println("err:", err)
			}
			fmt.Println(num)
		}
	} else {
		fmt.Println("find 0")
	}

	fmt.Println()
}
