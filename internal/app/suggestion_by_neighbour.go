package app

import (
	"fmt"
	"math"
)


type Suggestions_By_Neighbours struct {

	Suggestion [ 2 ] Property
	Properties_Envolved [ 4 ] Property

}

var Suggestions_point_6 []Suggestions_By_Neighbours;

func Init_all_suggestions_by_neighbour() {

	val, _ := GetOwnerGraph().Edges()

	for _, value := range val {

		cast := value.Properties.Data.(map[int]struct{})

		if len( cast ) < 4 {
			continue
		}

		handle_suggestions_group( cast )

	}

	fmt.Println( len( Suggestions_point_6 ) )

}

func isWithin10Percent(a, b float64) bool {
	diff := math.Abs(a - b)
	average := (a + b) / 2
	return diff/average < 0.10
}

func handle_suggestions_group( group map[ int ] struct {} ) {

	graph := GetGraph()

	keys := make([]int, 0, len(group))
	for k, _ := range group {
		keys = append(keys, k)
	}

	nei := [][]Property{}

	for index, _ := range keys {

		// Inner loop: from start up to current index
		for j := index + 1; j < len( keys ); j++ {

			edge, err := graph.Edge( keys[ index ], keys[ j ] )

			if err != nil {
				continue
			}

			arr := make([]Property, 2)
			arr[ 0 ] = edge.Source
			arr[ 1 ] = edge.Target

			nei = append( nei, arr )

		}
	}

	for index, neighbours := range nei {

		source := neighbours[ 0 ]
		target := neighbours[ 1 ]

		if( source.Owner == target.Owner ) {
			continue
		}

		for j := index + 1; j < len( neighbours ); j++ {
	
			source_2 := nei[ j ][ 0 ]
			target_2 := nei[ j ][ 1 ]		

			if( source_2.Owner == target_2.Owner ) {
				continue
			}

			if source.Owner != source_2.Owner && source.Owner != target_2.Owner {
				continue
			}

			if target.Owner != source_2.Owner && target.Owner != target_2.Owner {
				continue
			}

			final_1 := make([]Property, 2)
			final_2 := make([]Property, 2)

			if source.Owner != source_2.Owner {

				final_1[ 0 ] = source
				final_1[ 1 ] = source_2

				final_2[ 0 ] = target
				final_2[ 1 ] = target_2

			} else {

				final_1[ 0 ] = source
				final_1[ 1 ] = target_2

				final_2[ 0 ] = target
				final_2[ 1 ] = source_2


			}

			if isWithin10Percent( float64(final_1[ 0 ].ShapeArea), float64(final_1[ 1 ].ShapeArea) ) {

				suggestions := make( []Property, 2 )
				suggestions[ 0 ] = final_1[ 0 ]
				suggestions[ 1 ] = final_1[ 1 ]

				envolved := make( []Property, 4 )
				envolved[ 0 ] = final_1[ 0 ]
				envolved[ 1 ] = final_1[ 1 ]
				envolved[ 2 ] = final_2[ 0 ]
				envolved[ 3 ] = final_2[ 1 ]

				Suggestions_point_6 = append( Suggestions_point_6, Suggestions_By_Neighbours{ Suggestion: [2]Property(suggestions), Properties_Envolved: [4]Property(envolved) } )

			} 

			if isWithin10Percent( float64(final_2[ 0 ].ShapeArea), float64(final_2[ 1 ].ShapeArea) ) {

				suggestions := make( []Property, 2 )
				suggestions[ 0 ] = final_2[ 0 ]
				suggestions[ 1 ] = final_2[ 1 ]

				envolved := make( []Property, 4 )
				envolved[ 0 ] = final_1[ 0 ]
				envolved[ 1 ] = final_1[ 1 ]
				envolved[ 2 ] = final_2[ 0 ]
				envolved[ 3 ] = final_2[ 1 ]

				Suggestions_point_6 = append( Suggestions_point_6, Suggestions_By_Neighbours{ Suggestion: [2]Property(suggestions), Properties_Envolved: [4]Property(envolved) } )

			}

		}

	}

}
