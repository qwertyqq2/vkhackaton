import React from 'react';
import Posts from '../Posts/Posts';
import classes from './Profile.module.css';

const Profile = (props) => {
    return (
        <div className={classes.profile}>
            <div className={classes.avatar}>
                <div className={classes.avatarImg}>
                    <img src='mommymonkey.jpeg' alt='' />
                </div>
                <div className={classes.avatarDesc}>
                    <ul>
                        <li>
                            <p>Город: {props.account[0].city}</p>
                        </li>
                        <li>
                            <p>Возраст: {props.account[0].age}</p>
                        </li>
                    </ul>
                </div>
            </div>
            <div className={classes.description}>
                <p>
                Американский предприниматель, инженер и миллиардер. 
                Основатель, генеральный директор и главный инженер компании SpaceX; инвестор, генеральный директор и архитектор продукта компании Tesla; основатель The Boring Company; соучредитель Neuralink и OpenAI; владелец Twitter.
                </p>
            </div>
        </div>
    );
};

export default Profile;