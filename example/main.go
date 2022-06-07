package main

import (
	"context"
	"errors"
	"fmt"
	workerpool "github.com/fullaheadcoco/go-example-workerpool"
	"net/http"
	"time"
)

func main() {

	wp := workerpool.New(3)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go wp.Run(ctx)
	go checkJobResult(wp)

	http.HandleFunc("/insert", insert(wp))
	http.ListenAndServe(":8090", nil)
}

func insert(wp workerpool.WorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wp.Insert(newJob())
	}
}

func newJob() workerpool.Job {
	return workerpool.Job{
		workerpool.JobDescriptor{
			ID:    "1",
			JType: "anyType",
			Metadata: workerpool.JobMetadata{
				"craetedAt": time.Now(),
			},
		},
		func(ctx context.Context, args interface{}) (interface{}, error) {
			argVal, ok := args.(int)
			if !ok {
				return nil, errors.New("Error")
			}

			return argVal * 2, nil
		},
		2,
	}
}

func checkJobResult(wp workerpool.WorkerPool) {

	for {
		select {
		case r, ok := <-wp.Results():
			if !ok {
				continue
			}
			fmt.Printf("Result: %v\n", r)

		case <-wp.Done:
			fmt.Println("DoneDoneDoneDoneDoneDoneDone")
			return
		default:
		}
	}
}
