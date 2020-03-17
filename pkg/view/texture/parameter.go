package texture

import gl "github.com/adrianderstroff/pbr/pkg/core/gl"

// InternalFormat specifies the format of a pixel on the GPU.
type InternalFormat int32

// Format specifies the format of a pixel of a pixel of the image data.
type Format uint32

// PixelType specifies the data type of a single component of the image data.
type PixelType uint32

// List of supported internal formats.
const (
	INTERNALFORMAT_RED  InternalFormat = gl.RED
	INTERNALFORMAT_RG   InternalFormat = gl.RG
	INTERNALFORMAT_RGB  InternalFormat = gl.RGB
	INTERNALFORMAT_RGBA InternalFormat = gl.RGBA
)

// List of supported formats.
const (
	FORMAT_RED  Format = gl.RED
	FORMAT_RG   Format = gl.RG
	FORMAT_RGB  Format = gl.RGB
	FORMAT_RGBA Format = gl.RGBA
)

// List of supported pixel types.
const (
	PIXELTYPE_UNSIGNED_BYTE PixelType = gl.UNSIGNED_BYTE
	PIXELTYPE_UNSIGNED_INT  PixelType = gl.UNSIGNED_INT
	PIXELTYPE_FLOAT         PixelType = gl.FLOAT
)
