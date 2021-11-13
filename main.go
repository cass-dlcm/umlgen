package main

import (
	"encoding/json"
	"flag"
	"github.com/cass-dlcm/umlgen/lib"
	"io"
	"log"
	"math/rand"
	"os"
)

func getFlags() (string, string) {
	input := flag.String("input", "", "where to read the data from")
	output := flag.String("output", "", "where to save the file to")
	flag.Parse()
	return *input, *output
}

func main() {
	input, output := getFlags()
	var diagram lib.Diagram
	if input == "" {
		err := json.NewDecoder(os.Stdin).Decode(&diagram)
		if err != nil {
			log.Panic(err)
		}
	} else {
		inputFile, err := os.Open(input)
		if err != nil {
			log.Panic(err)
		}
		defer func(inputFile *os.File) {
			err := inputFile.Close()
			if err != nil {
				log.Panic(err)
			}
		}(inputFile)
		err = json.NewDecoder(inputFile).Decode(&diagram)
		if err != nil {
			log.Panic(err)
		}
	}
	var canvas io.Writer
	if output != "" {
		var err error
		canvas, err = os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			log.Panic(err)
		}
	} else {
		canvas = os.Stdout
	}
	rand.Seed(diagram.Seed)
	lib.Generate(canvas, diagram)
}
