package app

import "fmt"

// https://api.habitacao.net/graph/concelho/id_concelho para encontrar no postman o preço m² médio por freguesia

func main() {
	distritoID := "1629145" // Madeira

	concelhos := map[string]string{
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

	// Mapa no formato Concelho/Distrito/ID/
	mapaConcelhoDistrito := make(map[string]string)

	for nome, id := range concelhos {
		key := fmt.Sprintf("%s/%s/%s", nome, distritoID, id)
		mapaConcelhoDistrito[nome] = key
	}

	// Exemplo de impressão do mapa
	for nome, path := range mapaConcelhoDistrito {
		fmt.Printf("Concelho: %s => Caminho: %s\n", nome, path)
	}
}
