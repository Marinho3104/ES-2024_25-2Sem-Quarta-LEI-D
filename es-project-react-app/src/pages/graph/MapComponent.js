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

    useEffect(() => {
        const controller = new AbortController();

        const fetchData = async () => {
            try {
                setLoading(true);

                const cached = localStorage.getItem('geojson');
                if (cached) {
                    const parsed = JSON.parse(cached);
                    setGeojson(parsed);
                    setLoading(false);
                    return;
                }
                const response = await fetch('http://localhost:8080/api/prop', {
                    signal: controller.signal
                });

                if (!response.ok) throw new Error(`HTTP ${response.status}`);

                const data = await response.json();
                const convertedData = convertGeoJSON(data);

                setGeojson(convertedData);
                localStorage.setItem('geojson', JSON.stringify(convertedData));

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


    return (
        <div style={{ position: 'relative', width: '100%', height: '600px' }}>
            <Map
                {...viewState}
                onMove={evt => setViewState(evt.viewState)}
                mapLib={maplibregl}
                mapStyle="https://api.maptiler.com/maps/satellite/style.json?key=LSfqsa4izS6vDHlLAK4a"
            >
                {geojson && (
                    <Source
                        id="properties"
                        type="geojson"
                        data={geojson}
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
                    </Source>
                )}
            </Map>

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
                        <div>✅ {geojson.features.length} properties</div>
                        <div>📍 {viewState.longitude.toFixed(5)}, {viewState.latitude.toFixed(5)}</div>
                        <div>🔍 Zoom: {viewState.zoom.toFixed(1)}</div>
                    </>
                )}
            </div>
        </div>
    );
};

export default MapComponent;