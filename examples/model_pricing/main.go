// Package main demonstrates the Model Pricing API features of the kie-go SDK.
//
// This example shows how to:
// - List all model pricing with pagination
// - Search for specific model pricing by keyword
// - Filter by interface type
// - Retrieve all pricing records with rate limiting
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

	// Example 1: List first page of model pricing
	fmt.Println("=== List Model Pricing (Page 1) ===")
	listFirstPage(ctx, client)

	// Example 2: Search for specific models
	fmt.Println("\n=== Search Model Pricing ===")
	searchModels(ctx, client, "claude")

	// Example 3: Filter by interface type
	fmt.Println("\n=== Filter by Interface Type (chat) ===")
	filterByType(ctx, client, "chat")

	// Example 4: Get pricing summary
	fmt.Println("\n=== Pricing Summary ===")
	getPricingSummary(ctx, client)

	// Example 5: List all models with rate limiting (commented out - takes time)
	// fmt.Println("\n=== List All Models (with rate limiting) ===")
	// listAllModels(ctx, client)

	fmt.Println("\nDone!")
}

// listFirstPage demonstrates how to list the first page of model pricing.
func listFirstPage(ctx context.Context, client *kie.Client) {
	page, err := client.GetModelPricing(ctx, &kie.ModelPricingRequest{
		PageNum:  1,
		PageSize: 10,
	})
	if err != nil {
		log.Printf("Failed to get model pricing: %v", err)
		return
	}

	fmt.Printf("Total models: %d, Pages: %d, Current page: %d\n", page.Total, page.Pages, page.Current)
	fmt.Println("First 10 models:")
	for i, pricing := range page.Records {
		fmt.Printf("  %d. [%s] %s - %s %s (USD: $%s)\n",
			i+1,
			pricing.InterfaceType,
			pricing.ModelDescription,
			pricing.CreditPrice,
			pricing.CreditUnit,
			pricing.USDPrice,
		)
	}
}

// searchModels demonstrates how to search for specific model pricing.
func searchModels(ctx context.Context, client *kie.Client, keyword string) {
	page, err := client.SearchModelPricing(ctx, keyword, 1, 25)
	if err != nil {
		log.Printf("Failed to search model pricing: %v", err)
		return
	}

	fmt.Printf("Found %d models matching '%s':\n", page.Total, keyword)
	for _, pricing := range page.Records {
		fmt.Printf("  - [%s] %s: %s %s\n",
			pricing.Provider,
			pricing.ModelDescription,
			pricing.CreditPrice,
			pricing.CreditUnit,
		)
	}
}

// filterByType demonstrates how to filter model pricing by interface type.
func filterByType(ctx context.Context, client *kie.Client, interfaceType string) {
	page, err := client.GetModelPricing(ctx, &kie.ModelPricingRequest{
		PageNum:       1,
		PageSize:      10,
		InterfaceType: interfaceType,
	})
	if err != nil {
		log.Printf("Failed to get model pricing: %v", err)
		return
	}

	fmt.Printf("Found %d '%s' models:\n", page.Total, interfaceType)
	for _, pricing := range page.Records {
		discountInfo := ""
		if pricing.DiscountRate > 0 {
			discountInfo = fmt.Sprintf(" (%.0f%% off)", pricing.DiscountRate)
		}
		fmt.Printf("  - %s: %s %s%s\n",
			pricing.ModelDescription,
			pricing.CreditPrice,
			pricing.CreditUnit,
			discountInfo,
		)
	}
}

// getPricingSummary demonstrates how to get a summary of all pricing.
func getPricingSummary(ctx context.Context, client *kie.Client) {
	// Get first page to see total count
	page, err := client.GetModelPricing(ctx, &kie.ModelPricingRequest{
		PageNum:  1,
		PageSize: 1,
	})
	if err != nil {
		log.Printf("Failed to get model pricing: %v", err)
		return
	}

	fmt.Printf("Total pricing entries: %d\n", page.Total)

	// Count by interface type
	types := []string{"chat", "image", "video", "music"}
	for _, t := range types {
		typePage, err := client.GetModelPricing(ctx, &kie.ModelPricingRequest{
			PageNum:       1,
			PageSize:      1,
			InterfaceType: t,
		})
		if err != nil {
			continue
		}
		fmt.Printf("  - %s models: %d\n", t, typePage.Total)
	}
}

// listAllModels demonstrates how to fetch all model pricing with rate limiting.
// This uses random delays between requests to avoid triggering rate limits.
func listAllModels(ctx context.Context, client *kie.Client) {
	fmt.Println("Fetching all models with 1-3 second delays between pages...")

	allPricing, err := client.ListAllModelPricing(ctx,
		// Random delay between 1-3 seconds (default)
		kie.WithPricingInterval(1*time.Second, 3*time.Second),
		// Progress callback
		kie.WithPricingProgress(func(fetched, total int) {
			fmt.Printf("  Progress: %d/%d (%.1f%%)\n", fetched, total, float64(fetched)/float64(total)*100)
		}),
		// Page size (default is 100)
		kie.WithPricingPageSize(50),
	)
	if err != nil {
		log.Printf("Failed to list all models: %v", err)
		return
	}

	fmt.Printf("Successfully fetched %d model pricing entries\n", len(allPricing))

	// Show some statistics
	providers := make(map[string]int)
	for _, p := range allPricing {
		providers[p.Provider]++
	}
	fmt.Println("Models by provider:")
	for provider, count := range providers {
		fmt.Printf("  - %s: %d\n", provider, count)
	}
}
