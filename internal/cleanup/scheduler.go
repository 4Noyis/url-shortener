package cleanup

import (
	"log"
	"time"

	"github.com/4Noyis/url-shortener/internal/service"
)

type Scheduler struct {
	urlService *service.URLService
	interval   time.Duration
	stopChan   chan bool
}

func NewScheduler(urlService *service.URLService, interval time.Duration) *Scheduler {
	return &Scheduler{
		urlService: urlService,
		interval:   interval,
		stopChan:   make(chan bool),
	}
}

func (s *Scheduler) Start() {
	log.Printf("Starting URL cleanup scheduler with interval: %v", s.interval)
	
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run cleanup immediately on startup
	s.runCleanup()

	for {
		select {
		case <-ticker.C:
			s.runCleanup()
		case <-s.stopChan:
			log.Println("Stopping URL cleanup scheduler")
			return
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.stopChan)
}

func (s *Scheduler) runCleanup() {
	log.Println("Running expired URL cleanup...")
	
	deletedCount, err := s.urlService.CleanupExpiredURLs()
	if err != nil {
		log.Printf("Error during URL cleanup: %v", err)
		return
	}

	if deletedCount > 0 {
		log.Printf("Cleanup completed: deleted %d expired URLs", deletedCount)
	} else {
		log.Println("Cleanup completed: no expired URLs found")
	}
}