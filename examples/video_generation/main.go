// Example: Video generation using KIE SDK
//
// This example demonstrates video generation with the declarative architecture:
// 1. Text to Video generation
// 2. Image to Video generation
// 3. Video to Video transformation
// 4. Using predefined constants
// 5. Audio-based video generation (avatar)
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
	apiKey := os.Getenv("KIE_API_KEY")
	if apiKey == "" {
		log.Fatal("KIE_API_KEY environment variable is required")
	}

	client := kie.NewClient(
		apiKey,
		kie.WithTimeout(30*time.Second),
		kie.WithMaxRetries(3),
	)
	defer client.Close()

	ctx := context.Background()

	// Example 1: Text to Video with Builder
	fmt.Println("=== Example 1: Text to Video ===")
	if err := textToVideo(ctx, client); err != nil {
		log.Printf("Text to video failed: %v", err)
	}

	// Example 2: Image to Video
	fmt.Println("\n=== Example 2: Image to Video ===")
	if err := imageToVideo(ctx, client); err != nil {
		log.Printf("Image to video failed: %v", err)
	}

	// Example 3: Using Constants
	fmt.Println("\n=== Example 3: Using Constants ===")
	if err := generateWithConstants(ctx, client); err != nil {
		log.Printf("Constants example failed: %v", err)
	}

	// Example 4: Video to Video
	fmt.Println("\n=== Example 4: Video to Video ===")
	if err := videoToVideo(ctx, client); err != nil {
		log.Printf("Video to video failed: %v", err)
	}

	// Example 5: Audio-based Video (Avatar)
	fmt.Println("\n=== Example 5: Avatar with Audio ===")
	if err := audioToVideo(ctx, client); err != nil {
		log.Printf("Avatar generation failed: %v", err)
	}

	// Example 6: Video Upscale
	fmt.Println("\n=== Example 6: Video Upscale ===")
	if err := videoUpscale(ctx, client); err != nil {
		log.Printf("Video upscale failed: %v", err)
	}

	// Example 7: Upload + Generate workflow (recommended for local files)
	fmt.Println("\n=== Example 7: Upload + Generate Workflow ===")
	if err := uploadAndGenerateVideo(ctx, client); err != nil {
		log.Printf("Upload+Generate example failed: %v", err)
	}
}

func textToVideo(ctx context.Context, client *kie.Client) error {
	// Basic text to video generation with constants
	result, err := kie.Seedance15Pro.Request().
		Prompt("A robot dancing in the rain, cinematic lighting").
		AspectRatio(kie.Ratio16x9). // Use constant
		Duration(kie.Duration5s).   // Use constant
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	fmt.Printf("Task completed: %s\n", result.TaskID)
	fmt.Printf("Video URL: %s\n", result.URLs[0])
	return nil
}

func imageToVideo(ctx context.Context, client *kie.Client) error {
	// Transform an image into a video with constants
	result, err := kie.Seedance15ImageToVideo.Request().
		Prompt("Make the character wave hello").
		ImageURLs("https://example.com/character.jpg").
		Duration(kie.Duration5s). // Use constant
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	fmt.Printf("Video URL: %s\n", result.URLs[0])
	return nil
}

func generateWithConstants(ctx context.Context, client *kie.Client) error {
	// Using predefined constants for type safety
	result, err := kie.Wan26TextToVideo.Request().
		Prompt("A cinematic shot of a rocket launching").
		Set(kie.ParamDuration, kie.Duration10Sec).
		Set(kie.ParamResolution, kie.VideoResolution1080p).
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	fmt.Printf("Generated with constants: %s\n", result.URLs[0])

	// Alternative with Params map
	result2, err := kie.Generate(ctx, client, kie.BytedanceSeedance15Pro, kie.Params{
		kie.ParamPrompt:      "A peaceful garden with butterflies",
		kie.ParamAspectRatio: kie.Ratio16x9,
		kie.ParamDuration:    kie.Duration8Sec,
		kie.ParamResolution:  kie.VideoResolution720p,
	})

	if err != nil {
		return fmt.Errorf("generate with params: %w", err)
	}

	fmt.Printf("Generated with params: %s\n", result2.URLs[0])
	return nil
}

func videoToVideo(ctx context.Context, client *kie.Client) error {
	// Video transformation using Wan 2.6 Video to Video
	result, err := kie.Wan26VideoToVideo.Request().
		Prompt("Transform the video into an anime style").
		Set(kie.ParamVideoURLs, []string{"https://example.com/input-video.mp4"}).
		Set(kie.ParamDuration, kie.Duration5Sec).
		Set(kie.ParamResolution, kie.VideoResolution720p).
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("transform: %w", err)
	}

	fmt.Printf("Transformed video: %s\n", result.URLs[0])

	// Another option: Wan 2.2 Animate Move (transfer motion from video to image)
	result2, err := kie.Wan22AnimateMove.Request().
		Set(kie.ParamVideoURL, "https://example.com/motion-video.mp4").
		Set(kie.ParamImageURL, "https://example.com/character.png").
		Set(kie.ParamResolution, kie.VideoResolution480p).
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("animate: %w", err)
	}

	fmt.Printf("Animated video: %s\n", result2.URLs[0])
	return nil
}

