# KIE Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/lpg-it/kie-go.svg)](https://pkg.go.dev/github.com/lpg-it/kie-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/lpg-it/kie-go)](https://goreportcard.com/report/github.com/lpg-it/kie-go)

A high-performance, secure, and stable Go SDK for the [KIE AI](https://kie.ai) platform.

## Features

| Feature | Description |
|---------|-------------|
| **Security** | Secure API key handling, webhook signature verification |
| **Stability** | Circuit breaker, retry with backoff, rate limiting |
| **High Performance** | HTTP/2, connection pooling, zero-allocation buffers |
| **Batch Processing** | Concurrent task creation and waiting |
| **Webhooks** | Callback handler with HMAC signature verification |
| **Account Credits** | Check balance, threshold alerts |
| **Download URLs** | Temporary download links for generated files |
| **Image Generation** | Google Imagen4, Nano Banana Pro, Seedream, Grok Imagine, etc. |
| **Video Generation** | Seedance 1.5 Pro, Kling Video, Runway Gen3, Pika Video |

## Installation

```bash
go get github.com/lpg-it/kie-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    kie "github.com/lpg-it/kie-go"
)

func main() {
    client := kie.NewClient("your-api-key")
    defer client.Close()

    ctx := context.Background()

    // Method 1: Using RequestBuilder with constants (Recommended)
    result, err := kie.Image.NanoBananaPro.Request().
        Prompt("A futuristic city at sunset").
        AspectRatio(kie.Ratio16x9).      // Use predefined constant
        Set(kie.ParamResolution, kie.Resolution4K).
        Generate(ctx, client)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Generated:", result.URLs)

    // Method 2: Using Generate function with constants
    result, err = kie.Generate(ctx, client, kie.NanoBananaPro, kie.Params{
        kie.ParamPrompt:      "A beautiful sunset",
        kie.ParamAspectRatio: kie.Ratio16x9,
        kie.ParamResolution:  kie.Resolution4K,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Generated:", result.URLs)
}
```

---

## Complete API Reference

### Table of Contents

- [Client](#client)
  - [Creating a Client](#creating-a-client)
  - [Client Options](#client-options)
  - [Client Methods](#client-methods)
- [Task Operations](#task-operations)
  - [Creating Tasks](#creating-tasks)
  - [Waiting for Tasks](#waiting-for-tasks)
  - [Wait Options](#wait-options)
- [Batch Operations](#batch-operations)
  - [Batch Methods](#batch-methods)
  - [Batch Options](#batch-options)
- [Model Generation](#model-generation)
  - [Using RequestBuilder](#using-requestbuilder)
  - [Using Generate Function](#using-generate-function)
  - [Generate Options](#generate-options)
  - [RequestBuilder Methods](#requestbuilder-methods)
- [Image Models](#image-models)
- [Video Models](#video-models)
- [Model Registry](#model-registry)
- [File Upload](#file-upload)
  - [FileUploader Methods](#fileuploader-methods)
- [File Download](#file-download)
  - [Download URL Methods](#download-url-methods)
  - [Downloader](#downloader)
  - [Downloader Options](#downloader-options)
  - [Downloader Methods](#downloader-methods)
- [Account Credits](#account-credits)
- [Webhook Handler](#webhook-handler)
  - [WebhookHandler Methods](#webhookhandler-methods)
- [Circuit Breaker](#circuit-breaker)
  - [CircuitBreaker Methods](#circuitbreaker-methods)
- [Rate Limiting](#rate-limiting)
- [Observability](#observability)
- [Error Handling](#error-handling)
  - [Error Types](#error-types)
  - [Error Check Functions](#error-check-functions)
  - [Error Variables](#error-variables)
- [Types Reference](#types-reference)
- [Constants](#constants)

---

## Client

### Creating a Client

```go
func NewClient(apiKey string, opts ...Option) *Client
```

Creates a new KIE API client configured with HTTP/2, connection pooling, and automatic retry.

```go
// Basic usage
client := kie.NewClient("your-api-key")

// With options
client := kie.NewClient("your-api-key",
    kie.WithTimeout(60*time.Second),
    kie.WithMaxRetries(5),
    kie.WithRateLimit(10),
)
```

### Client Options

| Function | Description |
|----------|-------------|
| `WithBaseURL(url string)` | Set custom API base URL |
| `WithHTTPClient(client *http.Client)` | Set custom HTTP client |
| `WithTimeout(d time.Duration)` | Set request timeout |
| `WithRetry(config *RetryConfig)` | Configure retry behavior |
| `WithMaxRetries(n int)` | Set maximum retry attempts |
| `WithDebug(enabled bool)` | Enable debug logging |
| `WithTransport(config *TransportConfig)` | Configure HTTP transport |
| `WithCircuitBreaker(config *CircuitBreakerConfig)` | Enable circuit breaker |
| `WithRateLimit(rps int)` | Enable rate limiting |
| `WithMetrics(metrics Metrics)` | Enable metrics collection |
| `WithTracer(tracer Tracer)` | Enable distributed tracing |

### Client Methods

| Method | Description |
|--------|-------------|
| `Close() error` | Release resources held by the client |
| `CreateTask(ctx, req) (*CreateTaskResponse, error)` | Create a generation task |
| `GetTaskStatus(ctx, taskID) (*TaskInfo, error)` | Get task status |
| `WaitForTask(ctx, taskID, opts...) (*TaskInfo, error)` | Wait for task completion |
| `CreateTasksBatch(ctx, requests, opts...) <-chan BatchResult` | Create multiple tasks concurrently |
| `WaitForTasksBatch(ctx, taskIDs, opts...) <-chan BatchResult` | Wait for multiple tasks concurrently |
| `GetCredits(ctx) (int, error)` | Get current credit balance |
| `CheckCredits(ctx, threshold) (ok, balance, error)` | Check if credits meet threshold |
| `HasSufficientCredits(ctx, threshold) (bool, error)` | Simple credit sufficiency check |
| `GetDownloadURL(ctx, fileURL) (string, error)` | Get temporary download URL |
| `GetDownloadURLs(ctx, fileURLs) ([]string, error)` | Get multiple download URLs |
| `FileUploader() *FileUploader` | Get file uploader |
| `GenerateModel(ctx, modelID, input, opts) (*Result, error)` | Execute model generation (Generator interface) |

---

## Task Operations

### Creating Tasks

```go
func (c *Client) CreateTask(ctx context.Context, req *CreateTaskRequest) (*CreateTaskResponse, error)
```

```go
task, err := client.CreateTask(ctx, &kie.CreateTaskRequest{
    Model: "nano-banana-pro",
    Input: kie.Params{
        "prompt":       "A beautiful sunset",
        "aspect_ratio": "16:9",
    },
    CallbackURL: "https://your-server.com/webhook", // Optional
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("Task ID:", task.TaskID)
```

### Waiting for Tasks

```go
func (c *Client) GetTaskStatus(ctx context.Context, taskID string) (*TaskInfo, error)
func (c *Client) WaitForTask(ctx context.Context, taskID string, opts ...WaitOption) (*TaskInfo, error)
```

```go
// Non-blocking status check
status, err := client.GetTaskStatus(ctx, taskID)

// Blocking wait with exponential backoff
info, err := client.WaitForTask(ctx, taskID,
    kie.WithWaitTimeout(5*time.Minute),
    kie.WithPollInterval(time.Second),
    kie.WithMaxPollInterval(10*time.Second),
)
if err != nil {
    log.Fatal(err)
}

urls, _ := info.GetResultURLs()
fmt.Println("Results:", urls)
```

### Wait Options

| Function | Description |
|----------|-------------|
| `WithWaitTimeout(d time.Duration)` | Set maximum wait time |
| `WithPollInterval(d time.Duration)` | Set initial polling interval |
| `WithMaxPollInterval(d time.Duration)` | Set maximum polling interval |
| `WithProgressCallback(cb ProgressCallback)` | Set progress notification callback |

**Progress Callback:**

```go
// Track progress during long waits
info, err := client.WaitForTask(ctx, taskID,
    kie.WithProgressCallback(func(taskID string, elapsed time.Duration, pollCount int) {
        log.Printf("Task %s: waiting %.1fs (poll #%d)", taskID, elapsed.Seconds(), pollCount)
    }),
)
```

**Configuration Defaults:**

```go
func DefaultWaitConfig() *WaitConfig
// Returns: Timeout=10min, PollInterval=500ms, MaxPollInterval=5s, PollMultiplier=1.5
```

---

## Batch Operations

### Batch Methods

```go
func (c *Client) CreateTasksBatch(ctx context.Context, requests []*CreateTaskRequest, opts ...BatchOption) <-chan BatchResult
func (c *Client) WaitForTasksBatch(ctx context.Context, taskIDs []string, opts ...BatchOption) <-chan BatchResult
```

```go
// Create multiple tasks concurrently
requests := []*kie.CreateTaskRequest{
    {Model: "nano-banana-pro", Input: kie.Params{"prompt": "Cat"}},
    {Model: "nano-banana-pro", Input: kie.Params{"prompt": "Dog"}},
    {Model: "nano-banana-pro", Input: kie.Params{"prompt": "Bird"}},
}

for result := range client.CreateTasksBatch(ctx, requests, kie.WithConcurrency(5)) {
    if result.Error != nil {
        log.Printf("Task %d failed: %v", result.Index, result.Error)
        continue
    }
    log.Printf("Created task %d: %s", result.Index, result.TaskID)
}

// Wait for multiple tasks concurrently
taskIDs := []string{"task1", "task2", "task3"}
for result := range client.WaitForTasksBatch(ctx, taskIDs, kie.WithConcurrency(3)) {
    if result.Error != nil {
        log.Printf("Task %s failed: %v", result.TaskID, result.Error)
        continue
    }
    urls, _ := result.Info.GetResultURLs()
    log.Printf("Task %s completed: %v", result.TaskID, urls)
}
```

### Batch Options

| Function | Description |
|----------|-------------|
| `WithConcurrency(n int)` | Set concurrency level (default: 5) |
| `WithBatchTimeout(d time.Duration)` | Set timeout for batch operations |

**Configuration Defaults:**

```go
func DefaultBatchConfig() *BatchConfig
// Returns: Concurrency=5, Timeout=10min, PollInterval=500ms, MaxPollInterval=5s
```

### BatchProcessor (High-Level API)

For simpler batch processing, use the `BatchProcessor` wrapper:

```go
processor := kie.NewBatchProcessor(client,
    kie.WithConcurrency(10),
    kie.WithBatchTimeout(15*time.Minute),
)

// Create and wait for all tasks
results, err := processor.CreateAndWait(ctx, requests)
for _, r := range results {
    if r.Error != nil {
        log.Printf("Task %d failed: %v", r.Index, r.Error)
    } else {
        log.Printf("Task %d completed: %s", r.Index, r.Info.TaskID)
    }
}

// Process RequestBuilders directly
builders := []*model.RequestBuilder{
    kie.NanoBananaPro.Request().Prompt("A sunset"),
    kie.NanoBananaPro.Request().Prompt("A sunrise"),
}
results, err := processor.ProcessBuilders(ctx, builders)

// With progress callback
results, err := processor.ProcessAll(ctx, requests, func(completed, total int) {
    log.Printf("Progress: %d/%d", completed, total)
})
```

**BatchProcessor Methods:**

| Method | Description |
|--------|-------------|
| `CreateAndWait(ctx, requests) ([]BatchResult, error)` | Create and wait for all tasks |
| `ProcessBuilders(ctx, builders) ([]BatchResult, error)` | Process RequestBuilders |
| `ProcessAll(ctx, requests, progressCb) ([]BatchResult, error)` | Process with progress callback |

---

## Model Generation

### Using RequestBuilder

The recommended way to use models with full type safety and validation.

```go
// Build and execute a request with constants (recommended)
result, err := kie.Image.GoogleImagen4.Request().
    Prompt("A futuristic cityscape").
    NegativePrompt("blurry, dark").
    AspectRatio(kie.Ratio16x9).          // Use constant
    Set(kie.ParamSeed, 12345).
    Generate(ctx, client)

if err != nil {
    log.Fatal(err)
}
fmt.Println("URLs:", result.URLs)
fmt.Println("Cost Time:", result.CostTime)
```

### Using Generate Function

```go
func Generate(ctx context.Context, client *Client, m *Model, params Params, opts ...GenerateOption) (*Result, error)
```

```go
result, err := kie.Generate(ctx, client, kie.NanoBananaPro, kie.Params{
    kie.ParamPrompt:      "A beautiful sunset",
    kie.ParamAspectRatio: kie.Ratio16x9,
    kie.ParamResolution:  kie.Resolution4K,
},
    kie.WithGenTimeout(5*time.Minute),
    kie.WithGenCallback("https://your-server.com/webhook"),
)
```

### Generate Options

| Function | Description |
|----------|-------------|
| `WithGenTimeout(d time.Duration)` | Set generation timeout |
| `WithGenPollInterval(d time.Duration)` | Set polling interval |
| `WithGenMaxPollInterval(d time.Duration)` | Set max polling interval |
| `WithGenCallback(url string)` | Set callback URL for async mode |

### RequestBuilder Methods

| Method | Description |
|--------|-------------|
| `Set(name string, value any) *RequestBuilder` | Set any parameter |
| `Prompt(p string) *RequestBuilder` | Set prompt parameter |
| `NegativePrompt(p string) *RequestBuilder` | Set negative_prompt parameter |
| `AspectRatio(ar string) *RequestBuilder` | Set aspect_ratio parameter |
| `ImageURLs(urls ...string) *RequestBuilder` | Set image_urls parameter |
| `ImageURL(url string) *RequestBuilder` | Set single image_url parameter |
| `ImageInput(urls ...string) *RequestBuilder` | Set image_input parameter |
| `Seed(seed int) *RequestBuilder` | Set random seed |
| `Resolution(r string) *RequestBuilder` | Set resolution (1K, 2K, 4K) |
| `OutputFormat(f string) *RequestBuilder` | Set output format (png, jpg) |
| `Duration(d string) *RequestBuilder` | Set video duration |
| `Style(s string) *RequestBuilder` | Set generation style |
| `Strength(s float64) *RequestBuilder` | Set strength (0.0-1.0) |
| `ImageSize(size string) *RequestBuilder` | Set image_size parameter |
| `NumImages(n string) *RequestBuilder` | Set num_images parameter |
| `Params() Params` | Get current parameters |
| `Validate() error` | Validate parameters |
| `Model() *Model` | Get the model |
| `Generate(ctx, client, opts...) (*Result, error)` | Execute generation |

**Example with convenience methods:**

```go
result, err := kie.IdeogramV3Remix.Request().
    Prompt("A futuristic city").
    ImageURL("https://example.com/input.jpg").
    Strength(0.8).
    Style(kie.StyleRealistic).           // Use constant
    ImageSize(kie.SizeSquareHD).         // Use constant
    NumImages("2").
    Generate(ctx, client)
```

---

## Parameter Discovery

### Discovering Model Parameters

Each model has defined required and optional fields. You can query them programmatically:

```go
// Get model field information
model := kie.NanoBananaPro

// List required fields
fmt.Println("Required fields:")
for _, f := range model.RequiredFields() {
    fmt.Printf("  - %s (%s): %s\n", f.Name, f.Type, f.Description)
}

// List optional fields
fmt.Println("Optional fields:")
for _, f := range model.OptionalFields() {
    fmt.Printf("  - %s (%s): %s\n", f.Name, f.Type, f.Description)
    if len(f.EnumVals) > 0 {
        fmt.Printf("    Allowed values: %v\n", f.EnumVals)
    }
    if f.Default != nil {
        fmt.Printf("    Default: %v\n", f.Default)
    }
}
```

**Example output for NanoBananaPro:**

```
Required fields:
  - prompt (string): Text description of the image to generate

Optional fields:
  - image_input ([]string): Input images to transform or use as reference
  - aspect_ratio (enum): Aspect ratio of the generated image
    Allowed values: [1:1 2:3 3:2 3:4 4:3 4:5 5:4 9:16 16:9 21:9 auto]
    Default: 1:1
  - resolution (enum): Resolution of the generated image
    Allowed values: [1K 2K 4K]
    Default: 1K
  - output_format (enum): Format of the output image
    Allowed values: [png jpg]
    Default: png
```

### Model Methods for Discovery

| Method | Description |
|--------|-------------|
| `RequiredFields() []Field` | Get all required fields |
| `OptionalFields() []Field` | Get all optional fields |
| `AllFields() []Field` | Get all fields (required + optional) |

### Field Properties

| Property | Type | Description |
|----------|------|-------------|
| `Name` | `string` | Parameter name (use as key) |
| `Type` | `FieldType` | Type: string, int, []string, enum |
| `Description` | `string` | Human-readable description |
| `Required` | `bool` | Whether field is required |
| `EnumVals` | `[]string` | Allowed values for enum types |
| `Default` | `any` | Default value if not specified |
| `MaxLength` | `int` | Max length for strings |
| `MaxItems` | `int` | Max items for arrays |
| `Min` | `*float64` | Min value for numbers |
| `Max` | `*float64` | Max value for numbers |

---

## Parameter Constants

Use predefined constants for type safety and IDE autocompletion:

### Aspect Ratio Constants

```go
const (
    Ratio1x1   = "1:1"
    Ratio16x9  = "16:9"
    Ratio9x16  = "9:16"
    Ratio4x3   = "4:3"
    Ratio3x4   = "3:4"
    Ratio3x2   = "3:2"
    Ratio2x3   = "2:3"
    Ratio5x4   = "5:4"
    Ratio4x5   = "4:5"
    Ratio21x9  = "21:9"
    Ratio9x21  = "9:21"
    RatioAuto  = "auto"
)
```

### Resolution Constants

```go
const (
    Resolution1K = "1K"
    Resolution2K = "2K"
    Resolution4K = "4K"
)
```

### Output Format Constants

```go
const (
    FormatPNG  = "png"
    FormatJPG  = "jpg"
    FormatJPEG = "jpeg"
)
```

### Video Duration Constants

```go
const (
    Duration4s  = "4s"
    Duration5s  = "5s"
    Duration8s  = "8s"
    Duration10s = "10s"
)
```

### Video Duration Sec Constants (without "s" suffix)

```go
const (
    Duration4Sec  = "4"
    Duration5Sec  = "5"
    Duration6Sec  = "6"
    Duration8Sec  = "8"
    Duration10Sec = "10"
    Duration12Sec = "12"
    Duration15Sec = "15"
    Duration25Sec = "25"
)
```

### Video Resolution Mode Constants

```go
const (
    VideoResolution480p  = "480p"
    VideoResolution580p  = "580p"
    VideoResolution720p  = "720p"
    VideoResolution1080p = "1080p"
    VideoResolution512P  = "512P"
    VideoResolution768P  = "768P"
    VideoResolution1080P = "1080P"
)
```

### Quality Constants

```go
const (
    QualityBasic  = "basic"
    QualityMedium = "medium"
    QualityHigh   = "high"
)
```

### Parameter Name Constants

```go
const (
    // Common parameters
    ParamPrompt         = "prompt"
    ParamNegativePrompt = "negative_prompt"
    ParamAspectRatio    = "aspect_ratio"
    ParamImageSize      = "image_size"
    ParamResolution     = "resolution"
    ParamOutputFormat   = "output_format"
    ParamSeed           = "seed"

    // Image input parameters
    ParamImage      = "image"
    ParamImageURL   = "image_url"
    ParamImageURLs  = "image_urls"
    ParamImageInput = "image_input"
    ParamInputURLs  = "input_urls"

    // Video parameters
    ParamDuration   = "duration"
    ParamMode       = "mode"
    ParamVideoURL   = "video_url"
    ParamVideoURLs  = "video_urls"
    ParamAudioURL   = "audio_url"

    // Count and scale parameters
    ParamN     = "n"
    ParamScale = "scale"

    // Quality and style parameters
    ParamQuality        = "quality"
    ParamStyle          = "style"
    ParamRenderingSpeed = "rendering_speed"
    ParamNumImages      = "num_images"

    // Advanced parameters
    ParamCfgScale            = "cfg_scale"
    ParamGuidanceScale       = "guidance_scale"
    ParamStrength            = "strength"
    ParamEndImageURL         = "end_image_url"
    ParamTailImageURL        = "tail_image_url"
    ParamPromptOptimizer     = "prompt_optimizer"
    ParamCameraFixed         = "camera_fixed"
    ParamEnableSafetyChecker = "enable_safety_checker"
)
```

### Usage with Constants

```go
// Using constants for type safety
result, err := kie.NanoBananaPro.Request().
    Prompt("A beautiful sunset").
    Set(kie.ParamAspectRatio, kie.Ratio16x9).
    Set(kie.ParamResolution, kie.Resolution4K).
    Set(kie.ParamOutputFormat, kie.FormatPNG).
    Generate(ctx, client)

// Or with Params
result, err := kie.Generate(ctx, client, kie.NanoBananaPro, kie.Params{
    kie.ParamPrompt:       "A beautiful sunset",
    kie.ParamAspectRatio:  kie.Ratio16x9,
    kie.ParamResolution:   kie.Resolution4K,
})
```

---

## Image Models

Access via `kie.Image.*` namespace or directly as `kie.*`:

| Model Variable | Identifier | Category |
|----------------|------------|----------|
| `GoogleImagen4` | `google/imagen4` | Text-to-Image |
| `GoogleImagen4Fast` | `google/imagen4-fast` | Text-to-Image |
| `GoogleImagen4Ultra` | `google/imagen4-ultra` | Text-to-Image |
| `GoogleNanoBanana` | `google/nano-banana` | Text-to-Image |
| `GoogleNanoBananaEdit` | `google/nano-banana-edit` | Image-Edit |
| `NanoBananaPro` | `nano-banana-pro` | Text-to-Image |
| `GoogleNanoBanana2` | `nano-banana-2` | Text-to-Image |
| `GrokImagineTextToImage` | `grok-imagine/text-to-image` | Text-to-Image |
| `GrokImagineImageToImage` | `grok-imagine/image-to-image` | Image-to-Image |
| `GrokImagineUpscale` | `grok-imagine/upscale` | Upscale |
| `Seedream45TextToImage` | `seedream/4.5-text-to-image` | Text-to-Image |
| `Seedream45Edit` | `seedream/4.5-edit` | Image-Edit |
| `BytedanceSeedreamV4TextToImage` | `bytedance/seedream-v4-text-to-image` | Text-to-Image |
| `BytedanceSeedreamV4Edit` | `bytedance/seedream-v4-edit` | Image-Edit |
| `RecraftCrispUpscale` | `recraft/crisp-upscale` | Upscale |
| `RecraftRemoveBackground` | `recraft/remove-background` | Image-Edit |
| `TopazImageUpscale` | `topaz/image-upscale` | Upscale |
| `GptImage15ImageToImage` | `gpt-image/1.5-image-to-image` | Image-to-Image |
| `GptImage15TextToImage` | `gpt-image/1.5-text-to-image` | Text-to-Image |
| `ZImage` | `z-image` | Text-to-Image |
| `Flux2ProImageToImage` | `flux-2/pro-image-to-image` | Image-to-Image |
| `Flux2FlexImageToImage` | `flux-2/flex-image-to-image` | Image-to-Image |
| `Flux2FlexTextToImage` | `flux-2/flex-text-to-image` | Text-to-Image |
| `Flux2ProTextToImage` | `flux-2/pro-text-to-image` | Text-to-Image |
| `IdeogramV3Reframe` | `ideogram/v3-reframe` | Image-Edit |
| `IdeogramV3TextToImage` | `ideogram/v3-text-to-image` | Text-to-Image |
| `IdeogramV3Edit` | `ideogram/v3-edit` | Image-Edit |
| `IdeogramV3Remix` | `ideogram/v3-remix` | Image-to-Image |
| `BytedanceSeedream` | `bytedance/seedream` | Text-to-Image |
| `Wan27Image` | `wan/2-7-image` | Text-to-Image |
| `Wan27ImagePro` | `wan/2-7-image-pro` | Text-to-Image |
| `QwenImageToImage` | `qwen/image-to-image` | Image-to-Image |
| `QwenTextToImage` | `qwen/text-to-image` | Text-to-Image |
| `QwenImageEdit` | `qwen/image-edit` | Image-Edit |
| `Qwen2ImageEdit` | `qwen2/image-edit` | Image-Edit |

**Usage:**

```go
// Via namespace with constants (recommended)
result, _ := kie.Image.GoogleImagen4.Request().
    Prompt("A beautiful sunset").
    AspectRatio(kie.Ratio16x9).
    Generate(ctx, client)

// Direct access
result, _ := kie.GoogleImagen4.Request().
    Prompt("A beautiful sunset").
    AspectRatio(kie.Ratio16x9).
    Generate(ctx, client)
```

---

## Video Models

Access via `kie.Video.*` namespace or directly as `kie.*`:

| Model Variable | Identifier | Category |
|----------------|------------|----------|
| `Seedance15Pro` | `seedance/1.5-pro` | Text-to-Video |
| `Seedance15ImageToVideo` | `seedance/1.5-image-to-video` | Image-to-Video |
| `KlingVideo` | `kling/video` | Text-to-Video |
| `RunwayGen3` | `runway/gen3` | Text-to-Video |
| `PikaVideo` | `pika/video` | Text-to-Video |
| `GrokImagineImageToVideo` | `grok-imagine/image-to-video` | Image-to-Video |
| `GrokImagineTextToVideo` | `grok-imagine/text-to-video` | Text-to-Video |
| `Kling26MotionControl` | `kling-2.6/motion-control` | Image-to-Video |
| `Kling30MotionControl` | `kling-3.0/motion-control` | Image-to-Video |
| `BytedanceSeedance15Pro` | `bytedance/seedance-1.5-pro` | Text-to-Video |
| `BytedanceSeedance20Fast` | `bytedance/seedance-2-fast` | Text-to-Video |
| `BytedanceSeedance20` | `bytedance/seedance-2` | Text-to-Video |
| `Wan26TextToVideo` | `wan/2-6-text-to-video` | Text-to-Video |
| `Wan26ImageToVideo` | `wan/2-6-image-to-video` | Image-to-Video |
| `Wan26VideoToVideo` | `wan/2-6-video-to-video` | Video-to-Video |
| `Wan27TextToVideo` | `wan/2-7-text-to-video` | Text-to-Video |
| `Wan27ImageToVideo` | `wan/2-7-image-to-video` | Image-to-Video |
| `Wan27ReferenceToVideo` | `wan/2-7-r2v` | Text-to-Video |
| `Wan27VideoEdit` | `wan/2-7-videoedit` | Video-to-Video |
| `Kling26ImageToVideo` | `kling-2.6/image-to-video` | Image-to-Video |
| `Kling26TextToVideo` | `kling-2.6/text-to-video` | Text-to-Video |
| `BytedanceV1ProFastImageToVideo` | `bytedance/v1-pro-fast-image-to-video` | Image-to-Video |
| `Hailuo23ImageToVideoPro` | `hailuo/2-3-image-to-video-pro` | Image-to-Video |
| `Hailuo23ImageToVideoStandard` | `hailuo/2-3-image-to-video-standard` | Image-to-Video |
| `Sora2ProStoryboard` | `sora-2-pro-storyboard` | Image-to-Video |
| `Sora2ProTextToVideo` | `sora-2-pro-text-to-video` | Text-to-Video |
| `Sora2ProImageToVideo` | `sora-2-pro-image-to-video` | Image-to-Video |
| `Sora2Characters` | `sora-2-characters` | Text-to-Video |
| `SoraWatermarkRemover` | `sora-watermark-remover` | Video-to-Video |
| `Kling25TurboTextToVideoPro` | `kling/v2-5-turbo-text-to-video-pro` | Text-to-Video |
| `Kling25TurboImageToVideoPro` | `kling/v2-5-turbo-image-to-video-pro` | Image-to-Video |
| `Wan25ImageToVideo` | `wan/2-5-image-to-video` | Image-to-Video |
| `Wan25TextToVideo` | `wan/2-5-text-to-video` | Text-to-Video |
| `Wan22AnimateMove` | `wan/2-2-animate-move` | Video-to-Video |
| `Wan22AnimateReplace` | `wan/2-2-animate-replace` | Video-to-Video |
| `TopazVideoUpscale` | `topaz/video-upscale` | Video-to-Video |
| `InfinitalkFromAudio` | `infinitalk/from-audio` | Image-to-Video |
| `Wan22A14bSpeechToVideoTurbo` | `wan/2-2-a14b-speech-to-video-turbo` | Image-to-Video |
| `KlingV1AvatarStandard` | `kling/v1-avatar-standard` | Image-to-Video |
| `KlingAiAvatarV1Pro` | `kling/ai-avatar-v1-pro` | Image-to-Video |
| `Wan22A14bTextToVideoTurbo` | `wan/2-2-a14b-text-to-video-turbo` | Text-to-Video |
| `Wan22A14bImageToVideoTurbo` | `wan/2-2-a14b-image-to-video-turbo` | Image-to-Video |
| `KlingV21MasterImageToVideo` | `kling/v2-1-master-image-to-video` | Image-to-Video |
| `KlingV21Pro` | `kling/v2-1-pro` | Image-to-Video |
| `KlingV21Standard` | `kling/v2-1-standard` | Image-to-Video |
| `KlingV21MasterTextToVideo` | `kling/v2-1-master-text-to-video` | Text-to-Video |
| `BytedanceV1ProImageToVideo` | `bytedance/v1-pro-image-to-video` | Image-to-Video |
| `BytedanceV1LiteImageToVideo` | `bytedance/v1-lite-image-to-video` | Image-to-Video |
| `BytedanceV1ProTextToVideo` | `bytedance/v1-pro-text-to-video` | Text-to-Video |
| `BytedanceV1LiteTextToVideo` | `bytedance/v1-lite-text-to-video` | Text-to-Video |
| `Hailuo02TextToVideoStandard` | `hailuo/02-text-to-video-standard` | Text-to-Video |
| `Hailuo02ImageToVideoStandard` | `hailuo/02-image-to-video-standard` | Image-to-Video |
| `Hailuo02ImageToVideoPro` | `hailuo/02-image-to-video-pro` | Image-to-Video |
| `Hailuo02TextToVideoPro` | `hailuo/02-text-to-video-pro` | Text-to-Video |
| `Sora2ImageToVideo` | `sora-2-image-to-video` | Image-to-Video |
| `Sora2TextToVideo` | `sora-2-text-to-video` | Text-to-Video |

**Usage:**

```go
result, _ := kie.Video.Seedance15Pro.Request().
    Prompt("A dancing robot").
    Duration(kie.Duration10s).           // Use constant
    AspectRatio(kie.Ratio16x9).          // Use constant
    Generate(ctx, client)
```

---

## Model Registry

Functions for accessing models programmatically:

| Function | Description |
|----------|-------------|
| `GetImageModel(id string) *Model` | Get image model by identifier |
| `GetVideoModel(id string) *Model` | Get video model by identifier |
| `GetModel(id string) *Model` | Get any model by identifier |
| `AllImageModels() []*Model` | Get all image models |
| `AllVideoModels() []*Model` | Get all video models |
| `AllModels() []*Model` | Get all models |

```go
// Get model dynamically
model := kie.GetModel("nano-banana-pro")
if model != nil {
    result, _ := model.Request().Prompt("Hello").Generate(ctx, client)
}

// List all models
for _, m := range kie.AllModels() {
    fmt.Printf("Model: %s (%s)\n", m.Name, m.Identifier)
}
```

---

## File Upload

### FileUploader Methods

Get uploader via `client.FileUploader()`:

| Method | Description |
|--------|-------------|
| `UploadFile(ctx, filePath, opts) (*UploadResult, error)` | Upload local file |
| `UploadFromURL(ctx, fileURL, opts) (*UploadResult, error)` | Upload from remote URL |
| `UploadBase64(ctx, base64Data, opts) (*UploadResult, error)` | Upload base64 data |
| `UploadBytes(ctx, data, mimeType, opts) (*UploadResult, error)` | Upload raw bytes |

```go
uploader := client.FileUploader()

// Upload local file
result, err := uploader.UploadFile(ctx, "/path/to/image.jpg", &kie.UploadOptions{
    UploadPath: "images",
    FileName:   "my-photo.jpg",
})

// Upload from URL
result, err := uploader.UploadFromURL(ctx, "https://example.com/image.jpg", nil)

// Upload base64
result, err := uploader.UploadBase64(ctx, "data:image/png;base64,iVBORw0K...", nil)

// Upload bytes
result, err := uploader.UploadBytes(ctx, imageBytes, "image/jpeg", nil)
```

### Upload + Generate Workflow (Recommended)

For image-to-X and video-to-X models, upload your local files first:

```go
uploader := client.FileUploader()

// Step 1: Upload local image
uploadResult, err := uploader.UploadFile(ctx, "./my-photo.jpg", nil)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Uploaded: %s (expires: %s)\n", uploadResult.FileURL, uploadResult.ExpiresAt)

// Step 2: Use uploaded URL for image-to-video
videoResult, err := kie.Seedance15ImageToVideo.Request().
    Prompt("Make the character wave hello").
    ImageURLs(uploadResult.FileURL).  // Use the uploaded file URL
    Duration("5s").
    Generate(ctx, client)

// Step 3: For video-to-video, upload a video file
videoUpload, err := uploader.UploadFile(ctx, "./input.mp4", nil)
if err != nil {
    log.Fatal(err)
}

transformResult, err := kie.Wan26VideoToVideo.Request().
    Prompt("Transform into anime style").
    Set(kie.ParamVideoURLs, []string{videoUpload.FileURL}).
    Generate(ctx, client)

// Step 4: For avatar with audio, upload both image and audio
audioUpload, err := uploader.UploadFile(ctx, "./speech.mp3", nil)
avatarResult, err := kie.KlingV1AvatarStandard.Request().
    Set(kie.ParamImageURL, uploadResult.FileURL).
    Set(kie.ParamAudioURL, audioUpload.FileURL).
    Prompt("Professional speaking").
    Generate(ctx, client)
```

> **Note:** Uploaded files are temporary and will be deleted after 15 days. This is sufficient for generation tasks.

---

## File Download

### Download URL Methods

```go
func (c *Client) GetDownloadURL(ctx context.Context, fileURL string) (string, error)
func (c *Client) GetDownloadURLs(ctx context.Context, fileURLs []string) ([]string, error)
```

```go
// Single file
downloadURL, err := client.GetDownloadURL(ctx, "https://tempfile.1f6c...")

// Multiple files (concurrent)
downloadURLs, err := client.GetDownloadURLs(ctx, []string{
    "https://tempfile.1f6c...",
    "https://tempfile.2f7d...",
})
```

### Downloader

```go
func NewDownloader(opts ...DownloaderOption) *Downloader
```

### Downloader Options

| Function | Description |
|----------|-------------|
| `WithDownloadConcurrency(n int)` | Set concurrent downloads (default: 5) |
| `WithDownloadHTTPClient(client *http.Client)` | Set custom HTTP client |
| `WithKIEClient(client *Client)` | Enable automatic URL conversion |

### Downloader Methods

| Method | Description |
|--------|-------------|
| `DownloadFromTaskInfo(ctx, info, outputDir) ([]DownloadResult, error)` | Download all results from TaskInfo |
| `DownloadURLs(ctx, urls, outputDir, prefix) ([]DownloadResult, error)` | Download multiple URLs |
| `DownloadSingle(ctx, url, outputPath) error` | Download single file |

```go
// Create downloader with KIE client for auto URL conversion
downloader := kie.NewDownloader(
    kie.WithKIEClient(client),
    kie.WithDownloadConcurrency(10),
)

// Download from task info
results, err := downloader.DownloadFromTaskInfo(ctx, taskInfo, "./output/")
for _, r := range results {
    if r.Error != nil {
        log.Printf("Failed: %v", r.Error)
    } else {
        log.Printf("Downloaded: %s (%d bytes)", r.LocalPath, r.Size)
    }
}

// Download single file
err := downloader.DownloadSingle(ctx, url, "./output/image.png")
```

### Utility Functions

| Function | Description |
|----------|-------------|
| `IsKIETempFileURL(url string) bool` | Check if URL is a KIE temp file |

---

## Account Credits

```go
// Get current balance
credits, err := client.GetCredits(ctx)
fmt.Printf("Credits: %d\n", credits)

// Check if above threshold
ok, balance, err := client.CheckCredits(ctx, 100)
if !ok {
    log.Printf("Low credits: %d (need 100)", balance)
}

// Simple boolean check
if ok, _ := client.HasSufficientCredits(ctx, 50); !ok {
    log.Fatal("Insufficient credits")
}
```

---

## Webhook Handler

### Creating Handler

```go
func NewWebhookHandler(secretKey string) *WebhookHandler
```

### WebhookHandler Methods

| Method | Description |
|--------|-------------|
| `ParseRequest(r *http.Request) (*TaskInfo, error)` | Parse HTTP request |
| `ParsePayload(body []byte) (*TaskInfo, error)` | Parse raw payload |
| `Handler(callback func(*TaskInfo)) http.HandlerFunc` | Create HTTP handler |
| `HandlerWithError(callback func(*TaskInfo) error) http.HandlerFunc` | Create handler with error support |

```go
handler := kie.NewWebhookHandler("your-secret-key")

// Simple handler
http.HandleFunc("/webhook", handler.Handler(func(info *kie.TaskInfo) {
    log.Printf("Task %s completed: %s", info.TaskID, info.State)
    urls, _ := info.GetResultURLs()
    // Process results...
}))

// Handler with error support (returns 500 to trigger retry)
http.HandleFunc("/webhook", handler.HandlerWithError(func(info *kie.TaskInfo) error {
    if err := processTask(info); err != nil {
        return err // Will return 500, KIE will retry
    }
    return nil
}))

// Manual parsing
http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
    info, err := handler.ParseRequest(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // Process info...
    w.WriteHeader(http.StatusOK)
})
```

---

## Circuit Breaker

### Configuration

```go
func DefaultCircuitBreakerConfig() *CircuitBreakerConfig
func NewCircuitBreaker(config *CircuitBreakerConfig) *CircuitBreaker
```

**Configuration:**

```go
type CircuitBreakerConfig struct {
    FailureThreshold int           // Failures before opening (default: 5)
    SuccessThreshold int           // Successes to close (default: 2)
    Timeout          time.Duration // Wait before testing recovery (default: 30s)
}
```

### CircuitBreaker Methods

| Method | Description |
|--------|-------------|
| `Allow() bool` | Check if request should be allowed |
| `RecordSuccess()` | Record successful request |
| `RecordFailure()` | Record failed request |
| `State() CircuitState` | Get current state |
| `Reset()` | Reset to closed state |
| `Stats() CircuitBreakerStats` | Get statistics |

```go
// Enable via client option
client := kie.NewClient("api-key",
    kie.WithCircuitBreaker(&kie.CircuitBreakerConfig{
        FailureThreshold: 5,
        SuccessThreshold: 2,
        Timeout:          30 * time.Second,
    }),
)

// Manual usage
cb := kie.NewCircuitBreaker(kie.DefaultCircuitBreakerConfig())
if cb.Allow() {
    err := doRequest()
    if err != nil {
        cb.RecordFailure()
    } else {
        cb.RecordSuccess()
    }
}
```

**Circuit States:**

| Constant | Description |
|----------|-------------|
| `CircuitClosed` | Requests are allowed |
| `CircuitOpen` | Requests are blocked |
| `CircuitHalfOpen` | Testing if service recovered |

---

## Rate Limiting

```go
client := kie.NewClient("api-key",
    kie.WithRateLimit(10), // 10 requests per second
)
```

---

## Observability

### Metrics Interface

```go
type Metrics interface {
    IncCounter(name string, labels ...string)
    ObserveHistogram(name string, value float64, labels ...string)
}
```

### Tracer Interface

```go
type Tracer interface {
    Start(ctx context.Context, name string) (context.Context, Span)
}

type Span interface {
    SetAttribute(key string, value interface{})
    RecordError(err error)
    End()
}
```

```go
client := kie.NewClient("api-key",
    kie.WithMetrics(myMetricsCollector),
    kie.WithTracer(myTracer),
)
```

---

## Error Handling

### Error Types

| Type | Description |
|------|-------------|
| `*APIError` | Error returned by KIE API |
| `*TaskFailedError` | Task failed on server side |
| `*ResultValidationError` | Result data missing or invalid |

**APIError Fields:**

```go
type APIError struct {
    HTTPStatus int
    Code       int
    Message    string
    RequestID  string
    RetryAfter time.Duration
}
```

**TaskFailedError Fields:**

```go
type TaskFailedError struct {
    TaskID   string
    FailCode string
    FailMsg  string
}
```

### Error Check Functions

| Function | Description |
|----------|-------------|
| `IsRetryable(err error) bool` | Check if error can be retried |
| `IsAuthError(err error) bool` | Check if authentication error (401) |
| `IsRateLimitError(err error) bool` | Check if rate limit error (429) |
| `IsInsufficientFundError(err error) bool` | Check if insufficient funds (402) |
| `IsValidationError(err error) bool` | Check if validation error (422) |

```go
result, err := client.WaitForTask(ctx, taskID)
if err != nil {
    switch {
    case errors.Is(err, kie.ErrCircuitOpen):
        log.Fatal("Service unavailable, circuit breaker open")
    case kie.IsAuthError(err):
        log.Fatal("Invalid API key")
    case kie.IsRateLimitError(err):
        log.Println("Rate limited, will retry...")
    case kie.IsInsufficientFundError(err):
        log.Fatal("Please top up your account")
    case errors.Is(err, kie.ErrTaskFailed):
        var taskErr *kie.TaskFailedError
        errors.As(err, &taskErr)
        log.Printf("Task failed: %s - %s", taskErr.FailCode, taskErr.FailMsg)
    default:
        log.Printf("Unknown error: %v", err)
    }
}
```

### Error Variables

| Variable | Description |
|----------|-------------|
| `ErrNilInput` | Input cannot be nil |
| `ErrEmptyTaskID` | Task ID cannot be empty |
| `ErrTaskFailed` | Task failed (sentinel) |
| `ErrTimeout` | Operation timed out |
| `ErrMaxRetriesExceeded` | Max retries exceeded |
| `ErrCircuitOpen` | Circuit breaker is open |
| `ErrResultValidation` | Result validation failed (sentinel) |

---

## Types Reference

### Core Types

| Type | Description |
|------|-------------|
| `Client` | API client |
| `Model` | Model definition (re-exported from model package) |
| `Params` | Model input parameters (map[string]any) |
| `Result` | Generation result |
| `RequestBuilder` | Fluent request builder |
| `Field` | Model field definition |
| `Category` | Model category |

### Task Types

| Type | Description |
|------|-------------|
| `CreateTaskRequest` | Request to create task |
| `CreateTaskResponse` | Response from CreateTask |
| `TaskInfo` | Detailed task information |
| `TaskState` | Task state (waiting/success/fail) |
| `BatchResult` | Result of batch operation |

**TaskInfo Methods:**

| Method | Description |
|--------|-------------|
| `GetResultURLs() ([]string, error)` | Parse and get result URLs |

**TaskState Methods:**

| Method | Description |
|--------|-------------|
| `IsTerminal() bool` | Check if in terminal state |
| `IsSuccess() bool` | Check if succeeded |

### Upload Types

| Type | Description |
|------|-------------|
| `FileUploader` | File upload helper |
| `UploadOptions` | Upload configuration |
| `UploadResult` | Upload result with URLs |

### Download Types

| Type | Description |
|------|-------------|
| `Downloader` | File download helper |
| `DownloadResult` | Download result |

### Webhook Types

| Type | Description |
|------|-------------|
| `WebhookHandler` | Webhook request handler |
| `WebhookPayload` | Webhook callback structure |

### Configuration Types

| Type | Description |
|------|-------------|
| `RetryConfig` | Retry behavior configuration |
| `TransportConfig` | HTTP transport configuration |
| `CircuitBreakerConfig` | Circuit breaker configuration |
| `WaitConfig` | Wait behavior configuration |
| `BatchConfig` | Batch operation configuration |

### Configuration Factory Functions

| Function | Returns |
|----------|---------|
| `DefaultRetryConfig()` | Default retry config |
| `NoRetry()` | Config that disables retry |
| `DefaultTransportConfig()` | Default transport config |
| `HighConcurrencyTransportConfig()` | Config for 10k+ QPS |
| `DefaultCircuitBreakerConfig()` | Default circuit breaker config |
| `DefaultWaitConfig()` | Default wait config |
| `DefaultBatchConfig()` | Default batch config |

---

## Constants

### Version

```go
const Version = "0.1.0"
```

### Default Values

```go
const DefaultBaseURL = "https://api.kie.ai"
const DefaultTimeout = 30 * time.Second
```

### Task States

```go
const (
    TaskStateWaiting TaskState = "waiting"
    TaskStateSuccess TaskState = "success"
    TaskStateFail    TaskState = "fail"
)
```

### Circuit States

```go
const (
    CircuitClosed   CircuitState = iota
    CircuitOpen
    CircuitHalfOpen
)
```

### Error Codes

```go
const (
    ErrCodeSuccess          = 200
    ErrCodeBadRequest       = 400
    ErrCodeUnauthorized     = 401
    ErrCodeInsufficientFund = 402
    ErrCodeNotFound         = 404
    ErrCodeValidation       = 422
    ErrCodeRateLimit        = 429
    ErrCodeInternal         = 500
)
```

### Category Constants

```go
const (
    CategoryTextToImage  = model.CategoryTextToImage
    CategoryImageToImage = model.CategoryImageToImage
    CategoryImageEdit    = model.CategoryImageEdit
    CategoryUpscale      = model.CategoryUpscale
    CategoryTextToVideo  = model.CategoryTextToVideo
    CategoryImageToVideo = model.CategoryImageToVideo
    CategoryVideoToVideo = model.CategoryVideoToVideo
)
```

---

## License

MIT License
