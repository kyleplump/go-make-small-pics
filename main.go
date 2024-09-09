package main

import (
	"encoding/binary"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"os"

	// "encoding/json"
	// "bufio"
	"compress/zlib"
	// "unsafe"
	"image/color"
	"image/png"
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
	var stack Stack

	imageFile, err := os.Open("./test.jpeg");
	f, _ := os.Create("./compressed.gmis");
	w := zlib.NewWriter(f);
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

			if len(stack.items) == 0 {
			 	stack.Push(Color{r: pixel.r, g: pixel.g, b: pixel.b, a: pixel.a, run: 1 })
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
					 stack.Push(Color{r: pixel.r, g: pixel.g, b: pixel.b, a: pixel.a, run: 1 });
				}
			}
		}
	}

	fmt.Println("size of stack: ", len(stack.items))

	for _, color := range stack.items {
		arr := []uint32{
			color.r,
			color.g,
			color.b,
			color.a,
			color.run,
		}
		s := []byte{};
		// fmt.Println("writing color: ", color)
		for _, value := range arr {
			chunk := make([]byte, 4);
			binary.LittleEndian.PutUint32(chunk, value);
			s = append(s, chunk...);
		}
		w.Write(s)
		// stringifiedColor, _ := json.Marshal(color);
		// f.WriteString(string(stringifiedColor));
	}

	w.Close();
	f.Close();

	compressedFile, _ := os.Open("./compressed.gmis");

	compressedStats, _ := compressedFile.Stat();

	// reader := bufio.NewReader(compressedFile);

	fmt.Println("compressed file size:", compressedStats.Size());

	fmt.Println("uncompressing file ...");

	zlibReader, _ := zlib.NewReader(compressedFile);

	buf := make([]byte, 20);
	img := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y));
	var newStack Stack

	for {
		_, err := zlibReader.Read(buf);

		if err == io.EOF {
			break;
		}

		if err != nil {
			fmt.Println("errored");
			return;
		}
		r := binary.LittleEndian.Uint32(buf[0:4])
		g := binary.LittleEndian.Uint32(buf[4:8])
		b := binary.LittleEndian.Uint32(buf[8:12])
		a := binary.LittleEndian.Uint32(buf[12:16])
		run := binary.LittleEndian.Uint32(buf[16:20])
		// fmt.Println("pixel: ", color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		// fmt.Println("run: ", run)
		// fmt.Printf("Read %d bytes: %s\n", n, string(buf[:n]))
		c := Color{r: r, g: g, b: b, a: a, run: run}
		newStack.Push(c);
	}

	fmt.Println("decoded stack size: ", len(newStack.items))
	curRun, _ := stack.Pop();

	for i := bounds.Max.X; i > bounds.Min.X; i -- {
		for j := bounds.Max.Y; j > bounds.Min.Y; j -- {
			if curRun.run == 0 {
				curRun, _ = stack.Pop();
			}
			fmt.Println("color: ", curRun)
			img.Set(i, j, color.RGBA64{uint16(curRun.r), uint16(curRun.g), uint16(curRun.b), uint16(curRun.a)});
			curRun.run --;
		}
	}

	outfile, _ := os.Create("output.png");

	png.Encode(outfile, img);

	outfile.Close();


	// curr_color := make([]byte, compressedStats.Size() * 2);
	// zlibReader.Read(curr_color);
	// var result []Color

	// // Read 4 bytes at a time
	// for i := 0; i < len(curr_color); i += 20 {
	// 	chunk := curr_color[i : i+20]
	// 	color := Color{binary.LittleEndian.Uint32(chunk[0:4]), binary.LittleEndian.Uint32(chunk[4:8]), binary.LittleEndian.Uint32(chunk[8:12]), binary.LittleEndian.Uint32(chunk[12:16]), binary.LittleEndian.Uint32(chunk[16:20])}
	// 	result = append(result, color)
	// }

	// fmt.Println("result: ", result)





	// fmt.Println("bytes read: ", n)
	// fmt.Println("curr color size: ", len(curr_color))
	// page := int64(1)
	// curr_page := curr_color[0:20]
	// pixels_written := binary.LittleEndian.Uint32(curr_page[16:20]);
	// // fmt.Println("curr run: ", pixels_written)
	// img := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y));
	// fmt.Println("bounds: ", bounds)
	// for i := bounds.Min.X; i < bounds.Max.X; i ++ {
	// 	for j := bounds.Min.Y; j < bounds.Max.Y; j ++ {
	// 		fmt.Println("writing: ", i, j)
	// 		if(pixels_written > 0) {
	// 			r := binary.LittleEndian.Uint32(curr_page[0:4])
	// 			g := binary.LittleEndian.Uint32(curr_page[4:8])
	// 			b := binary.LittleEndian.Uint32(curr_page[8:12])
	// 			a := binary.LittleEndian.Uint32(curr_page[12:16])
	// 			fmt.Println("pixel: ", color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})

	// 			img.Set(i, j, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)});
	// 			pixels_written = pixels_written - 1;
	// 			// fmt.Println("pixels written: ", pixels_written)
	// 		} else {
	// 			fmt.Println("next page")
	// 			start := page * 20
	// 			end := start + 20

	// 			if end < (compressedStats.Size() - 20) {

	// 				curr_page = curr_color[start:end]
	// 				fmt.Println("start: ", start, ", end: ", end, ", size: ", compressedStats.Size())
	// 				pixels_written = binary.LittleEndian.Uint32(curr_page[(end - 4):end]);
	// 				page = page + 1
	// 			}

	// 		}
	// 	}
	// }

	// outfile, _ := os.Create("output.png");

	// png.Encode(outfile, img);

	// outfile.Close();
}
