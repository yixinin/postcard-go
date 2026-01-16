package postcard

import (
	"fmt"
	"reflect"
	"unicode/utf8"
)

type Deserializer struct {
	data []byte
	pos  int
}

func NewDeserializer(data []byte) *Deserializer {
	return &Deserializer{data: data, pos: 0}
}

func (d *Deserializer) popByte() (byte, error) {
	if d.pos >= len(d.data) {
		return 0, ErrDeserializeUnexpectedEnd
	}
	b := d.data[d.pos]
	d.pos++
	return b, nil
}

func (d *Deserializer) takeBytes(n int) ([]byte, error) {
	if d.pos+n > len(d.data) {
		return nil, ErrDeserializeUnexpectedEnd
	}
	result := d.data[d.pos : d.pos+n]
	d.pos += n
	return result, nil
}

func (d *Deserializer) DeserializeBool() (bool, error) {
	b, err := d.popByte()
	if err != nil {
		return false, err
	}
	switch b {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, ErrDeserializeBadBool
	}
}

func (d *Deserializer) DeserializeInt8() (int8, error) {
	b, err := d.popByte()
	if err != nil {
		return 0, err
	}
	return int8(b), nil
}

func (d *Deserializer) DeserializeInt16() (int16, error) {
	v, err := decodeVarintUint16(d.data, &d.pos)
	if err != nil {
		return 0, err
	}
	return zigzagDecodeInt16(v), nil
}

func (d *Deserializer) DeserializeInt32() (int32, error) {
	v, err := decodeVarintUint32(d.data, &d.pos)
	if err != nil {
		return 0, err
	}
	return zigzagDecodeInt32(v), nil
}

func (d *Deserializer) DeserializeInt64() (int64, error) {
	v, err := decodeVarintUint64(d.data, &d.pos)
	if err != nil {
		return 0, err
	}
	return zigzagDecodeInt64(v), nil
}

func (d *Deserializer) DeserializeInt() (int, error) {
	v, err := decodeVarintUint(d.data, &d.pos)
	if err != nil {
		return 0, err
	}
	return zigzagDecodeInt(v), nil
}

func (d *Deserializer) DeserializeUint8() (uint8, error) {
	b, err := d.popByte()
	if err != nil {
		return 0, err
	}
	return b, nil
}

func (d *Deserializer) DeserializeUint16() (uint16, error) {
	return decodeVarintUint16(d.data, &d.pos)
}

func (d *Deserializer) DeserializeUint32() (uint32, error) {
	return decodeVarintUint32(d.data, &d.pos)
}

func (d *Deserializer) DeserializeUint64() (uint64, error) {
	return decodeVarintUint64(d.data, &d.pos)
}

func (d *Deserializer) DeserializeUint() (uint, error) {
	return decodeVarintUint(d.data, &d.pos)
}

func (d *Deserializer) DeserializeFloat32() (float32, error) {
	return decodeFloat32LE(d.data, &d.pos)
}

func (d *Deserializer) DeserializeFloat64() (float64, error) {
	return decodeFloat64LE(d.data, &d.pos)
}

func (d *Deserializer) DeserializeString() (string, error) {
	sz, err := d.DeserializeUint()
	if err != nil {
		return "", err
	}
	if sz > 4 {
		bytes, err := d.takeBytes(int(sz))
		if err != nil {
			return "", err
		}
		if !utf8.Valid(bytes) {
			return "", ErrDeserializeBadUtf8
		}
		return string(bytes), nil
	}
	bytes, err := d.takeBytes(int(sz))
	if err != nil {
		return "", err
	}
	if !utf8.Valid(bytes) {
		return "", ErrDeserializeBadUtf8
	}
	return string(bytes), nil
}

func (d *Deserializer) DeserializeBytes() ([]byte, error) {
	sz, err := d.DeserializeUint()
	if err != nil {
		return nil, err
	}
	return d.takeBytes(int(sz))
}

