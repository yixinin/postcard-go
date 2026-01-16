package postcard

import (
	"encoding/binary"
	"math"
)

const (
	bitsPerByte      = 8
	bitsPerVarintByte = 7
)

func varintMax(size int) int {
	bits := size * bitsPerByte
	roundupBits := bits + (bitsPerVarintByte - 1)
	return roundupBits / bitsPerVarintByte
}

func maxOfLastByte(size int) uint8 {
	maxBits := size * 8
	extraBits := maxBits % 7
	return uint8((1 << extraBits) - 1)
}

func encodeVarintUint16(n uint16) []byte {
	maxLen := varintMax(2)
	buf := make([]byte, maxLen)
	if n < 128 {
		buf[0] = byte(n)
		return buf[:1]
	}
	value := n
	i := 0
	for value >= 128 {
		buf[i] = byte(value&0x7F) | 0x80
		value >>= 7
		i++
	}
	buf[i] = byte(value)
	return buf[:i+1]
}

func encodeVarintUint32(n uint32) []byte {
	maxLen := varintMax(4)
	buf := make([]byte, maxLen)
	if n < 128 {
		buf[0] = byte(n)
		return buf[:1]
	}
	value := n
	i := 0
	for value >= 128 {
		buf[i] = byte(value&0x7F) | 0x80
		value >>= 7
		i++
	}
	buf[i] = byte(value)
	return buf[:i+1]
}

func encodeVarintUint64(n uint64) []byte {
	maxLen := varintMax(8)
	buf := make([]byte, maxLen)
	if n < 128 {
		buf[0] = byte(n)
		return buf[:1]
	}
	value := n
	i := 0
	for value >= 128 {
		buf[i] = byte(value&0x7F) | 0x80
		value >>= 7
		i++
	}
	buf[i] = byte(value)
	return buf[:i+1]
}

func encodeVarintUint(n uint) []byte {
	if ^uint(0) == math.MaxUint64 {
		return encodeVarintUint64(uint64(n))
	}
	return encodeVarintUint32(uint32(n))
}

func decodeVarintUint16(data []byte, pos *int) (uint16, error) {
	var out uint16
	maxLen := varintMax(2)
	maxLast := maxOfLastByte(2)
	for i := 0; i < maxLen; i++ {
		if *pos >= len(data) {
			return 0, ErrDeserializeUnexpectedEnd
		}
		val := data[*pos]
		*pos++
		carry := uint16(val & 0x7F)
		out |= carry << (7 * i)
		if (val & 0x80) == 0 {
			if i == maxLen-1 && val > maxLast {
				return 0, ErrDeserializeBadVarint
			}
			return out, nil
		}
	}
	return 0, ErrDeserializeBadVarint
}

func decodeVarintUint32(data []byte, pos *int) (uint32, error) {
	var out uint32
	maxLen := varintMax(4)
	maxLast := maxOfLastByte(4)
	for i := 0; i < maxLen; i++ {
		if *pos >= len(data) {
			return 0, ErrDeserializeUnexpectedEnd
		}
		val := data[*pos]
		*pos++
		carry := uint32(val & 0x7F)
		out |= carry << (7 * i)
		if (val & 0x80) == 0 {
			if i == maxLen-1 && val > maxLast {
				return 0, ErrDeserializeBadVarint
			}
			return out, nil
		}
	}
	return 0, ErrDeserializeBadVarint
}

func decodeVarintUint64(data []byte, pos *int) (uint64, error) {
	var out uint64
	maxLen := varintMax(8)
	maxLast := maxOfLastByte(8)
	for i := 0; i < maxLen; i++ {
		if *pos >= len(data) {
			return 0, ErrDeserializeUnexpectedEnd
		}
		val := data[*pos]
		*pos++
		carry := uint64(val & 0x7F)
		out |= carry << (7 * i)
		if (val & 0x80) == 0 {
			if i == maxLen-1 && val > maxLast {
				return 0, ErrDeserializeBadVarint
			}
			return out, nil
		}
	}
	return 0, ErrDeserializeBadVarint
}

func decodeVarintUint(data []byte, pos *int) (uint, error) {
	if ^uint(0) == math.MaxUint64 {
		val, err := decodeVarintUint64(data, pos)
		return uint(val), err
	}
	val, err := decodeVarintUint32(data, pos)
	return uint(val), err
}

func encodeUint16LE(n uint16) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, n)
	return buf
}

func encodeUint32LE(n uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, n)
	return buf
}

func encodeUint64LE(n uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, n)
	return buf
}

func decodeUint16LE(data []byte, pos *int) (uint16, error) {
	if *pos+2 > len(data) {
		return 0, ErrDeserializeUnexpectedEnd
	}
	val := binary.LittleEndian.Uint16(data[*pos : *pos+2])
	*pos += 2
	return val, nil
}

func decodeUint32LE(data []byte, pos *int) (uint32, error) {
	if *pos+4 > len(data) {
		return 0, ErrDeserializeUnexpectedEnd
	}
	val := binary.LittleEndian.Uint32(data[*pos : *pos+4])
	*pos += 4
	return val, nil
}

func decodeUint64LE(data []byte, pos *int) (uint64, error) {
	if *pos+8 > len(data) {
		return 0, ErrDeserializeUnexpectedEnd
	}
	val := binary.LittleEndian.Uint64(data[*pos : *pos+8])
	*pos += 8
	return val, nil
}

func encodeFloat32LE(n float32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, math.Float32bits(n))
	return buf
}

func encodeFloat64LE(n float64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, math.Float64bits(n))
	return buf
}

func decodeFloat32LE(data []byte, pos *int) (float32, error) {
	if *pos+4 > len(data) {
		return 0, ErrDeserializeUnexpectedEnd
	}
	val := binary.LittleEndian.Uint32(data[*pos : *pos+4])
	*pos += 4
	return math.Float32frombits(val), nil
}

func decodeFloat64LE(data []byte, pos *int) (float64, error) {
	if *pos+8 > len(data) {
		return 0, ErrDeserializeUnexpectedEnd
	}
	val := binary.LittleEndian.Uint64(data[*pos : *pos+8])
	*pos += 8
	return math.Float64frombits(val), nil
}
