package cmd

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func Run() {
	var grayImageCodes []byte

	filePath := flag.String("p", `./images/tjImg_logo.png`, "image file path")
	scale := flag.Int("x", 10, "scale down the image x time")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
		return
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}

	bounds := img.Bounds()

	xMin, xMax := bounds.Min.X, bounds.Max.X
	yMin, yMax := bounds.Min.Y, bounds.Max.Y

	imageWidth := xMax - xMin
	imageHeight := yMax - yMin

	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r, g, b = r>>8, g>>8, b>>8 // convert 16 bit to 8 bit

			grayScale := int((0.299 * float64(r)) + (0.587 * float64(g)) + (0.114 * float64(b)))
			grayImageCodes = append(grayImageCodes, byte(grayScale))
			fmt.Printf("Color code : %v %v %v %v\n", r, g, b, grayScale)
		}
		grayImageCodes = append(grayImageCodes, byte('\n'))
	}

	colGroup := imageWidth / *scale
	rowGroup := imageHeight / *scale

	var rowCompressGrayimageCode []byte

	for row := range rowGroup {
		for col := range colGroup {
			grayValueSum := 0
			for s := range *scale {
				idx := row*col**scale + s

				fmt.Println("test", idx)
				if grayImageCodes[idx] == byte('\n') { // 10 mean \n new line
					break
				}
				grayValueSum += int(grayImageCodes[idx])
				rowCompressGrayimageCode = append(rowCompressGrayimageCode, byte(grayValueSum / *scale))
			}
			rowCompressGrayimageCode = append(rowCompressGrayimageCode, byte(grayValueSum / *scale), byte('\n'))
		}
	}
	fmt.Println(len(rowCompressGrayimageCode))

}
