package kie

import (
	"context"
	"net/http"
)

// GetCredits retrieves the current credit balance of the account.
//
// This method is useful for:
//   - Checking credits before starting generation tasks
//   - Monitoring credit consumption patterns
//   - Implementing credit threshold alerts
//
// Example:
//
//	credits, err := client.GetCredits(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Current credits: %d\n", credits)
func (c *Client) GetCredits(ctx context.Context) (int, error) {
	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+"/api/v1/chat/credit",
		nil,
	)
	if err != nil {
		return 0, wrapError("failed to create request", err)
	}

	// Execute request
	respBody, err := c.doRequest(ctx, httpReq)
	if err != nil {
		return 0, err
	}

	// Parse response - data is directly the credit value (int)
	credits, err := parseResponseData[int](respBody)
	if err != nil {
		return 0, err
	}

	return *credits, nil
}

// CheckCredits checks if the account has sufficient credits.
//
// This is a convenience method that combines GetCredits with a threshold check.
// It returns:
//   - ok: true if balance >= threshold
//   - balance: the current credit balance
//   - err: any error that occurred
//
// Example:
//
//	ok, balance, err := client.CheckCredits(ctx, 100)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if !ok {
//	    log.Printf("Warning: Low credits (%d), need at least 100", balance)
//	    return
//	}
//	// Proceed with operations
func (c *Client) CheckCredits(ctx context.Context, threshold int) (ok bool, balance int, err error) {
	balance, err = c.GetCredits(ctx)
	if err != nil {
		return false, 0, err
	}

	return balance >= threshold, balance, nil
}

// HasSufficientCredits is a simple check if credits are above threshold.
// Unlike CheckCredits, it only returns a boolean result.
//
// Example:
//
//	if ok, _ := client.HasSufficientCredits(ctx, 50); !ok {
//	    log.Fatal("Insufficient credits")
//	}
func (c *Client) HasSufficientCredits(ctx context.Context, threshold int) (bool, error) {
	ok, _, err := c.CheckCredits(ctx, threshold)
	return ok, err
}
