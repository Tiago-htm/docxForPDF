package main

import (
	"archive/zip"
	"encoding/xml"
	"path/filepath"
	"strings"
)

func GetRelationships(r *zip.ReadCloser) (map[string]string, error) {
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

func GetImgType(name string) string {
	ext := filepath.Ext(name)
	if len(ext) > 1 {
		return ext[1:]
	}

	return ""
}


func ConvertTwipsToMM(twips float64) float64 {
	const twipsPerInch = 1440.0
	const mmPerInch = 25.4
	return (twips / twipsPerInch) * mmPerInch
}

