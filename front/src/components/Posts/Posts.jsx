import React from "react";
import Post from "./Post/Post";
import s from "./Posts.module.css";
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import PrePost from "./Post/PrePost";

// eslint-disable-next-line
let htmlModule = require(`raw-loader!./htmlExample` + 1 + `.html`);
let html = htmlModule.default;

const Posts = () => {

    const getHtmlContent = () => {
        let content = [];

        for (let i = 1; i <= 3; i++) {
            // eslint-disable-next-line
            let htmlModule = require(`raw-loader!./htmlExample` + i + `.html`);
            let html = htmlModule.default;

            const item = html[i];
            content.push(i, item);
        }
        return content;
    };

    //let prePostMap = getHtmlContent.map( p => <PrePost source={getHtmlContent.id}/> );

    return (
        <div className={s.posts}>
            <PrePost source={html} />
            <Routes>
                
                <Route path='/post' element={<Post source={html}/>} />
            </Routes>
        </div>
    );
}

export default Posts;