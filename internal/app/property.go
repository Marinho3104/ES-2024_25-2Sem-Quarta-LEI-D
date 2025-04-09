package app

import "github.com/twpayne/go-geom"

type Property struct {
	id        int
	owner     int
	shapeArea float32
	freguesia string
	municipio string
	geometry  geom.MultiPolygon
}

func (p Property) Geometry() *geom.MultiPolygon {
	return p.geometry.Clone()
}
