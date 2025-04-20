import { useState } from 'react';
import SigmaGraph from './SigmaGraph';
import './graph.css';
import MapComponent from './MapComponent';

function Graph() {
  const [headerTabsNames] = useState(["Propriedades", "Proprietários"]);
  const [headerTabsState, setHeaderTabsState] = useState([true, false]);

  const graphHeadeTabOnClick = (e) => {
    const index = headerTabsNames.indexOf(e.target.innerText);
    setHeaderTabsState(headerTabsState.map((_, i) => i === index));
  };
  //  <SigmaGraph />
  return (
    <section className="graph-main-content">
      <section className="graph-content">
          <MapComponent />
      </section>

      <section className="graph-header-content">
        {headerTabsNames.map((tab, index) => (
          <section
            key={index}
            className={headerTabsState[index] ? "graph-header-tab-selected" : "graph-header-tab"}
            onClick={graphHeadeTabOnClick}
          >

            {tab}
          </section>
        ))}
      </section>

    </section>
  );
}

export default Graph;
