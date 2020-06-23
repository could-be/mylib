package util

import "errors"

// 封装位操作

// FLAGS &= ~( X | Y | Z )
// 1. 将特定的某几位对应的整数X, Y, Z使用或（|）运算组合成一个新的整数N；
// 2. 将新的整数N按位取反(~),得到新的整数M；
// 3. 以M为基，对FLAGS进行与(&)运算。
func ClearSpecialBit(flag *int, ops ...int) error {
	if len(ops) == 0 {
		return nil
	}

	var opSet int
	for _, op := range ops {
		if !IsAtomicOperation(op) {
			return errors.New("仅支持清除原子操作")
		}
		opSet |= op
	}

	// 清除指定操作集合
	*flag &^= opSet
	return nil
}

// 如果把一个整数减去1，再和 原整数做与运算，会把该整数最右边一个1变成0.如二进制1100，减去1后变为1011，1100和1011做位与运算是1000.把1100最右边的1变成了0
// 那么一个整数的二进制表示中有多少个1，就可以进行多少次这样的操作。代码如下：
func NumberOf1(x int) int {
	var cnt int
	for x != 0 {
		x &= (x - 1)
		cnt++
	}
	return cnt
}

// 判断是否是原子操作
// 原子操作的定义是 二进制表示时，内部，仅有一个1
func IsAtomicOperation(op int) bool {
	return NumberOf1(op) == 1
}
