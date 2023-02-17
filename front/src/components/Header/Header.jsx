import React from 'react';
import { NavLink, Route, Routes } from 'react-router-dom';
import Profile from '../Profile/Profile';
import classes from './Header.module.css';

const Header = () => {
    return (
    <header className={classes.header}>
        <div className={classes.vkImg}>
            <a href='/'><img src="vk.png" alt=''/></a>
        </div>
        <div>
            <NavLink to='/Profile'><img src="user.png" alt=''/></NavLink>
        </div>
        <div>
            <NavLink to='/Create'><img src="more.png" alt=''/></NavLink>
        </div>
    </header>
    );
}

export default Header;