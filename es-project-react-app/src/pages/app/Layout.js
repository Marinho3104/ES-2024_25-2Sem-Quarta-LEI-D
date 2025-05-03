
import { Outlet } from 'react-router-dom';
import Header from '../../components/header/header';
import './Layout.css';


function Layout() {

  const currentPage = window.location.href.split('/')[window.location.href.split('/').length - 1]

  const inicioIsSelected = currentPage === 'homepage' ? true : false
  const graphIsSelected = currentPage === 'graph' ? true : false
  const carregarIsSelected = currentPage === 'loadCSVAgain' ? true : false

  return (

    <section class="layout-content">

      <Header 
        inicioIsSelected={inicioIsSelected} 
        graphIsSelected={graphIsSelected} 
        carregarIsSelected={carregarIsSelected}
      />

      <Outlet />

    </section>

  );
  

}

export default Layout