func (d *Deserializer) DeserializeOption(v interface{}) error {
	b, err := d.popByte()
	if err != nil {
		return err
	}
	switch b {
	case 0:
		if v != nil {
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Ptr {
				rv.Elem().SetZero()
			}
		}
		return nil
	case 1:
		if v == nil {
			return fmt.Errorf("cannot deserialize into nil")
		}
		return d.DeserializeValue(v)
	default:
		return ErrDeserializeBadOption
	}
}

func (d *Deserializer) DeserializeSlice(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("expected pointer to slice, got %T", v)
	}
	
	slice := rv.Elem()
	
	sz, err := d.DeserializeUint()
	if err != nil {
		return err
	}
	
	if slice.IsNil() || slice.Cap() < int(sz) {
		slice.Set(reflect.MakeSlice(slice.Type(), int(sz), int(sz)))
	} else {
		slice.SetLen(int(sz))
	}
	
	for i := 0; i < int(sz); i++ {
		elem := slice.Index(i)
		if err := d.DeserializeValue(elem.Addr().Interface()); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deserializer) DeserializeArray(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Array {
		return fmt.Errorf("expected pointer to array, got %T", v)
	}
	
	arr := rv.Elem()
	
	for i := 0; i < arr.Len(); i++ {
		elem := arr.Index(i)
		if err := d.DeserializeValue(elem.Addr().Interface()); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deserializer) DeserializeMap(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Map {
		return fmt.Errorf("expected pointer to map, got %T", v)
	}
	
	m := rv.Elem()
	
	if m.IsNil() {
		m.Set(reflect.MakeMap(m.Type()))
	}
	
	sz, err := d.DeserializeUint()
	if err != nil {
		return err
	}
	
	keyType := m.Type().Key()
	elemType := m.Type().Elem()
	
	for i := 0; i < int(sz); i++ {
		key := reflect.New(keyType).Elem()
		if err := d.DeserializeValue(key.Addr().Interface()); err != nil {
			return err
		}
		value := reflect.New(elemType).Elem()
		if err := d.DeserializeValue(value.Addr().Interface()); err != nil {
			return err
		}
		m.SetMapIndex(key, value)
	}
	return nil
}

func (d *Deserializer) DeserializeStruct(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer, got %T", v)
	}
	
	val := rv.Elem()
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		val = val.Elem()
	}
	
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %T", v)
	}
	
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}
		fieldVal := val.Field(i)
		if fieldVal.CanAddr() {
			if err := d.DeserializeValue(fieldVal.Addr().Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *Deserializer) DeserializeEnum(variantIndex *uint32, value interface{}) error {
	v, err := d.DeserializeUint32()
	if err != nil {
		return err
	}
	if variantIndex != nil {
		*variantIndex = v
	}
	if value != nil {
		return d.DeserializeValue(value)
	}
	return nil
}

func (d *Deserializer) DeserializeValue(v interface{}) error {
	if v == nil {
		return d.DeserializeOption(nil)
	}
	
	rv := reflect.ValueOf(v)
	
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer, got %T", v)
	}
	
	if rv.IsNil() {
		rv.Set(reflect.New(rv.Type().Elem()))
	}
	
	val := rv.Elem()
	
	switch val.Kind() {
	case reflect.Bool:
		decoded, err := d.DeserializeBool()
		if err != nil {
			return err
		}
		val.SetBool(decoded)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch val.Kind() {
		case reflect.Int8:
			decoded, err := d.DeserializeInt8()
			if err != nil {
				return err
			}
			val.SetInt(int64(decoded))
		case reflect.Int16:
			decoded, err := d.DeserializeInt16()
			if err != nil {
				return err
			}
			val.SetInt(int64(decoded))
		case reflect.Int32:
			decoded, err := d.DeserializeInt32()
			if err != nil {
				return err
			}
			val.SetInt(int64(decoded))
		case reflect.Int64:
			decoded, err := d.DeserializeInt64()
			if err != nil {
				return err
			}
			val.SetInt(decoded)
		default:
			decoded, err := d.DeserializeInt()
			if err != nil {
				return err
			}
			val.SetInt(int64(decoded))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch val.Kind() {
		case reflect.Uint8:
			decoded, err := d.DeserializeUint8()
			if err != nil {
				return err
			}
			val.SetUint(uint64(decoded))
		case reflect.Uint16:
			decoded, err := d.DeserializeUint16()
			if err != nil {
				return err
			}
			val.SetUint(uint64(decoded))
		case reflect.Uint32:
			decoded, err := d.DeserializeUint32()
			if err != nil {
				return err
			}
			val.SetUint(uint64(decoded))
		case reflect.Uint64:
			decoded, err := d.DeserializeUint64()
			if err != nil {
				return err
			}
			val.SetUint(decoded)
		default:
			decoded, err := d.DeserializeUint()
			if err != nil {
				return err
			}
			val.SetUint(uint64(decoded))
		}
	case reflect.Float32:
		decoded, err := d.DeserializeFloat32()
		if err != nil {
			return err
		}
		val.SetFloat(float64(decoded))
	case reflect.Float64:
		decoded, err := d.DeserializeFloat64()
		if err != nil {
			return err
		}
		val.SetFloat(decoded)
	case reflect.String:
		decoded, err := d.DeserializeString()
		if err != nil {
			return err
		}
		val.SetString(decoded)
	case reflect.Slice:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			decoded, err := d.DeserializeBytes()
			if err != nil {
				return err
			}
			val.SetBytes(decoded)
		} else {
			return d.DeserializeSlice(v)
		}
	case reflect.Array:
		return d.DeserializeArray(v)
	case reflect.Map:
		return d.DeserializeMap(v)
	case reflect.Struct:
		return d.DeserializeStruct(v)
	default:
		return fmt.Errorf("unsupported type: %v", val.Kind())
	}
	return nil
}

