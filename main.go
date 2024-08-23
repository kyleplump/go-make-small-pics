package main

import (
	"fmt"
	"os"
	"image"
	_ "image/jpeg"

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

func (s *Stack) Push(v Color) {
	s.items = append(s.items, v);
}

func (s *Stack) Pop() (Color, bool) {
	popped_value := s.items[len(s.items) - 1]
	s.items = s.items[:len(s.items) - 1]
	return popped_value, true;
}

func main() {
	fmt.Println("start");
	var stack Stack

	imageFile, err := os.Open("./test.jpeg");
	f, _ := os.Create("./compressed.gmis");
	defer imageFile.Close();
	stats, _ := imageFile.Stat();

	fmt.Println("original file size: ", stats.Size());

	if err != nil {
		// handle error
		fmt.Println("error opening image");
	}

	imageData, _, err := image.Decode(imageFile)

	bounds := imageData.Bounds();

	for i := bounds.Min.X; i < bounds.Max.X; i ++ {
		for j := bounds.Min.Y; j < bounds.Max.Y; j ++ {
			r, g, b, a := imageData.At(i, j).RGBA()
			pixel := Color{r: r, g: g, b: b, a: a}

			// empty stack
			if len(stack.items) == 0 {
			 	stack.Push(pixel)
			} else {
				popped_value, _ := stack.Pop()
				matched := false

				if popped_value.r == r && popped_value.g == g && popped_value.b == b && popped_value.a == a {
					matched = true
				}

				if matched {
					popped_value.run = popped_value.run + 1
					 stack.Push(popped_value);
				} else {
					 stack.Push(popped_value);
					 stack.Push(pixel);
				}
			}
		}
	}

	for _, color := range stack.items {
		fmt.Fprintln(f, color)
	}


	done_stats, _ := f.Stat();

	fmt.Println("resulting size: ", done_stats.Size());
		f.Close();

}
