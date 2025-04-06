
import { Outlet } from 'react-router-dom';
import Header from '../../components/header/header';
import './Layout.css';


function Layout() {

  return (

    <section class="layout-content">

      <Header />

      <Outlet />

    </section>

  );
  

}

export default Layout
