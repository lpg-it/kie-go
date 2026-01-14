# Changelog

All notable changes to the KIE Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-01-11

### Added

#### Core Features
- High-performance HTTP client with HTTP/2 and connection pooling
- Task creation and status query APIs
- `WaitForTask` with exponential backoff polling
- Image models: Nano Banana Pro, Seedream 4.5 Edit
- Video model: Seedance 1.5 Pro

#### Reliability
- Circuit breaker for cascading failure protection
- Rate limiter with token bucket algorithm
- Retry with exponential backoff and jitter
- Result validation to prevent false negatives
- `Retry-After` header support for 429 responses

#### Developer Experience
- Mock client (`kietest` package) for unit testing
- Progress callback for real-time status updates
- Result downloader for auto-saving files
- Request middleware support
- Metrics and Tracer interfaces for observability
- Batch processing for concurrent operations
- Webhook handler with HMAC signature verification

#### Performance
- Zero-allocation buffer pooling
- 10k+ QPS high-concurrency configuration
- HTTP/2 multiplexing
