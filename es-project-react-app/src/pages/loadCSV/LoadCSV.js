
import './LoadCSV.css';
import ThemedButton from '../../components/themedButton/themedButton';

function LoadingCSV( {textContent} ) {

  return (
    <section class="loading-csv-content">

      <section class="text-content-limit">
        <p class="text-content"> {textContent} </p>
      </section>
      <ThemedButton buttonTextContent="Carregar" buttonOnClick={ () => window.location.href = '/homepage' } /> 

    </section>
  );

}

export default LoadingCSV;
