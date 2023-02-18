package main

import (
  "fmt"
  "github.com/julienschmidt/httprouter"
  //"github.com/julienschmidt/sse"
  "github.com/rs/cors"
  "github.com/kardianos/service"
  "net/http"
  "os"
  "sync"
  "time"
)

var (
  serviceIsRunning bool
  programIsRunning bool
  writingSync      sync.Mutex
)

const serviceName = "Medium service"
const serviceDescription = "Simple service, just for fun"

type program struct{}

func (p program) Start(s service.Service) error {
  fmt.Println(s.String() + " started")
  writingSync.Lock()
  serviceIsRunning = true
  writingSync.Unlock()
  go p.run()
  return nil
}

func (p program) Stop(s service.Service) error {
  writingSync.Lock()
  serviceIsRunning = false
  writingSync.Unlock()
  for programIsRunning {
    fmt.Println(s.String() + " stopping...")
    time.Sleep(1 * time.Second)
  }
  fmt.Println(s.String() + " stopped")
  return nil
}

func (p program) run() {
  router := httprouter.New()
  handler := cors.Default().Handler(router)
  //timer := sse.New()
  //router.ServeFiles("/js/*filepath", http.Dir("js"))
  //router.ServeFiles("/css/*filepath", http.Dir("css"))
  router.GET("/", serveHomepage)

  // router.POST("/get_time", getTime)
  router.POST("/create_post", createPost) // создать пост

  //router.POST("/post_comment", createComment) // создать коммент
  //router.GET("/get_new_post", getRecentPost) // получить новый пост
  //router.PATCH("/update_likes", updateLikes) // обновить счетчик лайков
  //router.Handler("GET", "/time", timer)

  //go streamTime(timer)
  err := http.ListenAndServe(":3001", handler)
  if err != nil {
    fmt.Println("Problem starting web server: " + err.Error())
    os.Exit(-1)
  }
}

func main() {
  serviceConfig := &service.Config{
    Name:        serviceName,
    DisplayName: serviceName,
    Description: serviceDescription,
  }
  prg := &program{}
  s, err := service.New(prg, serviceConfig)
  if err != nil {
    fmt.Println("Cannot create the service: " + err.Error())
  }
  err = s.Run()
  if err != nil {
    fmt.Println("Cannot start the service: " + err.Error())
  }
}