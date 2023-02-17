// load HTMLs
const loadHtmlContent = [];

for (let i = 1; i <= 3; i++) {
    // eslint-disable-next-line
    let htmlModule = require(`raw-loader!../data/htmlExample` + i + `.html`);
    let html = htmlModule.default;

    loadHtmlContent.push(html);
}
//

let state = {
  getHtmlContent: loadHtmlContent,
  comments: [
    { bindPostId: 1, id: 1, message: "Test comment" },
    { bindPostId: 2, id: 2, message: "T_T"}
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