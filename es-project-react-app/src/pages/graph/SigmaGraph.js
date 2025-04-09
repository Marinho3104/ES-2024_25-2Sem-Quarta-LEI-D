import React, { useEffect, useRef, useState } from 'react';
import Sigma from 'sigma';
import Graph from 'graphology';

const SigmaGraph = () => {
  const containerRef = useRef(null);

  useEffect(() => {
    fetch('http://localhost:8080/api/graph')
      .then(response => response.json())
      .then(data => {
        const { nodes, edges } = data;

        const graph = new Graph();

        nodes.forEach((node) => {
          graph.addNode(node.id, {
            label: node.label || node.id,
            x: node.x,
            y: node.y,
            size: 1,
            color: '#6FB1FC',
          });
        });

        edges.forEach((edge) => {
          graph.addEdgeWithKey(
            edge.id,
            edge.source,
            edge.target,
            {
              label: edge.label || '',
              color: 'rgba(150, 150, 150, 0.3)',
            }
          );
        });

        const sigmaInstance = new Sigma(graph, containerRef.current, {
          renderEdgeLabels: false, 
        });

        const camera = sigmaInstance.getCamera();

        camera.setState({
          ratio: 2,
          x: 0.5,
          y: 0.5 
        });


      })
      .catch(err => {
        console.error('Error fetching graph data:', err);
      });
  }, []); 

  return (

    <div ref={containerRef} className="sigma-container" />
  
  );
};

export default SigmaGraph;
