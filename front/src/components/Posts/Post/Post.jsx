import React from 'react';
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
        replace: ({attribs, children}) => {
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
                <div className={s.likeButton}>
                    <img src='like-button.png' alt='' />
                </div>
                <div className={s.commentButton}>
                    <img src='comment.png' alt='' />
                </div>
                <div className={s.arrow_back}>
                    <NavLink to=''><img src='back_arrow.png' alt='' /></NavLink>
                </div>
            </div>
            <div className={s.comments}>
                <Comments />
            </div>
        </div>
    );
}

export default Post;