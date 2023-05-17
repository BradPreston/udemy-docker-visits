package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-server:6379",
		Password: "",
		DB:       0,
	})
	err := rdb.Set(ctx, "visits", 0, 0).Err()
	if err != nil {
		log.Fatalf("could not set initial visits: %s", err)
	}

	mux := chi.NewRouter()

	mux.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val, err := rdb.Get(ctx, "visits").Result()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf("error getting visits: %s", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("number of visits is %s", val)))
		visits, err := strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf("error converting val to int: %s", err)))
			return
		}
		err = rdb.Set(ctx, "visits", visits+1, 0).Err()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf("error setting visits: %s", err)))
			return
		}
	}))

	fmt.Println("listening on port 8081")
	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("could not start app: %s", err)
	}
}
