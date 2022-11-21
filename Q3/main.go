package main

import (
	"log"
	"math/rand"
	"time"
)

// Constant time duration of 20 sec
const timeout time.Duration = time.Duration(20 * time.Second)

// Enum for Task Status
type TaskStatus string

const (
	StatusNew       TaskStatus = "new"
	StatusCompleted TaskStatus = "completed"
	StatusTimeout   TaskStatus = "timeout"
)

// Task Structure
type Task struct {
	Id          string
	IsCompleted bool
	Status      TaskStatus
}

// Structure of task data for processing in the queue
type TaskData struct {
	createdAt time.Time
	task      *Task
}

func main() {

	// Take some task datas for processing
	tasks := []Task{
		{Id: "1", IsCompleted: false, Status: StatusNew},
		{Id: "2", IsCompleted: false, Status: StatusNew},
		{Id: "3", IsCompleted: false, Status: StatusNew},
		{Id: "4", IsCompleted: false, Status: StatusNew},
		{Id: "5", IsCompleted: false, Status: StatusNew},
		{Id: "6", IsCompleted: false, Status: StatusNew},
		{Id: "7", IsCompleted: false, Status: StatusNew},
		{Id: "8", IsCompleted: false, Status: StatusNew},
		{Id: "9", IsCompleted: false, Status: StatusNew},
		{Id: "10", IsCompleted: false, Status: StatusNew},
	}

	// Go routine for taking a random integer in range of [0,n), where n is the length of tasks
	// and if the status is not completed, marked it as completed
	go func(tasks []Task) {
		for {
			idx := rand.Intn(len(tasks))
			if !tasks[idx].IsCompleted {
				tasks[idx].IsCompleted = true
			}
			time.Sleep(3 * time.Second)
		}
	}(tasks)

	// Taking one channel, we will be publishing the task data on that
	queue := make(chan *TaskData, 20)

	// Publishing task datas to the channel along with the createdAt time stamp
	for i := range tasks {

		queue <- &TaskData{createdAt: time.Now(), task: &tasks[i]}
	}

	// processing the queue, having one array of string called completedTasks, just to
	// show the tasks, that has been completed at a moment
	var completedTasks []string
	for len(queue) > 0 {
		temp := <-queue
		t := temp.task
		if t.IsCompleted {
			t.Status = StatusCompleted
			completedTasks = append(completedTasks, t.Id)
			log.Printf("task: %s is completed, [completed tasks are]: %v)", t.Id, completedTasks)
		} else {
			if time.Since(temp.createdAt) > timeout {
				log.Printf("Timeout, removing it from queu: %s", t.Id)
				t.Status = StatusTimeout
				continue
			}
			// if task is not completed then push back the task to the channel
			queue <- temp
		}
		time.Sleep(500 * time.Millisecond)
	}

	log.Println(tasks)
}
