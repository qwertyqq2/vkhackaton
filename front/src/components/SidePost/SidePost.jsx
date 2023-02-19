import React from 'react';
import s from './SidePost.module.css';
import { loadImages } from '../../redux/state';

const SidePost = () => {
    return (
        <div className={s.post}>
            <div className={s.title}>
                <p>
                    Most Popular
                </p>
            </div>
            <div className={s.pict}>
                {loadImages.map((image) => (
                    <img src={image} alt="" key={image.id} className={s.item} />
                ))}
            </div>
        </div>
    );
}

export default SidePost;