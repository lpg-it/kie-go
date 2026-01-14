package kie

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TaskState represents the state of a task.
type TaskState string

const (
	// TaskStateWaiting means the task is queued and waiting to be processed.
	TaskStateWaiting TaskState = "waiting"

	// TaskStateSuccess means the task completed successfully.
	TaskStateSuccess TaskState = "success"

	// TaskStateFail means the task failed.
	TaskStateFail TaskState = "fail"
)

// IsTerminal returns true if the task is in a terminal state.
func (s TaskState) IsTerminal() bool {
	return s == TaskStateSuccess || s == TaskStateFail
}

// IsSuccess returns true if the task completed successfully.
func (s TaskState) IsSuccess() bool {
	return s == TaskStateSuccess
}

// CreateTaskRequest represents a request to create a new task.
type CreateTaskRequest struct {
	// Model is the model identifier (e.g., "nano-banana-pro")
	Model string `json:"model"`

	// Input is the model-specific input parameters
	Input interface{} `json:"input"`

	// CallbackURL is an optional URL for completion notifications
	CallbackURL string `json:"callBackUrl,omitempty"`
}

// Validate validates the request.
func (r *CreateTaskRequest) Validate() error {
	if r.Model == "" {
		return fmt.Errorf("kie: model is required")
	}
	if r.Input == nil {
		return ErrNilInput
	}

	// If input implements Validator, validate it
	if v, ok := r.Input.(Validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validator is an interface for validatable inputs.
type Validator interface {
	Validate() error
}

// CreateTaskResponse represents the response from creating a task.
type CreateTaskResponse struct {
	// TaskID is the unique identifier for the created task
	TaskID string `json:"taskId"`
}

// TaskInfo represents detailed information about a task.
type TaskInfo struct {
	// TaskID is the unique identifier of the task
	TaskID string `json:"taskId"`

	// Model is the model name used for the task
	Model string `json:"model"`

	// State is the current state of the task
	State TaskState `json:"state"`

	// Param is the JSON-encoded request parameters
	Param string `json:"param"`

	// ResultJSON is the JSON-encoded result (available when state is success)
	// For image/video: {"resultUrls": ["url1", "url2", ...]}
	// For text: {"resultObject": {...}}
	ResultJSON string `json:"resultJson,omitempty"`

	// FailCode is the failure code (available when state is fail)
	FailCode string `json:"failCode,omitempty"`

	// FailMsg is the failure message (available when state is fail)
	FailMsg string `json:"failMsg,omitempty"`

	// CostTime is the processing duration in milliseconds
	CostTime int64 `json:"costTime,omitempty"`

	// CompleteTime is the completion timestamp in milliseconds
	CompleteTime int64 `json:"completeTime,omitempty"`

	// CreateTime is the creation timestamp in milliseconds
	CreateTime int64 `json:"createTime"`
}

// GetResultURLs parses ResultJSON and returns the result URLs.
// This is a convenience method for image/video results.
func (t *TaskInfo) GetResultURLs() ([]string, error) {
	if t.ResultJSON == "" {
		return nil, nil
	}

	var result struct {
		ResultURLs []string `json:"resultUrls"`
	}
	if err := json.Unmarshal([]byte(t.ResultJSON), &result); err != nil {
		return nil, fmt.Errorf("kie: failed to parse result: %w", err)
	}

	return result.ResultURLs, nil
}

// CreateTask creates a new generation task.
//
// The task is created asynchronously. Use GetTaskStatus or WaitForTask
// to check the task status and get results.
//
// Example:
//
//	task, err := client.CreateTask(ctx, &kie.CreateTaskRequest{
//	    Model: image.ModelGoogleNanaBananaPro,
//	    Input: &image.GoogleNanaBananaProInput{
//	        Prompt: "A beautiful sunset",
//	    },
//	})
func (c *Client) CreateTask(ctx context.Context, req *CreateTaskRequest) (*CreateTaskResponse, error) {
	if req == nil {
		return nil, ErrNilInput
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Serialize request using pooled buffer
	buf := c.getBuffer()
	defer c.putBuffer(buf)

	if err := json.NewEncoder(buf).Encode(req); err != nil {
		return nil, fmt.Errorf("kie: failed to encode request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/api/v1/jobs/createTask",
		bytes.NewReader(buf.Bytes()),
	)
	if err != nil {
		return nil, fmt.Errorf("kie: failed to create request: %w", err)
	}

	// GetBody for retry support
	bodyBytes := buf.Bytes()
	bodyCopy := make([]byte, len(bodyBytes))
	copy(bodyCopy, bodyBytes)
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(bodyCopy)), nil
	}

	// Execute request
	respBody, err := c.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	// Parse response
	resp, err := parseResponseData[CreateTaskResponse](respBody)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetTaskStatus retrieves the current status of a task.
//
// This method is useful for checking task progress without blocking.
// For blocking wait, use WaitForTask instead.
func (c *Client) GetTaskStatus(ctx context.Context, taskID string) (*TaskInfo, error) {
	if taskID == "" {
		return nil, ErrEmptyTaskID
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+"/api/v1/jobs/recordInfo?taskId="+taskID,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("kie: failed to create request: %w", err)
	}

	// Execute request
	respBody, err := c.doRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	// Parse response
	info, err := parseResponseData[TaskInfo](respBody)
	if err != nil {
		return nil, err
	}

	return info, nil
}
