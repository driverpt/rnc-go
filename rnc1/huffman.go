package rnc1

type HuffmanTable struct {
	Count int
	Leaves [32]HuffmanLeaf
}

type HuffmanLeaf struct {
	Code uint32
	CodeLength int
	Value int
}

func ReadHuffmanTable(bitstream *BitStream) HuffmanTable {
	num := int(bitstream.ReadBits(0x1F, 5))
	if num == 0 {
		return HuffmanTable{Count: 0}
	}
	leafLen := make([]int, 32)
	leafMax := 1

	result := HuffmanTable{
		Count: 0,
		Leaves: [32]HuffmanLeaf{
		},
	}

	for i := 0; i < num; i++ {
		leafLen[i] = int(bitstream.ReadBits(0x0F, 4))
		if leafMax < leafLen[i] {
			leafMax = leafLen[i]
		}
	}

	count := 0
	value := uint32(0)

	for i := 1; i <= leafMax; i++ {
		for j := 0; j < num; j++ {
			if leafLen[j] != i {
				continue
			}

			result.Leaves[count].Code = mirrorBits(value, i)
			result.Leaves[count].CodeLength = i
			result.Leaves[count].Value = j
			value++
			count++
		}

		value <<= 1
	}

	result.Count = count

	return result
}

func ReadHuffman(table *HuffmanTable, stream *BitStream) int32 {
	i := 0
	for i = 0; i < table.Count; i++ {
		mask := uint32 ((1 << table.Leaves[i].CodeLength) - 1)
		if stream.Peek(mask) == table.Leaves[i].Code {
			break
		}
	}

	if i == table.Count {
		return -1
	}

	stream.Advance(table.Leaves[i].CodeLength)
	value := uint32(table.Leaves[i].Value)

	if value >= uint32(2) {
		value = 1 << (value - 1)

		value |= stream.ReadBits(value-1, table.Leaves[i].Value-1)
	}

	return int32(value)
}


func mirrorBits(value uint32, bits int) uint32 {
	top := uint32(1 << (bits - 1))
	bottom := uint32(1)

	for top > bottom {
		mask := top | bottom
		masked := value & mask

		if masked != 0 && masked != mask {
			value ^= mask
		}

		top >>= 1
		bottom <<= 1
	}

	return value
}
