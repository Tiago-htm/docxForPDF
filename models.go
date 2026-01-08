package main

import "encoding/xml"

type Paragraph struct {
	Runs []Run `xml:"r"`
}
type Body struct {
	Paragraphs []Paragraph `xml:"p"`
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
}

type Relationship struct {
	ID     string `xml:"Id,attr"`
	Target string `xml:"Target,attr"`
}

type Relationships struct {
	XMLName xml.Name       `xml:"Relationships"`
	Items   []Relationship `xml:"Relationship"`
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
