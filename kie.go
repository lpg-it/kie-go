// Package kie provides a high-performance, secure, and stable Go SDK for the KIE AI platform.
//
// Key Features:
//   - Security: Secure API key handling, no credentials in logs
//   - Stability: Robust error handling, retry with exponential backoff
//   - High Performance: Connection pooling, zero-allocation JSON, buffer reuse
//
// Quick Start:
//
//	client := kie.NewClient("your-api-key")
//
//	task, err := client.CreateTask(ctx, &kie.CreateTaskRequest{
//	    Model: image.ModelGoogleNanaBananaPro,
//	    Input: &image.GoogleNanaBananaProInput{
//	        Prompt: "A beautiful sunset",
//	    },
//	})
//
//	result, err := client.WaitForTask(ctx, task.TaskID)
package kie

// Version is the current SDK version
const Version = "0.1.0"
