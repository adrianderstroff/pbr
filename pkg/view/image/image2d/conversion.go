package image2d

import (
	"encoding/binary"
	"math"
)

func float32SliceToUint8Slice(fls []float32) []uint8 {
	var bytes []byte
	for _, fl := range fls {
		bytes = append(bytes, float32ToBytes(fl)...)
	}
	return bytesToUint8(bytes)
}

func float64SliceToUint8Slice(fls []float64) []uint8 {
	var bytes []byte
	for _, fl := range fls {
		bytes = append(bytes, float64ToBytes(fl)...)
	}
	return bytesToUint8(bytes)
}

func float32ToBytes(fl float32) []byte {
	bits := math.Float32bits(fl)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func float64ToBytes(fl float64) []byte {
	bits := math.Float64bits(fl)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func bytesToUint8(bytes []byte) []uint8 {
	return bytes
}
