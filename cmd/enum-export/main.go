package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
)

type EnumSet struct {
	Name  string
	Enums []string
}

func toStringSlice[T dto.Enum](slice []T) []string {
	var strSlice = make([]string, len(slice))
	for i, s := range slice {
		strSlice[i] = s.String()
	}

	return strSlice
}

func main() {
	var enums = []EnumSet{
		{Name: "JobOrdinalAllValues", Enums: toStringSlice(dto.JobOrdinalAllValues)},
		{Name: "MediaOrdinalAllValues", Enums: toStringSlice(dto.MediaOrdinalAllValues)},
		{Name: "JobStatusAllValues", Enums: toStringSlice(model.JobStatusEnumAllValues)},
		{Name: "JobTypeAllValues", Enums: toStringSlice(model.JobTypeEnumAllValues)},
		{Name: "MediaTypeAllValues", Enums: toStringSlice(model.MediaTypeEnumAllValues)},
		{Name: "MediaRelationTypeAllValues", Enums: toStringSlice(model.MediaRelationTypeEnumAllValues)},
		{Name: "WSTopicAllValues", Enums: toStringSlice(dto.WSTopicAllValues)},
	}

	lines := []string{}
	for _, e := range enums {
		log.Printf("Generating type for %v", e.Name)
		lines = append(lines, fmt.Sprintf(
			`export type %v = "%v"`,
			e.Name,
			strings.Join(e.Enums, `" | "`),
		))
	}

	log.Print("Creating enum.d.ts")
	f, err := os.Create("./ts/enum.d.ts")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	text := strings.Join(lines, "\n")

	log.Print("Writing content to enum.d.ts")
	_, err = f.WriteString(text)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Done generating enums")
}
