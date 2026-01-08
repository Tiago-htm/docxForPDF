package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"

	//externas
	"github.com/jung-kurt/gofpdf"
)

func main() {

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.SetFont("Arial", "", 12)

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.AddPage()

	r, err := zip.OpenReader("teste.docx")
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

			var doc Document
			err = xml.Unmarshal(content, &doc)
			if err != nil {
				log.Fatal(err)
			}

			for _, para := range doc.Body.Paragraphs {
				for _, run := range para.Runs {

					if run.Drawing != nil {
						var idProcurado string

						if run.Drawing.Inline != nil {
							idProcurado = run.Drawing.Inline.Graphic.Data.Pic.BlipFill.Blip.Embed
						} else if run.Drawing.Anchor != nil {
							idProcurado = run.Drawing.Anchor.Graphic.Data.Pic.BlipFill.Blip.Embed
						}

						if idProcurado != "" {
							if nameArq, isRel := relacoes[idProcurado]; isRel {
								if arqZip, existZip := mapaFicheiros[nameArq]; existZip {

									tipoImg := GetImgType(nameArq)

									imgFile, err := arqZip.Open()
									if err == nil {
										info := pdf.RegisterImageOptionsReader(nameArq, gofpdf.ImageOptions{ImageType: tipoImg}, imgFile)
										imgFile.Close()

										if info != nil {
											_, h := info.Extent()
											x, y := pdf.GetX(), pdf.GetY()

											if y+h > 275 {
												pdf.AddPage()
												y = pdf.GetY()
											}

											pdf.ImageOptions(nameArq, x, y, 80.0, 0, false, gofpdf.ImageOptions{ImageType: tipoImg}, 0, "")
											pdf.SetY(y + h + 5)
										}
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

						pdf.Write(5, tr(t.Content))
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
