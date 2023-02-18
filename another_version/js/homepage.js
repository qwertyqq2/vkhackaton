const testButton = document.getElementById("test-button")
const buttonTime = document.getElementById("button-time")
const name = document.getElementById("name")
const askButton = document.getElementById("ask-button")
const askTime = document.getElementById("ask-time")
const formButton = document.getElementById("form-button")
const askPost = document.getElementById("post-added")
const form = document.querySelector('form');
const postCommentButton = document.getElementById("post-comment")

askButton.addEventListener("click", function () {
    let data = {
        Name: name.value,
        Time: new Date().toLocaleString("en-IE"),
    };
    fetch("/get_time", {
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            console.log(result)
            askTime.textContent = "Backend processing took " + result["Duration"] + " and ended at " + result["Time"]
        });
    }).catch((error) => {
        console.log(error)
    });
})

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
      fetch('/create_post', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
      })
      .then(response => {
        if (response.ok) {
          response.text().then(function (data) {
            let result = JSON.parse(data);
            console.log(result)
            askPost.textContent = "Post created at: " + result["Time"]
          });
        } else {
          throw new Error('Network response was not ok');
        }
      })
      .catch(error => {
        console.error('Error:', error);
      });
});

postCommentButton.addEventListener("click", function () {
  let data = {
      Name: name.value,
      Time: new Date().toLocaleString("en-IE"),
  }; // ??
  fetch("/post_comment", {
      headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
      },
      method: "POST",
      body: JSON.stringify(data)
  // }).then((response) => {
  //     response.text().then(function (data) {
  //         let result = JSON.parse(data);
  //         console.log(result)
  //         askTime.textContent = "Backend processing took " + result["Duration"] + " and ended at " + result["Time"]
  //     });
   }) 
    .catch((error) => {
      console.log(error)
  });
})

const time = new EventSource('/time');
time.addEventListener('time', (e) => {
    document.getElementById("actual-time").innerHTML = "Actual time using SSE: " + e.data;

}, false);