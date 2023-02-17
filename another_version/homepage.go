package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"html/template"
	"net/http"
	"time"
	"os"
)

type HomePage struct {
	Time string
}

type TimeDataInput struct {
	Name string
	Time string
}

type TimeDataOutput struct {
	Result   string
	Text     string
	Time     string
	Duration string
}

type FormData struct {
	Name string   `json:"name"`
	Message string `json:"message"`
	Interests []string `json:"interests"`
}

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
    html := "<html><head><title>" + formData.Name + "</title></head><body><h1>" + formData.Name + "</h1><p>" + formData.Message + "</p></body></html>"

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

	var responseData TimeDataOutput
	responseData.Result = "ok"
	responseData.Text = "everything went smooth"
	responseData.Time = time.Now().Format("02/01/2006, 15:04:05")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(responseData)
}

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
