import React, { useState } from 'react';
import s from './Post.module.css';
import parse, { domToReact } from 'html-react-parser';
import Comments from './Comments/Comments';
import { NavLink } from 'react-router-dom';

const PostItem = (props) => {
    return (
        <div className={s.postItem + ' ' + s.active}>
            <NavLink to={"/post/" + props.id}><Post /></NavLink>
        </div>
    );
}

const Post = (props) => {
    const html = props.source;

    const options = {
        replace: ({ attribs, children }) => {
            if (!attribs) {
                return;
            }

            if (attribs.class === 'title') {
                return (
                    <div className={s.title}>{domToReact(children, options)}</div>
                );
            }

            if (attribs.class === 'content') {
                return (
                    <div className={s.content}>{domToReact(children, options)}</div>
                );
            }
        }
    };

    const [showCommentElement, setShowCommentElement] = useState(false);

    const handleButtonClick = () => {
        setShowCommentElement(!showCommentElement);
    }

    const [likeElement, setLikeElement] = useState(true);

    const handleLikeClick = () => {
        setLikeElement(!likeElement);
    }

    return (
        <div className={s.post}>
            <div className={s.postName + ' ' + s.content}>
                {parse(html, options)}
            </div>
            <div className={s.imageContainer}>
                <img src='testpic.png' alt='' />
            </div>
            <div className={s.hl} />
            <div className={s.buttons}>
                <div className={s.likeButton} onClick={props.isLikePressed.h}>
                    {props.isLikePressed.l ? (
                        <img src='like.png' alt='' />
                    ) : 
                    (
                        <img src='like_active.png' alt='' />
                    )}
                </div>
                <div className={s.commentButton} onClick={handleButtonClick}>
                    <img src='comment.png' alt='' />
                </div>
                <div className={s.arrow_back}>
                    <NavLink to='/' onClick={props.isVisible.h}><img src='arrow_back.png' alt='' /></NavLink>
                </div>
            </div>
            {showCommentElement && (
                <div className={s.comments}>
                    <Comments comments={props.comments} postId={props.postId} />
                </div>
            )}
        </div>
    );
}

export default Post;