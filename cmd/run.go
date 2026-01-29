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

func rowCompress(imageMatrix [][]byte, scale int) [][]byte {
	// Row wise image gray code compression.
	var rowCompressGrayimageCode [][]byte

	for _, rowCodes := range imageMatrix {
		counter := 1
		avgRowValue := 0
		var tempCompressRowCodes []byte
		for _, rowValue := range rowCodes {
			if counter == scale+1 {
				tempCompressRowCodes = append(tempCompressRowCodes, byte(avgRowValue/scale))
				avgRowValue = 0
				counter = 1
			}
			avgRowValue += int(rowValue)
			counter++
		}

		if counter != 1 {
			tempCompressRowCodes = append(tempCompressRowCodes, byte(avgRowValue))
			avgRowValue = 0
			counter = 1
		}

		rowCompressGrayimageCode = append(rowCompressGrayimageCode, tempCompressRowCodes)
	}

	return rowCompressGrayimageCode
}

func colCopress(imageMatrix [][]byte, scale int) [][]byte {
	// column wise compress
	var colCompressGrayimageCode [][]byte
	totalColNumber := len(imageMatrix[0]) // col number in a image

	counter := 1
	tempColCompresionValues := make([]int, totalColNumber)
	for _, row := range imageMatrix {

		if counter == scale+1 {
			tempColCompresionValuesByte := make([]byte, totalColNumber)
			// convert int to byte slice
			for idx, value := range tempColCompresionValues {
				tempColCompresionValuesByte[idx] = byte(value / scale)
			}

			colCompressGrayimageCode = append(colCompressGrayimageCode, tempColCompresionValuesByte) // puting compress values

			counter = 1
			tempColCompresionValues = make([]int, totalColNumber)

			// scale + 1 row are ignore so i have to do this
			for idx, value := range row {
				tempColCompresionValues[idx] = tempColCompresionValues[idx] + int(value)
			}
			counter++
		} else {
			for idx, value := range row {
				tempColCompresionValues[idx] = tempColCompresionValues[idx] + int(value)
			}
			counter++
		}
	}

	if counter != 1 {
		tempColCompresionValuesByte := make([]byte, totalColNumber)

		for idx, value := range tempColCompresionValues {
			tempColCompresionValuesByte[idx] = byte(value / scale) // counter and scale both can use for division
		}

		colCompressGrayimageCode = append(colCompressGrayimageCode, tempColCompresionValuesByte) // puting compress values

	}

	return colCompressGrayimageCode

}

func imageCompression(imageMatrix [][]byte, scale int) [][]byte {
	var rowCompressGrayimageCode [][]byte

	var colCompressGrayimageCode [][]byte

	if scale != 1 {
		rowCompressGrayimageCode = rowCompress(imageMatrix, scale) // compress row wise scale time

		colCompressGrayimageCode = colCopress(rowCompressGrayimageCode, scale) // column wise compress

		return colCompressGrayimageCode
	}

	return imageMatrix
}

func imageToAscii(imageMatrix [][]byte) [][]string {
	chars := " .`'^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	lenChars := len(chars) - 1
	palletSize := len(imageMatrix[0])

	asciiArt := make([][]string, len(imageMatrix))

	for _, row := range imageMatrix {
		tempAsciiChars := make([]string, palletSize)

		for idx, v := range row {
			asciiIdx := (int(v) * lenChars) / 255
			tempAsciiChars[idx] = string(chars[asciiIdx])
		}

		asciiArt = append(asciiArt, tempAsciiChars)
	}
	return asciiArt
}

func Run() {

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

	var grayImageCodes [][]byte

	for y := yMin; y < yMax; y++ {
		var rowGrayCodes []byte
		for x := xMin; x < xMax; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r, g, b = r>>8, g>>8, b>>8 // convert 16 bit to 8 bit

			grayScale := int((0.299 * float64(r)) + (0.587 * float64(g)) + (0.114 * float64(b)))
			rowGrayCodes = append(rowGrayCodes, byte(grayScale))
			fmt.Printf("Color code : %v %v %v %v\n", r, g, b, grayScale)
		}
		grayImageCodes = append(grayImageCodes, rowGrayCodes)
	}

	compressImageCodes := imageCompression(grayImageCodes, *scale)

	asciiArt := imageToAscii(compressImageCodes)

	for _, row := range asciiArt {
		for _,v := range row {
			fmt.Print(v)
		}
		fmt.Println()
	}

}
