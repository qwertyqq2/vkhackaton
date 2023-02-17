import React from 'react';
import s from './Comment.module.css';

const Comment = (props) => {
    return (
        <div className={s.comment}>
            <div className={s.commentatorImg}>
                <img src='avatar.jpg' alt='' />
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