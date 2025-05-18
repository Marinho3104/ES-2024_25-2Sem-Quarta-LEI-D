import React, { useState, useEffect, useCallback } from 'react';
import { Map, Source, Layer } from 'react-map-gl/maplibre';
import maplibregl from 'maplibre-gl';
import 'maplibre-gl/dist/maplibre-gl.css';
import proj4 from 'proj4';

// Define custom projection
proj4.defs(
    'EPSG:5016',
    '+proj=utm +zone=28 +ellps=GRS80 +towgs84=0,0,0,0,0,0,0 +units=m +no_defs +type=crs'
);

const MapComponent = () => {
    // Map and data states
    const [geojson, setGeojson] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [viewState, setViewState] = useState({ longitude: -17.0, latitude: 32.763, zoom: 9 });

    // Administrative area filters
    const [admAreaData, setAdmAreaData] = useState(null);
    const [selectedDistrito, setSelectedDistrito] = useState('all');
    const [selectedMunicipio, setSelectedMunicipio] = useState('all');
    const [selectedFreguesia, setSelectedFreguesia] = useState('all');
    const [filteredPropertyIds, setFilteredPropertyIds] = useState([]);

    // Suggestion groups state
    const [suggestionGroups, setSuggestionGroups] = useState([]);
    const [activeGroupIndex, setActiveGroupIndex] = useState(null);
    const [showPropertyIds, setShowPropertyIds] = useState(true);

    // Convert coordinates from EPSG:5016 to WGS84
    const convertGeoJSON = useCallback((data) => {
        if (!data?.features) return data;
        return {
            ...data,
            features: data.features.map((feature) => {
                if (feature.geometry.type === 'MultiPolygon') {
                    const convertedCoordinates = feature.geometry.coordinates.map((polygon) =>
                        polygon.map((ring) =>
                            ring.map(([x, y]) => proj4('EPSG:5016', 'WGS84', [x, y]))
                        )
                    );
                    return { ...feature, geometry: { ...feature.geometry, coordinates: convertedCoordinates } };
                }
                return feature;
            }),
        };
    }, []);

    // Fetch admin area data
    const fetchAdmAreaData = useCallback(async () => {
        const ctrl = new AbortController();
        try {
            const resp = await fetch('http://localhost:8080/api/adm-area', { signal: ctrl.signal });
            if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
            const data = await resp.json();
            setAdmAreaData(data);
            // init all IDs
            const allIds = [];
            Object.values(data.distrito).forEach(d =>
                Object.values(d.municipio).forEach(m =>
                    Object.values(m.freguesia).forEach(f => allIds.push(...f.property_ids))
                )
            );
            setFilteredPropertyIds(allIds);
        } catch (e) {
            if (e.name !== 'AbortError') console.error(e);
        }
        return () => ctrl.abort();
    }, []);

    // Update filtered IDs by admin filters
    const updateFilteredPropertyIds = useCallback(() => {
        if (!admAreaData) return;
        let ids = [];
        Object.entries(admAreaData.distrito).forEach(([dKey, d]) => {
            if (selectedDistrito === 'all' || selectedDistrito === dKey) {
                Object.entries(d.municipio).forEach(([mKey, m]) => {
                    if (selectedMunicipio === 'all' || selectedMunicipio === mKey) {
                        Object.entries(m.freguesia).forEach(([fKey, f]) => {
                            if (selectedFreguesia === 'all' || selectedFreguesia === fKey) {
                                ids.push(...f.property_ids);
                            }
                        });
                    }
                });
            }
        });
        setFilteredPropertyIds(ids);
        // Clear any active suggestion when admin filters change
        setActiveGroupIndex(null);
    }, [admAreaData, selectedDistrito, selectedMunicipio, selectedFreguesia]);

    // Effect: refetch admin data on mount
    useEffect(() => { fetchAdmAreaData(); }, [fetchAdmAreaData]);

    // Effect: update when admin filter values change (only when no suggestion active)
    useEffect(() => {
        if (activeGroupIndex === null) {
            updateFilteredPropertyIds();
        }
    }, [selectedDistrito, selectedMunicipio, selectedFreguesia, updateFilteredPropertyIds, activeGroupIndex]);

    // Fetch suggestion groups
    useEffect(() => {
        const ctrl = new AbortController();
        (async () => {
            try {
                const resp = await fetch('http://localhost:8080/api/suggestions_by_neighbours', { signal: ctrl.signal });
                if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
                const data = await resp.json();
                setSuggestionGroups(data);
            } catch (e) {
                if (e.name !== 'AbortError') console.error(e);
            }
        })();
        return () => ctrl.abort();
    }, []);

    // Handle group click
    const handleGroupClick = (idx) => {
        setActiveGroupIndex(idx);
        const group = suggestionGroups[idx];
        const involvedIds = group.Properties_Envolved.map(p => p.Id);
        setFilteredPropertyIds(involvedIds);
    };

    const clearSelection = () => {
        setActiveGroupIndex(null);
        updateFilteredPropertyIds();
    };

    // Fetch properties
    useEffect(() => {
        const ctrl = new AbortController();
        (async () => {
            try {
                setLoading(true);
                const resp = await fetch('http://localhost:8080/api/prop', { signal: ctrl.signal });
                if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
                const data = await resp.json();
                setGeojson(convertGeoJSON(data));
            } catch (e) {
                if (e.name !== 'AbortError') {
                    console.error(e);
                    setError(e.message);
                }
            } finally { setLoading(false); }
        })();
        return () => ctrl.abort();
    }, [convertGeoJSON]);

    // Compute filtered GeoJSON
    const filteredGeojson = React.useMemo(() => {
        if (!geojson) return null;
        if (!filteredPropertyIds.length) return geojson;
        return {
            ...geojson,
            features: geojson.features.filter(f => filteredPropertyIds.includes(f.properties.id)),
        };
    }, [geojson, filteredPropertyIds]);

    // Admin select options
    const getDistritos = useCallback(() => (admAreaData ? Object.keys(admAreaData.distrito) : []), [admAreaData]);
    const getMunicipios = useCallback(() => {
        if (!admAreaData) return [];
        if (selectedDistrito === 'all') {
            const setM = new Set();
            Object.values(admAreaData.distrito).forEach(d => Object.keys(d.municipio).forEach(m => setM.add(m)));
            return Array.from(setM);
        }
        return Object.keys(admAreaData.distrito[selectedDistrito].municipio);
    }, [admAreaData, selectedDistrito]);
    const getFreguesias = useCallback(() => {
        if (!admAreaData) return [];
        if (selectedDistrito === 'all' || selectedMunicipio === 'all') {
            const setF = new Set();
            Object.values(admAreaData.distrito).forEach(d =>
                Object.values(d.municipio).forEach(m => Object.keys(m.freguesia).forEach(f => setF.add(f)))
            );
            return Array.from(setF);
        }
        return Object.keys(admAreaData.distrito[selectedDistrito].municipio[selectedMunicipio].freguesia);
    }, [admAreaData, selectedDistrito, selectedMunicipio]);

    return (
        <div style={{ display: 'flex', height: '100vh' }}>
            {/* Sidebar */}
            <div style={{ width: 300, overflowY: 'auto', borderRight: '1px solid #ccc', padding: 10, background: '#fafafa' }}>
                <h3>Suggestions</h3>
                <button onClick={clearSelection} disabled={activeGroupIndex === null} style={{ marginBottom: 10 }}>
                    Reset
                </button>
                {suggestionGroups.map((group, idx) => (
                    <div
                        key={idx}
                        onClick={() => handleGroupClick(idx)}
                        style={{
                            cursor: 'pointer',
                            padding: 8,
                            marginBottom: 8,
                            backgroundColor: idx === activeGroupIndex ? '#d0eaff' : '#fff',
                            borderRadius: 4,
                            boxShadow: '0 0 3px rgba(0,0,0,0.1)',
                        }}
                    >
                        {group.Suggestion.map((s, i) => (
                            <div key={s.Id} style={{ marginBottom: i < group.Suggestion.length - 1 ? 4 : 0 }}>
                                <strong>ID:</strong> {s.Id} &nbsp;
                                <strong>Owner:</strong> {s.Owner} &nbsp;
                                <strong>Freguesia:</strong> {s.Freguesia}
                            </div>
                        ))}
                    </div>
                ))}
            </div>

            {/* Map */}
            <div style={{ flexGrow: 1, position: 'relative' }}>
                <Map
                    {...viewState}
                    onMove={e => setViewState(e.viewState)}
                    mapLib={maplibregl}
                    mapStyle="https://api.maptiler.com/maps/satellite/style.json?key=LSfqsa4izS6vDHlLAK4a"
                >
                    {filteredGeojson && (
                        <Source id="properties" type="geojson" data={filteredGeojson} promoteId="id">
                            <Layer
                                id="properties-fill"
                                type="fill"
                                paint={{
                                    'fill-color': ['match', ['get', 'municipio'], 'Calheta', '#f28cb1', 'Funchal', '#8dd3c7', '#ffcc00'],
                                    'fill-opacity': 0.7,
                                    'fill-outline-color': '#000',
                                }}
                            />
                            <Layer id="properties-outline" type="line" paint={{ 'line-color': '#000', 'line-width': 1 }} />
                            {showPropertyIds && (
                                <Layer
                                    id="properties-label"
                                    type="symbol"
                                    layout={{ 'text-field': ['to-string', ['get', 'id']], 'text-size': 12, 'text-offset': [0, 0.5], 'text-anchor': 'top' }}
                                    paint={{ 'text-color': '#000', 'text-halo-color': '#fff', 'text-halo-width': 1 }}
                                />
                            )}
                        </Source>
                    )}
                </Map>

                {/* Filter controls */}
                <div style={{ position: 'absolute', top: 10, right: 10, background: 'rgba(255,255,255,0.9)', padding: 12, borderRadius: 4, zIndex: 1, width: 250, boxShadow: '0 0 10px rgba(0,0,0,0.1)' }}>
                    <h3 style={{ margin: '0 0 10px' }}>Filters</h3>
                    <label>Distrito:</label>
                    <select value={selectedDistrito} onChange={e => { setSelectedDistrito(e.target.value); setSelectedMunicipio('all'); setSelectedFreguesia('all'); }} style={{ width: '100%', padding: 5, marginBottom: 10 }}>
                        <option value="all">All</option>
                        {getDistritos().map(d => <option key={d} value={d}>{d}</option>)}
                    </select>
                    <label>Município:</label>
                    <select value={selectedMunicipio} onChange={e => { setSelectedMunicipio(e.target.value); setSelectedFreguesia('all'); }} style={{ width: '100%', padding: 5, marginBottom: 10 }}>
                        <option value="all">All</option>
                        {getMunicipios().map(m => <option key={m} value={m}>{m}</option>)}
                    </select>
                    <label>Freguesia:</label>
                    <select value={selectedFreguesia} onChange={e => setSelectedFreguesia(e.target.value)} style={{ width: '100%', padding: 5 }}>
                        <option value="all">All</option>
                        {getFreguesias().map(f => <option key={f} value={f}>{f}</option>)}
                    </select>
                    <div style={{ marginTop: 10 }}>
                        <label><input type="checkbox" checked={showPropertyIds} onChange={() => setShowPropertyIds(!showPropertyIds)} /> Show IDs</label>
                    </div>
                </div>

                {/* Debug overlay */}
                <div style={{ position: 'absolute', top: 10, left: 10, background: 'rgba(255,255,255,0.9)', padding: 8, borderRadius: 4, fontFamily: 'monospace', zIndex: 1 }}>
                    {loading && <div>🔄 Loading...</div>}
                    {error && <div style={{ color: 'red' }}>❌ {error}</div>}
                    {geojson && <div>✅ {filteredGeojson.features.length} / {geojson.features.length} properties</div>}
                    <div>📍 {viewState.longitude.toFixed(5)}, {viewState.latitude.toFixed(5)}</div>
                    <div>🔍 Zoom: {viewState.zoom.toFixed(1)}</div>
                </div>
            </div>
        </div>
    );
};

export default MapComponent;
