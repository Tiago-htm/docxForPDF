package main

import "encoding/xml"


type Paragraph struct {
	Properties ParagraphProperties `xml:"pPr"`
	Runs       []Run               `xml:"r"`
}
type Body struct {
	Paragraphs []Paragraph `xml:"p"`
	Sections SectionProperties `xml:"sectPr"`
}


type Document struct {
	XMLName xml.Name `xml:"document"`
	Body    Body     `xml:"body"`
}

type RunProperties struct {
	Bold   *string `xml:"b"`
	Italic *string `xml:"i"`
}

type Text struct {
	Content string `xml:",chardata"`
	Space   string `xml:"space,attr"`

}

type Run struct {
	Texts      []Text        `xml:"t"`
	Properties RunProperties `xml:"rPr"`
	Drawing    *Drawing      `xml:"drawing"`
	Tab 	   *struct{}      `xml:"tab"`
	
}

type Relationship struct {
	ID     string `xml:"Id,attr"`
	Target string `xml:"Target,attr"`
}

type Relationships struct {
	XMLName xml.Name       `xml:"Relationships"`
	Items   []Relationship `xml:"Relationship"`
}
type Justification struct {
    Val string `xml:"val,attr"`
}

type Tab struct {
	Val string `xml:"val,attr"`
	Pos float64 `xml:"pos,attr"`
}


type ParagraphProperties struct {
    Justification Justification `xml:"jc"` 
	Tabs 	   []Tab         `xml:"tabs"`
}

type SectionProperties struct {
	PageMargin struct {
		Left  float64 `xml:"left,attr"`
		Right float64 `xml:"right,attr"`
		Top   float64 `xml:"top,attr"`
		Bottom float64 `xml:"bottom,attr"`
	} `xml:"pgMar"`
}
	
// image

type Drawing struct {
	XMLName xml.Name `xml:"drawing"`
	Anchor  *Anchor  `xml:"anchor"`
	Inline  *Inline  `xml:"inline"`
}

type Inline struct {
	Graphic Graphic `xml:"graphic"`
}

type Anchor struct {
	Graphic Graphic `xml:"graphic"`
}

type Graphic struct {
	Data GraphicData `xml:"graphicData"`
}

type GraphicData struct {
	Pic Pic `xml:"pic"`
}

type Pic struct {
	BlipFill BlipFill `xml:"blipFill"`
}

type BlipFill struct {
	Blip Blip `xml:"blip"`
}

type Blip struct {
	Embed string `xml:"embed,attr"`
}
