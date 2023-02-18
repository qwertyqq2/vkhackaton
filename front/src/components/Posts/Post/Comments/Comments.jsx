import React, { useState } from 'react';
import Comment from './Comment/Comment';
import s from './Comments.module.css';
import { addComment } from '../../../../redux/state';

const Comments = (props) => {

    let commentators = [];

    props.comments.forEach(element => {
        if (element.bindPostId === props.postId) {
            commentators.push(element);
        }
    });

    let newCommentElem = React.createRef();

    let addComm = () => {
        console.log("It works");
        addComment({ postId: props.postId, id: 4, message: newCommentElem.current.value });
    };

    let onCommentChange = () => {
        let text = newCommentElem.current.value;
        console.log(text);
    }

    //to back
    const [inputValue, setInputValue] = useState('');

    const handleInputChange = (event) => {
        setInputValue(event.target.value);
    };

    const handleSubmit = (event) => {
        addComm();
        event.preventDefault();
        // do something with the input value
        const form = document.querySelector('form');

        console.log("1");

        form.addEventListener('submit', (event) => {
            event.preventDefault();

            const message = form.elements['message'].current.value;

            const data = {
                message: message,
            };
            fetch('http://localhost:3001/create_post', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            })
                // .then(response => {
                //     if (response.ok) {
                //         // handle the response from the Go backend here
                //     } else {
                //         throw new Error('Network response was not ok');
                //     }
                // })
                .catch(error => {
                    console.error('Error:', error);
                });
        });
    };
    //

    return (
        <div className={s.comments}>
            <div className={s.comment}>
                {commentators.map((elem) =>
                    <Comment comments={elem} />
                )}
            </div>
            <div className={s.createComment}>
                <form>
                    <div className={s.textArea}>
                        <textarea id="message" onChange={onCommentChange} ref={newCommentElem} ></textarea>
                    </div>
                </form>
                <div className={s.arrow}>
                    <img src='right-arrow.png' alt='' onClick={handleSubmit} />
                </div>
            </div>
        </div>
    );
}

export default Comments;