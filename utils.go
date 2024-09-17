package main

import (
	"image"
	"fmt"
	"encoding/binary"
	"os"
)

func ImageBounds(img *os.File) image.Rectangle {
	img.Seek(0, 0);
	imageData, _, err := image.Decode(img);

	if err != nil {
		fmt.Println("error: ", err, imageData);
	}

	bounds := imageData.Bounds();

	return bounds;
}

func DecodeColor(chunk []byte) Color {
	r := binary.LittleEndian.Uint16(chunk[0:2])
	g := binary.LittleEndian.Uint16(chunk[2:4])
	b := binary.LittleEndian.Uint16(chunk[4:6])
	a := binary.LittleEndian.Uint16(chunk[6:8])
	run := binary.LittleEndian.Uint16(chunk[8:10])

	color := Color{r: r, g: g, b: b, a: a, run: run}

	return color;
}
