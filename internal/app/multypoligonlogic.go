package app

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
	"github.com/twpayne/go-geom/xy/lineintersector"
)

func areMultiPolygonsNeighbors(mp1, mp2 *geom.MultiPolygon) bool {

	for i := 0; i < mp1.NumPolygons(); i++ {
		poly1 := mp1.Polygon(i)
		for j := 0; j < mp2.NumPolygons(); j++ {
			poly2 := mp2.Polygon(j)
			if doPolygonsTouch(poly1, poly2) {
				return true
			}
		}
	}

	return false
}

func doPolygonsTouch(poly1, poly2 *geom.Polygon) bool {
	if hasSharedBoundary(poly1, poly2) {
		return true
	}
	return false
}

func hasSharedBoundary(poly1, poly2 *geom.Polygon) bool {
	edges1 := extractEdges(poly1)
	edges2 := extractEdges(poly2)

	for _, e1 := range edges1 {
		for _, e2 := range edges2 {
			if xy.Distance(e1.start, e2.start) > 50 {
				if xy.Distance(e1.end, e2.end) > 50 {
					continue
				}

			}

			if !xy.DoLinesOverlap(e1.start, e1.end, e2.start, e2.end) {
				continue
			}
			result := lineintersector.LineIntersectsLine(lineintersector.RobustLineIntersector{}, e1.start, e1.end, e2.start, e2.end)
			if result.HasIntersection() {
				return true
			}

		}
	}

	return false
}

func extractEdges(poly *geom.Polygon) []struct{ start, end geom.Coord } {
	var edges []struct{ start, end geom.Coord }
	for i := 0; i < poly.NumLinearRings(); i++ {
		ring := poly.LinearRing(i)
		coords := ring.Coords()
		for j := 0; j < len(coords)-1; j++ {
			edges = append(edges, struct{ start, end geom.Coord }{
				start: coords[j],
				end:   coords[j+1],
			})
		}
	}
	return edges
}
