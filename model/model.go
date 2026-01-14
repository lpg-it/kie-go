package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Category represents the model category.
type Category string

const (
	CategoryTextToImage  Category = "text-to-image"
	CategoryImageToImage Category = "image-to-image"
	CategoryImageEdit    Category = "image-edit"
	CategoryUpscale      Category = "upscale"
	CategoryTextToVideo  Category = "text-to-video"
	CategoryImageToVideo Category = "image-to-video"
	CategoryVideoToVideo Category = "video-to-video"
)

// Model represents an AI model definition.
type Model struct {
	// Identifier is the model ID used in API calls.
	Identifier string

	// Name is the display name.
	Name string

	// Category is the model category.
	Category Category

	// Provider is the model provider.
	Provider string

	// Timeout is the default timeout for this model.
	Timeout time.Duration

	// Fields defines the input parameters.
	requiredFields []Field
	optionalFields []Field
}

// ModelOption configures a Model.
type ModelOption func(*Model)

// WithTimeout sets the model timeout.
func WithTimeout(d time.Duration) ModelOption {
	return func(m *Model) { m.Timeout = d }
}

// WithProvider sets the provider name.
func WithProvider(p string) ModelOption {
	return func(m *Model) { m.Provider = p }
}

// Define creates a new model definition.
func Define(id string, name string, category Category, opts ...ModelOption) *Model {
	m := &Model{
		Identifier: id,
		Name:       name,
		Category:   category,
		Timeout:    10 * time.Minute,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Required adds required fields to the model.
func (m *Model) Required(fields ...Field) *Model {
	for _, f := range fields {
		f.Required = true
		m.requiredFields = append(m.requiredFields, f)
	}
	return m
}

// Optional adds optional fields to the model.
func (m *Model) Optional(fields ...Field) *Model {
	for _, f := range fields {
		f.Required = false
		m.optionalFields = append(m.optionalFields, f)
	}
	return m
}

// RequiredFields returns all required fields for this model.
// Use this to discover what parameters are required.
func (m *Model) RequiredFields() []Field {
	result := make([]Field, len(m.requiredFields))
	copy(result, m.requiredFields)
	return result
}

// OptionalFields returns all optional fields for this model.
// Use this to discover what optional parameters are available.
func (m *Model) OptionalFields() []Field {
	result := make([]Field, len(m.optionalFields))
	copy(result, m.optionalFields)
	return result
}

// AllFields returns all fields (required + optional) for this model.
func (m *Model) AllFields() []Field {
	result := make([]Field, 0, len(m.requiredFields)+len(m.optionalFields))
	result = append(result, m.requiredFields...)
	result = append(result, m.optionalFields...)
	return result
}

// Validate validates parameters against the model schema.
func (m *Model) Validate(params Params) error {
	// Check required fields
	for _, f := range m.requiredFields {
		v, exists := params[f.Name]
		if !exists || v == nil {
			return fmt.Errorf("%s: %s is required", m.Identifier, f.Name)
		}
		if err := f.Validate(v); err != nil {
			return fmt.Errorf("%s: %w", m.Identifier, err)
		}
	}

	// Validate optional fields if present
	for _, f := range m.optionalFields {
		if v, exists := params[f.Name]; exists {
			if err := f.Validate(v); err != nil {
				return fmt.Errorf("%s: %w", m.Identifier, err)
			}
		}
	}

	return nil
}

// Params represents model input parameters.
type Params map[string]any

// Request creates a new request builder for this model.
func (m *Model) Request() *RequestBuilder {
	return &RequestBuilder{
		model:  m,
		params: make(Params),
	}
}

// RequestBuilder builds a request with fluent API.
type RequestBuilder struct {
	model  *Model
	params Params
}

// Set sets a parameter value.
func (b *RequestBuilder) Set(name string, value any) *RequestBuilder {
	b.params[name] = value
	return b
}

// Prompt is a convenience method for setting the prompt.
func (b *RequestBuilder) Prompt(p string) *RequestBuilder {
	return b.Set("prompt", p)
}

// NegativePrompt sets the negative prompt.
func (b *RequestBuilder) NegativePrompt(p string) *RequestBuilder {
	return b.Set("negative_prompt", p)
}

// AspectRatio sets the aspect ratio.
func (b *RequestBuilder) AspectRatio(ar string) *RequestBuilder {
	return b.Set("aspect_ratio", ar)
}

// ImageURLs sets the image URLs.
func (b *RequestBuilder) ImageURLs(urls ...string) *RequestBuilder {
	return b.Set("image_urls", urls)
}

// ImageURL sets a single image URL.
func (b *RequestBuilder) ImageURL(url string) *RequestBuilder {
	return b.Set("image_url", url)
}

// ImageInput sets the image input (for models that use image_input parameter).
func (b *RequestBuilder) ImageInput(urls ...string) *RequestBuilder {
	return b.Set("image_input", urls)
}

// Seed sets the random seed for reproducibility.
func (b *RequestBuilder) Seed(seed int) *RequestBuilder {
	return b.Set("seed", seed)
}

// Resolution sets the output resolution (e.g., "1K", "2K", "4K").
func (b *RequestBuilder) Resolution(r string) *RequestBuilder {
	return b.Set("resolution", r)
}

// OutputFormat sets the output format (e.g., "png", "jpg").
func (b *RequestBuilder) OutputFormat(f string) *RequestBuilder {
	return b.Set("output_format", f)
}

// Duration sets the video duration (e.g., "5s", "10s").
func (b *RequestBuilder) Duration(d string) *RequestBuilder {
	return b.Set("duration", d)
}

// Style sets the generation style.
func (b *RequestBuilder) Style(s string) *RequestBuilder {
	return b.Set("style", s)
}

// Strength sets the strength parameter (typically 0.0-1.0).
func (b *RequestBuilder) Strength(s float64) *RequestBuilder {
	return b.Set("strength", s)
}

// ImageSize sets the image size parameter.
func (b *RequestBuilder) ImageSize(size string) *RequestBuilder {
	return b.Set("image_size", size)
}

// NumImages sets the number of images to generate.
func (b *RequestBuilder) NumImages(n string) *RequestBuilder {
	return b.Set("num_images", n)
}

// Params returns the current parameters.
func (b *RequestBuilder) Params() Params {
	return b.params
}

// Validate validates the current parameters.
func (b *RequestBuilder) Validate() error {
	return b.model.Validate(b.params)
}

// Model returns the model for this builder.
func (b *RequestBuilder) Model() *Model {
	return b.model
}

// Generate executes the request and waits for completion.
// It validates parameters, creates the task, and waits for results.
func (b *RequestBuilder) Generate(ctx context.Context, client Generator, opts ...GenerateOption) (*Result, error) {
	// Validate parameters
	if err := b.Validate(); err != nil {
		return nil, err
	}

	// Build wait options from generate options
	cfg := DefaultGenerateConfig(b.model)
	ApplyOptions(cfg, opts...)
	timeout, pollInterval, maxPollInterval, _ := cfg.GetConfig()

	waitOpts := WaitOptions{
		Timeout:         timeout,
		PollInterval:    pollInterval,
		MaxPollInterval: maxPollInterval,
	}

	return client.GenerateModel(ctx, b.model.Identifier, b.params, waitOpts)
}

// Result represents the result of a generation.
type Result struct {
	TaskID   string
	URLs     []string
	CostTime int64
	Raw      json.RawMessage
}

// WaitOptions for waiting on task completion.
type WaitOptions struct {
	Timeout         time.Duration
	PollInterval    time.Duration
	MaxPollInterval time.Duration
}

// Generator is the interface that clients must implement to execute model requests.
type Generator interface {
	// GenerateModel creates and waits for a model task.
	GenerateModel(ctx context.Context, modelID string, input any, opts WaitOptions) (*Result, error)
}

// GenerateOption configures generation behavior.
type GenerateOption func(*generateConfig)

type generateConfig struct {
	timeout         time.Duration
	pollInterval    time.Duration
	maxPollInterval time.Duration
	callbackURL     string
}

// WithGenTimeout sets the generation timeout.
func WithGenTimeout(d time.Duration) GenerateOption {
	return func(c *generateConfig) { c.timeout = d }
}

// WithGenPollInterval sets the polling interval.
func WithGenPollInterval(d time.Duration) GenerateOption {
	return func(c *generateConfig) { c.pollInterval = d }
}

// WithGenMaxPollInterval sets the max polling interval.
func WithGenMaxPollInterval(d time.Duration) GenerateOption {
	return func(c *generateConfig) { c.maxPollInterval = d }
}

// WithGenCallback sets a callback URL for async notification.
func WithGenCallback(url string) GenerateOption {
	return func(c *generateConfig) { c.callbackURL = url }
}

// DefaultGenerateConfig returns the default generate config for a model.
func DefaultGenerateConfig(m *Model) *generateConfig {
	return &generateConfig{
		timeout:         m.Timeout,
		pollInterval:    500 * time.Millisecond,
		maxPollInterval: 5 * time.Second,
	}
}

// ApplyOptions applies options to a config.
func ApplyOptions(cfg *generateConfig, opts ...GenerateOption) {
	for _, opt := range opts {
		opt(cfg)
	}
}

// GetConfig returns config values.
func (c *generateConfig) GetConfig() (timeout, pollInterval, maxPollInterval time.Duration, callbackURL string) {
	return c.timeout, c.pollInterval, c.maxPollInterval, c.callbackURL
}