func audioToVideo(ctx context.Context, client *kie.Client) error {
	// Generate talking avatar from image and audio
	result, err := kie.KlingV1AvatarStandard.Request().
		Set(kie.ParamImageURL, "https://example.com/avatar.jpg").
		Set(kie.ParamAudioURL, "https://example.com/speech.mp3").
		Prompt("Professional speaking pose").
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("avatar: %w", err)
	}

	fmt.Printf("Avatar video: %s\n", result.URLs[0])

	// Alternative: Using Infinitalk
	result2, err := kie.InfinitalkFromAudio.Request().
		Set(kie.ParamImageURL, "https://example.com/person.jpg").
		Set(kie.ParamAudioURL, "https://example.com/narration.mp3").
		Prompt("Natural conversation").
		Set(kie.ParamResolution, kie.VideoResolution720p).
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("infinitalk: %w", err)
	}

	fmt.Printf("Infinitalk video: %s\n", result2.URLs[0])
	return nil
}

func videoUpscale(ctx context.Context, client *kie.Client) error {
	// Upscale video quality using Topaz
	result, err := kie.TopazVideoUpscale.Request().
		Set(kie.ParamVideoURL, "https://example.com/low-res-video.mp4").
		Set("upscale_factor", kie.Upscale2x).
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("upscale: %w", err)
	}

	fmt.Printf("Upscaled video: %s\n", result.URLs[0])
	return nil
}

// uploadAndGenerateVideo demonstrates the recommended workflow for image/video-to-video generation:
// 1. Upload local file to KIE temporary storage
// 2. Use the returned URL for video generation
//
// This is the best practice because:
// - Avoids needing to host files publicly
// - Uploaded files are valid for 15 days (enough for generation)
// - Works with all image-to-video and video-to-video models
func uploadAndGenerateVideo(ctx context.Context, client *kie.Client) error {
	uploader := client.FileUploader()

	// Step 1: Upload a local image for image-to-video
	fmt.Println("Step 1: Uploading local image...")
	imageUpload, err := uploader.UploadFile(ctx, "./character.jpg", &kie.UploadOptions{
		UploadPath: "video-inputs",
	})
	if err != nil {
		return fmt.Errorf("upload image: %w", err)
	}

	fmt.Printf("Image uploaded: %s\n", imageUpload.FileURL)
	fmt.Printf("Expires: %s\n", imageUpload.ExpiresAt)

	// Step 2: Generate video from uploaded image
	fmt.Println("\nStep 2: Generating video from image...")
	videoResult, err := kie.Seedance15ImageToVideo.Request().
		Prompt("Make the character smile and wave hello").
		ImageURLs(imageUpload.FileURL). // Use uploaded file URL
		Duration(kie.Duration5s).       // Use constant
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("generate video: %w", err)
	}

	fmt.Printf("Generated video: %s\n", videoResult.URLs[0])

	// Example: Upload video for video-to-video transformation
	fmt.Println("\n--- Video-to-Video workflow ---")

	// Upload a video file
	videoUpload, err := uploader.UploadFile(ctx, "./input-video.mp4", &kie.UploadOptions{
		UploadPath: "video-inputs",
	})
	if err != nil {
		return fmt.Errorf("upload video: %w", err)
	}

	fmt.Printf("Video uploaded: %s\n", videoUpload.FileURL)

	// Transform video style
	transformResult, err := kie.Wan26VideoToVideo.Request().
		Prompt("Transform into cinematic anime style").
		Set(kie.ParamVideoURLs, []string{videoUpload.FileURL}).
		Set(kie.ParamDuration, kie.Duration5Sec).
		Set(kie.ParamResolution, kie.VideoResolution720p).
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("transform video: %w", err)
	}

	fmt.Printf("Transformed video: %s\n", transformResult.URLs[0])

	// Example: Upload audio for avatar generation
	fmt.Println("\n--- Avatar with uploaded audio ---")

	audioUpload, err := uploader.UploadFile(ctx, "./speech.mp3", nil)
	if err != nil {
		return fmt.Errorf("upload audio: %w", err)
	}

	avatarResult, err := kie.KlingV1AvatarStandard.Request().
		Set(kie.ParamImageURL, imageUpload.FileURL). // Reuse uploaded image
		Set(kie.ParamAudioURL, audioUpload.FileURL). // Use uploaded audio
		Prompt("Professional presenter speaking").
		Generate(ctx, client)

	if err != nil {
		return fmt.Errorf("generate avatar: %w", err)
	}

	fmt.Printf("Avatar video: %s\n", avatarResult.URLs[0])

	return nil
}
