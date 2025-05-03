import './headerButton.css';

function HeaderButton({ buttonContent, isSelected}) {

  return (
    <button 
      class="header-button" 
      style={{ 
        "--backgroundColor": isSelected ? 'white' : 'transparent',
        "--color": isSelected ? '#101728' : 'white'
      }}
      onClick={() => {

        if( 'início' === buttonContent.toLowerCase() ) {
          window.location.href = '/homepage';
        }
        if( 'grafos' === buttonContent.toLowerCase() ) {
          window.location.href = `/graph`;
        }
        if( 'carregar' === buttonContent.toLowerCase() ) {
          window.location.href = `/loadCSVAgain`;
        }
      }}
    >
      <section class="header-button-content">{buttonContent}</section>
    </button>
  );

}

export default HeaderButton
