package postcard

import "errors"

var (
	ErrWontImplement             = errors.New("this is a feature that postcard will never implement")
	ErrNotYetImplemented         = errors.New("this is a feature that postcard intends to support, but does not yet")
	ErrSerializeBufferFull       = errors.New("the serialize buffer is full")
	ErrSerializeSeqLengthUnknown = errors.New("the length of a sequence must be known")
	ErrDeserializeUnexpectedEnd  = errors.New("hit the end of buffer, expected more data")
	ErrDeserializeBadVarint      = errors.New("found a varint that didn't terminate")
	ErrDeserializeBadBool        = errors.New("found a bool that wasn't 0 or 1")
	ErrDeserializeBadChar        = errors.New("found an invalid unicode char")
	ErrDeserializeBadUtf8        = errors.New("tried to parse invalid utf-8")
	ErrDeserializeBadOption      = errors.New("found an Option discriminant that wasn't 0 or 1")
	ErrDeserializeBadEnum        = errors.New("found an enum discriminant that was > u32::max_value()")
	ErrDeserializeBadEncoding    = errors.New("the original data was not well encoded")
	ErrDeserializeBadCrc         = errors.New("bad CRC while deserializing")
	ErrSerdeSerCustom            = errors.New("serde serialization error")
	ErrSerdeDeCustom             = errors.New("serde deserialization error")
	ErrCollectStrError           = errors.New("error while processing collect_str during serialization")
)
