package main

import (
	"fmt"
	// "os"
	// "image"
	// _ "image/jpeg"
)


type Color struct {
	r uint32
	g uint32
	b uint32
	a uint32
	run uint32
}

type Stack struct {
	items []Color
}

type Stackable interface {
	Push() []Color
	Pop() Color
}

func (s Stack) Push(v Color) Stack {
	s.items = append(s.items, v);
	return s;
}

func main() {
	fmt.Println("start");
	var stack Stack

	fmt.Println("stack: ", stack)
	color := Color{r: 123, g: 123, b:3254345, a: 1}
	stack = stack.Push(color)
	fmt.Println("stack 2: ", stack)

	// imageFile, err := os.Open("./test.jpeg");

	// if err != nil {
	// 	// handle error
	// 	fmt.Println("error opening image");
	// }

	// imageData, _, err := image.Decode(imageFile)

	// fmt.Println("image Data: ", imageData.Bounds())
	// bounds := imageData.Bounds();
	// var data []Color

	// for i := bounds.Min.X; i < bounds.Max.X; i ++ {
	// 	for j := bounds.Min.Y; j < bounds.Max.Y; j ++ {
	// 		r, g, b, a := imageData.At(i, j).RGBA()

	// 		fmt.Println("r: ", r)
	// 		fmt.Println("g: ", g)
	// 		fmt.Println("b: ", b)
	// 		fmt.Println("a: ", a)
	// 		// fmt.Println("point: ", imageData.At(i, j).RGBA())
	// 	}
	// }
}
