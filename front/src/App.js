import React from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import './App.css';
import Header from './components/Header/Header';
import Posts from './components/Posts/Posts';
import Profile from './components/Profile/Profile';
import Create from './components/Create/Create';

function App(props) {
  return (
    <BrowserRouter>
      <div className="app-wrapper">
        <Header />
        <div className='app-wrapper-content'>
          <Routes>
            <Route path='*' element={<Posts htmlContent={props.htmlContent} comments={props.comments}/>} />
            <Route path='/Profile' element={<Profile htmlContent={props.htmlContent} />} />
            <Route path='/Create' element={<Create />} />
          </Routes>
        </div>
        <div className='app-wrapper-s'>

        </div>
      </div>
    </BrowserRouter>
  );
}

export default App;
