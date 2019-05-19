package go_docker

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main()  {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)

	port := os.Getenv("PORT")

	srv := &http.Server{
		Handler:r,
		Addr: ":"+port,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	LogFileLocation := os.Getenv("LOG_FILE_LOCATION")

	if LogFileLocation != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   LogFileLocation,
			MaxSize:    500, //megabytes
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		})
	}

	go func() {
		log.Println("Starting server...")
		log.Println("Server running on port"+port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}

	}()

	//Gracefully shutdown
	waitForShutdown(srv)

}

func waitForShutdown(server *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<- interruptChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Shutting down")
	os.Exit(0)
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}
	log.Printf("Received request for %s\n", name)
	_, err := w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
	if err != nil {
		log.Fatal(err)
	}
}

