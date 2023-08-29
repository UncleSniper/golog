package golog

type BoolStack struct {
	fill uint8
	scalars []uint64
}

func(stack *BoolStack) Top() bool {
	index := len(stack.scalars) - 1
	scalar := stack.scalars[index]
	return scalar & uint64(1) != 0
}

func(stack *BoolStack) Push(top bool) {
	var scalar uint64
	if stack.fill == 0 || stack.fill == 64 {
		if top {
			scalar = 1
		}
		stack.scalars = append(stack.scalars, scalar)
		stack.fill = 1
	} else {
		index := len(stack.scalars) - 1
		scalar = stack.scalars[index] << 1
		if top {
			scalar |= uint64(1)
		}
		stack.scalars[index] = scalar
		stack.fill++
	}
}

func(stack *BoolStack) Pop() bool {
	index := len(stack.scalars) - 1
	scalar := stack.scalars[index]
	top := scalar & uint64(1) != 0
	if stack.fill == 1 {
		if index == 0 {
			stack.scalars = nil
			stack.fill = 0
		} else {
			stack.scalars = stack.scalars[:index]
			stack.fill = 64
		}
	} else {
		stack.scalars[index] = scalar >> 1
		stack.fill--
	}
	return top
}

func(stack *BoolStack) Replace(top bool) {
	index := len(stack.scalars) - 1
	scalar := stack.scalars[index]
	if top {
		scalar |= uint64(1)
	} else {
		scalar &^= uint64(1)
	}
	stack.scalars[index] = scalar
}

func(stack *BoolStack) IsEmpty() bool {
	return len(stack.scalars) == 0
}

func RepeatRune(r rune, times int) []rune {
	if times < 0 {
		return nil
	}
	buffer := make([]rune, times)
	for i := 0; i < times; i++ {
		buffer[i] = r
	}
	return buffer
}

type OrderRel uint

const (
	ORDR_GREATER_EQUAL OrderRel = iota
	ORDR_GREATER
	ORDR_EQUAL
	ORDR_LESS
	ORDR_LESS_EQUAL
)
