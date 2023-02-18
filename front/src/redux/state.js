// load HTMLs
const loadHtmlContent = [];

// for (let i = 1; i <= 3; i++) {
    // eslint-disable-next-line
    let htmlModule = require(`raw-loader!../data/` + 1 + `.html`);
    let html = htmlModule.default;

    loadHtmlContent.push(html);
// }
//

let state = {
  getHtmlContent: loadHtmlContent,
  comments: [
    { bindPostId: 1, id: 1, message: "Test comment" },
    { bindPostId: 2, id: 2, message: "T_T"},
    { bindPostId: 1, id: 3, message: "Second comm" }
  ],
  likes : [
    { bindPostId: 1, likePressed: false },
    { bindPostId: 2, likePressed: false }
  ],
  likeCount : [
    { bindPostId: 1, count: 0 },
    { bindPostId: 2, count: 0 }
  ],
  account : [
    { city: "Moscow", age: -1 }
  ],
};

let rerenderEntireTree = () => {};

export let addPost = () => {
    let newPost = {
        
    };

    rerenderEntireTree();
};

export const subscribe = (observer) => {
    rerenderEntireTree = observer;
};

export default state;