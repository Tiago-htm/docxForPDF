package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
)

func main() {
	r, error := zip.OpenReader("computa.docx")
	if error != nil {
		log.Fatal(error, `Error ao ler aquivo`)

	}

	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, error := f.Open()
			if error != nil {
				log.Fatal(error)
			}

			content, err := io.ReadAll(rc)
			rc.Close()

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(content))
		}
	}

}
