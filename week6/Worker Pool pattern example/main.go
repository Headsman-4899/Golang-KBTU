package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan string, results chan<- string) {
	for v := range jobs {
		fmt.Println("worker with id: ", id, " start job ", v)
		time.Sleep(2 * time.Second)
		results <- v + " is done."
	}
}

func main() {
	numberOfWokers := 3
	jobList := []string{"job1", "job2", "job3", "job4", "job5"}
	numberOfJobs := len(jobList)

	jobs := make(chan string, numberOfJobs)
	results := make(chan string, numberOfJobs)

	for i := 0; i < numberOfWokers; i++ {
		go worker(i, jobs, results)
	}

	for _, job := range jobList {
		jobs <- job
	}
	close(jobs)

	for a := 1; a <= numberOfJobs; a++ {
		fmt.Println(<-results)
	}

	fmt.Println("All jobs have been finished.")
}
