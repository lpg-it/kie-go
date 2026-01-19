package kie

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// ModelPricingRequest represents the request parameters for querying model pricing.
type ModelPricingRequest struct {
	// PageNum is the page number (1-based).
	PageNum int `json:"pageNum"`

	// PageSize is the number of records per page.
	PageSize int `json:"pageSize"`

	// ModelDescription is the search keyword for model description.
	// Leave empty to list all models.
	ModelDescription string `json:"modelDescription"`

	// InterfaceType filters by interface type (e.g., "chat", "image", "video", "music").
	// Leave empty to list all types.
	InterfaceType string `json:"interfaceType"`
}

// ModelPricing represents pricing information for a single model.
type ModelPricing struct {
	// ModelDescription describes the model, including its name and parameters.
	// Note: Different parameters (e.g., Input/Output tokens) are in separate records.
	ModelDescription string `json:"modelDescription"`

	// InterfaceType is the type of interface (e.g., "chat", "image", "video", "music").
	InterfaceType string `json:"interfaceType"`

	// Provider is the model provider (e.g., "Google", "Other", "ByteDance").
	Provider string `json:"provider"`

	// CreditPrice is the price in credits.
	CreditPrice string `json:"creditPrice"`

	// CreditUnit describes the unit for credit pricing (e.g., "per million tokens", "per image").
	CreditUnit string `json:"creditUnit"`

	// USDPrice is the price in USD.
	USDPrice string `json:"usdPrice"`

	// FalPrice is the fal.ai price for comparison (may be empty).
	FalPrice string `json:"falPrice"`

	// DiscountRate is the discount percentage.
	DiscountRate float64 `json:"discountRate"`

	// Anchor is the URL to the model page on kie.ai.
	Anchor string `json:"anchor"`

	// DiscountPrice indicates whether the price is discounted.
	DiscountPrice bool `json:"discountPrice"`
}

// ModelPricingPage represents a paginated response of model pricing.
type ModelPricingPage struct {
	// Records contains the list of model pricing entries.
	Records []ModelPricing `json:"records"`

	// Total is the total number of records.
	Total int `json:"total"`

	// Size is the page size.
	Size int `json:"size"`

	// Current is the current page number.
	Current int `json:"current"`

	// Pages is the total number of pages.
	Pages int `json:"pages"`

	// SearchCount indicates whether the search count is enabled.
	SearchCount bool `json:"searchCount"`

	// Orders contains the ordering information (usually empty).
	Orders []any `json:"orders"`
}

// GetModelPricing retrieves a paginated list of model pricing.
//
// This method allows you to:
//   - List all model pricing with pagination
//   - Search for specific models by description
//   - Filter by interface type (chat, image, video, music)
//
// Example - List all models:
//
//	page, err := client.GetModelPricing(ctx, &kie.ModelPricingRequest{
//	    PageNum:  1,
//	    PageSize: 25,
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, pricing := range page.Records {
//	    fmt.Printf("%s: %s %s\n", pricing.ModelDescription, pricing.CreditPrice, pricing.CreditUnit)
//	}
//
// Example - Search for a specific model:
//
//	page, err := client.GetModelPricing(ctx, &kie.ModelPricingRequest{
//	    PageNum:          1,
//	    PageSize:         25,
//	    ModelDescription: "claude",
//	})
//
// Example - Filter by interface type:
//
//	page, err := client.GetModelPricing(ctx, &kie.ModelPricingRequest{
//	    PageNum:       1,
//	    PageSize:      25,
//	    InterfaceType: "chat",
//	})
func (c *Client) GetModelPricing(ctx context.Context, req *ModelPricingRequest) (*ModelPricingPage, error) {
	// Set defaults if not provided
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 25
	}

	// Marshal request body
	body, err := json.Marshal(req)
	if err != nil {
		return nil, wrapError("failed to marshal request", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/client/v1/model-pricing/page",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, wrapError("failed to create request", err)
	}

	// Set GetBody for retry support
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}

	// Execute request
	respBody, err := c.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	// Parse response
	page, err := parseResponseData[ModelPricingPage](respBody)
	if err != nil {
		return nil, err
	}

	return page, nil
}

// SearchModelPricing is a convenience method to search for model pricing by description.
//
// This is equivalent to calling GetModelPricing with ModelDescription set.
//
// Example:
//
//	page, err := client.SearchModelPricing(ctx, "gemini", 1, 25)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Found %d models matching 'gemini'\n", page.Total)
func (c *Client) SearchModelPricing(ctx context.Context, keyword string, pageNum, pageSize int) (*ModelPricingPage, error) {
	return c.GetModelPricing(ctx, &ModelPricingRequest{
		PageNum:          pageNum,
		PageSize:         pageSize,
		ModelDescription: keyword,
	})
}

