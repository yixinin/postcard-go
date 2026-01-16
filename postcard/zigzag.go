package postcard

func zigzagEncodeInt16(n int16) uint16 {
	return uint16((n << 1) ^ (n >> 15))
}

func zigzagDecodeInt16(n uint16) int16 {
	return int16((n>>1)&0x7FFF) ^ -(int16(n & 0b1))
}

func zigzagEncodeInt32(n int32) uint32 {
	return uint32((n << 1) ^ (n >> 31))
}

func zigzagDecodeInt32(n uint32) int32 {
	return int32((n>>1)&0x7FFFFFFF) ^ -(int32(n & 0b1))
}

func zigzagEncodeInt64(n int64) uint64 {
	return uint64((n << 1) ^ (n >> 63))
}

func zigzagDecodeInt64(n uint64) int64 {
	return int64((n>>1)&0x7FFFFFFFFFFFFFFF) ^ -(int64(n & 0b1))
}

func zigzagEncodeInt(n int) uint {
	return uint(zigzagEncodeInt64(int64(n)))
}

func zigzagDecodeInt(n uint) int {
	return int(zigzagDecodeInt64(uint64(n)))
}
