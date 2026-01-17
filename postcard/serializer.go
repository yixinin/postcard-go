package postcard

import (
	"fmt"
	"reflect"
	"unicode/utf8"
)

type Serializer struct {
	buf []byte
}

func NewSerializer(buf []byte) *Serializer {
	return &Serializer{buf: buf}
}

func (s *Serializer) Result() ([]byte, error) {
	return s.buf, nil
}

func (s *Serializer) pushByte(b byte) error {
	s.buf = append(s.buf, b)
	return nil
}

func (s *Serializer) pushBytes(data []byte) error {
	s.buf = append(s.buf, data...)
	return nil
}

func (s *Serializer) pushVarintUint16(n uint16) error {
	s.buf = append(s.buf, encodeVarintUint16(n)...)
	return nil
}

func (s *Serializer) pushVarintUint32(n uint32) error {
	s.buf = append(s.buf, encodeVarintUint32(n)...)
	return nil
}

func (s *Serializer) pushVarintUint64(n uint64) error {
	s.buf = append(s.buf, encodeVarintUint64(n)...)
	return nil
}

func (s *Serializer) pushVarintUint(n uint) error {
	s.buf = append(s.buf, encodeVarintUint(n)...)
	return nil
}

func (s *Serializer) SerializeBool(v bool) error {
	if v {
		return s.pushByte(1)
	}
	return s.pushByte(0)
}

func (s *Serializer) SerializeInt8(v int8) error {
	return s.pushByte(uint8(v))
}

func (s *Serializer) SerializeInt16(v int16) error {
	zzv := zigzagEncodeInt16(v)
	return s.pushVarintUint16(zzv)
}

func (s *Serializer) SerializeInt32(v int32) error {
	zzv := zigzagEncodeInt32(v)
	return s.pushVarintUint32(zzv)
}

func (s *Serializer) SerializeInt64(v int64) error {
	zzv := zigzagEncodeInt64(v)
	return s.pushVarintUint64(zzv)
}

func (s *Serializer) SerializeInt(v int) error {
	zzv := zigzagEncodeInt(v)
	return s.pushVarintUint(zzv)
}

func (s *Serializer) SerializeUint8(v uint8) error {
	return s.pushByte(v)
}

func (s *Serializer) SerializeUint16(v uint16) error {
	return s.pushVarintUint16(v)
}

func (s *Serializer) SerializeUint32(v uint32) error {
	return s.pushVarintUint32(v)
}

func (s *Serializer) SerializeUint64(v uint64) error {
	return s.pushVarintUint64(v)
}

func (s *Serializer) SerializeUint(v uint) error {
	return s.pushVarintUint(v)
}

func (s *Serializer) SerializeVarInt(v Varint) error {
	return s.pushBytes(v.Encode())
}

func (s *Serializer) SerializeFloat32(v float32) error {
	s.buf = append(s.buf, encodeFloat32LE(v)...)
	return nil
}

func (s *Serializer) SerializeFloat64(v float64) error {
	s.buf = append(s.buf, encodeFloat64LE(v)...)
	return nil
}

func (s *Serializer) SerializeString(v string) error {
	if err := s.pushVarintUint(uint(len(v))); err != nil {
		return err
	}
	return s.pushBytes([]byte(v))
}

func (s *Serializer) SerializeBytes(v []byte) error {
	if err := s.pushVarintUint(uint(len(v))); err != nil {
		return err
	}
	return s.pushBytes(v)
}

func (s *Serializer) SerializeOption(v interface{}) error {
	if v == nil {
		return s.pushByte(0)
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return s.pushByte(0)
		}
		if err := s.pushByte(1); err != nil {
			return err
		}
		return s.SerializeValue(rv.Elem().Interface())
	}

	if err := s.pushByte(1); err != nil {
		return err
	}
	return s.SerializeValue(v)
}

