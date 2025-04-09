import { useEffect, useRef } from 'react';
import cytoscape from 'cytoscape';

function CytoscapeGraph() {
  const cyRef = useRef(null);

  useEffect(() => {
    if (cyRef.current) {
      cytoscape({
        container: cyRef.current,
        elements: [
          { data: { id: 'a' } },
          { data: { id: 'b' } },
          { data: { id: 'c' } },
          { data: { source: 'a', target: 'b' } },
          { data: { source: 'b', target: 'c' } }
        ],
        style: [
          {
            selector: 'node',
            style: {
              'background-color': '#0074D9',
              'label': 'data(id)',
              'color': '#fff',
              'text-valign': 'center',
              'text-halign': 'center',
              'font-size': '14px'
            }
          },
          {
            selector: 'edge',
            style: {
              'width': 2,
              'line-color': '#ccc',
              'target-arrow-color': '#ccc',
              'target-arrow-shape': 'triangle'
            }
          }
        ],
        layout: {
          name: 'grid',
          rows: 1
        }
      });
    }
  }, []);

  return <div ref={cyRef} className="cytoscape-container" />;
}

export default CytoscapeGraph;
