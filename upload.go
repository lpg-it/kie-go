package kie

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// FileUploader provides file upload capabilities for the KIE API.
// Files are uploaded to KIE's temporary storage and can be used as input
// for image-to-image, image-to-video, and other models.
//
// IMPORTANT: Uploaded files are temporary and will be automatically deleted after 3 days.
//
// Example:
//
//	uploader := client.FileUploader()
//
//	// Upload a local file
//	result, err := uploader.UploadFile(ctx, "/path/to/image.jpg", nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Use the uploaded file URL for image-to-video
//	video, err := video.GenerateSeedancePro(client, ctx, &video.SeedanceProInput{
//	    ImageURL: result.FileURL,
//	    Prompt:   "Make it move",
//	})
type FileUploader struct {
	client *Client
}

// FileUploader returns a FileUploader for the client.
func (c *Client) FileUploader() *FileUploader {
	return &FileUploader{client: c}
}

// UploadOptions configures file upload behavior.
type UploadOptions struct {
	// UploadPath is the directory path for the file (e.g., "images/user-uploads")
	UploadPath string

	// FileName is the custom file name. If empty, a random name is generated.
	// Note: Uploading with the same name overwrites the previous file.
	FileName string
}

// UploadResult represents the result of a file upload.
type UploadResult struct {
	// FileID is the unique identifier for the uploaded file
	FileID string `json:"fileId"`

	// FileName is the name of the uploaded file
	FileName string `json:"fileName"`

	// OriginalName is the original name before upload
	OriginalName string `json:"originalName"`

	// FileSize is the file size in bytes
	FileSize int64 `json:"fileSize"`

	// MimeType is the MIME type of the file
	MimeType string `json:"mimeType"`

	// UploadPath is the path where the file was uploaded
	UploadPath string `json:"uploadPath"`

	// FileURL is the URL to access the file (use this for API inputs)
	FileURL string `json:"fileUrl"`

	// DownloadURL is the direct download URL
	DownloadURL string `json:"downloadUrl"`

	// UploadTime is when the file was uploaded
	UploadTime string `json:"uploadTime"`

	// ExpiresAt is when the file will be deleted (3 days after upload)
	ExpiresAt string `json:"expiresAt"`
}

// uploadResponse is the API response structure.
type uploadResponse struct {
	Success bool         `json:"success"`
	Code    int          `json:"code"`
	Msg     string       `json:"msg"`
	Data    UploadResult `json:"data"`
}

// UploadFile uploads a local file to KIE's temporary storage.
// The returned FileURL can be used as input for image-to-image, image-to-video, etc.
//
// Example:
//
//	result, err := uploader.UploadFile(ctx, "/path/to/photo.jpg", &kie.UploadOptions{
//	    UploadPath: "images",
//	    FileName:   "my-photo.jpg",
//	})
func (u *FileUploader) UploadFile(ctx context.Context, filePath string, opts *UploadOptions) (*UploadResult, error) {
	// Open file
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("kie: failed to open file: %w", err)
	}
	defer f.Close()

	// Get file info
	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("kie: failed to stat file: %w", err)
	}

	// Create multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add file
	fileName := filepath.Base(filePath)
	if opts != nil && opts.FileName != "" {
		fileName = opts.FileName
	}

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("kie: failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, f); err != nil {
		return nil, fmt.Errorf("kie: failed to copy file data: %w", err)
	}

	// Add optional fields
	if opts != nil {
		if opts.UploadPath != "" {
			writer.WriteField("uploadPath", opts.UploadPath)
		}
		if opts.FileName != "" {
			writer.WriteField("fileName", opts.FileName)
		}
	}

	writer.Close()

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		u.client.baseURL+"/api/file-stream-upload", &buf)
	if err != nil {
		return nil, fmt.Errorf("kie: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+u.client.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute request
	resp, err := u.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kie: upload failed: %w", err)
	}
	defer resp.Body.Close()

	return u.parseResponse(resp, stat.Size())
}

