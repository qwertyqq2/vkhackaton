package main

import (
  "encoding/json"
  "fmt"
  "github.com/julienschmidt/httprouter"
  "github.com/julienschmidt/sse"
  "html/template"
  "net/http"
  "time"
  "strconv"
  "os"
)

type HomePage struct {
  Time string
}

type TimeDataInput struct {
  Name string
  Time string
}

type CommentInput struct {
  //Author string `json:"author"`
  Content string `json:"content"`
  //PostID int `json:"postID"`
}

type TimeDataOutput struct {
  Result   string
  Text     string
  Time     string
  Duration string
}

type PostDataOutput struct{
  Result string
  Time string
  Count int
}

type FormData struct {
  Name string   `json:"name"`
  Message string `json:"message"`
  Interests []string `json:"interests"`
}

type Comment struct {
  Author string 
  Content string 
  Time string
}

type ResponseData struct {
  Result string
  Text string
  Time string
}

var commentMap map[int][]Comment
var postsCount = 3
// func createComment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//   decoder := json.NewDecoder(r.Body)
//   var commentInput CommentInput
//   err := decoder.Decode(&commentInput)
//   if err != nil {
//     http.Error(w, err.Error(), http.StatusBadRequest)
//     return
//   }

//   // make html file from message and name
//   fmt.Println(commentInput.Author)
//   fmt.Println(commentInput.Content)
//   fmt.Println(commentInput.PostID)

//   var comment Comment;
//   comment.Author = commentInput.Author;
//   comment.Content = commentInput.Content;
//   comment.Time = time.Now().Format("02/01/2006, 15:04:05");

//   //writeContentHTML(commentInput.Author, commentInput.Content)
//   commentMap[commentInput.PostID].append(comment)
//   var responseData PostDataOutput
//   responseData.Result = "ok"
//   responseData.Text = "everything went smooth"
//   responseData.Time = time.Now().Format("02/01/2006, 15:04:05")
  
//   setupCORS(&w, r);

//   w.Header().Set("Content-Type", "application/json")
//   _ = json.NewEncoder(w).Encode(responseData)
// }

func createPost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  decoder := json.NewDecoder(r.Body)
  var formData FormData
  err := decoder.Decode(&formData)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  // make html file from message and name
  fmt.Println(formData.Message)
  fmt.Println(formData.Name)

    // Create the HTML string by wrapping the header and content in HTML tags
    html := "<div class ='title'>" + formData.Name + "</div>" + 
      "<div class ='content'>" + formData.Message + "</div>"

    // Create a new file to write to
	postsCount++
    // file, err := os.Create("../front/src/data/htmlExample" + strconv.Itoa(postsCount) + ".html" // for danil
	file, err := os.Create("../src/data/htmlExample" + strconv.Itoa(postsCount) + ".html" // for general
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }

    // Write the HTML string to the file
    _, err = file.WriteString(html)
    if err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }

    // Close the file
    err = file.Close()
    if err != nil {
        fmt.Println("Error closing file:", err)
        return
    }

    fmt.Println("HTML file created successfully!")

  var responseData PostDataOutput
  responseData.Result = "ok"
  responseData.Count = postsCount
  responseData.Time = time.Now().Format("02/01/2006, 15:04:05")
  
  setupCORS(&w, r);

  w.Header().Set("Content-Type", "application/json")
  _ = json.NewEncoder(w).Encode(responseData)
}

// func getRecentPost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//   /*var data TimeDataInput
//   err := json.NewDecoder(r.Body).Decode(&data) 
//   if err != nil {
//     fmt.Println(err.Error())
//     var responseData TimeDataOutput
//     responseData.Result = "no"
//     responseData.Text = "problem with user json data"
//     w.Header().Set("Content-Type", "application/json")
//     _ = json.NewEncoder(w).Encode(responseData)
//     return
//   } 
//   fmt.Println(data.Name)
//   fmt.Println(data.Time) */
//   time := time.Now()
//   var postData PostDataOutput
//   // как то получить недавний пост из хтмл
//   postData.Result = "ok"
//   postData.Content = "lorem ipsum dolor"
//   postData.Time = time.Now().Format("02/01/2006, 15:04:05")
//   postData.Author = "Vladimir Belyaev"
//   w.Header().Set("Content-Type", "application/json")
//   _ = json.NewEncoder(w).Encode(postData)
// }

func getTime(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
  var data TimeDataInput
  err := json.NewDecoder(request.Body).Decode(&data)
  if err != nil {
    fmt.Println(err.Error())
    var responseData TimeDataOutput
    responseData.Result = "no"
    responseData.Text = "problem with user json data"
    writer.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(writer).Encode(responseData)
    return
  }
  fmt.Println(data.Name)
  fmt.Println(data.Time)
  timer := time.Now()
  time.Sleep(1 * time.Second)
  end := time.Since(timer)
  fmt.Println("processing takes: " + end.String())
  var responseData TimeDataOutput
  responseData.Result = "ok"
  responseData.Text = "everything went smooth"
  responseData.Time = time.Now().Format("02/01/2006, 15:04:05")
  responseData.Duration = end.String()
  writer.Header().Set("Content-Type", "application/json")
  _ = json.NewEncoder(writer).Encode(responseData)
}

func serveHomepage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
  writingSync.Lock()
  programIsRunning = true
  writingSync.Unlock()
  var homepage HomePage
  homepage.Time = time.Now().Format("02/01/2006, 15:04:05")
  tmpl := template.Must(template.ParseFiles("html/homepage.html"))
  _ = tmpl.Execute(writer, homepage)
  writingSync.Lock()
  programIsRunning = false
  writingSync.Unlock()
}

func streamTime(timer *sse.Streamer) {
  fmt.Println("Streaming time started")
  for serviceIsRunning {
    timer.SendString("", "time", time.Now().Format("02/01/2006, 15:04:05"))
    time.Sleep(1 * time.Second)
  }
}

func setupCORS(w *http.ResponseWriter, r *http.Request) {
  (*w).Header().Set("Access-Control-Allow-Origin", "*")
  (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
  (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
  if r.Method == "OPTIONS" {
        (*w).WriteHeader(http.StatusOK)
        return
    }
}

func writeContentHTML(author string, content string) {
  // Create the HTML string by wrapping the header and content in HTML tags
    html := "<html><head><title>" + author + "</title></head><body><h1>" + content + "</h1><p>" + content + "</p></body></html>"

    // Create a new file to write to
    file, err := os.Create("html/mypage.html")
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }

    // Write the HTML string to the file
    _, err = file.WriteString(html)
    if err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }

    // Close the file
    err = file.Close()
    if err != nil {
        fmt.Println("Error closing file:", err)
        return
    }

    fmt.Println("HTML file created successfully!")
}