func Deserialize(data []byte, v interface{}) error {
	d := NewDeserializer(data)
	return d.DeserializeValue(v)
}

func DeserializeBool(data []byte) (bool, error) {
	d := NewDeserializer(data)
	return d.DeserializeBool()
}

func DeserializeInt8(data []byte) (int8, error) {
	d := NewDeserializer(data)
	return d.DeserializeInt8()
}

func DeserializeInt16(data []byte) (int16, error) {
	d := NewDeserializer(data)
	return d.DeserializeInt16()
}

func DeserializeInt32(data []byte) (int32, error) {
	d := NewDeserializer(data)
	return d.DeserializeInt32()
}

func DeserializeInt64(data []byte) (int64, error) {
	d := NewDeserializer(data)
	return d.DeserializeInt64()
}

func DeserializeInt(data []byte) (int, error) {
	d := NewDeserializer(data)
	return d.DeserializeInt()
}

func DeserializeUint8(data []byte) (uint8, error) {
	d := NewDeserializer(data)
	return d.DeserializeUint8()
}

func DeserializeUint16(data []byte) (uint16, error) {
	d := NewDeserializer(data)
	return d.DeserializeUint16()
}

func DeserializeUint32(data []byte) (uint32, error) {
	d := NewDeserializer(data)
	return d.DeserializeUint32()
}

func DeserializeUint64(data []byte) (uint64, error) {
	d := NewDeserializer(data)
	return d.DeserializeUint64()
}

func DeserializeUint(data []byte) (uint, error) {
	d := NewDeserializer(data)
	return d.DeserializeUint()
}

func DeserializeFloat32(data []byte) (float32, error) {
	d := NewDeserializer(data)
	return d.DeserializeFloat32()
}

func DeserializeFloat64(data []byte) (float64, error) {
	d := NewDeserializer(data)
	return d.DeserializeFloat64()
}

func DeserializeString(data []byte) (string, error) {
	d := NewDeserializer(data)
	return d.DeserializeString()
}

func DeserializeBytes(data []byte) ([]byte, error) {
	d := NewDeserializer(data)
	return d.DeserializeBytes()
}

func DeserializeRune(data []byte) (rune, error) {
	d := NewDeserializer(data)
	sz, err := d.DeserializeUint()
	if err != nil {
		return 0, err
	}
	if sz > 4 {
		return 0, ErrDeserializeBadChar
	}
	bytes, err := d.takeBytes(int(sz))
	if err != nil {
		return 0, err
	}
	r, _ := utf8.DecodeRune(bytes)
	if r == utf8.RuneError {
		return 0, ErrDeserializeBadChar
	}
	return r, nil
}
