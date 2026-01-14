// Package main demonstrates the Common API features of the kie-go SDK.
//
// This example shows how to:
// - Check account credit balance
// - Get temporary download URLs for generated files
// - Use the enhanced Downloader with automatic URL conversion
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

	client := kie.NewClient(apiKey)
	ctx := context.Background()

	// Example 1: Check credit balance
	fmt.Println("=== Credit Balance Check ===")
	checkCredits(ctx, client)

	// Example 2: Credit threshold check
	fmt.Println("\n=== Credit Threshold Check ===")
	checkCreditThreshold(ctx, client)

	// Example 3: Get download URL (commented out as it requires a real file URL)
	// fmt.Println("\n=== Download URL ===")
	// getDownloadURL(ctx, client)

	// Example 4: Use Downloader with KIE Client
	fmt.Println("\n=== Downloader with KIE Client ===")
	demonstrateDownloader(client)

	fmt.Println("\nDone!")
}

// checkCredits demonstrates how to check the current credit balance.
func checkCredits(ctx context.Context, client *kie.Client) {
	credits, err := client.GetCredits(ctx)
	if err != nil {
		log.Printf("Failed to get credits: %v", err)
		return
	}
	fmt.Printf("Current credit balance: %d\n", credits)
}

// checkCreditThreshold demonstrates how to check if credits are above a threshold.
func checkCreditThreshold(ctx context.Context, client *kie.Client) {
	threshold := 50

	ok, balance, err := client.CheckCredits(ctx, threshold)
	if err != nil {
		log.Printf("Failed to check credits: %v", err)
		return
	}

	if ok {
		fmt.Printf("✓ Sufficient credits: %d (threshold: %d)\n", balance, threshold)
	} else {
		fmt.Printf("⚠ Low credits: %d (need at least %d)\n", balance, threshold)
	}

	// Alternative: simple boolean check
	hasSufficient, _ := client.HasSufficientCredits(ctx, 100)
	fmt.Printf("Has 100+ credits: %v\n", hasSufficient)
}

// getDownloadURL demonstrates how to get a temporary download URL.
// Note: This requires a valid KIE-generated file URL.
func getDownloadURL(ctx context.Context, client *kie.Client) {
	// This is a placeholder URL - replace with an actual KIE-generated URL
	fileURL := "https://tempfile.1f6cxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxbd98"

	downloadURL, err := client.GetDownloadURL(ctx, fileURL)
	if err != nil {
		log.Printf("Failed to get download URL: %v", err)
		return
	}

	fmt.Printf("Original URL: %s\n", fileURL)
	fmt.Printf("Download URL: %s\n", downloadURL)
	fmt.Printf("Note: This URL is valid for 20 minutes\n")
}

// demonstrateDownloader shows how to create a Downloader with KIE Client integration.
func demonstrateDownloader(client *kie.Client) {
	// Create a Downloader with KIE Client for automatic URL conversion
	downloader := kie.NewDownloader(
		kie.WithKIEClient(client),
		kie.WithDownloadConcurrency(3),
	)

	fmt.Printf("Downloader created with KIE Client integration\n")
	fmt.Printf("- Automatic tempfile URL conversion: enabled\n")
	fmt.Printf("- Concurrent downloads: 3\n")

	// In a real scenario, you would use it like this:
	/*
		// Get task info from a completed generation
		taskInfo, err := client.GetTaskStatus(ctx, "your-task-id")
		if err != nil {
			log.Fatal(err)
		}

		// Download all results - tempfile URLs are automatically converted
		results, err := downloader.DownloadFromTaskInfo(ctx, taskInfo, "./output/")
		if err != nil {
			log.Fatal(err)
		}

		for _, r := range results {
			if r.Error != nil {
				log.Printf("Failed: %v", r.Error)
			} else {
				log.Printf("Downloaded: %s (%d bytes)", r.LocalPath, r.Size)
			}
		}
	*/

	// Demonstrate URL detection
	testURLs := []string{
		"https://tempfile.abc123xyz",
		"https://cdn.kie.ai/files/image.png",
	}

	for _, url := range testURLs {
		isTempFile := kie.IsKIETempFileURL(url)
		fmt.Printf("URL: %s -> Is tempfile: %v\n", url, isTempFile)
	}

	_ = downloader // suppress unused warning
}

// monitorCredits demonstrates a credit monitoring pattern.
// This is useful for long-running applications that need to track credit usage.
func monitorCredits(ctx context.Context, client *kie.Client) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	lowCreditThreshold := 50

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ok, balance, err := client.CheckCredits(ctx, lowCreditThreshold)
			if err != nil {
				log.Printf("Credit check failed: %v", err)
				continue
			}

			if !ok {
				log.Printf("⚠ LOW CREDITS ALERT: %d credits remaining", balance)
				// In a real application, you might send an email, Slack message, etc.
			} else {
				log.Printf("Credit check OK: %d credits", balance)
			}
		}
	}
}
