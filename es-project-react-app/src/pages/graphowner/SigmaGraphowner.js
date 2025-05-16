import React, { useEffect, useRef, useState } from 'react';
import Sigma from 'sigma';
import Graph from 'graphology';
import forceAtlas2 from "graphology-layout-forceatlas2";

const SigmaGraphowner = () => {
  const containerRef = useRef(null);

  useEffect(() => {
    fetch('http://localhost:8080/api/owner')
      .then(response => response.json())
      .then(data => {
        const { nodes, edges } = data;

        const graph = new Graph();

        nodes.forEach((node) => {
          graph.addNode(node.id, {
            label: node.label || node.id,
            size: 1,
              x: Math.random(),
              y: Math.random(),
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
          forceAtlas2.assign(graph, {
              iterations: 300,
              settings: {
                  gravity: 0.5,
                  scalingRatio: 5,
                  strongGravityMode: true,
                  barnesHutOptimize: true, // critical for performance
              }
          });

        const sigmaInstance = new Sigma(graph, containerRef.current, {
          renderEdgeLabels: false,
        });

        const camera = sigmaInstance.getCamera();

        camera.setState({
          ratio: 1,
          x: 0,
          y: 0
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

export default SigmaGraphowner;
