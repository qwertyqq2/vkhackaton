import React from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import './App.css';
import Header from './components/Header/Header';
import Posts from './components/Posts/Posts';
import Profile from './components/Profile/Profile';
import Create from './components/Create/Create';
import SidePost from './components/SidePost/SidePost';

function App(props) {

  return (
    <BrowserRouter>
      <div className="app-wrapper">
        <Header />
        <div className='app-wrapper-content'>
          <Routes>
            <Route path='/' element={<Posts
              htmlContent={props.store.getHtmlContent}
              comments={props.store.comments}
              likes={props.store.likes} />} />
            <Route path='/Profile' element={<Profile
              account={props.store.account} />} />
            <Route path='/Create' element={<Create />} />
          </Routes>
        </div>
        <div className='app-wrapper-s'>
          <Routes>
            <Route path='/' element={<SidePost />}/>
          </Routes>
        </div>
      </div>
    </BrowserRouter>
  );
}

export default App;
