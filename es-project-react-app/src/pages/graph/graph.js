import { useState } from 'react';
import './graph.css';

function Graph() {

  const [ headerTabsNames, setHeaderTabsNames ] = useState([ "Propriedades", "Proprietários" ]);
  const [ headerTabsState, setHeaderTabsState ] = useState([ true, false ]);


  const graphHeadeTabOnClick = ( e ) => {
    const index = headerTabsNames.indexOf( e.target.innerText );
    setHeaderTabsState( headerTabsState.map( ( tab, i ) => i === index ? true : false ) );
  }

  return (
    <section class="graph-main-content">

      <section class="graph-content">
      </section>

      <section class="graph-header-content">

        { headerTabsNames.map( ( tab, index ) => {
          
          return (
            <section key={ index } class={ headerTabsState[ index ] ? "graph-header-tab-selected" : "graph-header-tab" } onClick={( e ) => { graphHeadeTabOnClick( e ) }}>
              { tab } 
            </section>
          );

        })}

      </section>
    </section>
  );

}

export default Graph
