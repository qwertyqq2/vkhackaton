import React, { useState } from "react";
import Post from "./Post/Post";
import s from "./Posts.module.css";
import { Route, Routes } from 'react-router-dom';
import PrePost from "./Post/PrePost";

const Posts = (props) => {

    const [visibleElement, setVisibleElement] = useState(true);

    const handleVisibleElement = () => {
        setVisibleElement(!visibleElement);
    }

    const [likeElement, setLikeElement] = useState(true);

    const handleLikeClick = () => {
        setLikeElement(!likeElement);
    }

    // console.log(visibleElement);

    return (
        <div className={s.posts}>
            {props.htmlContent.map((item, ind) =>
                <PrePost
                    source={item} 
                    index={ind + 1} 
                    isVisible={{ v: visibleElement, h: handleVisibleElement }}
                    isLikePressed={{ l: likeElement, h: handleLikeClick }} />
            )}

            <Routes>
                {props.htmlContent.map((item, index) =>
                    <Route
                        path={'/post' + (index + 1)}
                        element={<Post source={item}
                            comments={props.comments}
                            postId={index + 1}
                            isVisible={{ v: visibleElement, h: handleVisibleElement }}
                            isLikePressed={{ l: likeElement, h: handleLikeClick }} />} />
                )}
                {/* <Route path='/post' element={<Post source={html} />} /> */}
                {/* className={!visibleElement ? '' : (' ' + s.hidden)} */}
            </Routes>

        </div>
    );
}

export default Posts;