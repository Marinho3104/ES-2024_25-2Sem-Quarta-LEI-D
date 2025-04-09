import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './pages/app/App';
import LandingPage from './pages/landingPage/LandingPage';
import LoadingCSV from './pages/loadCSV/LoadCSV';
import reportWebVitals from './tests/reportWebVitals';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <App />
);
// Removed StricMode cuz it was rendering the graph twice (annoying )
// <React.StrictMode>
// <App />
// {/* <LandingPage /> */}
// {/* <LoadingCSV /> */}
// </React.StrictMode> 

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
