package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

var outFormatF = flag.String("out", "jpeg", "output format: 'jpeg', 'png', 'gif'")
var supported = map[string]bool{
	"jpeg": true,
	"png":  true,
	"gif":  true,
}

func main() {
	flag.Parse()

	if err := convert(os.Stdin, os.Stdout, *outFormatF); err != nil {
		fmt.Fprintf(os.Stderr, "imgconv: %v\n", err)
		os.Exit(1)
	}
}

func convert(in io.Reader, out io.Writer, outFmt string) error {
	if !supported[outFmt] {
		return fmt.Errorf("%s is not supported", outFmt)
	}

	buf := bytes.NewBuffer([]byte{})
	buf.ReadFrom(in)
	br := bytes.NewReader(buf.Bytes())
	img, kind, err := image.Decode(br)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "input: %s / output: %s\n", kind, outFmt)
	if kind == outFmt {
		if _, e := br.Seek(io.SeekStart, io.SeekStart); e != nil {
			return e
		}
		_, e := io.Copy(out, br)
		return e
	}
	switch outFmt {
	case "jpeg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
	case "png":
		return png.Encode(out, img)
	case "gif":
		return gif.Encode(out, img, &gif.Options{})
	}
	panic("not reachable")
}
