/*
We will keep three channels, each for result, error and completion and will be performing
a simple mathematical operation, let's say summation of numbers from 1 to N( N is the input number)
*/
package main

import (
	"log"
	"time"
)

// Input structre, containing the array of integers, for which we have to find the summation

type Input struct {
	num []uint64
}

// structure of the input, we will keep input N and the total summtion as outputs
type output struct {
	inputNumber uint64
	output      uint64
}

// We will decleare all the required Methods in below interface
type Response interface {
	addResult(result output)
	setError(err error)
	setCompleted(status bool)
	subscribe(func(result output), func(err error), func(bool))
}

// Response struct conataing all those 3 channels
type FinalResponse struct {
	isCompleted      bool
	errorMsg         string
	results          []output
	errChannel       chan error
	completedChannel chan bool
	resChannel       chan output
}

// Method to find the summation of number from 1 to given n : like 1+2+3+4+ ... +n
func findSum(num uint64) (*output, error) {
	n1 := num * (num + 1)
	n1 /= 2
	return &output{inputNumber: num, output: (n1)}, nil
}

// Method for calling the summation function and publishing the respective data to repective channels
func execute(input Input) *FinalResponse {
	resp := &FinalResponse{
		resChannel:       make(chan output, 10),
		errChannel:       make(chan error),
		completedChannel: make(chan bool),
		results:          []output{},
	}

	go func(resp *FinalResponse) {
		for _, v := range input.num {
			res, err := findSum(v)
			if err != nil {
				log.Println("Error while finding the sum")
				resp.errChannel <- err
			} else {
				resp.resChannel <- *res
			}
		}
		close(resp.resChannel)
		// Operation has been finished so, sending true to the completedChannel
		resp.completedChannel <- true
	}(resp)
	return resp
}

// Function for appeding res to the Result array
func (r *FinalResponse) addResult(res output) {
	r.results = append(r.results, res)
}

// Setting error
func (r *FinalResponse) setError(err error) {
	r.errorMsg = err.Error()
}

func (r *FinalResponse) setCompleted(val bool) {
	if r.completedChannel == nil {
		r.completedChannel = make(chan bool)
	}
	r.isCompleted = val
}

// Function for subscribing respective channels and doing the things not synchronous (non-blocking way)
func (resp *FinalResponse) subscribe(Addfunction func(input output), Errfunction func(err error), Completefunction func(status bool)) {
	go resp.solve(Addfunction, Errfunction, Completefunction)
}

func (resp *FinalResponse) solve(Addfunction func(input output), Errfunction func(err error), Completefunction func(status bool)) {
	done := false
	for {
		select {
		case res, ok := <-resp.resChannel:
			{
				Addfunction(res)
				if !ok {
					done = true
					break
				}
			}
		case err := <-resp.errChannel:
			Errfunction(err)
		}

		if done {
			break
		}
	}
	Completefunction(<-resp.completedChannel)

}

func main() {
	input := Input{
		num: []uint64{1, 2, 5, 10, 20, 100, 1000, 5000, 10000},
	}
	resp := execute(input)
	log.Printf("%+v", resp)
	resp.subscribe(resp.addResult, resp.setError, resp.setCompleted)
	log.Println("Subscribed to the channels")
	for !resp.isCompleted {
		log.Println("listening to results ...")
		time.Sleep(1 * time.Second)
	}
	log.Printf("%+v", resp)
}
