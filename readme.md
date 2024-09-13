# Go make images smaller

- a Go program that takes a source jpeg, encodes it using run-length encoding, then compresses the RLE information using zlib into a custom file format (.gmis)
- then reads in the compressed file, and rebuilds the source jpeg as a png

this is a simple (sort of dumb) program just to practice starting to work with Go. i understand that in practice RLE is a terrible way of compressing images

pls don't make fun of my amatuer Go code i'm trying
