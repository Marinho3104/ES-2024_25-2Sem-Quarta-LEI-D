import ThemedButton from '../../components/themedButton/themedButton';
import './LandingPage.css';
import image from './../../images/oie_O5LqY2Q9XVJY.jpg'

function LandingPage({ textContent }) {

  return (
    <section class="landing-page-content">

      <section class="landing-page-header">

        <section class="landing-page-header-title">
          ES-2024_25-2Sem-Quarta-LEI-D
        </section>

      </section>

      <section class="landing-page-body">

        <section class="landing-page-body-left">

          <section class="landing-page-body-left-text-limit">
            <p class="landing-page-body-left-text">{textContent}</p>
          </section>

          <section class="landing-page-body-left-button-box">
            < ThemedButton buttonTextContent={"Começar"} buttonOnClick={ () => window.location.href = '/loadCSV' }/>
          </section>

        </section>

        <section class="landing-page-body-right">

          <section class="landing-page-body-right-image-box">

            <section class="landing-page-body-right-background-fade">

              <img 
                class="landing-page-body-right-image" 
                src={image}
              />

            </section>

          </section>

        </section>

      </section>

    </section>
  );

}

export default LandingPage;
