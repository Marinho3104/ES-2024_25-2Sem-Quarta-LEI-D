package tests

import "github.com/twpayne/go-geom"

type Property struct {
	id        int
	owner     int
	shapeArea float32
	freguesia string
	municipio string
	geometry geom.MultiPolygon 
}
