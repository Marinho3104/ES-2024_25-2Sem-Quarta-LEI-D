import { useState } from 'react';
import SigmaGraph from './SigmaGraph';
import SigmaGraphowner from "../graphowner/SigmaGraphowner";
import './graph.css';

function Graph() {
    const [headerTabsNames] = useState(["Propriedades", "Proprietários"]);
    const [selectedTabIndex, setSelectedTabIndex] = useState(0);

    const graphHeadeTabOnClick = (e) => {
        const index = headerTabsNames.indexOf(e.target.innerText);
        setSelectedTabIndex(index);
    };

    return (
        <section className="graph-main-content">
            <section className="graph-header-content">
                {headerTabsNames.map((tab, index) => (
                    <section
                        key={index}
                        className={selectedTabIndex === index ? "graph-header-tab-selected" : "graph-header-tab"}
                        onClick={graphHeadeTabOnClick}
                    >
                        {tab}
                    </section>
                ))}
            </section>

            <section className="graph-content">
                {selectedTabIndex === 0 && <SigmaGraph />}
                {selectedTabIndex === 1 && <SigmaGraphowner />}
            </section>
        </section>
    );
}

export default Graph;
