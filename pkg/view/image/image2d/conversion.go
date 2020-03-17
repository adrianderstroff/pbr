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

func float32ToBytes(fl float32) []byte {
	bits := math.Float32bits(fl)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func bytesToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func bytesToUint8(bytes []byte) []uint8 {
	return bytes
}
