/*
分别获取两张照片的像素点的值,转成uint8
将显示的图片的个位数去掉,然后将隐藏图片的的uint8值根据255/9的比例缩小,此时,隐藏图片的值一定是小于10的
将去除了显示图片的个位数的值和刚刚隐藏图片的值相加,这样影藏图片就成功影藏到了显示图片里
最后生成新的图片
因为只是个位数相加,也不会大幅改变原来图片
将两张图片从新的图片里剥离,只需要反着来就可以了
*/

package main

import (
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func main() {
	EncodeImage("a.png", "b.png", "c.png")
	DecodeImage("c.png", "d.png", "e.png")
}

// 将两张图片合成为一张图片
func EncodeImage(showImagePath, hideImagePath, outPutImagePath string) {

	showImage, err := os.Open(showImagePath)

	defer showImage.Close()

	if err != nil {
		panic(err)
	}

	showImageDecode, _, err := image.Decode(showImage)

	if err != nil {
		panic(err)
	}

	hideImage, err := os.Open(hideImagePath)

	defer hideImage.Close()

	if err != nil {
		panic(err)
	}

	hideImageDecode, _, err := image.Decode(hideImage)

	if err != nil {
		panic(err)
	}

	newImage := encodeImage(showImageDecode, hideImageDecode)

	targetFile, err := os.Create(outPutImagePath)

	defer targetFile.Close()

	if err != nil {
		panic(err)
	}

	if err := imageEncode(outPutImagePath, targetFile, newImage); err != nil {
		panic(err)
	}
}

// 将两张图片从一张图片里剥离
func DecodeImage(targetImagePath, showImagePath, hideImagePath string) {

	targetImage, err := os.Open(targetImagePath)

	if err != nil {
		panic(err)
	}

	defer targetImage.Close()

	targetImageDecode, _, err := image.Decode(targetImage)

	if err != nil {
		panic(err)
	}

	showImageRgba, hideImageRgba := decodeImage(targetImageDecode)

	showImage, err := os.Create(showImagePath)
	if err != nil {
		panic(err)
	}
	defer showImage.Close()

	hideImage, err := os.Create(hideImagePath)
	if err != nil {
		panic(err)
	}
	defer hideImage.Close()

	if err := imageEncode(showImagePath, showImage, showImageRgba); err != nil {
		panic(err)
	}

	if err := imageEncode(hideImagePath, hideImage, hideImageRgba); err != nil {
		panic(err)
	}

}

/*
图片编码
*/
func encodeImage(showImage, hideImage image.Image) *image.RGBA {

	// todo 将两张图片大小变成一样的

	outImageRgba := image.NewRGBA(showImage.Bounds())

	for x := 0; x < hideImage.Bounds().Dx(); x++ {

		for y := 0; y < hideImage.Bounds().Dy(); y++ {

			showPointRgba := showImage.At(x, y)

			showR, showG, showB, _ := showPointRgba.RGBA()

			showRUint8 := showR >> 8
			showGUint8 := showG >> 8
			showBUint8 := showB >> 8

			hidePointRgba := hideImage.At(x, y)

			hideR, hideG, hideB, _ := hidePointRgba.RGBA()

			hideRUint8 := hideR >> 8
			hideGUint8 := hideG >> 8
			hideBUint8 := hideB >> 8

			outValueFunc := func(showValue, hideValue uint8) uint8 {
				return showValue/10*10 + uint8(float32(hideValue)*9/255)
			}

			outG := outValueFunc(uint8(showRUint8), uint8(hideRUint8))
			outR := outValueFunc(uint8(showGUint8), uint8(hideGUint8))
			outB := outValueFunc(uint8(showBUint8), uint8(hideBUint8))

			outImageRgba.SetRGBA(x, y, color.RGBA{R: outR, G: outG, B: outB, A: 255})

		}
	}

	return outImageRgba
}

/*
图片解码
*/
func decodeImage(img image.Image) (*image.RGBA, *image.RGBA) {

	showImageRgba := image.NewRGBA(img.Bounds())

	hideImageRgba := image.NewRGBA(img.Bounds())

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {

			pointRgba := img.At(x, y)

			r, g, b, _ := pointRgba.RGBA()

			trans := func(value uint32) (uint8, uint8) {

				newValue := value >> 8

				// 舍去当前数的个位数
				showValue := newValue / 10 * 10

				// 将舍去的个位数转成float32 进行比例转换
				hideValue := float32(newValue-showValue) * 255 / 9

				// 最后转成uint8返回
				newShowValue := uint8(showValue)

				newHideValue := uint8(hideValue)

				return newShowValue, newHideValue
			}

			showR, hideR := trans(r)
			showG, hideG := trans(g)
			showB, hideB := trans(b)

			showImageRgba.SetRGBA(x, y, color.RGBA{R: showR, G: showG, B: showB, A: 255})

			hideImageRgba.SetRGBA(x, y, color.RGBA{R: hideR, G: hideG, B: hideB, A: 255})

		}
	}

	return showImageRgba, hideImageRgba
}

func imageEncode(fileName string, file *os.File, rgba *image.RGBA) error {

	// 将图片和扩展名分离
	stringSlice := strings.Split(fileName, ".")

	// 根据图片的扩展名来运用不同的处理
	switch stringSlice[len(stringSlice)-1] {
	case "jpg":
		return jpeg.Encode(file, rgba, nil)
	case "jpeg":
		return jpeg.Encode(file, rgba, nil)
	case "gif":
		return gif.Encode(file, rgba, nil)
	case "png":
		return png.Encode(file, rgba)
	default:
		panic("不支持的图片类型")
	}
}
