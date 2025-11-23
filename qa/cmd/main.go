package main

import (
	"log"
	"net/http"
	"os"

	"qa/internal/db"
	"qa/internal/handler"
	"qa/internal/model"
	"qa/internal/repo"
)

func main() {
	gdb, err := db.ConnectFromEnv()
	if err != nil {
		log.Fatal("db connect:", err)
	}

	if err := gdb.AutoMigrate(&model.Question{}, &model.Answer{}); err != nil {
		log.Fatal("migrate:", err)
	}

	r := repo.New(gdb)
	h := handler.NewHandler(r)
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, h.ServeMux()); err != nil {
		log.Fatal(err)
	}
}
