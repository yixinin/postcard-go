package postcard

import (
	"math"
	"reflect"
	"testing"
)

func TestVarintUint16(t *testing.T) {
	tests := []struct {
		name     string
		input    uint16
		expected []byte
	}{
		{"zero", 0, []byte{0x00}},
		{"127", 127, []byte{0x7F}},
		{"128", 128, []byte{0x80, 0x01}},
		{"16383", 16383, []byte{0xFF, 0x7F}},
		{"16384", 16384, []byte{0x80, 0x80, 0x01}},
		{"max", 65535, []byte{0xFF, 0xFF, 0x03}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := encodeVarintUint16(tt.input)
			if !reflect.DeepEqual(encoded, tt.expected) {
				t.Errorf("encodeVarintUint16(%d) = %v, want %v", tt.input, encoded, tt.expected)
			}

			pos := 0
			decoded, err := decodeVarintUint16(encoded, &pos)
			if err != nil {
				t.Errorf("decodeVarintUint16(%v) error = %v", encoded, err)
			}
			if decoded != tt.input {
				t.Errorf("decodeVarintUint16(%v) = %d, want %d", encoded, decoded, tt.input)
			}
		})
	}
}

func TestVarintUint32(t *testing.T) {
	tests := []struct {
		name     string
		input    uint32
		expected []byte
	}{
		{"zero", 0, []byte{0x00}},
		{"127", 127, []byte{0x7F}},
		{"128", 128, []byte{0x80, 0x01}},
		{"max", math.MaxUint32, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x0F}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := encodeVarintUint32(tt.input)
			if !reflect.DeepEqual(encoded, tt.expected) {
				t.Errorf("encodeVarintUint32(%d) = %v, want %v", tt.input, encoded, tt.expected)
			}

			pos := 0
			decoded, err := decodeVarintUint32(encoded, &pos)
			if err != nil {
				t.Errorf("decodeVarintUint32(%v) error = %v", encoded, err)
			}
			if decoded != tt.input {
				t.Errorf("decodeVarintUint32(%v) = %d, want %d", encoded, decoded, tt.input)
			}
		})
	}
}

func TestZigzagInt16(t *testing.T) {
	tests := []struct {
		input    int16
		expected uint16
	}{
		{0, 0},
		{-1, 1},
		{1, 2},
		{-2, 3},
		{2, 4},
		{32767, 65534},
		{-32768, 65535},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded := zigzagEncodeInt16(tt.input)
			if encoded != tt.expected {
				t.Errorf("zigzagEncodeInt16(%d) = %d, want %d", tt.input, encoded, tt.expected)
			}

			decoded := zigzagDecodeInt16(encoded)
			if decoded != tt.input {
				t.Errorf("zigzagDecodeInt16(%d) = %d, want %d", encoded, decoded, tt.input)
			}
		})
	}
}

func TestSerializeDeserializeBool(t *testing.T) {
	tests := []bool{true, false}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded, err := Serialize(tt)
			if err != nil {
				t.Fatalf("Serialize(%v) error = %v", tt, err)
			}

			var decoded bool
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if decoded != tt {
				t.Errorf("got %v, want %v", decoded, tt)
			}
		})
	}
}

func TestSerializeDeserializeUint8(t *testing.T) {
	tests := []uint8{0, 127, 255}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded, err := Serialize(tt)
			if err != nil {
				t.Fatalf("Serialize(%d) error = %v", tt, err)
			}

			var decoded uint8
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if decoded != tt {
				t.Errorf("got %d, want %d", decoded, tt)
			}
		})
	}
}

func TestSerializeDeserializeUint16(t *testing.T) {
	tests := []uint16{0, 127, 128, 16383, 16384, 65535}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded, err := Serialize(tt)
			if err != nil {
				t.Fatalf("Serialize(%d) error = %v", tt, err)
			}

			var decoded uint16
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if decoded != tt {
				t.Errorf("got %d, want %d", decoded, tt)
			}
		})
	}
}

func TestSerializeDeserializeInt16(t *testing.T) {
	tests := []int16{0, -1, 1, -64, 64, -32768, 32767}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded, err := Serialize(tt)
			if err != nil {
				t.Fatalf("Serialize(%d) error = %v", tt, err)
			}

			var decoded int16
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if decoded != tt {
				t.Errorf("got %d, want %d", decoded, tt)
			}
		})
	}
}

func TestSerializeDeserializeUint32(t *testing.T) {
	tests := []uint32{0, 127, 128, math.MaxUint32}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded, err := Serialize(tt)
			if err != nil {
				t.Fatalf("Serialize(%d) error = %v", tt, err)
			}

			var decoded uint32
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if decoded != tt {
				t.Errorf("got %d, want %d", decoded, tt)
			}
		})
	}
}

func TestSerializeDeserializeInt32(t *testing.T) {
	tests := []int32{0, -1, 1, math.MinInt32, math.MaxInt32}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded, err := Serialize(tt)
			if err != nil {
				t.Fatalf("Serialize(%d) error = %v", tt, err)
			}

			var decoded int32
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if decoded != tt {
				t.Errorf("got %d, want %d", decoded, tt)
			}
		})
	}
}

func TestSerializeDeserializeUint64(t *testing.T) {
	tests := []uint64{0, 127, 128, math.MaxUint64}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded, err := Serialize(tt)
			if err != nil {
				t.Fatalf("Serialize(%d) error = %v", tt, err)
			}

			var decoded uint64
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if decoded != tt {
				t.Errorf("got %d, want %d", decoded, tt)
			}
		})
	}
}

