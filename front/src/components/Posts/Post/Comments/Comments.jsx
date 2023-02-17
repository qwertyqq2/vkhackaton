import React from 'react';
import Comment from './Comment/Comment';
import s from './Comments.module.css';

const Comments = (props) => {

    let commentators = [];
    
    props.comments.forEach(element => {
        if (element.bindPostId === props.postId) {
            commentators.push(element);
        }
    });


    // for (let i = 1; i <= 3; i++) {
    //     // eslint-disable-next-line
    //     let htmlModule = require(`raw-loader!../data/htmlExample` + i + `.html`);
    //     let html = htmlModule.default;
    
    //     loadHtmlContent.push(html);
    // }

    return (
        <div className={s.comments}>
            <div className={s.comment}>
            {commentators.map((elem) =>
                <Comment comments={elem} />
            )}
                {/* <Comment /> */}
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