import React from 'react';
import { BrowserRouter, Route } from 'react-router-dom';
import './App.css';
import Header from './components/Header/Header';
import Posts from './components/Posts/Posts';

function App() {
  return (
    <BrowserRouter>
    <div className="app-wrapper">
      <Header />
      <div className='app-wrapper-content'>
        <Posts />
      </div>
      <div className='app-wrapper-s'>

      </div>
    </div>
    </BrowserRouter>
  );
}

export default App;
