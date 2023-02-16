import React from 'react';
import Comment from './Comment/Comment';
import s from './Comments.module.css';

const Comments = () => {
    return (
        <div className={s.comments}>
            <div className={s.comment}>
                <Comment />
            </div>
            <div className={s.createComment}>
                <div className={s.textArea}>
                    <textarea></textarea>
                </div>
                <div className={s.arrow}>
                    <img src='right-arrow.png' alt='' />
                </div>
            </div>
        </div>
    );
}

export default Comments;