// load HTMLs
const loadHtmlContent = [];

export const loadImages = ["love.png", "more.png", 
"photo1676760335.jpeg", "backgr.jpeg"];

export let postsCount = 100; // Четыре

let state = {
  getHtmlContent: loadHtmlContent,
  comments: [
    { bindPostId: 1, id: 1, message: "Test comment" },
    { bindPostId: 2, id: 2, message: "T_T" },
    { bindPostId: 1, id: 3, message: "Second comm" }
  ],
  likes: [
    
  ],
  account: [
    { city: "Moscow", age: -1 }
]
};

for (let i = 1; i <= postsCount; i++) {
  try {
  // eslint-disable-next-line
  let htmlModule = require(`raw-loader!../data/htmlExample` + i + `.html`);
  let html = htmlModule.default;

  loadHtmlContent.push(html);
  state.likes.push({ likePressed: false, count: 0 });
  } catch (error) {
    continue;
  }
}

let rerenderEntireTree = () => { };

export let addComment = (params) => {
  state.comments.push({ bindPostId: params.postId, id: params.id, message: params.message });

  //rerenderEntireTree(state);
}

export const subscribe = (observer) => {
  rerenderEntireTree = observer;
};

export default state;
