package workerpool

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/sounishnath003/customgo-mailer-service/internal/core"
	"github.com/sounishnath003/customgo-mailer-service/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type WorkerPool struct {
	concurrency int
	jobQueue    chan repository.JobQueue

	wg sync.WaitGroup
	mu sync.Mutex
	lo *slog.Logger
	co *core.Core
}

func NewWorkerPool(co *core.Core, concurrency int) *WorkerPool {
	return &WorkerPool{
		co: co,
		lo: slog.Default(),

		concurrency: concurrency,
		jobQueue:    make(chan repository.JobQueue, 10*concurrency),
	}
}

func (wp *WorkerPool) StartWorkers() {
	for i := 0; i < wp.concurrency; i++ {
		wp.wg.Add(1)
		go wp.worker()
		wp.lo.Info("[WORKERPOOL]:", "initilized.worker", i)
	}
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for job := range wp.jobQueue {
		wp.processJob(job)
	}
}

func (wp *WorkerPool) processJob(job repository.JobQueue) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	currentTime := time.Now()
	wp.lo.Info("[WORKERPOOL]: started processing new", "userEmail", job.UserEmailAddress, "job_type", job.JobType, "currentTime", currentTime)

	switch job.JobType {
	case repository.EXTRACT_CONTENT:
		// extract content from resume
		content, err := wp.co.ExtractResumeContentLLM(job.Payload.ResumeURL)
		if err != nil {
			break
		}
		job.Payload.ExtractedContent = content
		wp.lo.Info("[WORKERPOOL]: completed", "userEmail", job.UserEmailAddress, "job_type", job.JobType.String(), "timeElapsed", time.Since(currentTime))
		job.JobType = repository.GENERATE_PROFILE_SUMMARY
	case repository.GENERATE_PROFILE_SUMMARY:
		// generate profile summary
		wp.lo.Info("[WORKERPOOL]:", "userEmail", job.UserEmailAddress, "job_type", job.JobType, "completedAt", time.Now())
		summary, err := wp.co.GenerateProfileSummaryLLM(job.Payload.ExtractedContent)
		if err != nil {
			break
		}
		job.Payload.Summary = summary
		wp.lo.Info("[WORKERPOOL]: completed", "userEmail", job.UserEmailAddress, "job_type", job.JobType.String(), "timeElapsed", time.Since(currentTime))
		job.JobType = repository.UPDATE_RESUME_DOCUMENT
	case repository.UPDATE_RESUME_DOCUMENT:
		u, err := wp.co.DB.GetProfileByEmail(job.UserEmailAddress)
		if err != nil {
			return
		}
		u.ExtractedContent = job.Payload.ExtractedContent
		u.ProfileSummary = job.Payload.Summary
		if err = wp.co.DB.UpdateProfileInformation(u); err != nil {
			wp.lo.Error("error updating the users profile content", "userEmailAddress", job.UserEmailAddress, "error", err)
			break
		}
		wp.lo.Info("[WORKERPOOL]:", "userEmail", job.UserEmailAddress, "job_type", job.JobType.String(), "completedAt", time.Now())
		job.JobType = repository.EMAIL_NOTIFICATION
	case repository.EMAIL_NOTIFICATION:
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

	collection := wp.co.DB.Database("referrer").Collection("job_queues")
	m, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"userEmailAddress": job.UserEmailAddress, "status": "IN_PROGRESS"},
		bson.M{"$set": bson.M{
			"status":    job.Status,
			"jobType":   job.JobType,
			"updatedAt": job.UpdatedAt,
			"payload": bson.M{
				"resumeUrl":        job.Payload.ResumeURL,
				"extractedContent": job.Payload.ExtractedContent,
				"summary":          job.Payload.Summary,
			},
		}},
	)

	if m.MatchedCount == 0 || err != nil {
		wp.lo.Error("[WORKERPOOL]: failed to update job status in MongoDB:", "job", job, "error", err)
	}

	// If job is not completed, push it back to the job queue for the next stages
	if job.Status != "COMPLETED" && job.Status != "FAILED" {
		wp.jobQueue <- job
	}

}

func (wp *WorkerPool) ListenForThePendingJobs() {
	collection := wp.co.DB.Database("referrer").Collection("job_queues")

	for {
		var job repository.JobQueue
		err := collection.FindOneAndUpdate(
			context.TODO(),
			bson.M{"status": "PENDING"},
			bson.M{"$set": bson.M{"status": "IN_PROGRESS", "updatedAt": time.Now()}},
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		).Decode(&job)

		if err == nil {
			wp.jobQueue <- job
		} else if err != mongo.ErrNoDocuments {
			wp.lo.Error("[WORKERPOOL]: not able to pull pending jobs from job-queues:", "error", err)
		}

		time.Sleep(1 * time.Second)
	}
}

func (wp *WorkerPool) Wait() {
	defer close(wp.jobQueue)
	wp.wg.Wait()
}