// UploadFromURL uploads a file from a remote URL to KIE's temporary storage.
// This is useful for migrating files from other services.
//
// Example:
//
//	result, err := uploader.UploadFromURL(ctx, "https://example.com/image.jpg", nil)
func (u *FileUploader) UploadFromURL(ctx context.Context, fileURL string, opts *UploadOptions) (*UploadResult, error) {
	payload := map[string]string{
		"fileUrl": fileURL,
	}

	if opts != nil {
		if opts.UploadPath != "" {
			payload["uploadPath"] = opts.UploadPath
		}
		if opts.FileName != "" {
			payload["fileName"] = opts.FileName
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("kie: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		u.client.baseURL+"/api/file-url-upload", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("kie: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+u.client.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kie: upload failed: %w", err)
	}
	defer resp.Body.Close()

	return u.parseResponse(resp, 0)
}

// UploadBase64 uploads a base64-encoded file to KIE's temporary storage.
// The data should include the data URL prefix (e.g., "data:image/png;base64,...").
//
// Example:
//
//	result, err := uploader.UploadBase64(ctx, "data:image/png;base64,iVBORw0K...", nil)
func (u *FileUploader) UploadBase64(ctx context.Context, base64Data string, opts *UploadOptions) (*UploadResult, error) {
	payload := map[string]string{
		"base64Data": base64Data,
	}

	if opts != nil {
		if opts.UploadPath != "" {
			payload["uploadPath"] = opts.UploadPath
		}
		if opts.FileName != "" {
			payload["fileName"] = opts.FileName
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("kie: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		u.client.baseURL+"/api/file-base64-upload", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("kie: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+u.client.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kie: upload failed: %w", err)
	}
	defer resp.Body.Close()

	return u.parseResponse(resp, 0)
}

// UploadStream uploads data from an io.Reader to KIE's temporary storage using
// multipart/form-data binary transmission. This is the most flexible upload method,
// suitable for any data source (network streams, memory buffers, pipes, etc.).
//
// Unlike UploadBytes which uses Base64 encoding (33% size increase), UploadStream
// uses efficient binary transmission.
//
// The fileName parameter is required and should include the file extension.
//
// Example:
//
//	// Upload from HTTP response body (without saving to disk)
//	resp, _ := http.Get("https://example.com/image.jpg")
//	defer resp.Body.Close()
//	result, err := uploader.UploadStream(ctx, resp.Body, "image.jpg", nil)
//
//	// Upload from memory buffer
//	buf := bytes.NewReader(imageData)
//	result, err := uploader.UploadStream(ctx, buf, "photo.png", &kie.UploadOptions{
//	    UploadPath: "images",
//	})
//
//	// Upload from gzip reader
//	gzReader, _ := gzip.NewReader(compressedFile)
//	result, err := uploader.UploadStream(ctx, gzReader, "data.json", nil)
func (u *FileUploader) UploadStream(ctx context.Context, reader io.Reader, fileName string, opts *UploadOptions) (*UploadResult, error) {
	if reader == nil {
		return nil, fmt.Errorf("kie: reader cannot be nil")
	}
	if fileName == "" {
		return nil, fmt.Errorf("kie: fileName is required for stream upload")
	}

	// Create multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Use custom fileName from opts if provided
	if opts != nil && opts.FileName != "" {
		fileName = opts.FileName
	}

	// Add file part
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("kie: failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, reader); err != nil {
		return nil, fmt.Errorf("kie: failed to copy stream data: %w", err)
	}

	// Add optional fields
	if opts != nil {
		if opts.UploadPath != "" {
			if err := writer.WriteField("uploadPath", opts.UploadPath); err != nil {
				return nil, fmt.Errorf("kie: failed to write uploadPath field: %w", err)
			}
		}
		if opts.FileName != "" {
			if err := writer.WriteField("fileName", opts.FileName); err != nil {
				return nil, fmt.Errorf("kie: failed to write fileName field: %w", err)
			}
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("kie: failed to close multipart writer: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		u.client.baseURL+"/api/file-stream-upload", &buf)
	if err != nil {
		return nil, fmt.Errorf("kie: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+u.client.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute request
	resp, err := u.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kie: upload failed: %w", err)
	}
	defer resp.Body.Close()

	return u.parseResponse(resp, 0)
}

// UploadBytes uploads raw bytes to KIE's temporary storage.
// This is a convenience method for in-memory data.
//
// Note: This method uses Base64 encoding which increases data size by ~33%.
// For large data, consider using UploadStream instead for better efficiency.
//
// Example:
//
//	imageBytes, _ := io.ReadAll(imageReader)
//	result, err := uploader.UploadBytes(ctx, imageBytes, "image/jpeg", &kie.UploadOptions{
//	    FileName: "uploaded.jpg",
//	})
func (u *FileUploader) UploadBytes(ctx context.Context, data []byte, mimeType string, opts *UploadOptions) (*UploadResult, error) {
	// Convert to base64 data URL
	base64Data := fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data))
	return u.UploadBase64(ctx, base64Data, opts)
}

// parseResponse parses the upload API response.
func (u *FileUploader) parseResponse(resp *http.Response, expectedSize int64) (*UploadResult, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("kie: failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiResp uploadResponse
		if json.Unmarshal(body, &apiResp) == nil && apiResp.Msg != "" {
			return nil, newAPIError(resp.StatusCode, apiResp.Code, apiResp.Msg)
		}
		return nil, newAPIError(resp.StatusCode, resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var apiResp uploadResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("kie: failed to parse response: %w", err)
	}

	if !apiResp.Success {
		return nil, newAPIError(resp.StatusCode, apiResp.Code, apiResp.Msg)
	}

	return &apiResp.Data, nil
}
