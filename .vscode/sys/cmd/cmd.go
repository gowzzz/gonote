package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/codeskyblue/go-sh"
)

func main() {
	res, err := sh.Command("ls", "l").Output()
	if err != nil {
		fmt.Println("cmd err:", err)
	}
	fmt.Println("res:", res)
	return

	str := `
Timestamp                           : Thu Nov  7 15:43:06 2019
Driver Version                      : 430.26
CUDA Version                        : 10.2

Attached GPUs                       : 4
GPU 00000000:02:00.0
    Utilization
        Gpu                         : 0 %
        Memory                      : 0 %
        Encoder                     : 0 %
        Decoder                     : 0 %
    GPU Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %
    Memory Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %
    ENC Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %
    DEC Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %

GPU 00000000:03:00.0
    Utilization
        Gpu                         : 0 %
        Memory                      : 0 %
        Encoder                     : 0 %
        Decoder                     : 0 %
    GPU Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %
    Memory Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %
    ENC Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %
    DEC Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %

GPU 00000000:82:00.0
    Utilization
        Gpu                         : 0 %
        Memory                      : 0 %
        Encoder                     : 0 %
        Decoder                     : 0 %
    GPU Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %
    Memory Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %
    ENC Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %
    DEC Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %

GPU 00000000:83:00.0
    Utilization
        Gpu                         : 0 %
        Memory                      : 0 %
        Encoder                     : 0 %
        Decoder                     : 0 %
    GPU Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %
    Memory Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %
    ENC Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %
    DEC Utilization Samples
        Duration                    : 18446744073709.22 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %
        `
	buf := bufio.NewReader(bytes.NewReader([]byte(str)))
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				fmt.Println("eof")
				return
			}
			fmt.Println("err2:", err)
		}
		fmt.Println("line:", line)
	}
	// execCommand("ls", nil)
}
