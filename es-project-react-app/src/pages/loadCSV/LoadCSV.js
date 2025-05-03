import './LoadCSV.css';
import { useState } from 'react';
import ThemedButton from '../../components/themedButton/themedButton';

function LoadingCSV({ textContent }) {
  const [file, setFile] = useState(null);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleFileUpload = () => {
    const formData = new FormData();
    formData.append('file', file);

    fetch('http://localhost:8080/api/upload', {
      method: 'POST',
      body: formData,
    })
      .then(response => response.text())
      .then(result => {
        alert(result);
        window.location.href = '/homepage'; 
      })
      .catch((error) => {
        console.error('Error uploading file:', error);
      });
  };

  return (
    <section className="loading-csv-content">
      <section className="text-content-limit">
        <p className="text-content">{textContent}</p>
      </section>
      
      <input
        type="file"
        accept=".csv"
        onChange={handleFileChange}
        style={{ marginBottom: '20px' }}
      />
      
      <ThemedButton buttonTextContent="Carregar" buttonOnClick={handleFileUpload} />
    </section>
  );
}

export default LoadingCSV;
