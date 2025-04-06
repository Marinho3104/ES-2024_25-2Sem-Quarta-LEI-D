import logo from '../../images/logo.svg';
import './App.css';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import LandingPage from '../landingPage/LandingPage';
import LoadingCSV from '../loadCSV/LoadCSV';
import Layout from './Layout';
import Graph from '../graph/graph';

const textContentStandAloneLoadingCSV = "Como primeiro passo, solicitamos a disponibilização de um ficheiro CSV com os dados a carregar na plataforma."

function App() {
  return (
    <BrowserRouter>
      <section class="content">
        <Routes>

          <Route
            path="/"
          >

            <Route
              path="/loadCSV"
              element={<LoadingCSV textContent={textContentStandAloneLoadingCSV} />} 
            />

          </Route>


        <Route
          path="/"
          element={<Layout />}
        >
          <Route
            path='/graph'
            element={<Graph />}
          />

        </Route>

      </Routes>

      </section>
    </BrowserRouter>
  );
}

export default App;
