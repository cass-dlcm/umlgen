package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
)

func getFlags() (string, string) {
	input := flag.String("input", "", "where to read the data from")
	output := flag.String("output", "", "where to save the file to")
	flag.Parse()
	return *input, *output
}

type Attribute struct {
	Visibility string `json:"visibility,omitempty"`
	Name       string `json:"name"`
	Type       string `json:"type"`
}

type Argument struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Method struct {
	Visibilty string     `json:"visibilty,omitempty"`
	Name      string     `json:"name"`
	Args      []Argument `json:"args,omitempty"`
	Return    string     `json:"return,omitempty"`
}

type Class struct {
	Name       string      `json:"name" `
	Attributes []Attribute `json:"attributes,omitempty"`
	Methods    []Method    `json:"methods,omitempty"`
}

type Interaction struct {
	ClassAIndex int    `json:"class_a_index"`
	ClassBIndex int    `json:"class_b_index"`
	ClassAText  string `json:"class_a_text,omitempty"`
	ClassBText  string `json:"class_b_text,omitempty"`
}

type Diagram struct {
	Seed         int64         `json:"seed,omitempty"`
	Classes      []Class       `json:"classes"`
	Interactions []Interaction `json:"interactions,omitempty"`
}

func GetClassDimensions(c Class) (int, int) {
	longestStr := len(c.Name)
	height := 1
	for a := range c.Attributes {
		height++
		attrLen := len(c.Attributes[a].Name) + 1 + len(c.Attributes[a].Type)
		if c.Attributes[a].Visibility != "" {
			attrLen += 1
		}
		longestStr = int(math.Max(float64(longestStr), float64(attrLen)))
	}
	for m := range c.Methods {
		height++
		mLen := len(c.Methods[m].Name) + 2
		if c.Methods[m].Args != nil {
			mLen -= 2
			for a := range c.Methods[m].Args {
				mLen += len(c.Methods[m].Args[a].Type) + 3 + len(c.Methods[m].Args[a].Name)
			}
		}
		if c.Methods[m].Return != "" {
			mLen += len(c.Methods[m].Return) + 1
		}
		longestStr = int(math.Max(float64(longestStr), float64(mLen)))
	}
	return longestStr + 2, height + 1
}

