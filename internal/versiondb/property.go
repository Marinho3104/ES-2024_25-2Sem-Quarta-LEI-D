package versiondb

import (
	"errors"
	"strconv"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type Property struct {
	Id        int
	Owner     int
	ShapeArea float32
	Freguesia string
	Municipio string
	Geometry  geom.MultiPolygon
}

func createProperty(property []string) (*Property, error) {

	var record Property

	for j, field := range property {
		switch j {
		case 0:
			record.Id, _ = strconv.Atoi(field)
		case 1:
		case 2:
		case 3:
		case 4:
			var value, _ = strconv.ParseFloat(field, 32)
			record.ShapeArea = float32(value)
		case 5:

			convertedField, err := wkt.Unmarshal(field)

			if err != nil {
				break
			}

			record.Geometry =
				*geom.NewMultiPolygonFlat(convertedField.Layout(), convertedField.FlatCoords(), convertedField.Endss())

			geojson.Marshal(record.Geometry.Clone())
		case 6:
			record.Owner, _ = strconv.Atoi(field)
		case 7:
			record.Freguesia = field
		case 8:
			record.Municipio = field
		case 9:
		default:
			panic("Unreconized field")
		}
	}

	// Do not add wrong data
	if record.ShapeArea == 0 || record.Geometry.Bounds().IsEmpty() {
		return nil, errors.New("invalid geometry")
	}
	return &record, nil

}