// ListAllModelPricing retrieves all model pricing by iterating through all pages.
//
// This method automatically handles pagination and returns all records.
// By default, it adds a random delay (1-3 seconds) between requests to avoid rate limiting.
//
// Options:
//   - WithPricingInterval(min, max): Set custom delay range between requests
//   - WithPricingProgress(callback): Receive progress updates
//   - WithPricingPageSize(size): Set the page size for each request
//
// Example - Basic usage:
//
//	allPricing, err := client.ListAllModelPricing(ctx)
//
// Example - With custom interval:
//
//	allPricing, err := client.ListAllModelPricing(ctx,
//	    kie.WithPricingInterval(2*time.Second, 5*time.Second),
//	)
//
// Example - With progress callback:
//
//	allPricing, err := client.ListAllModelPricing(ctx,
//	    kie.WithPricingProgress(func(fetched, total int) {
//	        fmt.Printf("Progress: %d/%d\n", fetched, total)
//	    }),
//	)
//
// Example - No delay (use with caution):
//
//	allPricing, err := client.ListAllModelPricing(ctx,
//	    kie.WithPricingInterval(0, 0),
//	)
func (c *Client) ListAllModelPricing(ctx context.Context, opts ...ListAllPricingOption) ([]ModelPricing, error) {
	// Apply default options
	options := &listAllPricingOptions{
		minInterval: 1 * time.Second,
		maxInterval: 3 * time.Second,
		pageSize:    100,
	}
	for _, opt := range opts {
		opt(options)
	}

	var allRecords []ModelPricing
	pageNum := 1

	for {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return allRecords, ctx.Err()
		default:
		}

		page, err := c.GetModelPricing(ctx, &ModelPricingRequest{
			PageNum:  pageNum,
			PageSize: options.pageSize,
		})
		if err != nil {
			return allRecords, err
		}

		allRecords = append(allRecords, page.Records...)

		// Call progress callback if set
		if options.progressFn != nil {
			options.progressFn(len(allRecords), page.Total)
		}

		// Check if we've retrieved all records
		if pageNum >= page.Pages || len(page.Records) == 0 {
			break
		}

		pageNum++

		// Add random delay between requests (skip for last page)
		if options.maxInterval > 0 {
			delay := randomDuration(options.minInterval, options.maxInterval)
			select {
			case <-ctx.Done():
				return allRecords, ctx.Err()
			case <-time.After(delay):
			}
		}
	}

	return allRecords, nil
}

// listAllPricingOptions holds options for ListAllModelPricing.
type listAllPricingOptions struct {
	minInterval time.Duration
	maxInterval time.Duration
	pageSize    int
	progressFn  func(fetched, total int)
}

// ListAllPricingOption is a functional option for ListAllModelPricing.
type ListAllPricingOption func(*listAllPricingOptions)

// WithPricingInterval sets the random delay range between paginated requests.
// This helps avoid rate limiting when fetching all pricing data.
//
// Default: 1-3 seconds
// Set both to 0 to disable delay (use with caution).
//
// Example:
//
//	client.ListAllModelPricing(ctx, kie.WithPricingInterval(2*time.Second, 5*time.Second))
func WithPricingInterval(min, max time.Duration) ListAllPricingOption {
	return func(o *listAllPricingOptions) {
		o.minInterval = min
		o.maxInterval = max
	}
}

// WithPricingProgress sets a callback function to receive progress updates.
// The callback is called after each page is fetched with the current count
// of fetched records and the total count.
//
// Example:
//
//	client.ListAllModelPricing(ctx, kie.WithPricingProgress(func(fetched, total int) {
//	    fmt.Printf("Fetched %d of %d records\n", fetched, total)
//	}))
func WithPricingProgress(fn func(fetched, total int)) ListAllPricingOption {
	return func(o *listAllPricingOptions) {
		o.progressFn = fn
	}
}

// WithPricingPageSize sets the page size for each request.
// Larger page sizes mean fewer requests but more data per request.
//
// Default: 100
//
// Example:
//
//	client.ListAllModelPricing(ctx, kie.WithPricingPageSize(50))
func WithPricingPageSize(size int) ListAllPricingOption {
	return func(o *listAllPricingOptions) {
		if size > 0 {
			o.pageSize = size
		}
	}
}

// randomDuration returns a random duration between min and max.
func randomDuration(min, max time.Duration) time.Duration {
	if max <= min {
		return min
	}
	return min + time.Duration(rand.Int63n(int64(max-min)))
}
