import React from 'react';
import s from './Comment.module.css';

const Comment = () => {
    return (
        <div className={s.comment}>
            <div className={s.commentatorImg}>
                <img src='avatar.jpg' alt='' />
            </div>
            <div className={s.commentatorText}>
                <p>
                    aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaffffffffffffffffffff
                    sddddddddddddddddddddddddddaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaff
                    ffffffffffffffffff
                    sdddddddddddddddddddddddddd
                    aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaffffffffffffffffffff
                    sdddddddddddddddddddddddddd
                </p>
            </div>
        </div>
    );
}

export default Comment;