package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// https://api.habitacao.net/graph/distrito/1629145 para encontrar no postman o preço m² médio por freguesia
// https://api.habitacao.net/graph/distrito/1629145 para encontrar no postman o preço m² médio por freguesia

// Data para consulta: 201804

var BASE_URL = "https://api.habitacao.net/graph"

var concelhos_codes = map[ string ] string {
	"Machico":         "8421411",
	"Santa Cruz":      "8421412",
	"Funchal":         "8421413",
	"Câmara de Lobos": "8421414",
	"Ribeira Brava":   "8421415",
	"Ponta do Sol":    "8421416",
	"Santana":         "8421417",
	"São Vicente":     "8421418",
	"Porto Moniz":     "8421419",
	"Calheta":         "8421420",
	"Porto Santo":     "8435154",
}

// distritos concelhos/municipio e freguesia

var prices_distritos = map[ string ] any {}
var prices_concelhos = map[ string ] map[ string ] any {}

func Get_average_price_by_distrito_by_name( name string ) any {
	return prices_distritos[ name ];
}

func Get_average_price_by_concelhos_by_name( name string, concelhos string ) any {
	return prices_concelhos[ name ][ concelhos ];
}

func Init_prices_by_concelho_and_fregusesia() {

	distritos := get_info_from_url( BASE_URL + "/distrito/1629145" )
	handle_distritos_info( distritos, prices_distritos )

	for key, value := range concelhos_codes {

		prices_concelhos[ key ] = map[ string ] any {}

		url := BASE_URL + "/concelho/" + value

		concelhos := get_info_from_url( url )
		handle_distritos_info( concelhos, prices_concelhos[ key ] )

	}

	check_for_NAN_zones()

}

func check_for_NAN_zones() {

	area := GetAdministrativeArea()

	for key_municipio, value := range area.Distritos[ "Ilha da Madeira" ].Municipios {

		_, exists := prices_distritos[ key_municipio ]
		if ! exists {
			continue
		}

		for key_freguesias := range value.Freguesias {

			_, exists := prices_concelhos[ key_municipio ][ key_freguesias ]
			if ! exists {
				prices_concelhos[ key_municipio ][ key_freguesias ] = prices_distritos[ key_municipio ]
			}

		}

	}

	// fmt.Println( area )
	// fmt.Println( prices_distritos )

}

func get_info_from_url( url string ) []map[ string ] any {

	response, err := http.Get( url )
	if err != nil {
		fmt.Println( "Error when fetching data" )
		return []map[ string ] any {}
	}

	bytes, err := io.ReadAll( response.Body )
	if err != nil {
		fmt.Println( "Error when decoding data" )
		return []map[ string ] any {}
	}
	
	var info []map[string] any;

	err = json.Unmarshal([]byte(bytes), &info)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return []map[ string ] any {}
	}

	return info

}

func handle_distritos_info( distritos_info []map[ string ] any, target_map map[ string ] any ) {

	if len( distritos_info ) == 0 {
		return;
	}


	for _, item := range distritos_info {

		for key, value := range item {

			if key == "time" {
				continue
			}

			_, exists := target_map[ key ]
			if ! exists {
				target_map[ key ] = nil;
			}

			if value == nil {
				continue
			}

			target_map[ key ] = value

		}

	}

}


// var distritoID = "1629145" // Madeira

// func Get_maneira_municipio_by_name() {

	// // Mapa no formato Concelho/Distrito/ID/
	// mapaConcelhoDistrito := make(map[string]string)
	//
	// for nome, id := range concelhos {
	// 	key := fmt.Sprintf("%s/%s/%s", nome, distritoID, id)
	// 	mapaConcelhoDistrito[nome] = key
	// }
	//
	// // Exemplo de impressão do mapa
	// for nome, path := range mapaConcelhoDistrito {
	// 	fmt.Printf("Concelho: %s => Caminho: %s\n", nome, path)
	// }
// }