func ClassGen(canvas io.Writer, x, y int, c Class) {
	longestStr := len(c.Name)
	height := 1
	for a := range c.Attributes {
		height++
		attrLen := len(c.Attributes[a].Name) + 1 + len(c.Attributes[a].Type)
		if c.Attributes[a].Visibility != "" {
			attrLen += 1
		}
		longestStr = int(math.Max(float64(longestStr), float64(attrLen)))
	}
	for m := range c.Methods {
		height++
		mLen := len(c.Methods[m].Name) + 2
		if c.Methods[m].Args != nil {
			mLen -= 2
			for a := range c.Methods[m].Args {
				mLen += len(c.Methods[m].Args[a].Type) + 3 + len(c.Methods[m].Args[a].Name)
			}
		}
		if c.Methods[m].Return != "" {
			mLen += len(c.Methods[m].Return) + 1
		}
		longestStr = int(math.Max(float64(longestStr), float64(mLen)))
	}
	heightOffset := 0
	if _, err := canvas.Write([]byte(fmt.Sprintf("\t<rect x=\"%dem\" y=\"%dem\" width=\"%dem\" height=\"%dem\" style=\"stroke-width:1;stroke:rgb(0,0,0);fill:rgb(255,255,255)\" />\n", x, y, longestStr+2, height+1))); err != nil {
		log.Panic(err)
	}
	heightOffset++
	for r := range c.Name {
		if _, err := canvas.Write([]byte(fmt.Sprintf("\t<text x=\"%dem\" y=\"%dem\" style=\"text-size:16em\" >%c</text>\n", x+1+r, y+heightOffset, c.Name[r]))); err != nil {
			log.Panic(err)
		}
	}
	heightOffset++
	if _, err := canvas.Write([]byte(fmt.Sprintf("\t<line x1=\"%dem\" y1=\"%dem\" x2=\"%dem\" y2=\"%dem\" style=\"stroke:rgb(0,0,0);stroke-width:1\" />\n", x, y+heightOffset-1, longestStr+2+x, y+heightOffset-1))); err != nil {
		log.Panic(err)
	}
	for a := range c.Attributes {
		vis := '\000'
		switch c.Attributes[a].Visibility {
		case "protected":
			vis = '*'
		case "public":
			vis = '+'
		case "private":
			vis = '-'
		}
		if vis != '\000' {
			if _, err := canvas.Write([]byte(fmt.Sprintf("\t<text x=\"%dem\" y=\"%dem\" style=\"text-size:16em\" >%c</text>\n", x+1, y+heightOffset, vis))); err != nil {
				log.Panic(err)
			}
		}
		for r := range c.Attributes[a].Name {
			if vis != '\000' {
				if _, err := canvas.Write([]byte(fmt.Sprintf("\t<text x=\"%dem\" y=\"%dem\" style=\"text-size:16em\" >%c</text>\n", x+2+r, y+heightOffset, c.Attributes[a].Name[r]))); err != nil {
					log.Panic(err)
				}
			} else {
				if _, err := canvas.Write([]byte(fmt.Sprintf("\t<text x=\"%dem\" y=\"%dem\" style=\"text-size:16em\" >%c</text>\n", x+1+r, y+heightOffset, c.Attributes[a].Name[r]))); err != nil {
					log.Panic(err)
				}
			}
		}
		for r := range c.Attributes[a].Type {
			if vis != '\000' {
				if _, err := canvas.Write([]byte(fmt.Sprintf("\t<text x=\"%dem\" y=\"%dem\" style=\"text-size:16em\" >%c</text>\n", x+3+r+len(c.Attributes[a].Name), y+heightOffset, c.Attributes[a].Type[r]))); err != nil {
					log.Panic(err)
				}
			} else {
				if _, err := canvas.Write([]byte(fmt.Sprintf("\t<text x=\"%dem\" y=\"%dem\" style=\"text-size:16em\" >%c</text>\n", x+2+r+len(c.Attributes[a].Name), y+heightOffset, c.Attributes[a].Type[r]))); err != nil {
					log.Panic(err)
				}
			}
		}
		heightOffset++
	}
	if _, err := canvas.Write([]byte(fmt.Sprintf("\t<line x1=\"%dem\" y1=\"%dem\" x2=\"%dem\" y2=\"%dem\" style=\"stroke:rgb(0,0,0);stroke-width:1\" />\n", x, y+heightOffset-1, longestStr+2+x, y+heightOffset-1))); err != nil {
		log.Panic(err)
	}
	for m := range c.Methods {
		vis := '\000'
		switch c.Methods[m].Visibilty {
		case "protected":
			vis = '*'
		case "public":
			vis = '+'
		case "private":
			vis = '-'
		}
		length := x + 1
		if vis != '\000' {
			if _, err := canvas.Write([]byte(fmt.Sprintf("\t<text x=\"%dem\" y=\"%dem\" style=\"text-size:16em\" >%c</text>\n", length, y+heightOffset, vis))); err != nil {
				log.Panic(err)
			}
			length++
		}
		for r := range c.Methods[m].Name {
			if _, err := canvas.Write([]byte(fmt.Sprintf("\t<text x=\"%dem\" y=\"%dem\" style=\"text-size:16em\" >%c</text>\n", length+r, y+heightOffset, c.Methods[m].Name[r]))); err != nil {
				log.Panic(err)
			}
		}
		length += len(c.Methods[m].Name)
		if _, err := canvas.Write([]byte(fmt.Sprintf("\t<text x=\"%dem\" y=\"%dem\" style=\"text-size:16em\">(</text>\n", length, y+heightOffset))); err != nil {
			log.Panic(err)
		}
		length++
		for a := range c.Methods[m].Args {
			for r := range c.Methods[m].Args[a].Type {
				if _, err := canvas.Write([]byte(fmt.Sprintf("\t<text x=\"%dem\" y=\"%dem\" style=\"text-size:16em\" >%c</text>\n", length+r, y+heightOffset, c.Methods[m].Args[a].Type[r]))); err != nil {
					log.Panic(err)
				}
			}
			length += len(c.Methods[m].Args[a].Type) + 1
			for r := range c.Methods[m].Args[a].Name {
				if _, err := canvas.Write([]byte(fmt.Sprintf("\t<text x=\"%dem\" y=\"%dem\" style=\"text-size:16em\" >%c</text>\n", length+r, y+heightOffset, c.Methods[m].Args[a].Name[r]))); err != nil {
					log.Panic(err)
				}
			}
			length += len(c.Methods[m].Args[a].Name)
			if a < len(c.Methods[m].Args)-1 {
				if _, err := canvas.Write([]byte(fmt.Sprintf("\t<text x=\"%dem\" y=\"%dem\" style=\"text-size:16em\" >,</text>\n", length, y+heightOffset))); err != nil {
					log.Panic(err)
				}
				length += 2
			}
		}
		if _, err := canvas.Write([]byte(fmt.Sprintf("\t<text x=\"%dem\" y=\"%dem\" style=\"text-size:16em\">)</text>\n", length, y+heightOffset))); err != nil {
			log.Panic(err)
		}
		length += 2
		for r := range c.Methods[m].Return {
			if _, err := canvas.Write([]byte(fmt.Sprintf("\t<text x=\"%dem\" y=\"%dem\" style=\"text-size:16em\" >%c</text>\n", length+r, y+heightOffset, c.Methods[m].Return[r]))); err != nil {
				log.Panic(err)
			}
		}
		heightOffset++
	}
}

