package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationEmail NotificationType = "email"
	NotificationSMS   NotificationType = "sms"
)

type NotificationTask struct {
	UserID  uuid.UUID
	Message string
	Type    NotificationType
}

type NotificationWorker struct {
	taskChan    chan NotificationTask
	workerCount int
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewNotificationWorker(workerCount int) *NotificationWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &NotificationWorker{
		taskChan:    make(chan NotificationTask, 100),
		workerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (w *NotificationWorker) Start() {
	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.worker(i)
	}
	log.Printf("Notification worker started with %d workers", w.workerCount)
}

func (w *NotificationWorker) worker(id int) {
	defer w.wg.Done()
	log.Printf("Worker %d started", id)

	for {
		select {
		case task := <-w.taskChan:
			w.processNotification(task, id)
		case <-w.ctx.Done():
			log.Printf("Worker %d stopping", id)
			return
		}
	}
}

func (w *NotificationWorker) processNotification(task NotificationTask, workerID int) {
	// Simulate notification processing
	log.Printf("[Worker %d] Processing notification: Type=%s, UserID=%s, Message=%s",
		workerID, task.Type, task.UserID, task.Message)

	// Simulate network delay
	time.Sleep(100 * time.Millisecond)

	// In a real implementation, this would:
	// - Send email via SMTP service
	// - Send SMS via SMS gateway
	// - Push notification via FCM/APNS
	log.Printf("[Worker %d] Notification sent successfully", workerID)
}

func (w *NotificationWorker) Enqueue(task NotificationTask) {
	select {
	case w.taskChan <- task:
		// Task queued successfully
	case <-time.After(1 * time.Second):
		log.Printf("Warning: Notification queue full, dropping task for user %s", task.UserID)
	}
}

func (w *NotificationWorker) Shutdown() {
	log.Println("Shutting down notification worker...")
	w.cancel()

	// Wait for all workers to finish processing current tasks
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All notification workers stopped")
	case <-time.After(10 * time.Second):
		log.Println("Timeout waiting for workers to stop")
	}
}

