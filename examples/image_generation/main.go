// Example: Basic image generation using KIE SDK
//
// This example demonstrates the declarative architecture:
// 1. Create a client with authentication
// 2. Use declarative model definitions with fluent Builder API
// 3. Use simple Params map for quick usage
// 4. Use predefined constants for type safety
// 5. Discover model parameters programmatically
//
// Usage:
//
//	KIE_API_KEY=your_api_key go run main.go
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
	// Get API key from environment
	apiKey := os.Getenv("KIE_API_KEY")
	if apiKey == "" {
		log.Fatal("KIE_API_KEY environment variable is required")
	}

	// Create client with options
	client := kie.NewClient(
		apiKey,
		kie.WithTimeout(30*time.Second),
		kie.WithMaxRetries(3),
	)
	defer client.Close()

	ctx := context.Background()

	// Example 1: Using Builder pattern (recommended)
	fmt.Println("=== Example 1: Builder Pattern ===")
	if err := generateWithBuilder(ctx, client); err != nil {
		log.Printf("Builder example failed: %v", err)
	}

	// Example 2: Using Params map (quick & simple)
	fmt.Println("\n=== Example 2: Params Map ===")
	if err := generateWithParams(ctx, client); err != nil {
		log.Printf("Params example failed: %v", err)
	}

	// Example 3: Using predefined constants
	fmt.Println("\n=== Example 3: Using Constants ===")
	if err := generateWithConstants(ctx, client); err != nil {
		log.Printf("Constants example failed: %v", err)
	}

	// Example 4: Parameter discovery
	fmt.Println("\n=== Example 4: Parameter Discovery ===")
	discoverModelParameters()

	// Example 5: Image editing
	fmt.Println("\n=== Example 5: Image Editing ===")
	if err := editImage(ctx, client); err != nil {
		log.Printf("Edit example failed: %v", err)
	}

	// Example 6: Upscale image
	fmt.Println("\n=== Example 6: Upscale Image ===")
	if err := upscaleImage(ctx, client); err != nil {
		log.Printf("Upscale example failed: %v", err)
	}

	// Example 7: Upload + Generate workflow (recommended for local files)
	fmt.Println("\n=== Example 7: Upload + Generate Workflow ===")
	if err := uploadAndGenerate(ctx, client); err != nil {
		log.Printf("Upload+Generate example failed: %v", err)
	}
}

func generateWithBuilder(ctx context.Context, client *kie.Client) error {
	// Fluent Builder API with constants - recommended approach
	result, err := kie.GoogleImagen4.Request().
		Prompt("A futuristic city at sunset with flying cars").
		AspectRatio(kie.Ratio16x9). // Use constant instead of "16:9"
		NegativePrompt("blurry, low quality").
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	fmt.Printf("Task completed: %s\n", result.TaskID)
	fmt.Printf("Generated %d image(s):\n", len(result.URLs))
	for i, url := range result.URLs {
		fmt.Printf("  %d: %s\n", i+1, url)
	}
	return nil
}

func generateWithParams(ctx context.Context, client *kie.Client) error {
	// Params map with constants - type-safe key names
	result, err := kie.Generate(ctx, client, kie.GoogleImagen4, kie.Params{
		kie.ParamPrompt:      "A serene mountain landscape at dawn",
		kie.ParamAspectRatio: kie.Ratio16x9,
	})

	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	fmt.Printf("Task completed: %s\n", result.TaskID)
	fmt.Printf("Generated: %v\n", result.URLs)
	return nil
}

func generateWithConstants(ctx context.Context, client *kie.Client) error {
	// Using predefined constants for type safety and IDE support
	result, err := kie.NanoBananaPro.Request().
		Prompt("A beautiful sunset over the ocean").
		Set(kie.ParamAspectRatio, kie.Ratio16x9).
		Set(kie.ParamResolution, kie.Resolution4K).
		Set(kie.ParamOutputFormat, kie.FormatPNG).
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	fmt.Printf("Generated with constants: %s\n", result.URLs[0])

	// Alternative: Using constants with Params map
	result2, err := kie.Generate(ctx, client, kie.NanoBananaPro, kie.Params{
		kie.ParamPrompt:       "A starry night sky",
		kie.ParamAspectRatio:  kie.Ratio9x16,
		kie.ParamResolution:   kie.Resolution2K,
		kie.ParamOutputFormat: kie.FormatJPG,
	})

	if err != nil {
		return fmt.Errorf("generate with params: %w", err)
	}

	fmt.Printf("Generated with params constants: %s\n", result2.URLs[0])
	return nil
}

