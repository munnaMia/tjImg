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
	var asciiImage []byte

	filePath := flag.String("p", `./images/tjImg_logo.png`, "image file path")
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

	fmt.Println("Red Green Blue Alpha")
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r, g, b = r>>8, g>>8, b>>8 // convert 16 bit to 8 bit

			grayScale := int((0.299 * float64(r)) + (0.587 * float64(g)) + (0.114 * float64(b)))
			asciiImage = append(asciiImage, byte(grayScale))
			fmt.Printf("Color code : %v %v %v %v\n", r, g, b, grayScale)
		}
		asciiImage = append(asciiImage, '\n')
	}

	fmt.Println(asciiImage)
}
