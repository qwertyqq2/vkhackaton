import React from 'react';
import { NavLink } from 'react-router-dom';
import classes from './Header.module.css';

const Header = () => {
    return <header className={classes.header}>
        <NavLink to=''><img src="vk.png" alt=''/></NavLink>
    </header>
}

export default Header;