func discoverModelParameters() {
	// Discover what parameters a model accepts
	model := kie.NanoBananaPro

	fmt.Printf("Model: %s (%s)\n", model.Name, model.Identifier)

	// List required fields
	fmt.Println("\nRequired fields:")
	for _, f := range model.RequiredFields() {
		fmt.Printf("  - %s (%s): %s\n", f.Name, f.Type, f.Description)
	}

	// List optional fields
	fmt.Println("\nOptional fields:")
	for _, f := range model.OptionalFields() {
		fmt.Printf("  - %s (%s): %s\n", f.Name, f.Type, f.Description)
		if len(f.EnumVals) > 0 {
			fmt.Printf("    Allowed values: %v\n", f.EnumVals)
		}
		if f.Default != nil {
			fmt.Printf("    Default: %v\n", f.Default)
		}
	}

	// Also works with video models
	fmt.Println("\n--- Video Model Example ---")
	videoModel := kie.Seedance15Pro
	fmt.Printf("Model: %s\n", videoModel.Name)
	fmt.Printf("Required: %d fields, Optional: %d fields\n",
		len(videoModel.RequiredFields()),
		len(videoModel.OptionalFields()))
}

func editImage(ctx context.Context, client *kie.Client) error {
	// Image editing with Seedream 4.5
	result, err := kie.Seedream45Edit.Request().
		Prompt("Change the dress color to bright blue").
		ImageURLs("https://example.com/your-image.jpg").
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("edit: %w", err)
	}

	fmt.Printf("Edited image: %s\n", result.URLs[0])
	return nil
}

func upscaleImage(ctx context.Context, client *kie.Client) error {
	// Upscale an image using Topaz
	result, err := kie.TopazImageUpscale.Request().
		Set(kie.ParamImageURL, "https://example.com/low-res-image.jpg").
		Set("upscale_factor", kie.Upscale4x).
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("upscale: %w", err)
	}

	fmt.Printf("Upscaled image: %s\n", result.URLs[0])
	return nil
}

// uploadAndGenerate demonstrates the recommended workflow for image-to-X generation:
// 1. Upload local file to KIE temporary storage
// 2. Use the returned URL for model generation
//
// This is the best practice for local files because:
// - Avoids needing to host the image publicly
// - Uploaded files are valid for 15 days (enough time for generation)
// - Works seamlessly with all image-to-X and video-to-X models
func uploadAndGenerate(ctx context.Context, client *kie.Client) error {
	uploader := client.FileUploader()

	// Step 1: Upload a local image file
	fmt.Println("Step 1: Uploading local image...")
	uploadResult, err := uploader.UploadFile(ctx, "./my-photo.jpg", &kie.UploadOptions{
		UploadPath: "images",
		FileName:   "my-photo.jpg", // Optional: custom filename
	})
	if err != nil {
		return fmt.Errorf("upload: %w", err)
	}

	fmt.Printf("Uploaded: %s\n", uploadResult.FileURL)
	fmt.Printf("Expires: %s\n", uploadResult.ExpiresAt)

	// Step 2: Use the uploaded URL for image editing
	fmt.Println("\nStep 2: Generating edited image...")
	editResult, err := kie.IdeogramV3Remix.Request().
		Prompt("Transform into a watercolor painting style").
		ImageURL(uploadResult.FileURL). // Use uploaded file URL
		Style(kie.StyleRealistic).
		Strength(0.7).
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("remix: %w", err)
	}

	fmt.Printf("Remixed image: %s\n", editResult.URLs[0])

	// Alternative: Upload from URL (if you have a remote image)
	fmt.Println("\n--- Alternative: Upload from URL ---")
	urlUploadResult, err := uploader.UploadFromURL(ctx, "https://example.com/remote-image.jpg", nil)
	if err != nil {
		return fmt.Errorf("url upload: %w", err)
	}
	fmt.Printf("URL upload result: %s\n", urlUploadResult.FileURL)

	// Alternative: Upload base64 data (for in-memory images)
	/*
		base64Result, err := uploader.UploadBase64(ctx, "data:image/png;base64,iVBORw0K...", nil)
		if err != nil {
			return fmt.Errorf("base64 upload: %w", err)
		}
		fmt.Printf("Base64 upload result: %s\n", base64Result.FileURL)
	*/

	return nil
}
