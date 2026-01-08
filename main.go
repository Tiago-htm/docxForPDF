package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"path/filepath"

	//externas
	"github.com/jung-kurt/gofpdf"
)

func main() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	r, err := zip.OpenReader("computa.docx")
	if err != nil {
		log.Fatal(err, `err ao ler arquivo`)
	}
	defer r.Close()

	relacoes, err := GetRelationships(r)

	if err != nil {
		log.Fatal("Erro ao carregar relações:", err)
	}

	// re :=  regexp.MustCompile(`^word/media/.*\.(png|jpe?g|gif)$`)

	mapaFicheiros := make(map[string]*zip.File)

	for _, f := range r.File {
		mapaFicheiros[f.Name] = f
	}

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
				for _, run := range para.Runs {

					if run.Drawing != nil && run.Drawing.Inline != nil {
						getDrawingId := run.Drawing.Inline.Graphic.Data.Pic.BlipFill.Blip.Embed
						fmt.Println("Tentando processar imagemID: ", getDrawingId)
						if nameArq, isRel := relacoes[getDrawingId]; isRel {
							fmt.Println("Caminho: ", nameArq)
							if arqZip, existZip := mapaFicheiros[nameArq]; existZip {
								extensao := filepath.Ext(nameArq)
								tipoImg := GetImgType(extensao)
								imgFile, err := arqZip.Open()
								if err == nil {
									info := pdf.RegisterImageOptionsReader(nameArq, gofpdf.ImageOptions{ImageType: tipoImg}, imgFile)
									imgFile.Close()
									if info == nil {
										fmt.Println("Não foi possivel ler o nome do arquivo", nameArq)
									} else {

										_, h := info.Extent()
										x, y := pdf.GetX(), pdf.GetY()
										pdf.ImageOptions(nameArq, x, y, 80.0, 0, false, gofpdf.ImageOptions{ImageType: tipoImg}, 0, "")
										pdf.SetY(y + h + 5)
									}

								}
							}

						}

					}

					style := ""

					if run.Properties.Bold != nil {
						style += "B"
					}

					if run.Properties.Italic != nil {
						style += "I"
					}

					pdf.SetFont("Arial", style, 12)

					for _, t := range run.Texts {
						pdf.Write(5, t.Content)
					}
				}

				pdf.Ln(7)

			}
		}
	}

	err = pdf.OutputFileAndClose("arquivo.pdf")
	if err != nil {
		panic(err)
	}

	fmt.Println("Escrita terminada....")
}
