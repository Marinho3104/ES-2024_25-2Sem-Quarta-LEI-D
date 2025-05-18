package app

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dhconnelly/rtreego"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type Property struct {
	Id        int
	Owner     int
	ShapeArea float32
	Freguesia string
	Municipio string
	Distrito  string
	Geometry  geom.MultiPolygon
	Rect      *rtreego.Rect //to build corretly RTree
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

		case 6:
			record.Owner, _ = strconv.Atoi(field)
		case 7:
			record.Freguesia = field
		case 8:
			record.Municipio = field
		case 9:
			record.Distrito = field
		default:
			panic("Unreconized field")
		}
	}

	// Do not add wrong data
	if record.ShapeArea == 0 || record.Geometry.Bounds().IsEmpty() {
		return nil, errors.New("invalid geometry")
	}

	rect, _ := BoundsToRect(&record.Geometry)
	record.Rect = rect

	return &record, nil

}

func BoundsToRect(multipolygon *geom.MultiPolygon) (*rtreego.Rect, error) {
	bounds := multipolygon.Bounds()

	minx := bounds.Min(0)
	miny := bounds.Min(1)
	maxx := bounds.Max(0)
	maxy := bounds.Max(1)

	width := maxx - minx
	height := maxy - miny

	// Create an R-tree Rect
	pt := rtreego.Point{minx, miny}
	rect, err := rtreego.NewRect(pt, []float64{width, height})
	if err != nil {
		return nil, fmt.Errorf("failed to create bounding box: %v", err)
	}

	return &rect, nil
}

func propertyEquals(a, b Property) bool {
	return a.Id == b.Id
}

func (obj *Property) Bounds() rtreego.Rect {
	return *obj.Rect
}
