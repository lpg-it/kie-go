// Example: Batch operations using KIE SDK
//
// This example demonstrates batch processing capabilities:
// 1. Creating multiple tasks concurrently
// 2. Waiting for multiple tasks concurrently
// 3. Using BatchProcessor for simplified workflows
// 4. Progress tracking
//
// Usage:
//
//	KIE_API_KEY=your_api_key go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	kie "github.com/lpg-it/kie-go"
)

func main() {
	apiKey := os.Getenv("KIE_API_KEY")
	if apiKey == "" {
		log.Fatal("KIE_API_KEY environment variable is required")
	}

	client := kie.NewClient(
		apiKey,
		kie.WithTimeout(30*time.Second),
		kie.WithMaxRetries(3),
	)
	defer client.Close()

	ctx := context.Background()

	// Example 1: Create multiple tasks concurrently
	fmt.Println("=== Example 1: Batch Task Creation ===")
	batchCreate(ctx, client)

	// Example 2: Wait for multiple tasks concurrently
	fmt.Println("\n=== Example 2: Batch Wait ===")
	batchWait(ctx, client)

	// Example 3: Using BatchProcessor
	fmt.Println("\n=== Example 3: BatchProcessor ===")
	useBatchProcessor(ctx, client)

	// Example 4: Process RequestBuilders with BatchProcessor
	fmt.Println("\n=== Example 4: Process Builders ===")
	processBuilders(ctx, client)
}

func batchCreate(ctx context.Context, client *kie.Client) {
	// Prepare multiple requests
	prompts := []string{
		"A peaceful mountain lake at sunrise",
		"A futuristic cityscape at night",
		"A magical forest with glowing mushrooms",
		"An underwater coral reef scene",
		"A desert oasis under starry sky",
	}

	requests := make([]*kie.CreateTaskRequest, len(prompts))
	for i, prompt := range prompts {
		requests[i] = &kie.CreateTaskRequest{
			Model: "nano-banana-pro",
			Input: kie.Params{
				kie.ParamPrompt:      prompt,
				kie.ParamAspectRatio: kie.Ratio16x9,
				kie.ParamResolution:  kie.Resolution2K,
			},
		}
	}

	// Create all tasks concurrently (with concurrency limit)
	taskIDs := []string{}
	for result := range client.CreateTasksBatch(ctx, requests, kie.WithConcurrency(3)) {
		if result.Error != nil {
			log.Printf("Task %d failed to create: %v", result.Index, result.Error)
			continue
		}
		fmt.Printf("Created task %d: %s\n", result.Index, result.TaskID)
		taskIDs = append(taskIDs, result.TaskID)
	}

	fmt.Printf("\nCreated %d tasks successfully\n", len(taskIDs))
}

func batchWait(ctx context.Context, client *kie.Client) {
	// In a real scenario, these would be actual task IDs from previous creation
	// Here we demonstrate the API usage pattern
	taskIDs := []string{
		"task-id-1",
		"task-id-2",
		"task-id-3",
	}

	fmt.Printf("Waiting for %d tasks...\n", len(taskIDs))

	// Wait for all tasks concurrently
	successCount := 0
	for result := range client.WaitForTasksBatch(ctx, taskIDs, kie.WithConcurrency(3)) {
		if result.Error != nil {
			log.Printf("Task %s failed: %v", result.TaskID, result.Error)
			continue
		}

		urls, err := result.Info.GetResultURLs()
		if err != nil {
			log.Printf("Task %s has no URLs: %v", result.TaskID, err)
			continue
		}

		successCount++
		fmt.Printf("Task %s completed with %d result(s)\n", result.TaskID, len(urls))
	}

	fmt.Printf("\n%d/%d tasks completed successfully\n", successCount, len(taskIDs))
}

func useBatchProcessor(ctx context.Context, client *kie.Client) {
	// Create a BatchProcessor for simplified batch operations
	processor := kie.NewBatchProcessor(client,
		kie.WithConcurrency(5),
		kie.WithBatchTimeout(15*time.Minute),
	)

	// Prepare requests
	requests := []*kie.CreateTaskRequest{
		{
			Model: "nano-banana-pro",
			Input: kie.Params{
				kie.ParamPrompt:      "A cute cat wearing a hat",
				kie.ParamAspectRatio: kie.Ratio1x1,
			},
		},
		{
			Model: "nano-banana-pro",
			Input: kie.Params{
				kie.ParamPrompt:      "A dog playing in the park",
				kie.ParamAspectRatio: kie.Ratio1x1,
			},
		},
		{
			Model: "nano-banana-pro",
			Input: kie.Params{
				kie.ParamPrompt:      "A bird flying over mountains",
				kie.ParamAspectRatio: kie.Ratio16x9,
			},
		},
	}

	// CreateAndWait - creates all tasks and waits for completion
	results, err := processor.CreateAndWait(ctx, requests)
	if err != nil {
		log.Printf("Batch processing failed: %v", err)
		return
	}

	// Process results
	for _, r := range results {
		if r.Error != nil {
			log.Printf("Task %d failed: %v", r.Index, r.Error)
			continue
		}
		fmt.Printf("Task %d completed: %s\n", r.Index, r.Info.TaskID)
		urls, _ := r.Info.GetResultURLs()
		if len(urls) > 0 {
			fmt.Printf("  URL: %s\n", urls[0])
		}
	}
}

func processBuilders(ctx context.Context, client *kie.Client) {
	// Create a BatchProcessor
	processor := kie.NewBatchProcessor(client, kie.WithConcurrency(3))

	// Use RequestBuilders for type-safe parameter setting
	builders := []*kie.RequestBuilder{
		kie.NanoBananaPro.Request().
			Prompt("A sunset over the ocean").
			Set(kie.ParamAspectRatio, kie.Ratio16x9).
			Set(kie.ParamResolution, kie.Resolution4K),
		kie.NanoBananaPro.Request().
			Prompt("A sunrise over mountains").
			Set(kie.ParamAspectRatio, kie.Ratio9x16).
			Set(kie.ParamResolution, kie.Resolution2K),
		kie.GoogleImagen4.Request().
			Prompt("A fantasy castle in the clouds").
			AspectRatio(kie.Ratio1x1), // Use constant
	}

	// ProcessBuilders handles conversion and execution
	results, err := processor.ProcessBuilders(ctx, builders)
	if err != nil {
		log.Printf("Processing failed: %v", err)
		return
	}

	// Display results
	successCount := 0
	for _, r := range results {
		if r.Error == nil {
			successCount++
			fmt.Printf("Builder %d succeeded: %s\n", r.Index, r.Info.TaskID)
		}
	}

	fmt.Printf("\n%d/%d builders processed successfully\n", successCount, len(builders))

	// With progress callback
	fmt.Println("\nProcessing with progress tracking...")
	results2, err := processor.ProcessAll(ctx, []*kie.CreateTaskRequest{
		{Model: "nano-banana-pro", Input: kie.Params{kie.ParamPrompt: "Test 1"}},
		{Model: "nano-banana-pro", Input: kie.Params{kie.ParamPrompt: "Test 2"}},
	}, func(completed, total int) {
		fmt.Printf("Progress: %d/%d (%.0f%%)\n", completed, total, float64(completed)/float64(total)*100)
	})

	if err != nil {
		log.Printf("Processing failed: %v", err)
		return
	}

	fmt.Printf("Completed %d tasks\n", len(results2))
}
