package main

import "fmt"
import "container/list"

func main() {
	s := []int{3, 2, 1, 4, 6, 5, 8, 2, 9}
	fmt.Println(s)
	sonicSortRecursive(s, 0, len(s)-1)
	// sonicSortStack(s, 0, len(s)-1)
	fmt.Println(s)
}

// insertSort 认为从左侧第一个元素开始是有序的，从右侧的第二个元素开始和左侧有序的所有元素比较并插入到合适位置，依次向右执行到最后
func insertSort(sl []int) {
	// 最优：已经是有序的了，只进行了依次向右遍历O(n);最差:完全逆序，向右遍历一遍，而且每个元素都向左比较O(n^2)

	// 下标从0~ len-1  假设排序从小到大排序
	// 我们先假设下标是0的元素是已排序的（因为只有一个元素）
	// 假设下标是1的元素且向右的所有元素都是未排序的
	// 我们从下标为1的元素一次2向右遍历直到最后一个元素，每次遍历保证被遍历的元素本身以及向左都是已排序的。

	// 被遍历到的元素，向左遍历；被遍历到的元素向右移动一个
	// 终止向左有两种可能：
	// 		1.遇到比他小的元素了
	// 		2.遇到下标是0的元素了（到头了）
	unsorted_suf := 1
	slen := len(sl)
	for i := unsorted_suf; i <= slen-1; i++ {
		tmp := sl[i] //先备份当前处理的未排序的元素，留在后边插入
		for j := i; j >= 0; j-- {
			if j == 0 || sl[j-1] < tmp {
				//左边没有元素了
				sl[j] = tmp
				break
			}
			sl[j] = sl[j-1]
		}

	}
	fmt.Println("slen:", slen)

}

// shellSort 是在 shellSort 基础上，让数组大体有序，然后最后完全有序。之前是2的倍数，有人提出使用基数
func shellSort(sl []int) {
	slen := len(sl)
	inter := slen / 2 //4   =>2
	for {
		tmp := slen / inter
		for i := 0; i < inter; i++ {
			start := tmp * i
			end := tmp * (i + 1)
			insertSort(sl[start:end])
		}
		if inter <= 1 {
			break
		}
		inter /= 2
	}
}

// selectSort 是从左到右 遍历，用当前遍历到的元素向右比较，找到最小的一个元素与之交换。
func selectSort(sl []int) {
	//第一个循环 O(n),第二个循环O(n)=>O(n^2)
	slen := len(sl)
	for i := 0; i < slen; i++ {
		minSuffix := i
		for j := i; j < slen; j++ {
			if sl[minSuffix] > sl[j] {
				minSuffix = j
			}
		}
		if minSuffix != i {
			sl[minSuffix], sl[i] = sl[i], sl[minSuffix]
		}
	}
}

// bubbleSort 冒泡 遍历数组，向右依次相互比较，大的右移，第一次遍历保证了最右边是最大的。依次遍历直到最左边
func bubbleSort(sl []int) {
	slen := len(sl)
	slen--
	for {
		for j := 0; j < slen; j++ {
			if sl[j] > sl[j+1] {
				sl[j], sl[j+1] = sl[j+1], sl[j]
			}
		}
		slen--
		if slen == 0 {
			break
		}
	}
}

func sonicSortRecursive(sl []int, start, end int) {
	povitIndex := partitionRecursive(sl, start, end)
	if povitIndex > start {
		sonicSortRecursive(sl, start, povitIndex-1)
	}
	if povitIndex+1 < end {
		sonicSortRecursive(sl, povitIndex+1, end)
	}
}

type Povit struct {
	Start int
	End   int
}

func sonicSortStack(sl []int, start, end int) {
	list := list.New()
	list.PushBack(Povit{
		Start: start,
		End:   end,
	})
	// 先入基础栈
	for {
		fmt.Println("list.Len():", list.Len())
		if list.Len() == 0 {
			break
		}
		// 然后循环取栈，根据需求入栈，然后取栈。。。直到为空
		e := list.Back()
		list.Remove(e)
		if p, ok := e.Value.(Povit); !ok {
			panic("type err")
		} else {
			povitIndex := partitionRecursive(sl, p.Start, p.End)
			if povitIndex > p.Start {
				list.PushBack(Povit{
					Start: p.Start,
					End:   povitIndex - 1,
				})
			}
			if povitIndex+1 < p.End {
				list.PushBack(Povit{
					Start: povitIndex + 1,
					End:   p.End,
				})
			}
		}
	}
}

func partitionRecursive(sl []int, start, end int) (index int) {
	povit := sl[start]
	left, right := start, end
	for {
		if left >= right {
			break
		}
		// right
		for {
			if left < right && sl[right] > povit {
				right--
			} else {
				break
			}
		}
		// left
		for {
			if left < right && sl[left] <= povit {
				left++
			} else {
				break
			}

		}

		//如果此时left和right重合了，就没必要交换了，让基准点和重合点值交换,且就让重合点作为下一个基准点返回
		if left < right {
			sl[right], sl[left] = sl[left], sl[right]
		}
	}
	sl[left], sl[start] = povit, sl[left] //选的sl[0]作为的基准点
	index = left
	return
}
