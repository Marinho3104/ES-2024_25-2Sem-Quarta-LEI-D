DEFINE TABLE property SCHEMALESS;
DEFINE TABLE neighbor_property;

DEFINE FUNCTION fn::boundingBox($polygon:  geometry ) {
    LET $coords = $polygon.coordinates;
    LET $flat = array::flatten(array::flatten($coords));
    LET $x = $flat.map(|$v| $v[0] );
    LET $y = $flat.map(|$v| $v[1] );
    RETURN {
        "x_max" :  array::max($x),
        "x_min" :  array::min($x),
        "y_max" :  array::max($y),
        "y_min" :  array::min($y),
    };
};

