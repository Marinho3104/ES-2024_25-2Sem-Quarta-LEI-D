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

// InsertProperty inserts a Property into the hierarchical structure of Districts, Municipalities, and Parishes (Freguesias).
// It initializes missing levels in the hierarchy as needed and appends the Property ID to the appropriate Freguesia.
func InsertProperty(prop Property) {
	if area.Distritos == nil {
		area.Distritos = make(map[string]*Distrito)
	}

	dist, ok := area.Distritos[prop.Distrito]
	if !ok {
		dist = &Distrito{
			Municipios: make(map[string]*Municipio),
		}
		area.Distritos[prop.Distrito] = dist
	}

	mun, ok := dist.Municipios[prop.Municipio]
	if !ok {
		mun = &Municipio{
			Freguesias: make(map[string]*Freguesia),
		}
		dist.Municipios[prop.Municipio] = mun
	}

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
