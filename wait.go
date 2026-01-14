package kie

import (
	"context"
	"sync"
	"time"
)

// WaitConfig configures the WaitForTask behavior.
type WaitConfig struct {
	// Timeout is the maximum time to wait for the task to complete.
	// If zero, defaults to 10 minutes.
	Timeout time.Duration

	// PollInterval is the initial interval between status checks.
	// If zero, defaults to 500ms.
	PollInterval time.Duration

	// MaxPollInterval is the maximum interval between status checks.
	// The interval increases exponentially up to this maximum.
	// If zero, defaults to 5 seconds.
	MaxPollInterval time.Duration

	// PollMultiplier is the factor by which the poll interval increases.
	// If zero, defaults to 1.5.
	PollMultiplier float64

	// ProgressCallback is called during polling to report progress.
	// Parameters: taskID, elapsed time since start, poll count.
	ProgressCallback ProgressCallback
}

// ProgressCallback is called during task polling to report progress.
// Use this to display progress to users during long-running tasks.
//
// Parameters:
//   - taskID: the task being polled
//   - elapsed: time since polling started
//   - pollCount: number of polls completed (starts at 1)
//
// Example:
//
//	info, err := client.WaitForTask(ctx, taskID,
//	    kie.WithProgressCallback(func(taskID string, elapsed time.Duration, pollCount int) {
//	        log.Printf("Task %s: waiting %.1fs (poll #%d)", taskID, elapsed.Seconds(), pollCount)
//	    }),
//	)
type ProgressCallback func(taskID string, elapsed time.Duration, pollCount int)

// DefaultWaitConfig returns the default wait configuration.
func DefaultWaitConfig() *WaitConfig {
	return &WaitConfig{
		Timeout:         10 * time.Minute,
		PollInterval:    500 * time.Millisecond,
		MaxPollInterval: 5 * time.Second,
		PollMultiplier:  1.5,
	}
}

// WaitOption configures WaitForTask behavior.
type WaitOption func(*WaitConfig)

// WithWaitTimeout sets the maximum wait timeout.
func WithWaitTimeout(d time.Duration) WaitOption {
	return func(c *WaitConfig) {
		c.Timeout = d
	}
}

// WithPollInterval sets the initial polling interval.
func WithPollInterval(d time.Duration) WaitOption {
	return func(c *WaitConfig) {
		c.PollInterval = d
	}
}

// WithMaxPollInterval sets the maximum polling interval.
func WithMaxPollInterval(d time.Duration) WaitOption {
	return func(c *WaitConfig) {
		c.MaxPollInterval = d
	}
}

// WithProgressCallback sets a callback for progress notifications.
// The callback is invoked after each poll with:
//   - taskID: the task ID being waited on
//   - elapsed: time since the wait started
//   - pollCount: number of polls completed
func WithProgressCallback(cb ProgressCallback) WaitOption {
	return func(c *WaitConfig) {
		c.ProgressCallback = cb
	}
}

// WaitForTask polls until the task completes or the context is canceled.
// It uses exponential backoff for polling.
//
// Example:
//
//	info, err := client.WaitForTask(ctx, taskID,
//	    kie.WithWaitTimeout(5*time.Minute),
//	    kie.WithPollInterval(time.Second),
//	)
func (c *Client) WaitForTask(ctx context.Context, taskID string, opts ...WaitOption) (*TaskInfo, error) {
	if taskID == "" {
		return nil, ErrEmptyTaskID
	}

	cfg := DefaultWaitConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	// Apply timeout if not zero
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}

	interval := cfg.PollInterval
	startTime := time.Now()
	pollCount := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		pollCount++
		info, err := c.GetTaskStatus(ctx, taskID)
		if err != nil {
			return nil, err
		}

		// Call progress callback if configured
		if cfg.ProgressCallback != nil {
			cfg.ProgressCallback(taskID, time.Since(startTime), pollCount)
		}

		if info.State.IsTerminal() {
			if info.State.IsSuccess() {
				return info, nil
			}
			return info, &TaskFailedError{
				TaskID:   info.TaskID,
				FailCode: info.FailCode,
				FailMsg:  info.FailMsg,
			}
		}

		// Wait before next poll
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(interval):
		}

		// Increase interval with exponential backoff
		interval = time.Duration(float64(interval) * cfg.PollMultiplier)
		if interval > cfg.MaxPollInterval {
			interval = cfg.MaxPollInterval
		}
	}
}

