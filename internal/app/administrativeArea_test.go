package app

import (
	"testing"

	"github.com/dhconnelly/rtreego"
	"github.com/twpayne/go-geom"
)

func TestInsertProperty(t *testing.T) {
	area = AdministrativeArea{}

	mp := geom.NewMultiPolygon(geom.XY)
	rectPtr, err := rtreego.NewRect(rtreego.Point{0.0, 0.0}, []float64{1.0, 1.0})
	if err != nil {
		t.Fatalf("failed to create rect: %v", err)
	}

	prop := Property{
		Id:        101,
		Owner:     1,
		ShapeArea: 123.45,
		Freguesia: "FreguesiaA",
		Municipio: "MunicipioA",
		Distrito:  "DistritoA",
		Geometry:  *mp,
		Rect:      &rectPtr, // use pointer directly, no dereference!
	}

	InsertProperty(prop)

	a := GetAdministrativeArea()

	dist, ok := a.Distritos["DistritoA"]
	if !ok {
		t.Fatalf("Expected Distrito 'DistritoA' to exist")
	}

	mun, ok := dist.Municipios["MunicipioA"]
	if !ok {
		t.Fatalf("Expected Municipio 'MunicipioA' to exist")
	}

	freg, ok := mun.Freguesias["FreguesiaA"]
	if !ok {
		t.Fatalf("Expected Freguesia 'FreguesiaA' to exist")
	}

	if len(freg.PropertyIDs) != 1 || freg.PropertyIDs[0] != 101 {
		t.Errorf("Expected PropertyIDs to contain 101, got %v", freg.PropertyIDs)
	}
}

