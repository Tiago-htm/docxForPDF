package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Body struct {
	Paragraphs []Paragraph `xml:"p"`
}
type Document struct {
	XMLName xml.Name `xml:"document"`
	Body    Body     `xml:"body"`
}

type Text struct {
	Content string `xml:",chardata"`
}

type Paragraph struct {
	Texts []Text `xml:"r>t"`
}

func main() {
	r, err := zip.OpenReader("computa.docx")
	if err != nil {
		log.Fatal(err, `err ao ler aquivo`)

	}

	defer r.Close()

	var todasPalavras []string

	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}

			content, err := io.ReadAll(rc)
			rc.Close()

			if err != nil {
				log.Fatal(err)
			}
			var doc Document
			err = xml.Unmarshal(content, &doc)

			if err != nil {
				log.Fatal(err)
			}

			for _, para := range doc.Body.Paragraphs {
				for _, text := range para.Texts {
					palavras := strings.Fields(text.Content)
					todasPalavras = append(todasPalavras, palavras...)
				}
			}

		}
	}

	arquivo, err := os.Create("palavras.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer arquivo.Close()

	for _, palavra := range todasPalavras {
		_, err := arquivo.WriteString(palavra + "\n")
		if err != nil {
			log.Fatal(err)
		}

	}

	fmt.Println("Escrita terminada....")

}
