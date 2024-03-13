package img

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)



// SubImage 裁剪图片
func SubImage(imgByte []byte) (subImgList [][]byte, err error) {
	reader := bytes.NewReader(imgByte)

	img, _, err := image.Decode(reader)
	if err != nil {
		return
	}

	var (
		// 获取图片的宽度和高度
		bounds = img.Bounds()
		width  = bounds.Max.X
		height = bounds.Max.Y

		// 保存剪切后的四个图片
		topLeftBuf     = new(bytes.Buffer)
		topRightBuf    = new(bytes.Buffer)
		bottomLeftBuf  = new(bytes.Buffer)
		bottomRightBuf = new(bytes.Buffer)
	)

	// 剪切四个角落的图片
	topLeft := image.NewRGBA(image.Rect(0, 0, width/2, height/2))
	draw.Draw(topLeft, topLeft.Bounds(), img, image.Point{}, draw.Src)
	if err = png.Encode(topLeftBuf, topLeft); err != nil {
		return
	}
	subImgList = append(subImgList, topLeftBuf.Bytes())

	topRight := image.NewRGBA(image.Rect(0, 0, width/2, height/2))
	draw.Draw(topRight, topRight.Bounds(), img, image.Point{X: width / 2}, draw.Src)
	if err = png.Encode(topRightBuf, topRight); err != nil {
		return
	}
	subImgList = append(subImgList, topRightBuf.Bytes())

	bottomLeft := image.NewRGBA(image.Rect(0, 0, width/2, height/2))
	draw.Draw(bottomLeft, bottomLeft.Bounds(), img, image.Point{Y: height / 2}, draw.Src)
	if err = png.Encode(bottomLeftBuf, bottomLeft); err != nil {
		return
	}
	subImgList = append(subImgList, bottomLeftBuf.Bytes())

	bottomRight := image.NewRGBA(image.Rect(0, 0, width/2, height/2))
	draw.Draw(bottomRight, bottomRight.Bounds(), img, image.Point{X: width / 2, Y: height / 2}, draw.Src)
	if err = png.Encode(bottomRightBuf, bottomRight); err != nil {
		return
	}
	subImgList = append(subImgList, bottomRightBuf.Bytes())
	return
}

// SubImageByGap 根据间隔裁剪图片
func SubImageByGap(imgByte []byte) (subImgList [][]byte, err error) {
	reader := bytes.NewReader(imgByte)

	imgBg, _, err := image.Decode(reader)
	if err != nil {
		return
	}

	var (
		// 获取图片的宽度和高度
		bounds   = imgBg.Bounds()
		bgWidth  = bounds.Max.X
		bgHeight = bounds.Max.Y
	)

	// 识别小图的水平间隔
	var (
		lastGapY  int
		rawImages []*image.RGBA // 裁剪出的小图
	)
	for y := 0; y < bgHeight; y++ { // 遍历图片的高度
		gapWidth := 0                  // 间隔宽度
		for x := 0; x < bgWidth; x++ { // 遍历图片的宽度
			at := imgBg.At(x, y) // 获取像素点的颜色
			if !isWhite(at) {    // 判断是否是白色
				break
			}
			gapWidth++ // 间隔宽度+1
		}
		if gapWidth == bgWidth { // 间隔宽度等于图片宽度，说明是白色间隔
			if y-lastGapY <= 1 { // 排除连续的白色间隔
				lastGapY = y
				continue
			}
			// 裁剪出小图
			smallImg := image.NewRGBA(image.Rect(0, 0, bgWidth, y-lastGapY))
			draw.Draw(smallImg, smallImg.Bounds(), imgBg, image.Point{Y: lastGapY}, draw.Src)
			rawImages = append(rawImages, smallImg)
			lastGapY = y
		}
	}

	var smallImages []*image.RGBA   // 裁剪出的小图
	for _, img := range rawImages { // 遍历竖行图
		var lastGapX int
		// 重新获取图片的宽度和高度
		bgWidth = img.Bounds().Max.X
		bgHeight = img.Bounds().Max.Y
		for x := 0; x < bgWidth; x++ { // 遍历图片的宽度
			gapHeight := 0                  // 间隔高度
			for y := 0; y < bgHeight; y++ { // 遍历图片的高度
				at := img.At(x, y) // 获取像素点的颜色
				if !isWhite(at) {  // 判断是否是白色
					break
				}
				gapHeight++ // 间隔高度+1
			}
			if gapHeight == bgHeight { // 间隔高度等于图片高度，说明是白色间隔
				if x-lastGapX == 1 { // 排除连续的白色间隔
					lastGapX = x
					continue
				}
				// 裁剪出竖行图
				smallImg := image.NewRGBA(image.Rect(0, 0, x-lastGapX, bgHeight))
				draw.Draw(smallImg, smallImg.Bounds(), img, image.Point{X: lastGapX}, draw.Src)
				smallImages = append(smallImages, smallImg)
				lastGapX = x
			}
		}
		// 添加最后一个列图
		if lastGapX < bgWidth {
			smallImg := image.NewRGBA(image.Rect(0, 0, bgWidth-lastGapX, bgHeight))
			draw.Draw(smallImg, smallImg.Bounds(), img, image.Point{X: lastGapX}, draw.Src)
			smallImages = append(smallImages, smallImg)
		}
	}
	// 保存裁剪出的小图
	for _, img := range smallImages {
		if fullWhiteImg(img) {
			continue
		}
		imgBuf := new(bytes.Buffer)
		if err = png.Encode(imgBuf, img); err != nil {
			continue
		}
		subImgList = append(subImgList, imgBuf.Bytes())
	}
	return
}

// fullWhiteImg 判断是否是全白图片
// img 图片
// is 是否是全白图片: true 是, false 否
func fullWhiteImg(img *image.RGBA) (is bool) {
	bounds := img.Bounds()
	bgWidth := bounds.Max.X
	bgHeight := bounds.Max.Y
	for y := 0; y < bgHeight; y++ { // 遍历图片的高度
		for x := 0; x < bgWidth; x++ { // 遍历图片的宽度
			at := img.At(x, y) // 获取像素点的颜色
			if !isWhite(at) {  // 判断是否是白色
				return
			}
		}
	}
	is = true
	return
}

// isWhite 判断是否是白色
func isWhite(c color.Color) bool {
	maxColor := color.RGBA{R: 255, G: 255, B: 255, A: 255} // 白色的最大值
	minColor := color.RGBA{R: 240, G: 240, B: 240, A: 240} // 白色的最小值

	r, g, b, a := c.RGBA()
	var rBool, gBool, bBool, aBool bool
	//fmt.Printf("byte(r): %d, byte(g): %d, byte(b): %d, byte(a): %d\n", byte(r), byte(g), byte(b), byte(a))
	if byte(r) <= maxColor.R && byte(r) >= minColor.R {
		rBool = true
	}
	if byte(g) <= maxColor.G && byte(g) >= minColor.G {
		gBool = true
	}
	if byte(b) <= maxColor.B && byte(b) >= minColor.B {
		bBool = true
	}
	if byte(a) <= maxColor.A && byte(a) >= minColor.A {
		aBool = true
	}
	return rBool && gBool && bBool && aBool
}
