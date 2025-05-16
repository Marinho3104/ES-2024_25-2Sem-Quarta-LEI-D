import { useState } from 'react';
import SigmaGraphowner from './SigmaGraphowner';
import './graphowner.css';

function Graphowner() {
  const [headerTabsNames] = useState(["Propriedades", "Proprietários"]);
  const [headerTabsState, setHeaderTabsState] = useState([true, false]);

  const graphHeadeTabOnClick = (e) => {
    const index = headerTabsNames.indexOf(e.target.innerText);
    setHeaderTabsState(headerTabsState.map((_, i) => i === index));
  };
  //  
  return (
    <section className="graph-main-content">
      <section className="graph-content">
          <SigmaGraphowner />
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

export default Graphowner;
