import React from "react";
import s from "./Post.module.css";
import parse, { domToReact } from 'html-react-parser';
import { NavLink } from "react-router-dom";

const PrePost = (props) => {
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

    return (
        <div className={s.post + (props.isVisible.v ? '' : (' ' + s.hidden))}>
            <NavLink to={'/post' + props.index}>
                <div onClick={props.isVisible.h}>
                    {parse(html, options)}
                </div>
            </NavLink>
            <div className={s.buttons}>
                <div className={s.likeButton} onClick={props.isLikePressed.h}>
                    {props.isLikePressed.l ? (
                        <img src='like.png' alt='' />
                    ) : 
                    (
                        <img src='like_active.png' alt='' />
                    )}
                    {props.likeCount.count}
                </div>
            </div>
        </div>
    );
}

export default PrePost;