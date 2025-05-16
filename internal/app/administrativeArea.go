package app

type AdministrativeArea struct {
	Distritos map[string]*Distrito `json:"distrito"`
}

type Distrito struct {
	Municipios map[string]*Municipio `json:"municipio"`
}

type Municipio struct {
	Freguesias map[string]*Freguesia `json:"freguesia"`
}

type Freguesia struct {
	PropertyIDs []int `json:"property_ids"`
}

var area AdministrativeArea

func InsertProperty(prop Property) {
	// Ensure Distrito map exists
	if area.Distritos == nil {
		area.Distritos = make(map[string]*Distrito)
	}

	// Get or create Distrito
	
	dist, ok := area.Distritos[prop.Distrito]
	if !ok {
		dist = &Distrito{
			Municipios: make(map[string]*Municipio),
		}
		area.Distritos[prop.Distrito] = dist
	}

	// Get or create Municipio
	mun, ok := dist.Municipios[prop.Municipio]
	if !ok {
		mun = &Municipio{
			Freguesias: make(map[string]*Freguesia),
		}
		dist.Municipios[prop.Municipio] = mun
	}

	// Get or create Freguesia
	freg, ok := mun.Freguesias[prop.Freguesia]
	if !ok {
		freg = &Freguesia{
			PropertyIDs: []int{},
		}
		mun.Freguesias[prop.Freguesia] = freg

	}
	freg.PropertyIDs = append(freg.PropertyIDs, prop.Id)
}

func GetAdministrativeArea() AdministrativeArea {
	return area
}
