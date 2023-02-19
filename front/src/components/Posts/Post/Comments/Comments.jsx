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
<<<<<<< HEAD
        console.log("It works");
=======
        // console.log("It works");
>>>>>>> c9e4d4e02ae8cc9163044fa9830434ff55c3f139
        addComment({ postId: props.postId, id: 4, message: newCommentElem.current.value });
    };

    let onCommentChange = () => {
        let text = newCommentElem.current.value;
<<<<<<< HEAD
        console.log(text);
=======
        // console.log(text);
>>>>>>> c9e4d4e02ae8cc9163044fa9830434ff55c3f139
    }

    //to back
    const [inputValue, setInputValue] = useState('');

    const handleInputChange = (event) => {
        setInputValue(event.target.value);
    };

    const handleSubmit = (event) => {
        event.preventDefault();
        // do something with the input value
        const form = document.querySelector('form');
<<<<<<< HEAD

        console.log("1");

=======

        // console.log("1");

>>>>>>> c9e4d4e02ae8cc9163044fa9830434ff55c3f139
        form.addEventListener('submit', (event) => {
            event.preventDefault();

            const message = form.elements['message'].current.value;

            const data = {
                message: message,
            };
            fetch('http://localhost:3001/create_comment', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            })
<<<<<<< HEAD
                .then(response => {
                    if (response.ok) {
                        // handle the response from the Go backend here
                    } else {
                        throw new Error('Network response was not ok');
                    }
                })
=======
                // .then(response => {
                //     if (response.ok) {
                //         // handle the response from the Go backend here
                //     } else {
                //         throw new Error('Network response was not ok');
                //     }
                // })
>>>>>>> c9e4d4e02ae8cc9163044fa9830434ff55c3f139
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
<<<<<<< HEAD
            <div className={s.createComment}>
<<<<<<< HEAD
                <form onSubmit ={handleSubmit}>
                    <div className={s.textArea}>
                        <textarea id="message" onChange={onCommentChange} ref={newCommentElem} ></textarea>
                    </div>
                <div className={s.arrow}>
                    {/* <img src='right-arrow.png' alt='' onClick={handleSubmit} /> */}
                    <input type = "submit" value="Submit" onClick={handleInputChange}/>
                </div>
=======
                <form onSubmit={handleSubmit}>
                    <div className={s.textArea}>
                        <textarea id="message" onChange={onCommentChange} ref={newCommentElem} ></textarea>
                    </div>

                    <div className={s.arrow}>
                        <img src='right-arrow.png' alt='' onClick={addComm} />
                        {/* <input type="submit" value="Submit" onClick={addComm} /> */}
                    </div>
>>>>>>> refs/remotes/origin/main
                </form>
            </div>
=======
            <div className={s.hl} />

            <form onSubmit={handleSubmit} className={s.createComment}>
                <div className={s.textArea}>
                    <textarea id="message" onChange={onCommentChange} ref={newCommentElem} ></textarea>
                </div>

                <div className={s.arrow}>
                    <img src='right-arrow.png' alt='' onClick={addComm} />
                    {/* <input type="submit" value="Submit" onClick={addComm} /> */}
                </div>
            </form>

>>>>>>> c9e4d4e02ae8cc9163044fa9830434ff55c3f139
        </div>
    );
}

export default Comments;