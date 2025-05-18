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

    const fetchAdmAreaData = useCallback(async () => {
        const controller = new AbortController();
        try {
            const response = await fetch('http://localhost:8080/api/adm-area', { signal: controller.signal });
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            const data = await response.json();
            setAdmAreaData(data);

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
            if (err.name !== 'AbortError') console.error('Error fetching administrative area data:', err);
        }
        return () => controller.abort();
    }, []);

    const updateFilteredPropertyIds = useCallback(() => {
        if (!admAreaData) return;
        let newFilteredIds = [];
        Object.keys(admAreaData.distrito).forEach(distrito => {
            if (selectedDistrito === 'all' || selectedDistrito === distrito) {
                Object.keys(admAreaData.distrito[distrito].municipio).forEach(municipio => {
                    if (selectedMunicipio === 'all' || selectedMunicipio === municipio) {
                        Object.keys(admAreaData.distrito[distrito].municipio[municipio].freguesia).forEach(freguesia => {
                            if (selectedFreguesia === 'all' || selectedFreguesia === freguesia) {
                                newFilteredIds.push(...admAreaData.distrito[distrito].municipio[municipio].freguesia[freguesia].property_ids);
                            }
                        });
                    }
                });
            }
        });
        setFilteredPropertyIds(newFilteredIds);
    }, [admAreaData, selectedDistrito, selectedMunicipio, selectedFreguesia]);

    useEffect(() => {
        if (admAreaData) updateFilteredPropertyIds();
    }, [admAreaData, selectedDistrito, selectedMunicipio, selectedFreguesia, updateFilteredPropertyIds]);

    useEffect(() => { fetchAdmAreaData(); }, [fetchAdmAreaData]);

    useEffect(() => {
        const controller = new AbortController();
        const fetchData = async () => {
            try {
                setLoading(true);
                const response = await fetch('http://localhost:8080/api/prop', { signal: controller.signal });
                if (!response.ok) throw new Error(`HTTP ${response.status}`);
                const data = await response.json();
                setGeojson(convertGeoJSON(data));
            } catch (err) {
                if (err.name !== 'AbortError') {
                    console.error('Error:', err);
                    setError(err.message);
                    setGeojson({
                        type: "FeatureCollection",
                        features: [{
                            type: "Feature",
                            geometry: {
                                type: "MultiPolygon",
                                coordinates: [[[[-17.903, 32.763], [-17.902, 32.763], [-17.902, 32.764], [-17.903, 32.764], [-17.903, 32.763]]]]
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

    const getFilteredGeojson = useCallback(() => {
        if (!geojson || !filteredPropertyIds.length) return geojson;
        return {
            ...geojson,
            features: geojson.features.filter(feature => filteredPropertyIds.includes(feature.properties.id))
        };
    }, [geojson, filteredPropertyIds]);

    const getAverageArea = useCallback(() => {
        if (!geojson || !filteredPropertyIds.length) return 0;
        const filteredFeatures = geojson.features.filter(feature => filteredPropertyIds.includes(feature.properties.id));
        if (filteredFeatures.length === 0) return 0;
        const totalArea = filteredFeatures.reduce((sum, feature) => sum + (feature.properties.area || 0), 0);
        return totalArea / filteredFeatures.length;
    }, [geojson, filteredPropertyIds]);

    const getDistritos = useCallback(() => admAreaData ? Object.keys(admAreaData.distrito) : [], [admAreaData]);
    const getMunicipios = useCallback(() => {
        if (!admAreaData || selectedDistrito === 'all') {
            const allMunicipios = new Set();
            if (admAreaData) {
                Object.values(admAreaData.distrito).forEach(d => {
                    Object.keys(d.municipio).forEach(m => allMunicipios.add(m));
                });
            }
            return Array.from(allMunicipios);
        }
        return Object.keys(admAreaData.distrito[selectedDistrito].municipio);
    }, [admAreaData, selectedDistrito]);

    const getFreguesias = useCallback(() => {
        if (!admAreaData) return [];
        if (selectedDistrito === 'all' || selectedMunicipio === 'all') {
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

    const handleDistritoChange = (e) => {
        setSelectedDistrito(e.target.value);
        setSelectedMunicipio('all');
        setSelectedFreguesia('all');
    };

    const handleMunicipioChange = (e) => {
        setSelectedMunicipio(e.target.value);
        setSelectedFreguesia('all');
    };

    const handleFreguesiaChange = (e) => {
        setSelectedFreguesia(e.target.value);
    };

    const togglePropertyIds = () => setShowPropertyIds(!showPropertyIds);
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
                    <Source id="properties" type="geojson" data={filteredGeojson} promoteId="id">
                        <Layer
                            id="properties-fill"
                            type="fill"
                            paint={{
                                'fill-color': [
                                    'match', ['get', 'municipio'],
                                    'Calheta', '#f28cb1',
                                    'Câmara de Lobos', '#66c2a5',
                                    'Funchal', '#8dd3c7',
                                    'Machico', '#ffd92f',
                                    'Ponta do Sol', '#e78ac3',
                                    'Porto Moniz', '#a6d854',
                                    'Ribeira Brava', '#fc8d62',
                                    'Santa Cruz', '#b3b3b3',
                                    'Santana', '#80b1d3',
                                    'São Vicente', '#fb8072',
                                    'Porto Santo', '#cab2d6',
                                    '#ffcc00'
                                ],
                                'fill-opacity': 0.7,
                                'fill-outline-color': '#000'
                            }}
                        />
                        <Layer id="properties-outline" type="line" paint={{ 'line-color': '#000', 'line-width': 1 }} />
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

            {/* Filters Panel */}
            <div style={{
                position: 'absolute', top: 10, right: 10, backgroundColor: 'rgba(255,255,255,0.9)',
                padding: 12, borderRadius: 4, zIndex: 1, width: '250px', boxShadow: '0 0 10px rgba(0,0,0,0.1)'
            }}>
                <h3 style={{ margin: '0 0 10px 0' }}>Filters</h3>
                <div style={{ marginBottom: '10px' }}>
                    <label>Distrito:</label>
                    <select value={selectedDistrito} onChange={handleDistritoChange} style={{ width: '100%' }}>
                        <option value="all">All Distritos</option>
                        {getDistritos().map(d => <option key={d} value={d}>{d}</option>)}
                    </select>
                </div>
                <div style={{ marginBottom: '10px' }}>
                    <label>Município:</label>
                    <select value={selectedMunicipio} onChange={handleMunicipioChange} style={{ width: '100%' }}>
                        <option value="all">All Municípios</option>
                        {getMunicipios().map(m => <option key={m} value={m}>{m}</option>)}
                    </select>
                </div>
                <div style={{ marginBottom: '15px' }}>
                    <label>Freguesia:</label>
                    <select value={selectedFreguesia} onChange={handleFreguesiaChange} style={{ width: '100%' }}>
                        <option value="all">All Freguesias</option>
                        {getFreguesias().map(f => <option key={f} value={f}>{f}</option>)}
                    </select>
                </div>
                <div>
                    <label style={{ display: 'flex', alignItems: 'center' }}>
                        <input type="checkbox" checked={showPropertyIds} onChange={togglePropertyIds} style={{ marginRight: '8px' }} />
                        Show Property IDs
                    </label>
                </div>
                {filteredPropertyIds.length > 0 && (
                    <>
                        <div style={{ marginTop: '15px', fontSize: '0.9em' }}>
                            <strong>Properties:</strong> {filteredPropertyIds.length}
                        </div>
                        <div style={{ marginTop: '5px', fontSize: '0.9em' }}>
                            <strong>Avg. Area:</strong> {getAverageArea().toFixed(2)} m²
                        </div>
                    </>
                )}
            </div>

            {/* Debug Overlay */}
            <div style={{
                position: 'absolute', top: 10, left: 10, backgroundColor: 'rgba(255,255,255,0.9)',
                padding: 8, borderRadius: 4, zIndex: 1, fontFamily: 'monospace'
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
