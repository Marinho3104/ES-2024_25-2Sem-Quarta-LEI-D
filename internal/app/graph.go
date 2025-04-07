package app

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func ReadFile() {

	fmt.Println("Loading file...")

	file, err := os.Open("../../assets/madeira.csv")

	if err != nil {
		fmt.Println("Error while trying to opening the file")
		fmt.Println(err)
		file.Close()
		return

	}

	fmt.Println("File loaded successfuly!!")
	fmt.Println()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	var records [][]string

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		records = append(records, record)
	}

	CreatePropertyList(records)
	file.Close()

}

func CreatePropertyList(data [][]string) {
	var propertyList []Property

	for i, line := range data {
		// Ignoring first line (header)
		if i > 0 {
			var record Property

			for j, field := range line {
				switch j {
				case 0:
					record.id, _ = strconv.Atoi(field)
				case 1:
				case 2:
				case 3:
				case 4:
					var value, _ = strconv.ParseFloat(field, 32)
					record.shapeArea = float32(value)
				case 5:
					record.geometry = field
				case 6:
					record.owner, _ = strconv.Atoi(field)
				case 7:
					record.freguesia = field
				case 8:
					record.municipio = field
				case 9:
				default:
					panic("Unreconized field")
				}
			}

			propertyList = append(propertyList, record)
		}
	}

	fmt.Println(propertyList)
}
