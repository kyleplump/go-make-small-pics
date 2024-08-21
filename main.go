package main

import (
	"fmt"
	"os"
	"image"
	_ "image/jpeg"
	// "bufio"
)

func main() {
	fmt.Println("hello world");
	imageFile, err := os.Open("./test.jpeg");

	if err != nil {
		// handle error
		fmt.Println("error opening image");
	}

	imageData, _, err := image.Decode(imageFile)

	fmt.Println("image Data: ", imageData.Bounds())
	bounds := imageData.Bounds();

	for i := bounds.Min.X; i < bounds.Max.X; i ++ {
		for j := bounds.Min.Y; j < bounds.Max.Y; j ++ {
			r, g, b, a := imageData.At(i, j).RGBA()
			fmt.Println("r: ", r)
			fmt.Println("g: ", g)
			fmt.Println("b: ", b)
			fmt.Println("a: ", a)
			// fmt.Println("point: ", imageData.At(i, j).RGBA())
		}
	}
}
