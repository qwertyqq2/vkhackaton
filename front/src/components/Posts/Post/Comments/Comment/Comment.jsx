import React from 'react';
import { NavLink } from 'react-router-dom';
import s from './Comment.module.css';

const Comment = (props) => {
    return (
        <div className={s.comment}>
            <div className={s.commentatorImg}>
                <NavLink to="/Profile"><img src='mommymonkey.jpeg' alt='' /></NavLink>
            </div>
            <div className={s.commentatorText}>
                <p>
                    { props.comments.message }
                </p>
            </div>
        </div>
    );
}

export default Comment;