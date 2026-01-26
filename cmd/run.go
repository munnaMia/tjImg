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
			r, g, b, a := img.At(x, y).RGBA()
			fmt.Printf("Color code : %v %v %v %v\n", r, g, b, a)
		}
	}
}
