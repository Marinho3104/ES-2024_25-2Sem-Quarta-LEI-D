package app

import "fmt"

// CalcOfArea computes the average area of properties at a specified administrative level (distrito, municipio, or freguesia).
// It takes a name and level as parameters, retrieves relevant properties for the level, and calculates the average area.
// Returns the computed average area and an error if no properties are found or the name is invalid.
func CalcOfArea(name string, level int) (float32, error) {
	g := GetGraph()
	var totalArea float32
	var count int

	for dName, distrito := range area.Distritos {
		if level == 0 && dName == name {
			// Distrito level
			for _, municipio := range distrito.Municipios {
				for _, freguesia := range municipio.Freguesias {
					for _, propID := range freguesia.PropertyIDs {
						prop, err := g.Vertex(propID)
						if err != nil {
							continue
						}
						totalArea += prop.ShapeArea
						count++
					}
				}
			}
			if count > 0 {
				return totalArea / float32(count), nil
			}
			return 0, fmt.Errorf("no properties found in distrito %s", name)
		}

		for mName, municipio := range distrito.Municipios {
			if level == 1 && mName == name {
				// Municipio level
				for _, freguesia := range municipio.Freguesias {
					for _, propID := range freguesia.PropertyIDs {
						prop, err := g.Vertex(propID)
						if err != nil {
							continue
						}
						totalArea += prop.ShapeArea
						count++
					}
				}
				if count > 0 {
					return totalArea / float32(count), nil
				}
				return 0, fmt.Errorf("no properties found in municipio %s", name)
			}

			for fName, freguesia := range municipio.Freguesias {
				if level == 2 && fName == name {
					// Freguesia level
					for _, propID := range freguesia.PropertyIDs {
						prop, err := g.Vertex(propID)
						if err != nil {
							continue
						}
						totalArea += prop.ShapeArea
						count++
					}
					if count > 0 {
						return totalArea / float32(count), nil
					}
					return 0, fmt.Errorf("no properties found in freguesia %s", name)
				}
			}
		}
	}

	return 0, fmt.Errorf("name '%s' not found at specified level", name)
}
