package operation

import (
	"encoding/json"
	"fmt"
	"operationmanage/operation/config"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"

	// "operationmanage/operation" //顺便初始化配置文件
	"os/exec"
)

func GetSysInfo() (gpuInfo, memoryInfo, cpusInfo, diskInfo []byte, errCode int, errMsg string) {
	gpu, err := GetGpuInfo()
	if err != nil {
		return nil, nil, nil, nil, ERROR_OPT_SERVER, GetMsg(ERROR_OPT_SERVER)
	}
	gpuInfo, err = json.Marshal(gpu)
	if err != nil {
		return nil, nil, nil, nil, ERROR_OPT_SERVER, GetMsg(ERROR_OPT_GPU_JSON)
	}
	memory, cpus, disk := SysInfo()
	memoryInfo, err = json.Marshal(memory)
	if err != nil {
		return nil, nil, nil, nil, ERROR_OPT_SERVER, GetMsg(ERROR_OPT_MEM_JSON)
	}
	cpusInfo, err = json.Marshal(cpus)
	if err != nil {
		return nil, nil, nil, nil, ERROR_OPT_SERVER, GetMsg(ERROR_OPT_CPU_JSON)
	}
	diskInfo, err = json.Marshal(disk)
	if err != nil {
		return nil, nil, nil, nil, ERROR_OPT_SERVER, GetMsg(ERROR_OPT_DISK_JSON)
	}

	return
}

type GpuInfo struct {
	Name      string
	Num       int
	Used      int
	Threshold int
}
type GpuInfos []GpuInfo

func GetGpuInfo() (GpuInfos, error) {
	// 先获取gpu数量，再
	cmdOutput, err := exec.Command("/bin/bash", "-c", "nvidia-smi -L").Output()
	if err != nil {
		fmt.Println("cmd err:", err)
		return nil, err
	}
	lines, err := splitByLine(cmdOutput)
	if err != nil {
		fmt.Println("splitByLine err:", err)
		return nil, err
	}
	var gpuInfos GpuInfos
	for _, line := range lines {
		var gpuInfo = GpuInfo{Threshold: config.GetConf().Services.Common.GpuThreshold}
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
		cmdstr := "nvidia-smi -d UTILIZATION -q -i " + splitOutput[1]
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
		for _, line := range lines {
			fmt.Println("line:", line)
			if strings.Contains(line, "Gpu") {
				resp := regexp.MustCompile("[0-9]+").FindString(line)
				fmt.Println("resp:", resp)
				if len(resp) > 0 {
					gpuInfo.Used, err = strconv.Atoi(resp)
					if err != nil {
						fmt.Println("err:", err)
						break
					}
				}
			}
		}
		gpuInfos = append(gpuInfos, gpuInfo)
	}
	return gpuInfos, nil
}

type SysUseInfo struct {
	Totol     uint64  //M
	Free      uint64  //M
	Used      uint64  //M
	Usage     float64 //M
	Threshold int
}

func SysInfo() (*SysUseInfo, []SysUseInfo, *SysUseInfo) {
	v, _ := mem.VirtualMemory()
	fmt.Printf("        Mem       : %v MB  Free: %v MB Used:%v Usage:%f%%\n", v.Total/1024/1024, v.Available/1024/1024, v.Used/1024/1024, v.UsedPercent)
	var memoryInfo = &SysUseInfo{
		Totol:     v.Total / 1024 / 1024,
		Free:      v.Available / 1024 / 1024,
		Used:      v.Used / 1024 / 1024,
		Usage:     v.UsedPercent,
		Threshold: config.GetConf().Services.Common.MemoryThreshold,
	}
	// c, _ := cpu.Info()
	cc, _ := cpu.Percent(time.Second, false)
	fmt.Printf("        CPU:  %f\n", cc)
	var cpuInfos []SysUseInfo
	for _, c := range cc {
		cpuInfos = append(cpuInfos, SysUseInfo{
			Usage:     c,
			Threshold: config.GetConf().Services.Common.CpuThreshold,
		})
	}
	// fmt.Printf("        CPU Used    : used %f%% \n", cc[0])
	d, _ := disk.Usage("/")
	fmt.Printf("        HD        : %v GB  Free: %v GB Usage:%f%%\n", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)
	var diskInfo = &SysUseInfo{
		Totol:     d.Total / 1024 / 1024,
		Free:      d.Free / 1024 / 1024,
		Used:      d.Used / 1024 / 1024,
		Usage:     d.UsedPercent,
		Threshold: config.GetConf().Services.Common.DiskThreshold,
	}
	return memoryInfo, cpuInfos, diskInfo
	// n, _ := host.Info()
	// nv, _ := net.IOCounters(true)
	// boottime, _ := host.BootTime()
	// btime := time.Unix(int64(boottime), 0).Format("2006-01-02 15:04:05")
	// if len(c) > 1 {
	// 	for _, sub_cpu := range c {
	// 		modelname := sub_cpu.ModelName
	// 		cores := sub_cpu.Cores
	// 		fmt.Printf("        CPU       : %v   %v cores \n", modelname, cores)
	// 	}
	// } else {
	// 	sub_cpu := c[0]
	// 	modelname := sub_cpu.ModelName
	// 	cores := sub_cpu.Cores
	// 	fmt.Printf("        CPU       : %v   %v cores \n", modelname, cores)
	// }
	// fmt.Printf("        Network: %v bytes / %v bytes\n", nv[0].BytesRecv, nv[0].BytesSent)
	// fmt.Printf("        SystemBoot:%v\n", btime)
	// fmt.Printf("        OS        : %v(%v)   %v  \n", n.Platform, n.PlatformFamily, n.PlatformVersion)
	// fmt.Printf("        Hostname  : %v  \n", n.Hostname)
}
