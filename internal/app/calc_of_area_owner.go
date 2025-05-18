package app

import "fmt"

// CalcOfAreaOwner calculates the average area of properties owned by a given owner based on the specified administrative level.
// It takes the owner's name and level as inputs and returns the calculated average area or an error if the operation fails.
func CalcOfAreaOwner(name string, level int) (float32, error) {
	a := agregateProprities()
	fmt.Println("A: ", a)
	return 0, nil
}

func agregateProprities() float32 {
	i := 0
	totalArea := float32(0.0)
	g := GetGraph()
	o := GetOwnerGraph()
	adjMap, _ := g.AdjacencyMap()
	c := 0

	//ownerAdjMap, _ := o.AdjacencyMap()
	cleared := make(map[int]struct{})
	for v := range adjMap {
		if _, ok := cleared[v]; ok {
			fmt.Println("colided ", v)
			c += 1
			continue
		}
		p, _ := g.Vertex(v)
		edge, err := o.Edge(p.Owner, p.Owner)
		if err != nil {
			totalArea += p.ShapeArea
			i++
			continue
		} else {

			m := edge.Properties.Data.(map[int]struct{})

			localtotalArea := float32(0.0)
			fmt.Println("----------------_")
			fmt.Println("current ", v)
			fmt.Println("----------------_")
			for key, _ := range m {
				lp, _ := g.Vertex(key)
				localtotalArea += lp.ShapeArea
				cleared[key] = struct{}{}
				fmt.Println("key ", key)
			}
			totalArea += localtotalArea / float32(len(m))
			i++

		}

	}
	fmt.Println(c)
	return totalArea / float32(i)
}
