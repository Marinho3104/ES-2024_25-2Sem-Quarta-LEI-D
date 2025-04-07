import HeaderButton from '../headerButton/headerButton';
import './header.css';

function Header({ inicioIsSelected, graphIsSelected, carregarIsSelected}) {

  return (
    <section class="header-content">

      <HeaderButton buttonContent="Início" isSelected={inicioIsSelected}/>
      <HeaderButton buttonContent="Grafos" isSelected={graphIsSelected}/>
      <HeaderButton buttonContent="Carregar" isSelected={carregarIsSelected}/>

    </section>
  );

}

export default Header; 
