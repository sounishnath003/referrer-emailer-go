package core

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/sounishnath003/customgo-mailer-service/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type WorkerPool struct {
	Concurrency int
	JobQueue    chan repository.JobQueue

	wg       sync.WaitGroup
	mu       sync.Mutex
	lo       *slog.Logger
	dbClient *repository.MongoDBClient
}

func NewWorkerPool(dbClient *repository.MongoDBClient, concurrency, bufferSize int) *WorkerPool {
	return &WorkerPool{
		Concurrency: concurrency,
		JobQueue:    make(chan repository.JobQueue, bufferSize),
		dbClient:    dbClient,
		lo:          slog.Default(),
	}
}

func (wp *WorkerPool) StartWorkers() {
	for i := 0; i < wp.Concurrency; i++ {
		wp.wg.Add(1)
		go wp.worker()
		wp.lo.Info("[WORKERPOOL]:", "initilized.worker", i)
	}
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for job := range wp.JobQueue {
		wp.processJob(job)
	}
}

func (wp *WorkerPool) processJob(job repository.JobQueue) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	currentTime := time.Now()
	wp.lo.Info("[WORKERPOOL]: started processing new", "userEmail", job.UserEmailAddress, "job_type", job.JobType, "currentTime", currentTime)

	switch job.JobType {
	case "EXTRACT_CONTENT":
		// extract content from resume
		
		wp.lo.Info("[WORKERPOOL]:", "userEmail", job.UserEmailAddress, "job_type", job.JobType, "completedAt", time.Now())
		job.JobType = "GENERATE_PROFILE_SUMMARY"
	case "GENERATE_PROFILE_SUMMARY":
		// generate profile summary
		wp.lo.Info("[WORKERPOOL]:", "userEmail", job.UserEmailAddress, "job_type", job.JobType, "completedAt", time.Now())
		job.JobType = "RESUME_TO_JSON"
	case "RESUME_TO_JSON":
		// generate the JSON struct for the resume
		wp.lo.Info("[WORKERPOOL]:", "userEmail", job.UserEmailAddress, "job_type", job.JobType, "completedAt", time.Now())
		job.JobType = "UPDATE_RESUME_DOCUMENT"
	case "UPDATE_RESUME_DOCUMENT":
		wp.lo.Info("[WORKERPOOL]:", "userEmail", job.UserEmailAddress, "job_type", job.JobType, "completedAt", time.Now())
		job.JobType = "EMAIL_NOTIFICATION"
	case "EMAIL_NOTIFICATION":
		// send the email to user
		wp.lo.Info("[WORKERPOOL]:", "userEmail", job.UserEmailAddress, "job_type", job.JobType, "status", "COMPLETED", "completedAt", time.Now())
		job.Status = "COMPLETED"
	default:
		wp.lo.Error("[WORKERPOOL]: unknown", "job type", job.JobType, "error", fmt.Errorf("job type not recognized: email=%s, job_type=%s\n", job.UserEmailAddress, job.JobType))
		job.Status = "FAILED"
	}

	wp.lo.Info("[WORKERPOOL]: completed processing", "userEmail", job.UserEmailAddress, "job_type", job.JobType, "totalElapsedTime", time.Since(currentTime))

	// Update job status in MongoDB
	job.UpdatedAt = time.Now()

	collection := wp.dbClient.Database("referrer").Collection("job_queues")
	m, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"userEmailAddress": job.UserEmailAddress, "status": "IN_PROGRESS"},
		bson.M{"$set": bson.M{"status": job.Status, "jobType": job.JobType, "updatedAt": job.UpdatedAt}},
	)

	if m.MatchedCount == 0 || err != nil {
		wp.lo.Error("[WORKERPOOL]: failed to update job status in MongoDB:", "job", job, "error", err)
	}

	// If job is not completed, push it back to the job queue for the next stages
	if job.Status != "COMPLETED" && job.Status != "FAILED" {
		wp.JobQueue <- job
	}

}

func (wp *WorkerPool) ListenForThePendingJobs() {
	collection := wp.dbClient.Database("referrer").Collection("job_queues")

	for {
		var job repository.JobQueue
		err := collection.FindOneAndUpdate(
			context.TODO(),
			bson.M{"status": "PENDING"},
			bson.M{"$set": bson.M{"status": "IN_PROGRESS", "updatedAt": time.Now()}},
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		).Decode(&job)

		if err == nil {
			wp.JobQueue <- job
		} else if err != mongo.ErrNoDocuments {
			wp.lo.Error("[WORKERPOOL]: not able to pull pending jobs from job-queues:", "error", err)
		}

		time.Sleep(1 * time.Second)
	}
}

func (wp *WorkerPool) Wait() {
	defer close(wp.JobQueue)
	wp.wg.Wait()
}