// BatchResult represents the result of a single task in a batch operation.
type BatchResult struct {
	Index  int
	TaskID string
	Info   *TaskInfo
	Error  error
}

// BatchOption configures batch operations.
type BatchOption func(*BatchConfig)

// BatchConfig configures batch operations.
type BatchConfig struct {
	Concurrency     int
	Timeout         time.Duration
	PollInterval    time.Duration
	MaxPollInterval time.Duration
}

// DefaultBatchConfig returns the default batch configuration.
func DefaultBatchConfig() *BatchConfig {
	return &BatchConfig{
		Concurrency:     5,
		Timeout:         10 * time.Minute,
		PollInterval:    500 * time.Millisecond,
		MaxPollInterval: 5 * time.Second,
	}
}

// WithConcurrency sets the concurrency level for batch operations.
func WithConcurrency(n int) BatchOption {
	return func(c *BatchConfig) {
		if n > 0 {
			c.Concurrency = n
		}
	}
}

// WithBatchTimeout sets the timeout for batch operations.
func WithBatchTimeout(d time.Duration) BatchOption {
	return func(c *BatchConfig) {
		c.Timeout = d
	}
}

// WaitForTasksBatch waits for multiple tasks concurrently.
// Results are returned through the channel as they complete.
func (c *Client) WaitForTasksBatch(ctx context.Context, taskIDs []string, opts ...BatchOption) <-chan BatchResult {
	cfg := DefaultBatchConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	results := make(chan BatchResult, len(taskIDs))

	go func() {
		defer close(results)

		sem := make(chan struct{}, cfg.Concurrency)

		for i, taskID := range taskIDs {
			select {
			case <-ctx.Done():
				return
			case sem <- struct{}{}:
			}

			go func(idx int, id string) {
				defer func() { <-sem }()

				info, err := c.WaitForTask(ctx, id,
					WithWaitTimeout(cfg.Timeout),
					WithPollInterval(cfg.PollInterval),
					WithMaxPollInterval(cfg.MaxPollInterval),
				)

				select {
				case results <- BatchResult{Index: idx, TaskID: id, Info: info, Error: err}:
				case <-ctx.Done():
				}
			}(i, taskID)
		}

		// Wait for all goroutines
		for i := 0; i < cfg.Concurrency && i < len(taskIDs); i++ {
			sem <- struct{}{}
		}
	}()

	return results
}

// CreateTasksBatch creates multiple tasks concurrently.
// Results are returned through the channel as they complete.
func (c *Client) CreateTasksBatch(ctx context.Context, requests []*CreateTaskRequest, opts ...BatchOption) <-chan BatchResult {
	cfg := DefaultBatchConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	results := make(chan BatchResult, len(requests))

	go func() {
		defer close(results)

		sem := make(chan struct{}, cfg.Concurrency)
		var wg sync.WaitGroup

		for i, req := range requests {
			select {
			case <-ctx.Done():
				// Wait for in-flight operations before returning
				wg.Wait()
				return
			case sem <- struct{}{}:
			}

			wg.Add(1)
			go func(idx int, r *CreateTaskRequest) {
				defer func() {
					<-sem
					wg.Done()
				}()

				resp, err := c.CreateTask(ctx, r)
				var taskID string
				if resp != nil {
					taskID = resp.TaskID
				}

				select {
				case results <- BatchResult{Index: idx, TaskID: taskID, Error: err}:
				case <-ctx.Done():
					// Context canceled, don't send result
				}
			}(i, req)
		}

		// Wait for all goroutines to finish
		wg.Wait()
	}()

	return results
}
