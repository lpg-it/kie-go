package kie

import (
	"context"
	"sync"

	"github.com/lpg-it/kie-go/model"
)

// BatchProcessor provides high-level batch task processing.
// It wraps the lower-level batch methods with a simpler API.
//
// Example:
//
//	processor := kie.NewBatchProcessor(client,
//	    kie.WithConcurrency(10),
//	    kie.WithBatchTimeout(15*time.Minute),
//	)
//
//	results, err := processor.CreateAndWait(ctx, requests)
type BatchProcessor struct {
	client *Client
	config *BatchConfig
}

// NewBatchProcessor creates a new BatchProcessor.
func NewBatchProcessor(client *Client, opts ...BatchOption) *BatchProcessor {
	cfg := DefaultBatchConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return &BatchProcessor{
		client: client,
		config: cfg,
	}
}

// CreateAndWait creates multiple tasks and waits for all to complete.
// Returns all results in order, including any errors.
//
// This is a convenience method that combines CreateTask and WaitForTask
// for batch processing.
//
// Example:
//
//	results, err := processor.CreateAndWait(ctx, []*kie.CreateTaskRequest{
//	    {Model: "nano-banana-pro", Input: map[string]any{"prompt": "A sunset"}},
//	    {Model: "nano-banana-pro", Input: map[string]any{"prompt": "A sunrise"}},
//	})
//	for _, r := range results {
//	    if r.Error != nil {
//	        log.Printf("Task %d failed: %v", r.Index, r.Error)
//	    } else {
//	        log.Printf("Task %d completed: %s", r.Index, r.Info.TaskID)
//	    }
//	}
func (p *BatchProcessor) CreateAndWait(ctx context.Context, requests []*CreateTaskRequest) ([]BatchResult, error) {
	if len(requests) == 0 {
		return nil, nil
	}

	// Create all tasks first
	taskIDs := make([]string, len(requests))
	createErrors := make([]error, len(requests))

	// Use channel-based approach for creation
	createResults := p.client.CreateTasksBatch(ctx, requests,
		WithConcurrency(p.config.Concurrency),
	)

	for result := range createResults {
		if result.Error != nil {
			createErrors[result.Index] = result.Error
		} else {
			taskIDs[result.Index] = result.TaskID
		}
	}

	// Wait for all successfully created tasks
	var validTaskIDs []string
	taskIndexMap := make(map[string]int) // map taskID -> original index
	for i, id := range taskIDs {
		if id != "" {
			validTaskIDs = append(validTaskIDs, id)
			taskIndexMap[id] = i
		}
	}

	// Prepare final results
	results := make([]BatchResult, len(requests))

	// Copy create errors
	for i, err := range createErrors {
		if err != nil {
			results[i] = BatchResult{Index: i, Error: err}
		}
	}

	// Wait for valid tasks
	if len(validTaskIDs) > 0 {
		waitResults := p.client.WaitForTasksBatch(ctx, validTaskIDs,
			WithConcurrency(p.config.Concurrency),
			WithBatchTimeout(p.config.Timeout),
		)

		for result := range waitResults {
			originalIndex := taskIndexMap[result.TaskID]
			results[originalIndex] = BatchResult{
				Index:  originalIndex,
				TaskID: result.TaskID,
				Info:   result.Info,
				Error:  result.Error,
			}
		}
	}

	return results, nil
}

// ProcessBuilders creates and waits for tasks from RequestBuilders.
// This is a convenience method for processing multiple model requests.
//
// Example:
//
//	builders := []*model.RequestBuilder{
//	    kie.NanoBananaPro.Request().Prompt("A sunset"),
//	    kie.NanoBananaPro.Request().Prompt("A sunrise"),
//	}
//	results, err := processor.ProcessBuilders(ctx, builders)
func (p *BatchProcessor) ProcessBuilders(ctx context.Context, builders []*model.RequestBuilder) ([]BatchResult, error) {
	if len(builders) == 0 {
		return nil, nil
	}

	// Convert builders to CreateTaskRequests
	requests := make([]*CreateTaskRequest, len(builders))
	for i, b := range builders {
		// Validate parameters first
		if err := b.Validate(); err != nil {
			// Return early with validation error
			results := make([]BatchResult, len(builders))
			results[i] = BatchResult{Index: i, Error: err}
			return results, nil
		}

		requests[i] = &CreateTaskRequest{
			Model: b.Model().Identifier,
			Input: b.Params(),
		}
	}

	return p.CreateAndWait(ctx, requests)
}

// ProcessAll creates tasks, waits for completion, and collects all results.
// Unlike channel-based methods, this blocks until all tasks complete.
//
// Parameters:
//   - requests: slice of task requests to process
//   - progressCb: optional callback for progress (can be nil)
//
// Example:
//
//	results, err := processor.ProcessAll(ctx, requests,
//	    func(completed, total int) {
//	        log.Printf("Progress: %d/%d", completed, total)
//	    },
//	)
func (p *BatchProcessor) ProcessAll(ctx context.Context, requests []*CreateTaskRequest, progressCb func(completed, total int)) ([]BatchResult, error) {
	if len(requests) == 0 {
		return nil, nil
	}

	results := make([]BatchResult, len(requests))
	var completed int
	var mu sync.Mutex

	// Process in batches with concurrency control
	sem := make(chan struct{}, p.config.Concurrency)
	var wg sync.WaitGroup

	for i, req := range requests {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		case sem <- struct{}{}:
		}

		wg.Add(1)
		go func(idx int, r *CreateTaskRequest) {
			defer func() {
				<-sem
				wg.Done()
			}()

			// Create task
			resp, err := p.client.CreateTask(ctx, r)
			if err != nil {
				mu.Lock()
				results[idx] = BatchResult{Index: idx, Error: err}
				completed++
				if progressCb != nil {
					progressCb(completed, len(requests))
				}
				mu.Unlock()
				return
			}

			// Wait for completion
			info, err := p.client.WaitForTask(ctx, resp.TaskID,
				WithWaitTimeout(p.config.Timeout),
				WithPollInterval(p.config.PollInterval),
				WithMaxPollInterval(p.config.MaxPollInterval),
			)

			mu.Lock()
			results[idx] = BatchResult{
				Index:  idx,
				TaskID: resp.TaskID,
				Info:   info,
				Error:  err,
			}
			completed++
			if progressCb != nil {
				progressCb(completed, len(requests))
			}
			mu.Unlock()
		}(i, req)
	}

	wg.Wait()
	return results, nil
}

// Config returns the current batch configuration.
func (p *BatchProcessor) Config() *BatchConfig {
	return p.config
}
