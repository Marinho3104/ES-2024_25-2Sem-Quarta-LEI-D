package app

import (
	"testing"

	"github.com/twpayne/go-geom"
)

func TestAreMultiPolygonsNeighbors_Touching(t *testing.T) {
	poly1 := geom.NewPolygonFlat(geom.XY, []float64{
		0, 0, 2, 0, 2, 2, 0, 2, 0, 0, // square
	}, []int{10})

	poly2 := geom.NewPolygonFlat(geom.XY, []float64{
		2, 0, 4, 0, 4, 2, 2, 2, 2, 0, // adjacent square (touching at edge)
	}, []int{10})

	mp1 := geom.NewMultiPolygon(geom.XY)
	_ = mp1.Push(poly1)
	mp2 := geom.NewMultiPolygon(geom.XY)
	_ = mp2.Push(poly2)

	if !areMultiPolygonsNeighbors(mp1, mp2) {
		t.Errorf("Expected touching polygons to be neighbors")
	}
}

func TestAreMultiPolygonsNeighbors_NotTouching(t *testing.T) {
	poly1 := geom.NewPolygonFlat(geom.XY, []float64{
		0, 0, 1, 0, 1, 1, 0, 1, 0, 0,
	}, []int{10})

	poly2 := geom.NewPolygonFlat(geom.XY, []float64{
		5, 5, 6, 5, 6, 6, 5, 6, 5, 5, // far away
	}, []int{10})

	mp1 := geom.NewMultiPolygon(geom.XY)
	_ = mp1.Push(poly1)
	mp2 := geom.NewMultiPolygon(geom.XY)
	_ = mp2.Push(poly2)

	if areMultiPolygonsNeighbors(mp1, mp2) {
		t.Errorf("Expected non-touching polygons NOT to be neighbors")
	}
}