func (s *Serializer) SerializeSlice(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("expected slice, got %v", val.Kind())
	}

	if err := s.pushVarintUint(uint(val.Len())); err != nil {
		return err
	}

	for i := 0; i < val.Len(); i++ {
		if err := s.SerializeValue(val.Index(i).Interface()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Serializer) SerializeArray(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Array {
		return fmt.Errorf("expected array, got %v", val.Kind())
	}

	for i := 0; i < val.Len(); i++ {
		if err := s.SerializeValue(val.Index(i).Interface()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Serializer) SerializeMap(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Map {
		return fmt.Errorf("expected map, got %v", val.Kind())
	}

	if err := s.pushVarintUint(uint(val.Len())); err != nil {
		return err
	}

	keys := val.MapKeys()
	for _, key := range keys {
		if err := s.SerializeValue(key.Interface()); err != nil {
			return err
		}
		value := val.MapIndex(key)
		if err := s.SerializeValue(value.Interface()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Serializer) SerializeStruct(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %v", val.Kind())
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}
		if err := s.SerializeValue(val.Field(i).Interface()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Serializer) SerializeEnum(variantIndex uint32, value interface{}) error {
	if err := s.pushVarintUint32(variantIndex); err != nil {
		return err
	}
	if value != nil {
		return s.SerializeValue(value)
	}
	return nil
}

func (s *Serializer) SerializeValue(v interface{}) error {
	if v == nil {
		return s.SerializeOption(nil)
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Bool:
		return s.SerializeBool(val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch val.Kind() {
		case reflect.Int8:
			return s.SerializeInt8(int8(val.Int()))
		case reflect.Int16:
			return s.SerializeInt16(int16(val.Int()))
		case reflect.Int32:
			return s.SerializeInt32(int32(val.Int()))
		case reflect.Int64:
			return s.SerializeInt64(val.Int())
		default:
			return s.SerializeInt(int(val.Int()))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch val.Kind() {
		case reflect.Uint8:
			return s.SerializeUint8(uint8(val.Uint()))
		case reflect.Uint16:
			return s.SerializeUint16(uint16(val.Uint()))
		case reflect.Uint32:
			return s.SerializeUint32(uint32(val.Uint()))
		case reflect.Uint64:
			if val.Type() == reflect.TypeOf(Varint(0)) {
				return s.SerializeVarInt(Varint(val.Uint()))
			}
			return s.SerializeUint64(val.Uint())
		default:
			return s.SerializeUint(uint(val.Uint()))
		}
	case reflect.Float32:
		return s.SerializeFloat32(float32(val.Float()))
	case reflect.Float64:
		return s.SerializeFloat64(val.Float())
	case reflect.String:
		return s.SerializeString(val.String())
	case reflect.Slice:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			return s.SerializeBytes(val.Bytes())
		}
		return s.SerializeSlice(v)
	case reflect.Array:
		return s.SerializeArray(v)
	case reflect.Map:
		return s.SerializeMap(v)
	case reflect.Struct:
		return s.SerializeStruct(v)
	case reflect.Ptr:
		if val.IsNil() {
			return s.SerializeOption(nil)
		}
		return s.SerializeValue(val.Elem().Interface())
	default:
		return fmt.Errorf("unsupported type: %v", val.Kind())
	}
}

func Serialize(v interface{}) ([]byte, error) {
	s := NewSerializer(nil)
	if err := s.SerializeValue(v); err != nil {
		return nil, err
	}
	return s.Result()
}

func SerializeToSlice(v interface{}, buf []byte) ([]byte, error) {
	s := NewSerializer(buf[:0])
	if err := s.SerializeValue(v); err != nil {
		return nil, err
	}
	return s.Result()
}

func SerializeString(v string) ([]byte, error) {
	s := NewSerializer(nil)
	if err := s.SerializeString(v); err != nil {
		return nil, err
	}
	return s.Result()
}

func SerializeBytes(v []byte) ([]byte, error) {
	s := NewSerializer(nil)
	if err := s.SerializeBytes(v); err != nil {
		return nil, err
	}
	return s.Result()
}

func SerializeRune(r rune) ([]byte, error) {
	buf := make([]byte, 4)
	n := utf8.EncodeRune(buf, r)
	return SerializeBytes(buf[:n])
}
