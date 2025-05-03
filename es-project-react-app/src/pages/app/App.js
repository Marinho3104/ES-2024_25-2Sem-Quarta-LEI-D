import logo from '../../images/logo.svg';
import './App.css';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import LandingPage from '../landingPage/LandingPage';
import LoadingCSV from '../loadCSV/LoadCSV';
import Layout from './Layout';
import Graph from '../graph/graph';
import LoadCSVAgain from '../loadCSVAgain/loadCSVAgain';
import HomePage from '../homepage/homepage';
import Map from '../map/map';
const textContentStandAloneLoadingCSV = "Como primeiro passo, solicitamos a disponibilização de um ficheiro CSV com os dados a carregar na plataforma."
const textContentLoadingCSV = "A importação de um novo ficheiro CSV irá substituir o ficheiro atual e eliminar todos os dados previamente carregados."
const textContentLandingPage = "Junta Terras, Cria Futuro"

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

            <Route
              path="/"
              element={<LandingPage textContent={textContentLandingPage}/>}
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
          <Route
              path='/map'
              element={<Map />}
          />

          <Route
            path='/homepage'
            element={<HomePage />}
          />

          <Route
            path='/loadCSVAgain'
            element={<LoadCSVAgain textContent={textContentLoadingCSV}/>}
          />

        </Route>

      </Routes>

      </section>
    </BrowserRouter>
  );
}

export default App;
