
import './themedButton.css';

function ThemedButton( { buttonTextContent, buttonOnClick } ) {

  return (
    <button class="themed-button" onClick={() => buttonOnClick()}>
     {buttonTextContent} 
    </button>
  );

}

export default ThemedButton;
