package main

import (
    "github.com/ajstarks/svgo"
    //"flag"
    "strconv"
    "os"
    "bufio"
    "fmt"
    "math"
    //"strings"
    "path/filepath"
)

var border_width uint32 = 1;
	
func check(e error) {
    if e != nil {
        panic(e)
    }
}

func nextPerfectSquare(num uint32) (uint32, uint64) {
	s := uint32( math.Ceil(math.Sqrt(float64(num))) )
	sq := uint64(s)*uint64(s)
	return s, sq
}

func main() {
	//Parse command line arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: input_file output_directory")
	}
	var input = os.Args[1];
	var output = os.Args[2];


	//Load target dictionary file
	fmt.Printf("Opening %s...\n", input)
	f, err := os.Open(input)
    check(err)
    defer f.Close()

    var min_digits float64 = 0.0;

    codes := make([]uint64,0)

	scanner := bufio.NewScanner(f)
    for scanner.Scan() {
	    code, err := strconv.ParseUint(scanner.Text(), 0, 64);
	    check(err)
	    codes = append(codes, code)
	    //fmt.Println(code)
	    t := math.Ceil(math.Log2(float64(code)))
	    if t > min_digits {
	    	min_digits = t;
	    }
    }
    fmt.Printf("%d codes found\n", len(codes))
    pattern_width, num_bits := nextPerfectSquare(uint32(min_digits))

    //fmt.Printf("Minimum binary digits: %d\n", int(min_digits))
    fmt.Printf("Detected pattern size: %dx%d\n", pattern_width, pattern_width)

	_, serr := os.Stat(output)
    if os.IsNotExist(serr) {
    	os.MkdirAll(output, os.ModePerm)
    }

    fmt.Printf("Writing files to %s...", output)
    for c := 0; c < len(codes); c++ {
		path := filepath.Join(output, strconv.Itoa(c) + ".svg")
		out, ferr := os.Create(path)
		check(ferr)

		//Generate svg

		idx := num_bits-1;
		pixel_size := uint32(50);

		num := codes[c]


		//Save to disk at target direcory
	    width := int(pixel_size*(pattern_width + border_width*2))
	    height := int(pixel_size*(pattern_width + border_width*2))
	    canvas := svg.New(out)
	    canvas.Start(width, height)
	    canvas.Square(0, 0, width, "fill:black")
		for j := 0; j < int(pattern_width); j++ {
			for i := 0; i < int(pattern_width); i++ {
				mask := uint64(1) << idx
				x := (i+int(border_width))*int(pixel_size)
				y := (j+int(border_width))*int(pixel_size)
				//fmt.Println(num & mask)
				if (num & mask) != 0 {
					canvas.Square(x, y, int(pixel_size), "fill:white")
				}
				idx--
			}
		}	
	    canvas.End()
	    out.Close();
    }
    fmt.Println("done!")
}