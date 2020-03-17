package image2d

// GetR returns the red value of the pixel at (x,y).
func (img *Image2D) GetR(x, y int) uint8 {
	idx := img.getIdx(x, y)
	return img.data[idx]
}

// GetG returns the green value of the pixel at (x,y).
func (img *Image2D) GetG(x, y int) uint8 {
	idx := img.getIdx(x, y)
	return img.data[idx+1]
}

// GetB returns the blue value of the pixel at (x,y).
func (img *Image2D) GetB(x, y int) uint8 {
	idx := img.getIdx(x, y)
	return img.data[idx+2]
}

// GetA returns the alpha value of the pixel at (x,y).
func (img *Image2D) GetA(x, y int) uint8 {
	idx := img.getIdx(x, y)
	return img.data[idx+3]
}

// GetRGB returns the RGB values of the pixel at (x,y).
func (img *Image2D) GetRGB(x, y int) (uint8, uint8, uint8) {
	idx := img.getIdx(x, y)
	return img.data[idx],
		img.data[idx+1],
		img.data[idx+2]
}

// GetRGBA returns the RGBA value of the pixel at (x,y).
func (img *Image2D) GetRGBA(x, y int) (uint8, uint8, uint8, uint8) {
	idx := img.getIdx(x, y)
	return img.data[idx],
		img.data[idx+1],
		img.data[idx+2],
		img.data[idx+3]
}

// SetR sets the red value of the pixel at (x,y).
func (img *Image2D) SetR(x, y int, r uint8) {
	idx := img.getIdx(x, y)
	img.data[idx] = r
}

// SetG sets the green value of the pixel at (x,y).
func (img *Image2D) SetG(x, y int, g uint8) {
	idx := img.getIdx(x, y)
	img.data[idx+1] = g
}

// SetB sets the blue value of the pixel at (x,y).
func (img *Image2D) SetB(x, y int, b uint8) {
	idx := img.getIdx(x, y)
	img.data[idx+2] = b
}

// SetA sets the alpha value of the pixel at (x,y).
func (img *Image2D) SetA(x, y int, a uint8) {
	idx := img.getIdx(x, y)
	img.data[idx+3] = a
}

// SetRGB sets the RGB values of the pixel at (x,y).
func (img *Image2D) SetRGB(x, y int, r, g, b uint8) {
	idx := img.getIdx(x, y)
	img.data[idx] = r
	img.data[idx+1] = g
	img.data[idx+2] = b
}

// SetRGBA sets the RGBA values of the pixel at (x,y).
func (img *Image2D) SetRGBA(x, y int, r, g, b, a uint8) {
	idx := img.getIdx(x, y)
	img.data[idx] = r
	img.data[idx+1] = g
	img.data[idx+2] = b
	img.data[idx+3] = a
}