func main() {
	input, output := getFlags()
	var diagram Diagram
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
	width := 0
	height := 0
	classWidths := []int{}
	classHeights := []int{}
	for c := range diagram.Classes {
		classWidth, classHeight := GetClassDimensions(diagram.Classes[c])
		classWidths = append(classWidths, classWidth)
		classHeights = append(classHeights, classHeight)
		width += classWidth + 5
		height += classHeight + 5
	}
	if _, err := canvas.Write([]byte("<?xml version=\"1.0\"?>\n")); err != nil {
		log.Panic(err)
	}
	if _, err := canvas.Write([]byte(fmt.Sprintf("<svg width=\"%dem\" height=\"%dem\"\nxmlns=\"http://www.w3.org/2000/svg\"\nxmlns:xlink=\"http://www.w3.org/1999/xlink\">\n\t<rect x=\"0\" y=\"0\" width=\"%dem\" height=\"%dem\" style=\"fill:rgb(255,255,255)\" />\n", width, height, width, height))); err != nil {
		log.Panic(err)
	}
	classX := []int{}
	classY := []int{}
	for true {
		satisfied := true
		for c := range diagram.Classes {
			classX = append(classX, rand.Intn(width-classWidths[c]-2)+1)
			classY = append(classY, rand.Intn(height-classHeights[c]-2)+1)
			cont := true
			for index := 0; index < c && cont; index++ {
				if (classX[c] < classX[index]+classWidths[index]+1 && classX[c] > classX[index]-1 && ((classY[c] < classY[index]+classHeights[index]+1 && classY[c] > classY[index]-1)  || (classY[c]+classHeights[c] > classY[index]-1 && classY[c] < classY[index]+classHeights[index]+1))) || (classX[c]+classWidths[c] > classX[index]-1 && classX[c] < classX[index]+classWidths[index]+1  && ((classY[c] < classY[index]+classHeights[index]+1 && classY[c] > classY[index]-1)  || (classY[c]+classHeights[c] > classY[index]-1 && classY[c] < classY[index]+classHeights[index]+1))) {
					classX = []int{}
					classY = []int{}
					cont = false
					satisfied = false
				}
			}
		}
		if satisfied {
			break
		}
	}
	for c := range diagram.Classes {
		ClassGen(canvas, classX[c], classY[c], diagram.Classes[c])
	}
	if _, err := canvas.Write([]byte("</svg>\n")); err != nil {
		log.Panic(err)
	}
}
