import React, { useState, useRef } from 'react';
import classes from './Create.module.css';

const Create = () => {

    const [inputValue, setInputValue] = useState('');

    const handleInputChange = (event) => {
        setInputValue(event.target.value);
    };

    const handleSubmit = (event) => {
        event.preventDefault();
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
            fetch('http://localhost:3001/create_post', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            })
                .then(response => {
                    if (response.ok) {
                        // handle the response from the Go backend here
                    } else {
                        throw new Error('Network response was not ok');
                    }
                })
                .catch(error => {
                    console.error('Error:', error);
                });
        });
    };

    return (
        <div>
            <form onSubmit={handleSubmit}>
                <label for="name">Name:</label>
                <input type="text" id="name" name="name" onChange={handleInputChange} /><br />

                <label for="message">Message:</label>
                <textarea id="message" name="message" onChange={handleInputChange} ></textarea><br />

                <label for="interests">Interests:</label>
                <input type="checkbox" id="Minting" name="interests" value="Minting" onChange={handleInputChange} />
                <label for="Minting">Minting</label>
                <input type="checkbox" id="Blockhain" name="interests" value="Blockhain" onChange={handleInputChange} />
                <label for="Blockhain">Blockhain</label>
                <input type="checkbox" id="Economics" name="interests" value="Economics" onChange={handleInputChange}/>
                <label for="Economics">Economics</label><br />

                <input type="submit" value="Submit" onChange={handleInputChange} />
            </form>
        </div>
    );
}

export default Create;
