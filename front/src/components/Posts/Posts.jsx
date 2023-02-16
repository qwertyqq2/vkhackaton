import React from "react";
import Post from "./Post/Post";
import s from "./Posts.module.css";
import { Route, Routes } from 'react-router-dom';
import PrePost from "./Post/PrePost";

const Posts = (props) => {

    return (
        <div className={s.posts}>
            {props.htmlContent.map((item, ind) =>
                <PrePost isVisible={true} source={item} index={ind + 1} />
            )}

            <Routes>
                {props.htmlContent.map((item, index) =>
                    <Route path={'/post' + (index + 1)} element={<Post source={item} />} />
                )}
                {/* <Route path='/post' element={<Post source={html} />} /> */}
            </Routes>

        </div>
    );
}

export default Posts;