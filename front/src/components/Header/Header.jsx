import React from 'react';
import { NavLink, Route, Routes } from 'react-router-dom';
import Profile from '../Profile/Profile';
import classes from './Header.module.css';

const Header = () => {
    return (
    <header className={classes.header}>
        <div className={classes.vkImg}>
            <a href='/'><img src="photo1676760339.jpeg" alt=''/></a>
        </div>
        <div>
            <NavLink to='/Profile'><img src="photo1676760637.jpeg" alt=''/></NavLink>
        </div>
        <div>
            <NavLink to='/Create'><img src="create_post.jpeg" alt=''/></NavLink>
        </div>
    </header>
    );
}

export default Header;