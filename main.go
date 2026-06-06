package main

import (
	"fmt"
	"time"
)

type Order struct {
	Uuid     string
	Type     string
	CookTime int
	IsReady  bool
}

const tick = 1 * time.Second

func worker(id int, jobs <-chan Order, results chan<- Order) {
	for job := range jobs {

		sleep := time.Duration(job.CookTime) * tick

		fmt.Printf("Cook %d prepares %s (%s) - %v\n", id, job.Uuid, job.Type, sleep)

		job.IsReady = true

		results <- job
	}
}

func createKitchen(numWorkers int, jobs <-chan Order, results chan<- Order) {
	for i := 1; i <= numWorkers; i++ {
		go worker(i, jobs, results)
	}
}

func runKitchen(label string, numWorkers int, orders []Order) {
	fmt.Printf("===== %s: %d cooks, %d orders =====\n ", label, numWorkers, len(orders))
	start := time.Now()

	jobs := make(chan Order, len(orders))
	results := make(chan Order, len(orders))

	createKitchen(numWorkers, jobs, results)

	for _, o := range orders {
		jobs <- o
	}
	close(jobs)

	for range orders {
		<-results
	}
	fmt.Printf("===== ready in %v =====\n\n", time.Since(start))
}

func main() {

	orders := []Order{
		{Uuid: "order-1", Type: "pizza", CookTime: 7},
		{Uuid: "order-2", Type: "salade", CookTime: 2},
		{Uuid: "order-3", Type: "soup", CookTime: 5},
		{Uuid: "order-4", Type: "sandwich", CookTime: 4},
		{Uuid: "order-5", Type: "desert", CookTime: 8},
		{Uuid: "order-6", Type: "french fries", CookTime: 3},
	}

	runKitchen("1 cook", 1, orders)
	runKitchen("3 cooks", 3, orders)
	runKitchen("6 cooks", 6, orders)
}
