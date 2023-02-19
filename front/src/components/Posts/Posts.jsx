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

    return (
        <div className={s.posts}>
            {props.htmlContent.map((item, ind) =>
                <PrePost
                    source={item}
                    index={ind + 1}
                    isVisible={{ v: visibleElement, h: handleVisibleElement }}
                    likes={props.likes[ind]} />
            )}

            <Routes>
                {props.htmlContent.map((item, index) =>
                    <Route
                        path={'/post' + (index + 1)}
                        element={<Post source={item}
                            comments={props.comments}
                            postId={index + 1}
                            isVisible={{ v: visibleElement, h: handleVisibleElement }}
                            likes={props.likes[index]} />} />
                )}
            </Routes>
        </div>
    );
}

export default Posts;