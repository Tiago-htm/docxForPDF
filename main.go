package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"time"
	"strings"

	//externas
	"github.com/jung-kurt/gofpdf"
)

func main() {
	inicio := time.Now()


	pdf := gofpdf.New("P", "mm", "A4", "")


	
	pdf.SetFont("Arial", "", 12)
	pdf.AddPage()

	tr := pdf.UnicodeTranslatorFromDescriptor("")
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
			
			var doc Document
			err = xml.Unmarshal(content, &doc)
			margens := doc.Body.Sections.PageMargin


			
			pdf.SetMargins(ConvertTwipsToMM(margens.Left), ConvertTwipsToMM(margens.Top), ConvertTwipsToMM(margens.Right))
			pdf.SetAutoPageBreak(true, ConvertTwipsToMM(margens.Bottom))

			if err != nil {
				log.Fatal(err)
			}
			
			for _, para := range doc.Body.Paragraphs {
				var align string
				switch para.Properties.Justification.Val {
						case "center": align = "C"
						case "right": align = "R"
						case "both": align = "J"
						default: align = "L"
				}

				widthMargin := 210 - ConvertTwipsToMM(float64(doc.Body.Sections.PageMargin.Left)) - ConvertTwipsToMM(float64(doc.Body.Sections.PageMargin.Right))

				var paragraphCompleted string

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

				if run.Tab != nil {
					paragraphCompleted += "       " 
				}
					
					for _, t := range run.Texts  {
						paragraphCompleted += t.Content
					}
				}
				textoFinal := strings.TrimSpace(paragraphCompleted)
				if len(textoFinal) == 0 {
					pdf.Ln(5)
					
				}
				
				pdf.MultiCell(widthMargin, 5, tr(paragraphCompleted), "", align, false)

				pdf.Ln(2)
			}
		}
	}

	err = pdf.OutputFileAndClose("arquivo.pdf")
	if err != nil {
		panic(err)
	}
	duracao := time.Since(inicio)
	fmt.Printf("Tempo decorrido: %v\n", duracao)
	fmt.Println("Escrita terminada....")

}