func TestSerializeDeserializeInt64(t *testing.T) {
	tests := []int64{0, -1, 1, math.MinInt64, math.MaxInt64}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded, err := Serialize(tt)
			if err != nil {
				t.Fatalf("Serialize(%d) error = %v", tt, err)
			}

			var decoded int64
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if decoded != tt {
				t.Errorf("got %d, want %d", decoded, tt)
			}
		})
	}
}

func TestSerializeDeserializeFloat32(t *testing.T) {
	tests := []float32{0.0, -0.0, 1.0, -1.0, 3.14159, -32.005859375}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded, err := Serialize(tt)
			if err != nil {
				t.Fatalf("Serialize(%f) error = %v", tt, err)
			}

			var decoded float32
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if decoded != tt {
				t.Errorf("got %f, want %f", decoded, tt)
			}
		})
	}
}

func TestSerializeDeserializeFloat64(t *testing.T) {
	tests := []float64{0.0, -0.0, 1.0, -1.0, 3.141592653589793, -32.005859375}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded, err := Serialize(tt)
			if err != nil {
				t.Fatalf("Serialize(%f) error = %v", tt, err)
			}

			var decoded float64
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if decoded != tt {
				t.Errorf("got %f, want %f", decoded, tt)
			}
		})
	}
}

func TestSerializeDeserializeString(t *testing.T) {
	tests := []string{"", "hello", "hello, postcard!", "你好世界"}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded, err := Serialize(tt)
			if err != nil {
				t.Fatalf("Serialize(%q) error = %v", tt, err)
			}

			var decoded string
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if decoded != tt {
				t.Errorf("got %q, want %q", decoded, tt)
			}
		})
	}
}

func TestSerializeDeserializeBytes(t *testing.T) {
	tests := [][]byte{
		{},
		{0x01, 0x02, 0x03},
		{0x01, 0x00, 0x20, 0x30},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded, err := Serialize(tt)
			if err != nil {
				t.Fatalf("Serialize(%v) error = %v", tt, err)
			}

			var decoded []byte
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if !reflect.DeepEqual(decoded, tt) {
				t.Errorf("got %v, want %v", decoded, tt)
			}
		})
	}
}

func TestSerializeDeserializeSlice(t *testing.T) {
	tests := [][]int{
		{},
		{1, 2, 3},
		{1, 2, 3, 4, 5},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded, err := Serialize(tt)
			if err != nil {
				t.Fatalf("Serialize(%v) error = %v", tt, err)
			}

			var decoded []int
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if !reflect.DeepEqual(decoded, tt) {
				t.Errorf("got %v, want %v", decoded, tt)
			}
		})
	}
}

func TestSerializeDeserializeStruct(t *testing.T) {
	type BasicStruct struct {
		A uint16
		B uint8
		C uint64
		D uint32
	}

	tests := []BasicStruct{
		{0xABCD, 0xFE, 0x1234_4321_ABCD_DCBA, 0xACAC_ACAC},
		{0, 0, 0, 0},
		{65535, 255, math.MaxUint64, math.MaxUint32},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded, err := Serialize(tt)
			if err != nil {
				t.Fatalf("Serialize(%v) error = %v", tt, err)
			}

			var decoded BasicStruct
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if !reflect.DeepEqual(decoded, tt) {
				t.Errorf("got %v, want %v", decoded, tt)
			}
		})
	}
}

func TestSerializeDeserializeMap(t *testing.T) {
	tests := []map[string]int{
		{},
		{"a": 1, "b": 2},
		{"hello": 42, "world": 123},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			encoded, err := Serialize(tt)
			if err != nil {
				t.Fatalf("Serialize(%v) error = %v", tt, err)
			}

			var decoded map[string]int
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if !reflect.DeepEqual(decoded, tt) {
				t.Errorf("got %v, want %v", decoded, tt)
			}
		})
	}
}

func TestSerializeDeserializePointer(t *testing.T) {
	tests := []struct {
		name  string
		input *int
	}{
		{"nil", nil},
		{"value", func() *int { v := 42; return &v }()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := Serialize(tt.input)
			if err != nil {
				t.Fatalf("Serialize(%v) error = %v", tt.input, err)
			}

			var decoded int
			err = Deserialize(encoded, &decoded)
			if err != nil {
				t.Fatalf("Deserialize(%v) error = %v", encoded, err)
			}

			if tt.input == nil {
				if decoded != 0 {
					t.Errorf("got %d, want 0", decoded)
				}
			} else {
				if decoded != *tt.input {
					t.Errorf("got %d, want %d", decoded, *tt.input)
				}
			}
		})
	}
}

func TestVarintBoundaryCanon(t *testing.T) {
	x := uint32(math.MaxUint32)
	encoded := encodeVarintUint32(x)
	expected := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x0F}
	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("encodeVarintUint32(%d) = %v, want %v", x, encoded, expected)
	}

	pos := 0
	decoded, err := decodeVarintUint32(encoded, &pos)
	if err != nil {
		t.Fatalf("decodeVarintUint32(%v) error = %v", encoded, err)
	}
	if decoded != x {
		t.Errorf("decodeVarintUint32(%v) = %d, want %d", encoded, decoded, x)
	}

	badEncoded := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x1F}
	pos = 0
	_, err = decodeVarintUint32(badEncoded, &pos)
	if err != ErrDeserializeBadVarint {
		t.Errorf("decodeVarintUint32(%v) error = %v, want %v", badEncoded, err, ErrDeserializeBadVarint)
	}
}
