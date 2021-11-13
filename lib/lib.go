package lib

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
)

type Point struct {
	x float64
	y float64
}

type ClassSpace struct {
	topLeft  Point
	botRight Point
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

func ClassGen(canvas io.Writer, p Point, c Class) {
	x := int(p.x)
	y := int(p.y)
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

func AddInteraction(canvas io.Writer, cA, cB ClassSpace) {
	p1, p2 := Point{(cA.topLeft.x + cA.botRight.x) / 2, (cA.topLeft.y + cA.botRight.y) / 2}, Point{(cB.topLeft.x + cB.botRight.x) / 2, (cB.topLeft.y + cB.botRight.y) / 2}
	if _, err := canvas.Write([]byte(fmt.Sprintf("\t<line x1=\"%fem\" y1=\"%fem\" x2=\"%fem\" y2=\"%fem\" style=\"stroke:rgb(0,0,0);stroke-width:1\" />\n", p1.x, p1.y, p2.x, p2.y))); err != nil {
		log.Panic(err)
	}
	if p1.y - p2.y == 0 {
		arrowHead := Point{(p1.x + p2.x) / 2, (p1.y + p2.y) / 2}
		for true {
			var tempArrowHead Point
			if arrowHead.x > p2.x {
				tempArrowHead = Point{arrowHead.x - 1,arrowHead.y}
			} else {
				tempArrowHead = Point{arrowHead.x + 1,arrowHead.y}
			}
			if !tempArrowHead.isInsideRect(cB) {
				arrowHead = tempArrowHead
			} else {
				break
			}
		}
		var pointL, pointR Point
		if p2.x > p1.x {
			log.Println("a")
			pointL = Point{
				arrowHead.x + math.Sin(2*math.Pi/3),
				arrowHead.y - math.Cos(2*math.Pi/3),
			}
			pointR = Point{
				arrowHead.x - math.Sin(math.Pi/3),
				arrowHead.y + math.Cos(math.Pi/3),
			}
		} else {
			log.Println("b")
			pointL = Point{
				arrowHead.x + math.Sin(2*math.Pi/3),
				arrowHead.y - math.Cos(math.Pi/3),
			}
			pointR = Point{
				arrowHead.x + math.Sin(math.Pi/3),
				arrowHead.y + math.Cos(math.Pi/3),
			}
		}
		if _, err := canvas.Write([]byte(fmt.Sprintf("\t<line x1=\"%fem\" y1=\"%fem\" x2=\"%fem\" y2=\"%fem\" style=\"stroke:rgb(0,0,0);stroke-width:1\" />\n", pointL.x, pointL.y, arrowHead.x, arrowHead.y))); err != nil {
			log.Panic(err)
		}
		if _, err := canvas.Write([]byte(fmt.Sprintf("\t<line x1=\"%fem\" y1=\"%fem\" x2=\"%fem\" y2=\"%fem\" style=\"stroke:rgb(0,0,0);stroke-width:1\" />\n", pointR.x, pointR.y, arrowHead.x, arrowHead.y))); err != nil {
			log.Panic(err)
		}
		return
	}
	arrowHead := Point{(p1.x + p2.x) / 2, (p1.y + p2.y) / 2}
	slope := (p2.y-p1.y)/(p2.x-p1.x)
	for true {
		var tempArrowHead Point
		if arrowHead.x > p2.x {
			tempArrowHead = Point{arrowHead.x - .5, arrowHead.y - slope/2}
		} else {
			tempArrowHead = Point{arrowHead.x + .5,arrowHead.y + slope/2}
		}
		if !tempArrowHead.isInsideRect(cB) {
			arrowHead = tempArrowHead
		} else {
			break
		}
	}
	var pointL, pointR Point
	if p2.y - p1.y < 0 {
		if p2.x - p1.x < 0 {
			pointL = Point{
				arrowHead.x + math.Sin(math.Atan(slope)+math.Pi/3),
				arrowHead.y - math.Cos(math.Atan(slope)+math.Pi/3),
			}
			pointR = Point{
				arrowHead.x - math.Sin(math.Atan(slope)-math.Pi/3),
				arrowHead.y + math.Cos(math.Atan(slope)-math.Pi/3),
			}
		} else {
			pointL = Point{
				arrowHead.x + math.Sin(math.Atan(slope)-math.Pi/3),
				arrowHead.y - math.Cos(math.Atan(slope)-math.Pi/3),
			}
			pointR = Point{
				arrowHead.x - math.Sin(math.Atan(slope)+math.Pi/3),
				arrowHead.y + math.Cos(math.Atan(slope)+math.Pi/3),
			}
		}
	} else {
		if p2.x - p1.x < 0 {
			pointL = Point{
				arrowHead.x + math.Sin(math.Atan(slope)+math.Pi/3),
				arrowHead.y - math.Cos(math.Atan(slope)+math.Pi/3),
			}
			pointR = Point{
				arrowHead.x - math.Sin(math.Atan(slope)-math.Pi/3),
				arrowHead.y + math.Cos(math.Atan(slope)-math.Pi/3),
			}
		} else {
			pointL = Point{
				arrowHead.x + math.Sin(math.Atan(slope)-math.Pi/3),
				arrowHead.y - math.Cos(math.Atan(slope)-math.Pi/3),
			}
			pointR = Point{
				arrowHead.x - math.Sin(math.Atan(slope)+math.Pi/3),
				arrowHead.y - math.Cos(math.Atan(slope)+math.Pi/3),
			}
		}
	}
	if _, err := canvas.Write([]byte(fmt.Sprintf("\t<line x1=\"%fem\" y1=\"%fem\" x2=\"%fem\" y2=\"%fem\" style=\"stroke:rgb(0,0,0);stroke-width:1\" />\n", pointL.x, pointL.y, arrowHead.x, arrowHead.y))); err != nil {
		log.Panic(err)
	}
	if _, err := canvas.Write([]byte(fmt.Sprintf("\t<line x1=\"%fem\" y1=\"%fem\" x2=\"%fem\" y2=\"%fem\" style=\"stroke:rgb(0,0,0);stroke-width:1\" />\n", pointR.x, pointR.y, arrowHead.x, arrowHead.y))); err != nil {
		log.Panic(err)
	}
}

func (p *Point) isInsideRect(rect ClassSpace) bool {
	return p.x <= rect.botRight.x && p.y <= rect.botRight.y && p.x >= rect.topLeft.x && p.y >= rect.topLeft.y
}

func Generate(canvas io.Writer, diagram Diagram) {
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
	classSpaces := []ClassSpace{}
	for true {
		satisfied := true
		for c := range diagram.Classes {
			topLeft := Point{float64(rand.Intn(width-classWidths[c]-2) + 1), float64(rand.Intn(height-classHeights[c]-2) + 1)}
			botRight := Point{topLeft.x + float64(classWidths[c]), topLeft.y + float64(classHeights[c])}
			classSpaces = append(classSpaces, ClassSpace{topLeft, botRight})
			log.Println(len(classSpaces))
			cont := true
			for index := 0; index < c && cont; index++ {
				log.Printf("c: %d, index: %d", c, index)
				if (classSpaces[c].topLeft.x < classSpaces[index].botRight.x+1 && classSpaces[c].topLeft.x > classSpaces[index].topLeft.x-1 && ((classSpaces[c].topLeft.y < classSpaces[index].botRight.y+1 && classSpaces[c].topLeft.y > classSpaces[index].topLeft.y-1) || (classSpaces[c].botRight.y > classSpaces[index].topLeft.y-1 && classSpaces[c].topLeft.y < classSpaces[index].botRight.y+1))) || (classSpaces[c].botRight.x > classSpaces[index].topLeft.x && classSpaces[c].topLeft.x < classSpaces[index].botRight.x+1 && ((classSpaces[c].topLeft.y < classSpaces[index].botRight.y+1 && classSpaces[c].topLeft.y > classSpaces[index].topLeft.y-1) || (classSpaces[c].botRight.y > classSpaces[index].topLeft.y-1 && classSpaces[c].topLeft.y < classSpaces[index].botRight.y+1))) {
					classSpaces = []ClassSpace{}
					cont = false
					satisfied = false
				}
			}
			if cont == false {
				break
			}
		}
		if satisfied {
			break
		}
	}
	for i := range diagram.Interactions {
		AddInteraction(canvas, classSpaces[diagram.Interactions[i].ClassAIndex], classSpaces[diagram.Interactions[i].ClassBIndex])
	}
	for c := range diagram.Classes {
		ClassGen(canvas, classSpaces[c].topLeft, diagram.Classes[c])
	}
	if _, err := canvas.Write([]byte("</svg>\n")); err != nil {
		log.Panic(err)
	}
}
