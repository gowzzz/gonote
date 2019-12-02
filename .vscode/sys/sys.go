package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
	// "operationmanage/operation" //顺便初始化配置文件
	"os/exec"
)

//
func splitByLine(in []byte) ([]string, error) {
	buf := bufio.NewReader(bytes.NewReader(in))
	var lines []string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil && err != io.EOF {
			fmt.Println("err2:", err)
			return nil, err
		}
		lines = append(lines, line)
		if err == io.EOF { //读取结束，会报EOF
			fmt.Println("eof")
			break
		}
	}
	return lines, nil
}

type GpuInfo struct {
	Name      string
	Num       int
	Used      int
	Threshold int
}
type GpuInfos []GpuInfo

func GetGpuInfo() {
	// 先获取gpu数量，再
	cmdOutput, err := exec.Command("/bin/bash", "-c", "nvidia-smi -L").Output()
	if err != nil {
		fmt.Println("cmd err:", err)
		return
	}
	lines, err := splitByLine(cmdOutput)
	if err != nil {
		fmt.Println("splitByLine err:", err)
		return
	}
	var gpuInfos GpuInfos
	for _, line := range lines {
		var gpuInfo = GpuInfo{}
		splitNOutput := strings.SplitN(line, ":", 2) //GPU 0 |rest
		if len(splitNOutput) == 0 {
			break
		}
		gpuInfo.Name = splitNOutput[0]
		splitOutput := strings.Split(gpuInfo.Name, " ") //GPU 0=>GPU|0
		if len(splitOutput) < 2 {
			break
		}
		gpuInfo.Num, err = strconv.Atoi(splitOutput[1])
		if err != nil {
			fmt.Println("Num Atoi err:", err)
			break
		}
		cmdstr := "nvidia-smi -d MEMORY -q -i " + splitOutput[1]
		cmdOutput, err := exec.Command("/bin/bash", "-c", cmdstr).Output()
		if err != nil {
			fmt.Println("cmd2 err:", err)
			break
		}
		fmt.Println(cmdstr+":", string(cmdOutput))
		lines, err := splitByLine(cmdOutput)
		if err != nil {
			fmt.Println("SplitByLine2 err:", err)
			break
		}
		for i := 0; i < len(lines); i++ {
			line := lines[i]
			fmt.Println("line:", line)
			if strings.Contains(line, "FB Memory Usage") {
				total := regexp.MustCompile("[0-9]+").FindString(lines[i+1])
				used := regexp.MustCompile("[0-9]+").FindString(lines[i+2])
				free := regexp.MustCompile("[0-9]+").FindString(lines[i+3])
				fmt.Println("total:", total)
				fmt.Println("used:", used)
				fmt.Println("free:", free)
				break

			}
		}
		gpuInfos = append(gpuInfos, gpuInfo)
	}
	return
}

type SysUseInfo struct {
	Totol     uint64  //M
	Free      uint64  //M
	Used      uint64  //M
	Usage     float64 //M
	Threshold int
}

func SysInfo() {
	v, _ := mem.VirtualMemory()
	fmt.Printf("        Mem       : %v MB  Free: %v MB Used:%v Usage:%f%%\n", v.Total/1024/1024, v.Available/1024/1024, v.Used/1024/1024, v.UsedPercent)

	// c, _ := cpu.Info()
	cc, _ := cpu.Percent(time.Second, false)
	fmt.Printf("        CPU:  %f\n", cc)
	var cpuInfos []SysUseInfo
	for _, c := range cc {
		cpuInfos = append(cpuInfos, SysUseInfo{
			Usage: c,
		})
	}
	// fmt.Printf("        CPU Used    : used %f%% \n", cc[0])
	// d, _ := disk.Usage("/")
	// fmt.Printf("        HD        : %v GB  Free: %v GB Usage:%f%%\n", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)
	dps, _ := disk.Partitions(true)
	for k, v := range dps {
		fmt.Printf("        Disk num:  %d  name:  %s\n", k, v.Mountpoint)
		d, err := disk.Usage(v.Mountpoint)
		if err != nil {
			fmt.Println(v.Mountpoint+" err:", err)
			continue
		}
		fmt.Printf("        HD        : %v GB  Free: %v GB Usage:%f%%\n", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)

	}
	return

}

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
func main() {
	GetGpuInfo()
	SysInfo()
}
