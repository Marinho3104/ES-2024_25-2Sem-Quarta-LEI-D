package app

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

var globalPropertyList []Property

// readFile reads data from a CSV file, parses it, and returns the records as a slice of string slices or an error.
func readFile() ([][]string, error) {

	fmt.Println("Loading file...")

	file, err := os.Open("../../assets/madeira_corrected.csv")

	if err != nil {
		return nil, err
	}
	defer file.Close()

	fmt.Println("File loaded successfuly!!")

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

	return records, nil
}

func GetPropertigesList() []Property {
	if globalPropertyList == nil {
		globalPropertyList = createPropertyList()
	}
	return globalPropertyList
}

// createPropertyList reads property data from a file and converts it into a list of Property objects while skipping invalid entries.
func createPropertyList() []Property {
	a := 0
	var propertyList []Property
	data, err := readFile()
	if err != nil {
		fmt.Println("Erro trying reading the file")
		fmt.Println(err)
		return nil
	}
	for i, line := range data {
		if i > 0 {
			prop, err := createProperty(line)
			if err != nil {
				continue
			}
			propertyList = append(propertyList, *prop)
			if prop.Freguesia == "NA" {
				a += 1
			}
			InsertProperty(*prop)
		}

	}
	fmt.Println("Total of invalid entries: ", a)
	return propertyList
}
