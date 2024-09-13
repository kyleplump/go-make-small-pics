package main

import (
	"encoding/binary"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"os"
	"compress/zlib"
	"image/color"
	"image/png"
)


type Color struct {
	r uint16
	g uint16
	b uint16
	a uint16
	run uint16
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

func compressImage(imgFile *os.File) *os.File {
	var stack Stack
	compressedFile, err := os.Create("./compressed.gmis");
	imageData, _, err := image.Decode(imgFile);
	bounds := imageData.Bounds();

	if err != nil {
		// todo
		fmt.Println("error decoding image file");
	}

	writer := zlib.NewWriter(compressedFile);

	defer writer.Close();

	for i := bounds.Min.X; i < bounds.Max.X; i ++ {
		for j := bounds.Min.Y; j < bounds.Max.Y; j ++ {
			r, g, b, a := imageData.At(i, j).RGBA()
			pixel := Color{r: uint16(r), g: uint16(g), b: uint16(b), a: uint16(a)}

			if len(stack.items) == 0 {
			 	stack.Push(Color{r: pixel.r, g: pixel.g, b: pixel.b, a: pixel.a, run: 1 })
			} else {
				popped_value, _ := stack.Pop()
				matched := false

				if popped_value.r == uint16(r) && popped_value.g == uint16(g) && popped_value.b == uint16(b) && popped_value.a == uint16(a) {
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

	for _, color := range stack.items {
		arr := []uint16{
			color.r,
			color.g,
			color.b,
			color.a,
			color.run,
		}
		s := []byte{};

		for _, value := range arr {
			chunk := make([]byte, 2);
			binary.LittleEndian.PutUint16(chunk, value);
			s = append(s, chunk...);
		}
		writer.Write(s)
	}

	return compressedFile;
}

func rebuildFile(compressedFile *os.File, bounds image.Rectangle) {
	zlibReader, _ := zlib.NewReader(compressedFile);


	img := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y));
	var newStack Stack

	for {
		buf := make([]byte, 10);
		n, err := zlibReader.Read(buf);

		if err == io.EOF {
			lastChunk := buf[:10];

			r := binary.LittleEndian.Uint16(lastChunk[0:2])
			g := binary.LittleEndian.Uint16(lastChunk[2:4])
			b := binary.LittleEndian.Uint16(lastChunk[4:6])
			a := binary.LittleEndian.Uint16(lastChunk[6:8])
			run := binary.LittleEndian.Uint16(lastChunk[8:10])

			c := Color{r: r, g: g, b: b, a: a, run: run}
			newStack.Push(c);
			break;
		}

		if err != nil {
			fmt.Println("error: ", err)
			return;
		}

		if n != 10 {
			if n < 10 {
				tmpBuf := make([]byte, 10 - n);
				zlibReader.Read(tmpBuf);

				if n == 2 {
					r := binary.LittleEndian.Uint16(buf[0:2])
					g := binary.LittleEndian.Uint16(tmpBuf[0:2])
					b := binary.LittleEndian.Uint16(tmpBuf[2:4])
					a := binary.LittleEndian.Uint16(tmpBuf[4:8])
					run := binary.LittleEndian.Uint16(tmpBuf[6:8])

					c := Color{r: r, g: g, b: b, a: a, run: run}
					newStack.Push(c);
					continue;

				} else if n == 4 {
					r := binary.LittleEndian.Uint16(buf[0:2])
					g := binary.LittleEndian.Uint16(buf[2:4])
					b := binary.LittleEndian.Uint16(tmpBuf[0:2])
					a := binary.LittleEndian.Uint16(tmpBuf[2:4])
					run := binary.LittleEndian.Uint16(tmpBuf[4:6])

					c := Color{r: r, g: g, b: b, a: a, run: run}
					newStack.Push(c);
					continue;

				} else if n == 6 {
					r := binary.LittleEndian.Uint16(buf[0:2])
					g := binary.LittleEndian.Uint16(buf[2:4])
					b := binary.LittleEndian.Uint16(buf[4:6])
					a := binary.LittleEndian.Uint16(tmpBuf[0:2])
					run := binary.LittleEndian.Uint16(tmpBuf[2:4])

					c := Color{r: r, g: g, b: b, a: a, run: run}
					newStack.Push(c);
					continue;

				} else if n == 8 {
					r := binary.LittleEndian.Uint16(buf[0:2])
					g := binary.LittleEndian.Uint16(buf[2:4])
					b := binary.LittleEndian.Uint16(buf[4:6])
					a := binary.LittleEndian.Uint16(buf[6:8])
					run := binary.LittleEndian.Uint16(tmpBuf[0:2])

					c := Color{r: r, g: g, b: b, a: a, run: run}
					newStack.Push(c);
					continue;
				}
			} else {
				chunk := buf[:10];
				buf = buf[10:];

				r := binary.LittleEndian.Uint16(chunk[0:2])
				g := binary.LittleEndian.Uint16(chunk[2:4])
				b := binary.LittleEndian.Uint16(chunk[4:6])
				a := binary.LittleEndian.Uint16(chunk[6:8])
				run := binary.LittleEndian.Uint16(chunk[8:10])

				c := Color{r: r, g: g, b: b, a: a, run: run}
				newStack.Push(c);
				continue;
			}
		}

		lastChunk := buf[:10];

		r := binary.LittleEndian.Uint16(lastChunk[0:2])
		g := binary.LittleEndian.Uint16(lastChunk[2:4])
		b := binary.LittleEndian.Uint16(lastChunk[4:6])
		a := binary.LittleEndian.Uint16(lastChunk[6:8])
		run := binary.LittleEndian.Uint16(lastChunk[8:10])

		c := Color{r: r, g: g, b: b, a: a, run: run}
		newStack.Push(c);
	}

	curRun, _ := newStack.Pop();

	for i := bounds.Max.X; i > bounds.Min.X; i -- {
		for j := bounds.Max.Y; j > bounds.Min.Y; j -- {
			if curRun.run == 0 {
				curRun, _ = newStack.Pop();
			}
			img.Set(i, j, color.RGBA64{curRun.r, curRun.g, curRun.b, curRun.a});
			curRun.run --;
		}
	}

	outfile, _ := os.Create("output.png");
	png.Encode(outfile, img);

	outfile.Close();
}

func getImageBounds(img *os.File) image.Rectangle {
	img.Seek(0, 0);
	imageData, _, err := image.Decode(img);

	if err != nil {
		fmt.Println("error: ", err, imageData);
	}

	bounds := imageData.Bounds();

	return bounds;
}

func main() {

	imageFile, err := os.Open("./test_input.jpeg");

	if err != nil {
		// todo
		fmt.Println("error opening image file");
	}

	defer imageFile.Close();

	// print input file size
	stats, _ := imageFile.Stat();
	fmt.Println("original file size: ", stats.Size());

	compressedFile := compressImage(imageFile);
	compressedFile.Seek(0, 0);

	compressedStats, err := compressedFile.Stat();

	if err != nil {
		// todo
	}

	// report compressed size
	fmt.Println("compressed file size:", compressedStats.Size());
	fmt.Println("uncompressing file ...");

	bounds := getImageBounds(imageFile);
	rebuildFile(compressedFile, bounds);
}
