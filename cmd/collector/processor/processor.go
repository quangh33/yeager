package collector

import (
	"context"
	"log"
	"sync"
	"yeager/pkg/model"
	"yeager/pkg/storage"
)

type SpanProcessor struct {
	storage    storage.SpanWriter
	queue      chan *model.Span
	numWorkers int
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

func NewSpanProcessor(storage storage.SpanWriter, queueSize int, numWorkers int) *SpanProcessor {
	sp := &SpanProcessor{
		storage:    storage,
		queue:      make(chan *model.Span, queueSize),
		numWorkers: numWorkers,
		stopCh:     make(chan struct{}),
	}

	sp.startWorkers()
	return sp
}

func (sp *SpanProcessor) startWorkers() {
	log.Printf("Starting %d collector workers...", sp.numWorkers)
	for i := 0; i < sp.numWorkers; i++ {
		sp.wg.Add(1)
		go sp.workerLoop(i)
	}
}

func (sp *SpanProcessor) workerLoop(id int) {
	defer sp.wg.Done()
	for {
		select {
		case span := <-sp.queue:
			log.Printf("[Worker %d] Pick up span: %s", id, span.SpanID)
			if err := sp.storage.WriteSpan(context.Background(), span); err != nil {
				log.Printf("[Worker %d] Failed to save span: %v", id, err)
			}
		case <-sp.stopCh:
			for {
				select {
				case span := <-sp.queue:
					sp.storage.WriteSpan(context.Background(), span)
				default:
					return
				}
			}
		}
	}
}

func (sp *SpanProcessor) ProcessSpan(span *model.Span) error {
	select {
	case sp.queue <- span:
		return nil
	default:
		log.Println("Collector queue full, dropping span!")
		return nil
	}
}

func (sp *SpanProcessor) Close() {
	close(sp.stopCh)
	sp.wg.Wait()
	log.Println("All collector workers stopped.")
}
