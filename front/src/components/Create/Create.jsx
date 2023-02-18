import React, { useState } from 'react';
import s from './Create.module.css';
import { postsCount } from '../../redux/state';

const Create = () => {

    const [inputValue, setInputValue] = useState('');

    const handleInputChange = (event) => {
        setInputValue(event.target.value);

    };

    const handleSubmit = (event) => {
        event.preventDefault();
        debugger;
        // do something with the input value
        const form = document.querySelector('form');

        console.log("1");

        form.addEventListener('submit', (event) => {
            event.preventDefault();

            const name = form.elements['name'].value;
            const message = form.elements['message'].value;
            const interests = [];

            const checkboxes = form.querySelectorAll('input[type="checkbox"]');
            checkboxes.forEach((checkbox) => {
                if (checkbox.checked) {
                    interests.push(checkbox.value);
                }
            });
            const data = {
                name: name,
                message: message,
                interests: interests
            };
            debugger;
            fetch('http://localhost:3001/create_post', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            })
                .then(response => {
                    if (response.ok) {
                        debugger;
                        let temp = JSON.parse(data);
                        debugger;
                        console.log(postsCount);
                        postsCount = temp["Count"];
                    } else {
                        debugger;
                        throw new Error('Network response was not ok');
                    }
                })
                .catch(error => {
                    debugger;
                    console.error('Error:', error);
                });
        });
    };

    return (
        <div className={s.createPost}>
            <form onSubmit={handleSubmit}>
                <div className={s.postName}>
                    <label for="name">Name:</label><br />
                    <input type="text" id="name" name="name" onChange={handleInputChange} /><br />
                </div>

                <div className={s.postContent}>
                    <label for="message">Message:</label><br />
                    <textarea id="message" name="message" onChange={handleInputChange} ></textarea><br />
                </div>

                <div className={s.hz}>
                    <label for="interests">Interests:</label>
                    <ul>
                        <li>
                            <input type="checkbox" id="Minting" name="interests" value="Minting" onChange={handleInputChange} />
                            <label for="Minting">Minting</label>
                        </li>
                        <li>
                            <input type="checkbox" id="Blockhain" name="interests" value="Blockhain" onChange={handleInputChange} />
                            <label for="Blockhain">Blockhain</label>
                        </li>
                        <li>
                            <input type="checkbox" id="Economics" name="interests" value="Economics" onChange={handleInputChange} />
                            <label for="Economics">Economics</label><br />
                        </li>
                    </ul>
                </div>

                <div className={s.submit}>
                    <input type="submit" value="Submit" onChange={handleInputChange} />
                </div>
            </form>
        </div>
    );
}

export default Create;
