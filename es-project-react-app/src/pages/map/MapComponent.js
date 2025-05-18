import React, { useState, useEffect, useCallback } from 'react';
import { Map, Source, Layer } from 'react-map-gl/maplibre';
import maplibregl from 'maplibre-gl';
import 'maplibre-gl/dist/maplibre-gl.css';
import proj4 from 'proj4';

proj4.defs("EPSG:5016","+proj=utm +zone=28 +ellps=GRS80 +towgs84=0,0,0,0,0,0,0 +units=m +no_defs +type=crs");

const MapComponent = () => {
    const [geojson, setGeojson] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [viewState, setViewState] = useState({
        longitude: -17.00,
        latitude: 32.763,
        zoom: 9
    });
    const [admAreaData, setAdmAreaData] = useState(null);
    const [showPropertyIds, setShowPropertyIds] = useState(true);
    const [selectedDistrito, setSelectedDistrito] = useState('all');
    const [selectedMunicipio, setSelectedMunicipio] = useState('all');
    const [selectedFreguesia, setSelectedFreguesia] = useState('all');
    const [filteredPropertyIds, setFilteredPropertyIds] = useState([]);
    const convertGeoJSON = useCallback((data) => {
        if (!data?.features) return data;

        return {
            ...data,
            features: data.features.map(feature => {
                if (feature.geometry.type === 'MultiPolygon') {
                    const convertedCoordinates = feature.geometry.coordinates.map(polygon =>
                        polygon.map(ring =>
                            ring.map(([x, y]) => proj4('EPSG:5016', 'WGS84', [x, y]))
                        )
                    );
                    return {
                        ...feature,
                        geometry: {
                            ...feature.geometry,
                            coordinates: convertedCoordinates
                        }
                    };
                }
                return feature;
            })
        };
    }, []);

    // Function to fetch administrative area data
    const fetchAdmAreaData = useCallback(async () => {
        const controller = new AbortController();
        try {
            const response = await fetch('http://localhost:8080/api/adm-area', {
                signal: controller.signal
            });

            if (!response.ok) throw new Error(`HTTP ${response.status}`);

            const data = await response.json();
            setAdmAreaData(data);

            // Initialize with all property IDs
            const allPropertyIds = [];
            Object.keys(data.distrito).forEach(distrito => {
                Object.keys(data.distrito[distrito].municipio).forEach(municipio => {
                    Object.keys(data.distrito[distrito].municipio[municipio].freguesia).forEach(freguesia => {
                        allPropertyIds.push(...data.distrito[distrito].municipio[municipio].freguesia[freguesia].property_ids);
                    });
                });
            });
            setFilteredPropertyIds(allPropertyIds);

        } catch (err) {
            if (err.name !== 'AbortError') {
                console.error('Error fetching administrative area data:', err);
            }
        }
        return () => controller.abort();
    }, []);

    // Function to update filtered property IDs based on selected filters
    const updateFilteredPropertyIds = useCallback(() => {
        if (!admAreaData) return;

        let newFilteredIds = [];

        Object.keys(admAreaData.distrito).forEach(distrito => {
            if (selectedDistrito === 'all' || selectedDistrito === distrito) {
                Object.keys(admAreaData.distrito[distrito].municipio).forEach(municipio => {
                    if (selectedMunicipio === 'all' || selectedMunicipio === municipio) {
                        Object.keys(admAreaData.distrito[distrito].municipio[municipio].freguesia).forEach(freguesia => {
                            if (selectedFreguesia === 'all' || selectedFreguesia === freguesia) {
                                newFilteredIds.push(
                                    ...admAreaData.distrito[distrito].municipio[municipio].freguesia[freguesia].property_ids
                                );
                            }
                        });
                    }
                });
            }
        });

        setFilteredPropertyIds(newFilteredIds);
    }, [admAreaData, selectedDistrito, selectedMunicipio, selectedFreguesia]);

    // Effect to update filters when selection changes
    useEffect(() => {
        if (admAreaData) {
            updateFilteredPropertyIds();
        }
    }, [admAreaData, selectedDistrito, selectedMunicipio, selectedFreguesia, updateFilteredPropertyIds]);

    // Effect to fetch administrative area data
    useEffect(() => {
        fetchAdmAreaData();
    }, [fetchAdmAreaData]);

    useEffect(() => {
        const controller = new AbortController();

        const fetchData = async () => {
            try {
                setLoading(true);
                const response = await fetch('http://localhost:8080/api/prop', {
                    signal: controller.signal
                });

                const response_2 = await fetch('http://localhost:8080/api/suggestions_by_neighbours');
                if (!response.ok) throw new Error(`HTTP ${response.status}`);

                var body_2 = await response_2.json()
                
                // body_2 = body_2.map(
                //   data => data.Properties_Envolved
                // )
                var prop = body_2[ 0 ].Properties_Envolved.map( data => data.Id )
                console.log( body_2[ 0 ].Suggestion )

                const data = await response.json();

                data.features = data.features.filter(
                  data => prop.includes( data.properties.id )
                )

                const convertedData = convertGeoJSON(data);

                // console.log( data.features[ 0 ].properties.id )

                setGeojson(convertedData);

            } catch (err) {
                if (err.name !== 'AbortError') {
                    console.error('Error:', err);
                    setError(err.message);

                    // Fallback test data (Madeira coordinates)
                    setGeojson({
                        type: "FeatureCollection",
                        features: [{
                            type: "Feature",
                            geometry: {
                                type: "MultiPolygon",
                                coordinates: [
                                    [[
                                        [-17.903, 32.763], [-17.902, 32.763],
                                        [-17.902, 32.764], [-17.903, 32.764],
                                        [-17.903, 32.763]  // Closing point
                                    ]]
                                ]
                            },
                            properties: {
                                area: 439.68985,
                                freguesia: "Fajã da Ovelha",
                                id: 2671,
                                municipio: "Calheta"
                            }
                        }]
                    });
                }
            } finally {
                setLoading(false);
            }
        };

        fetchData();
        return () => controller.abort();
    }, [convertGeoJSON]);


    // Filter properties based on filteredPropertyIds
    const getFilteredGeojson = useCallback(() => {
        if (!geojson || !filteredPropertyIds.length) return geojson;

        return {
            ...geojson,
            features: geojson.features.filter(feature => 
                filteredPropertyIds.includes(feature.properties.id)
            )
        };
    }, [geojson, filteredPropertyIds]);

    // Get unique distritos from admAreaData
    const getDistritos = useCallback(() => {
        if (!admAreaData) return [];
        return Object.keys(admAreaData.distrito);
    }, [admAreaData]);

    // Get municipios for selected distrito
    const getMunicipios = useCallback(() => {
        if (!admAreaData || selectedDistrito === 'all') {
            // Get all municipios if 'all' is selected
            const allMunicipios = new Set();
            if (admAreaData) {
                Object.keys(admAreaData.distrito).forEach(distrito => {
                    Object.keys(admAreaData.distrito[distrito].municipio).forEach(municipio => {
                        allMunicipios.add(municipio);
                    });
                });
            }
            return Array.from(allMunicipios);
        }

        return Object.keys(admAreaData.distrito[selectedDistrito].municipio);
    }, [admAreaData, selectedDistrito]);

    // Get freguesias for selected municipio
    const getFreguesias = useCallback(() => {
        if (!admAreaData) return [];

        if (selectedDistrito === 'all' || selectedMunicipio === 'all') {
            // Get all freguesias if 'all' is selected for distrito or municipio
            const allFreguesias = new Set();
            Object.keys(admAreaData.distrito).forEach(distrito => {
                if (selectedDistrito === 'all' || selectedDistrito === distrito) {
                    Object.keys(admAreaData.distrito[distrito].municipio).forEach(municipio => {
                        if (selectedMunicipio === 'all' || selectedMunicipio === municipio) {
                            Object.keys(admAreaData.distrito[distrito].municipio[municipio].freguesia).forEach(freguesia => {
                                allFreguesias.add(freguesia);
                            });
                        }
                    });
                }
            });
            return Array.from(allFreguesias);
        }

        return Object.keys(admAreaData.distrito[selectedDistrito].municipio[selectedMunicipio].freguesia);
    }, [admAreaData, selectedDistrito, selectedMunicipio]);

    // Handle distrito selection change
    const handleDistritoChange = (e) => {
        const newDistrito = e.target.value;
        setSelectedDistrito(newDistrito);
        setSelectedMunicipio('all');
        setSelectedFreguesia('all');
    };

    // Handle municipio selection change
    const handleMunicipioChange = (e) => {
        setSelectedMunicipio(e.target.value);
        setSelectedFreguesia('all');
    };

    // Handle freguesia selection change
    const handleFreguesiaChange = (e) => {
        setSelectedFreguesia(e.target.value);
    };

    // Toggle property IDs visibility
    const togglePropertyIds = () => {
        setShowPropertyIds(!showPropertyIds);
    };

    const filteredGeojson = getFilteredGeojson();

    return (
        <div style={{ position: 'relative', width: '100%', height: '100%' }}>
            <Map
                {...viewState}
                onMove={evt => setViewState(evt.viewState)}
                mapLib={maplibregl}
                mapStyle="https://api.maptiler.com/maps/satellite/style.json?key=LSfqsa4izS6vDHlLAK4a"
            >
                {filteredGeojson && (
                    <Source
                        id="properties"
                        type="geojson"
                        data={filteredGeojson}
                        promoteId="id"
                    >
                        <Layer
                            id="properties-fill"
                            type="fill"
                            paint={{
                                'fill-color': [
                                    'match',
                                    ['get', 'municipio'],
                                    'Calheta', '#f28cb1',
                                    'Funchal', '#8dd3c7',
                                    '#ffcc00'  // Default color
                                ],
                                'fill-opacity': 0.7,
                                'fill-outline-color': '#000'
                            }}
                        />
                        <Layer
                            id="properties-outline"
                            type="line"
                            paint={{
                                'line-color': '#000',
                                'line-width': 1
                            }}
                        />
                        {showPropertyIds && (
                            <Layer
                                id="properties-label"
                                type="symbol"
                                layout={{
                                    'text-field': ['to-string', ['get', 'id']],
                                    'text-size': 12,
                                    'text-offset': [0, 0.5],
                                    'text-anchor': 'top'
                                }}
                                paint={{
                                    'text-color': '#000000',
                                    'text-halo-color': '#ffffff',
                                    'text-halo-width': 1
                                }}
                            />
                        )}
                    </Source>
                )}
            </Map>

            {/* Filter controls */}
            <div style={{
                position: 'absolute',
                top: 10,
                right: 10,
                backgroundColor: 'rgba(255,255,255,0.9)',
                padding: 12,
                borderRadius: 4,
                zIndex: 1,
                width: '250px',
                boxShadow: '0 0 10px rgba(0,0,0,0.1)'
            }}>
                <h3 style={{ margin: '0 0 10px 0' }}>Filters</h3>

                <div style={{ marginBottom: '10px' }}>
                    <label style={{ display: 'block', marginBottom: '5px' }}>Distrito:</label>
                    <select 
                        value={selectedDistrito} 
                        onChange={handleDistritoChange}
                        style={{ width: '100%', padding: '5px' }}
                    >
                        <option value="all">All Distritos</option>
                        {getDistritos().map(distrito => (
                            <option key={distrito} value={distrito}>{distrito}</option>
                        ))}
                    </select>
                </div>

                <div style={{ marginBottom: '10px' }}>
                    <label style={{ display: 'block', marginBottom: '5px' }}>Município:</label>
                    <select 
                        value={selectedMunicipio} 
                        onChange={handleMunicipioChange}
                        style={{ width: '100%', padding: '5px' }}
                    >
                        <option value="all">All Municípios</option>
                        {getMunicipios().map(municipio => (
                            <option key={municipio} value={municipio}>{municipio}</option>
                        ))}
                    </select>
                </div>

                <div style={{ marginBottom: '15px' }}>
                    <label style={{ display: 'block', marginBottom: '5px' }}>Freguesia:</label>
                    <select 
                        value={selectedFreguesia} 
                        onChange={handleFreguesiaChange}
                        style={{ width: '100%', padding: '5px' }}
                    >
                        <option value="all">All Freguesias</option>
                        {getFreguesias().map(freguesia => (
                            <option key={freguesia} value={freguesia}>{freguesia}</option>
                        ))}
                    </select>
                </div>

                <div>
                    <label style={{ display: 'flex', alignItems: 'center', cursor: 'pointer' }}>
                        <input 
                            type="checkbox" 
                            checked={showPropertyIds} 
                            onChange={togglePropertyIds}
                            style={{ marginRight: '8px' }}
                        />
                        Show Property IDs
                    </label>
                </div>

                {filteredPropertyIds.length > 0 && (
                    <div style={{ marginTop: '15px', fontSize: '0.9em' }}>
                        <strong>Properties:</strong> {filteredPropertyIds.length}
                    </div>
                )}
            </div>

            {/* Debug overlay */}
            <div style={{
                position: 'absolute',
                top: 10,
                left: 10,
                backgroundColor: 'rgba(255,255,255,0.9)',
                padding: 8,
                borderRadius: 4,
                zIndex: 1,
                fontFamily: 'monospace'
            }}>
                {loading && <div>🔄 Loading...</div>}
                {error && <div style={{ color: 'red' }}>❌ {error}</div>}
                {geojson && (
                    <>
                        <div>✅ {filteredGeojson?.features?.length || 0} / {geojson.features.length} properties</div>
                        <div>📍 {viewState.longitude.toFixed(5)}, {viewState.latitude.toFixed(5)}</div>
                        <div>🔍 Zoom: {viewState.zoom.toFixed(1)}</div>
                    </>
                )}
            </div>
        </div>
    );
};

export default MapComponent;
