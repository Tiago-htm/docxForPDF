package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
    "strings"
	"path/filepath"

     //externas
    "github.com/jung-kurt/gofpdf"

)


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
	Bold *string `xml:"b"`
	Italic *string `xml:"i"`
}

type Text struct {
	Content string `xml:",chardata"`
	Space 	string `xml:"space,attr"`
}


type Run struct {
	Texts []Text `xml:"t"`
	Properties RunProperties `xml:"rPr"`
	Drawing * Drawing `xml:"drawing"`
}


type Relationship struct{
	ID 	string `xml:"Id,attr"`
	Target string `xml:"Target,attr"`
}

type Relationships struct{
	XMLName	xml.Name `xml:"Relationships"`
	Items []Relationship `xml:"Relationship"`
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

func getRelationships(r *zip.ReadCloser) (map[string]string, error) {
    relMap := make(map[string]string) 

    for _, f := range r.File { 
        if f.Name == "word/_rels/document.xml.rels" {
            rc, err := f.Open()
            if err != nil {
                return nil, err
            }
            defer rc.Close()

            var rels Relationships
            err = xml.NewDecoder(rc).Decode(&rels)
            if err != nil {
                return nil, err
            }

            for _, rel := range rels.Items {
                target := rel.Target
              
                if !strings.HasPrefix(target, "word/") && !strings.HasPrefix(target, "/") {
                    target = "word/" + target
                }
                relMap[rel.ID] = target
            }
        }
    }
    return relMap, nil
}




func main() {
	pdf :=  gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()


	r, err := zip.OpenReader("1.docx")
	if err != nil {
		log.Fatal(err, `err ao ler arquivo`)
	}
	defer r.Close()

	relacoes, err := getRelationships(r)


	if err != nil {
		log.Fatal("Erro ao carregar relações:", err)
	}

	// re :=  regexp.MustCompile(`^word/media/.*\.(png|jpe?g|gif)$`)

	mapaFicheiros := make(map[string]*zip.File)

	for _, f :=  range r.File {
		mapaFicheiros[f.Name] = f
	}

	for _, f := range r.File {

		// if re.MatchString(f.Name){
		// 	imgFile, err := f.Open()
			
		// 	if err != nil {
		// 		log.Println("Error ao pegar imagem", err)
		// 		continue
		// 	}

		// 	info :=  pdf.RegisterImageOptionsReader(f.Name, gofpdf.ImageOptions{ImageType: "png"}, imgFile)
			
		//      imgFile.Close()

		// 	if info == nil{
		// 		log.Println("Erro ao registrar", err)
		// 		continue
		// 	}
		// 	info =  pdf.GetImageInfo(f.Name)

		
		// 	_,alturaCalculada := info.Extent()
		// 	larguraDesejada := 100.0 

		// 	_, alturaCalculada = pdf.GetImageInfo(f.Name).Extent()
		
		// 	x, y := pdf.GetX(), pdf.GetY()
		
		// 	pdf.ImageOptions(f.Name, x, y, larguraDesejada, 0, false, gofpdf.ImageOptions{ImageType: ""},0,"")


		// 	pdf.RegisterImageOptionsReader(f.Name, gofpdf.ImageOptions{ImageType: "png"}, imgFile)

		// 	pdf.SetY(y + alturaCalculada + 5) 

		// }
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
    
    						//idProcurado := run.Drawing.RID
							getDrawingId := run.Drawing.Inline.Graphic.Data.Pic.BlipFill.Blip.Embed
							fmt.Println("Tentando processar imagemID: ", getDrawingId)	
							if nameArq, isRel := relacoes[getDrawingId]; isRel{
								fmt.Println("Caminho: ", nameArq)

								if arqZip, existZip := mapaFicheiros[nameArq]; existZip{
									extensao := filepath.Ext(nameArq)
									tipoImg :=  ""
									if len(extensao) > 1 {
										tipoImg = extensao[1:]
									}

									imgFile, err := arqZip.Open()
									if err == nil {
											info := pdf.RegisterImageOptionsReader(nameArq, gofpdf.ImageOptions{ImageType: tipoImg}, imgFile)
											imgFile.Close()
											if info == nil {
												fmt.Println("Não foi possivel ler o nome do arquivo", nameArq)
											} else {

												w, h := info.Extent()
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

	// arquivo, err := os.Create("palavras.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer arquivo.Close()

	// for _, linha := range linhas {
	// 	if linha != "" {
	// 		_, err := arquivo.WriteString(linha + "\n")
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	}
	// }

	err = pdf.OutputFileAndClose("arquivo.pdf")
	if err != nil {
		panic(err)
	}

	fmt.Println("Escrita terminada....")
}