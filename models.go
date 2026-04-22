// Package kie provides the KIE AI SDK for Go.
//
// This package re-exports commonly used types and functions from subpackages
// for convenient single-import usage.
//
// Example:
//
//	import "github.com/lpg-it/kie-go"
//
//	client := kie.NewClient(apiKey)
//	result, err := kie.Image.GoogleImagen4.Request().
//	    Prompt("A sunset over the ocean").
//	    Generate(ctx, client)
package kie

import (
	"context"
	"encoding/json"

	"github.com/lpg-it/kie-go/image"
	"github.com/lpg-it/kie-go/model"
	"github.com/lpg-it/kie-go/video"
)

// ================================================================================
// Re-exported types from model package
// ================================================================================

// Model is re-exported from model package.
type Model = model.Model

// Params is re-exported from model package.
type Params = model.Params

// Result is re-exported from model package.
type Result = model.Result

// RequestBuilder is re-exported from model package.
type RequestBuilder = model.RequestBuilder

// Field is re-exported from model package.
type Field = model.Field

// Category is re-exported from model package.
type Category = model.Category

// Category constants
const (
	CategoryTextToImage  = model.CategoryTextToImage
	CategoryImageToImage = model.CategoryImageToImage
	CategoryImageEdit    = model.CategoryImageEdit
	CategoryUpscale      = model.CategoryUpscale
	CategoryTextToVideo  = model.CategoryTextToVideo
	CategoryImageToVideo = model.CategoryImageToVideo
	CategoryVideoToVideo = model.CategoryVideoToVideo
)

// ================================================================================
// Model Namespaces - Access models by category
// ================================================================================

// ImageModels provides access to all image generation models.
var Image = struct {
	GoogleImagen4                  *Model
	GoogleImagen4Fast              *Model
	GoogleImagen4Ultra             *Model
	GoogleNanoBanana               *Model
	GoogleNanoBananaEdit           *Model
	NanoBananaPro                  *Model
	GoogleNanoBanana2              *Model
	GrokImagineTextToImage         *Model
	GrokImagineImageToImage        *Model
	GrokImagineUpscale             *Model
	Seedream45TextToImage          *Model
	Seedream45Edit                 *Model
	BytedanceSeedreamV4TextToImage *Model
	BytedanceSeedreamV4Edit        *Model
	RecraftCrispUpscale            *Model
	RecraftRemoveBackground        *Model
	TopazImageUpscale              *Model
	GptImage15ImageToImage         *Model
	GptImage15TextToImage          *Model
	GptImage2ImageToImage          *Model
	GptImage2TextToImage           *Model
	ZImage                         *Model
	Flux2ProImageToImage           *Model
	Flux2FlexImageToImage          *Model
	Flux2FlexTextToImage           *Model
	Flux2ProTextToImage            *Model
	IdeogramV3Reframe              *Model
	IdeogramV3TextToImage          *Model
	IdeogramV3Edit                 *Model
	IdeogramV3Remix                *Model
	BytedanceSeedream              *Model
	Wan27Image                     *Model
	Wan27ImagePro                  *Model
	QwenImageToImage               *Model
	QwenTextToImage                *Model
	QwenImageEdit                  *Model
	Qwen2ImageEdit                 *Model
}{
	GoogleImagen4:                  image.GoogleImagen4,
	GoogleImagen4Fast:              image.GoogleImagen4Fast,
	GoogleImagen4Ultra:             image.GoogleImagen4Ultra,
	GoogleNanoBanana:               image.GoogleNanoBanana,
	GoogleNanoBananaEdit:           image.GoogleNanoBananaEdit,
	NanoBananaPro:                  image.NanoBananaPro,
	GoogleNanoBanana2:              image.GoogleNanoBanana2,
	GrokImagineTextToImage:         image.GrokImagineTextToImage,
	GrokImagineImageToImage:        image.GrokImagineImageToImage,
	GrokImagineUpscale:             image.GrokImagineUpscale,
	Seedream45TextToImage:          image.Seedream45TextToImage,
	Seedream45Edit:                 image.Seedream45Edit,
	BytedanceSeedreamV4TextToImage: image.BytedanceSeedreamV4TextToImage,
	BytedanceSeedreamV4Edit:        image.BytedanceSeedreamV4Edit,
	RecraftCrispUpscale:            image.RecraftCrispUpscale,
	RecraftRemoveBackground:        image.RecraftRemoveBackground,
	TopazImageUpscale:              image.TopazImageUpscale,
	GptImage15ImageToImage:         image.GptImage15ImageToImage,
	GptImage15TextToImage:          image.GptImage15TextToImage,
	GptImage2ImageToImage:          image.GptImage2ImageToImage,
	GptImage2TextToImage:           image.GptImage2TextToImage,
	ZImage:                         image.ZImage,
	Flux2ProImageToImage:           image.Flux2ProImageToImage,
	Flux2FlexImageToImage:          image.Flux2FlexImageToImage,
	Flux2FlexTextToImage:           image.Flux2FlexTextToImage,
	Flux2ProTextToImage:            image.Flux2ProTextToImage,
	IdeogramV3Reframe:              image.IdeogramV3Reframe,
	IdeogramV3TextToImage:          image.IdeogramV3TextToImage,
	IdeogramV3Edit:                 image.IdeogramV3Edit,
	IdeogramV3Remix:                image.IdeogramV3Remix,
	BytedanceSeedream:              image.BytedanceSeedream,
	Wan27Image:                     image.Wan27Image,
	Wan27ImagePro:                  image.Wan27ImagePro,
	QwenImageToImage:               image.QwenImageToImage,
	QwenTextToImage:                image.QwenTextToImage,
	QwenImageEdit:                  image.QwenImageEdit,
	Qwen2ImageEdit:                 image.Qwen2ImageEdit,
}

// VideoModels provides access to all video generation models.
var Video = struct {
	Seedance15Pro                  *Model
	Seedance15ImageToVideo         *Model
	KlingVideo                     *Model
	RunwayGen3                     *Model
	PikaVideo                      *Model
	GrokImagineImageToVideo        *Model
	GrokImagineTextToVideo         *Model
	Kling26MotionControl           *Model
	Kling30MotionControl           *Model
	BytedanceSeedance15Pro         *Model
	BytedanceSeedance20Fast        *Model
	BytedanceSeedance20            *Model
	Wan26TextToVideo               *Model
	Wan26ImageToVideo              *Model
	Wan26VideoToVideo              *Model
	Wan27TextToVideo               *Model
	Wan27ImageToVideo              *Model
	Wan27ReferenceToVideo          *Model
	Wan27VideoEdit                 *Model
	Kling26ImageToVideo            *Model
	Kling26TextToVideo             *Model
	BytedanceV1ProFastImageToVideo *Model
	Hailuo23ImageToVideoPro        *Model
	Hailuo23ImageToVideoStandard   *Model
	Sora2ProStoryboard             *Model
	Sora2ProTextToVideo            *Model
	Sora2ProImageToVideo           *Model
	Sora2Characters                *Model
	SoraWatermarkRemover           *Model
	Kling25TurboTextToVideoPro     *Model
	Kling25TurboImageToVideoPro    *Model
	Wan25ImageToVideo              *Model
	Wan25TextToVideo               *Model
	Wan22AnimateMove               *Model
	Wan22AnimateReplace            *Model
	TopazVideoUpscale              *Model
	InfinitalkFromAudio            *Model
	Wan22A14bSpeechToVideoTurbo    *Model
	KlingV1AvatarStandard          *Model
	KlingAiAvatarV1Pro             *Model
	Wan22A14bTextToVideoTurbo      *Model
	Wan22A14bImageToVideoTurbo     *Model
	KlingV21MasterImageToVideo     *Model
	KlingV21Pro                    *Model
	KlingV21Standard               *Model
	KlingV21MasterTextToVideo      *Model
	BytedanceV1ProImageToVideo     *Model
	BytedanceV1LiteImageToVideo    *Model
	BytedanceV1ProTextToVideo      *Model
	BytedanceV1LiteTextToVideo     *Model
	Hailuo02TextToVideoStandard    *Model
	Hailuo02ImageToVideoStandard   *Model
	Hailuo02ImageToVideoPro        *Model
	Hailuo02TextToVideoPro         *Model
	Sora2ImageToVideo              *Model
	Sora2TextToVideo               *Model
}{
	Seedance15Pro:                  video.Seedance15Pro,
	Seedance15ImageToVideo:         video.Seedance15ImageToVideo,
	KlingVideo:                     video.KlingVideo,
	RunwayGen3:                     video.RunwayGen3,
	PikaVideo:                      video.PikaVideo,
	GrokImagineImageToVideo:        video.GrokImagineImageToVideo,
	GrokImagineTextToVideo:         video.GrokImagineTextToVideo,
	Kling26MotionControl:           video.Kling26MotionControl,
	Kling30MotionControl:           video.Kling30MotionControl,
	BytedanceSeedance15Pro:         video.BytedanceSeedance15Pro,
	BytedanceSeedance20Fast:        video.BytedanceSeedance20Fast,
	BytedanceSeedance20:            video.BytedanceSeedance20,
	Wan26TextToVideo:               video.Wan26TextToVideo,
	Wan26ImageToVideo:              video.Wan26ImageToVideo,
	Wan26VideoToVideo:              video.Wan26VideoToVideo,
	Wan27TextToVideo:               video.Wan27TextToVideo,
	Wan27ImageToVideo:              video.Wan27ImageToVideo,
	Wan27ReferenceToVideo:          video.Wan27ReferenceToVideo,
	Wan27VideoEdit:                 video.Wan27VideoEdit,
	Kling26ImageToVideo:            video.Kling26ImageToVideo,
	Kling26TextToVideo:             video.Kling26TextToVideo,
	BytedanceV1ProFastImageToVideo: video.BytedanceV1ProFastImageToVideo,
	Hailuo23ImageToVideoPro:        video.Hailuo23ImageToVideoPro,
	Hailuo23ImageToVideoStandard:   video.Hailuo23ImageToVideoStandard,
	Sora2ProStoryboard:             video.Sora2ProStoryboard,
	Sora2ProTextToVideo:            video.Sora2ProTextToVideo,
	Sora2ProImageToVideo:           video.Sora2ProImageToVideo,
	Sora2Characters:                video.Sora2Characters,
	SoraWatermarkRemover:           video.SoraWatermarkRemover,
	Kling25TurboTextToVideoPro:     video.Kling25TurboTextToVideoPro,
	Kling25TurboImageToVideoPro:    video.Kling25TurboImageToVideoPro,
	Wan25ImageToVideo:              video.Wan25ImageToVideo,
	Wan25TextToVideo:               video.Wan25TextToVideo,
	Wan22AnimateMove:               video.Wan22AnimateMove,
	Wan22AnimateReplace:            video.Wan22AnimateReplace,
	TopazVideoUpscale:              video.TopazVideoUpscale,
	InfinitalkFromAudio:            video.InfinitalkFromAudio,
	Wan22A14bSpeechToVideoTurbo:    video.Wan22A14bSpeechToVideoTurbo,
	KlingV1AvatarStandard:          video.KlingV1AvatarStandard,
	KlingAiAvatarV1Pro:             video.KlingAiAvatarV1Pro,
	Wan22A14bTextToVideoTurbo:      video.Wan22A14bTextToVideoTurbo,
	Wan22A14bImageToVideoTurbo:     video.Wan22A14bImageToVideoTurbo,
	KlingV21MasterImageToVideo:     video.KlingV21MasterImageToVideo,
	KlingV21Pro:                    video.KlingV21Pro,
	KlingV21Standard:               video.KlingV21Standard,
	KlingV21MasterTextToVideo:      video.KlingV21MasterTextToVideo,
	BytedanceV1ProImageToVideo:     video.BytedanceV1ProImageToVideo,
	BytedanceV1LiteImageToVideo:    video.BytedanceV1LiteImageToVideo,
	BytedanceV1ProTextToVideo:      video.BytedanceV1ProTextToVideo,
	BytedanceV1LiteTextToVideo:     video.BytedanceV1LiteTextToVideo,
	Hailuo02TextToVideoStandard:    video.Hailuo02TextToVideoStandard,
	Hailuo02ImageToVideoStandard:   video.Hailuo02ImageToVideoStandard,
	Hailuo02ImageToVideoPro:        video.Hailuo02ImageToVideoPro,
	Hailuo02TextToVideoPro:         video.Hailuo02TextToVideoPro,
	Sora2ImageToVideo:              video.Sora2ImageToVideo,
	Sora2TextToVideo:               video.Sora2TextToVideo,
}

// ================================================================================
// Direct Model Access (shortcuts)
// ================================================================================

// Image Models - Direct access without namespace
var (
	GoogleImagen4                  = image.GoogleImagen4
	GoogleImagen4Fast              = image.GoogleImagen4Fast
	GoogleImagen4Ultra             = image.GoogleImagen4Ultra
	GoogleNanoBanana               = image.GoogleNanoBanana
	GoogleNanoBananaEdit           = image.GoogleNanoBananaEdit
	NanoBananaPro                  = image.NanoBananaPro
	GoogleNanoBanana2              = image.GoogleNanoBanana2
	GrokImagineTextToImage         = image.GrokImagineTextToImage
	GrokImagineImageToImage        = image.GrokImagineImageToImage
	GrokImagineUpscale             = image.GrokImagineUpscale
	Seedream45TextToImage          = image.Seedream45TextToImage
	Seedream45Edit                 = image.Seedream45Edit
	BytedanceSeedreamV4TextToImage = image.BytedanceSeedreamV4TextToImage
	BytedanceSeedreamV4Edit        = image.BytedanceSeedreamV4Edit
	RecraftCrispUpscale            = image.RecraftCrispUpscale
	RecraftRemoveBackground        = image.RecraftRemoveBackground
	TopazImageUpscale              = image.TopazImageUpscale
	GptImage15ImageToImage         = image.GptImage15ImageToImage
	GptImage15TextToImage          = image.GptImage15TextToImage
	GptImage2ImageToImage          = image.GptImage2ImageToImage
	GptImage2TextToImage           = image.GptImage2TextToImage
	ZImage                         = image.ZImage
	Flux2ProImageToImage           = image.Flux2ProImageToImage
	Flux2FlexImageToImage          = image.Flux2FlexImageToImage
	Flux2FlexTextToImage           = image.Flux2FlexTextToImage
	Flux2ProTextToImage            = image.Flux2ProTextToImage
	IdeogramV3Reframe              = image.IdeogramV3Reframe
	IdeogramV3TextToImage          = image.IdeogramV3TextToImage
	IdeogramV3Edit                 = image.IdeogramV3Edit
	IdeogramV3Remix                = image.IdeogramV3Remix
	BytedanceSeedream              = image.BytedanceSeedream
	Wan27Image                     = image.Wan27Image
	Wan27ImagePro                  = image.Wan27ImagePro
	QwenImageToImage               = image.QwenImageToImage
	QwenTextToImage                = image.QwenTextToImage
	QwenImageEdit                  = image.QwenImageEdit
	Qwen2ImageEdit                 = image.Qwen2ImageEdit
)

// Video Models - Direct access without namespace
var (
	Seedance15Pro                  = video.Seedance15Pro
	Seedance15ImageToVideo         = video.Seedance15ImageToVideo
	KlingVideo                     = video.KlingVideo
	RunwayGen3                     = video.RunwayGen3
	PikaVideo                      = video.PikaVideo
	GrokImagineImageToVideo        = video.GrokImagineImageToVideo
	GrokImagineTextToVideo         = video.GrokImagineTextToVideo
	Kling26MotionControl           = video.Kling26MotionControl
	Kling30MotionControl           = video.Kling30MotionControl
	BytedanceSeedance15Pro         = video.BytedanceSeedance15Pro
	BytedanceSeedance20Fast        = video.BytedanceSeedance20Fast
	BytedanceSeedance20            = video.BytedanceSeedance20
	Wan26TextToVideo               = video.Wan26TextToVideo
	Wan26ImageToVideo              = video.Wan26ImageToVideo
	Wan26VideoToVideo              = video.Wan26VideoToVideo
	Wan27TextToVideo               = video.Wan27TextToVideo
	Wan27ImageToVideo              = video.Wan27ImageToVideo
	Wan27ReferenceToVideo          = video.Wan27ReferenceToVideo
	Wan27VideoEdit                 = video.Wan27VideoEdit
	Kling26ImageToVideo            = video.Kling26ImageToVideo
	Kling26TextToVideo             = video.Kling26TextToVideo
	BytedanceV1ProFastImageToVideo = video.BytedanceV1ProFastImageToVideo
	Hailuo23ImageToVideoPro        = video.Hailuo23ImageToVideoPro
	Hailuo23ImageToVideoStandard   = video.Hailuo23ImageToVideoStandard
	Sora2ProStoryboard             = video.Sora2ProStoryboard
	Sora2ProTextToVideo            = video.Sora2ProTextToVideo
	Sora2ProImageToVideo           = video.Sora2ProImageToVideo
	Sora2Characters                = video.Sora2Characters
	SoraWatermarkRemover           = video.SoraWatermarkRemover
	Kling25TurboTextToVideoPro     = video.Kling25TurboTextToVideoPro
	Kling25TurboImageToVideoPro    = video.Kling25TurboImageToVideoPro
	Wan25ImageToVideo              = video.Wan25ImageToVideo
	Wan25TextToVideo               = video.Wan25TextToVideo
	Wan22AnimateMove               = video.Wan22AnimateMove
	Wan22AnimateReplace            = video.Wan22AnimateReplace
	TopazVideoUpscale              = video.TopazVideoUpscale
	InfinitalkFromAudio            = video.InfinitalkFromAudio
	Wan22A14bSpeechToVideoTurbo    = video.Wan22A14bSpeechToVideoTurbo
	KlingV1AvatarStandard          = video.KlingV1AvatarStandard
	KlingAiAvatarV1Pro             = video.KlingAiAvatarV1Pro
	Wan22A14bTextToVideoTurbo      = video.Wan22A14bTextToVideoTurbo
	Wan22A14bImageToVideoTurbo     = video.Wan22A14bImageToVideoTurbo
	KlingV21MasterImageToVideo     = video.KlingV21MasterImageToVideo
	KlingV21Pro                    = video.KlingV21Pro
	KlingV21Standard               = video.KlingV21Standard
	KlingV21MasterTextToVideo      = video.KlingV21MasterTextToVideo
	BytedanceV1ProImageToVideo     = video.BytedanceV1ProImageToVideo
	BytedanceV1LiteImageToVideo    = video.BytedanceV1LiteImageToVideo
	BytedanceV1ProTextToVideo      = video.BytedanceV1ProTextToVideo
	BytedanceV1LiteTextToVideo     = video.BytedanceV1LiteTextToVideo
	Hailuo02TextToVideoStandard    = video.Hailuo02TextToVideoStandard
	Hailuo02ImageToVideoStandard   = video.Hailuo02ImageToVideoStandard
	Hailuo02ImageToVideoPro        = video.Hailuo02ImageToVideoPro
	Hailuo02TextToVideoPro         = video.Hailuo02TextToVideoPro
	Sora2ImageToVideo              = video.Sora2ImageToVideo
	Sora2TextToVideo               = video.Sora2TextToVideo
)

// ================================================================================
// Model Registry Functions
// ================================================================================

// GetImageModel returns an image model by identifier.
func GetImageModel(id string) *Model {
	return image.Get(id)
}

// GetVideoModel returns a video model by identifier.
func GetVideoModel(id string) *Model {
	return video.Get(id)
}

// GetModel returns any model by identifier.
func GetModel(id string) *Model {
	if m := image.Get(id); m != nil {
		return m
	}
	return video.Get(id)
}

// AllImageModels returns all image models.
func AllImageModels() []*Model {
	return image.Models
}

// AllVideoModels returns all video models.
func AllVideoModels() []*Model {
	return video.Models
}

// AllModels returns all registered models.
func AllModels() []*Model {
	all := make([]*Model, 0, len(image.Models)+len(video.Models))
	all = append(all, image.Models...)
	all = append(all, video.Models...)
	return all
}

// ================================================================================
// Generate Functions
// ================================================================================

// GenerateOption configures generation behavior.
type GenerateOption = model.GenerateOption

// Generation option functions
var (
	WithGenTimeout         = model.WithGenTimeout
	WithGenPollInterval    = model.WithGenPollInterval
	WithGenMaxPollInterval = model.WithGenMaxPollInterval
	WithGenCallback        = model.WithGenCallback
)

// Generate executes a model with the given parameters.
func Generate(ctx context.Context, client *Client, m *Model, params Params, opts ...GenerateOption) (*Result, error) {
	// Validate parameters
	if err := m.Validate(params); err != nil {
		return nil, err
	}

	// Apply options
	cfg := model.DefaultGenerateConfig(m)
	model.ApplyOptions(cfg, opts...)
	timeout, pollInterval, maxPollInterval, callbackURL := cfg.GetConfig()

	// Create task
	task, err := client.CreateTask(ctx, &CreateTaskRequest{
		Model:       m.Identifier,
		Input:       params,
		CallbackURL: callbackURL,
	})
	if err != nil {
		return nil, err
	}

	// Async mode
	if callbackURL != "" {
		return &Result{TaskID: task.TaskID}, nil
	}

	// Wait for completion
	info, err := client.WaitForTask(ctx, task.TaskID,
		WithWaitTimeout(timeout),
		WithPollInterval(pollInterval),
		WithMaxPollInterval(maxPollInterval),
	)
	if err != nil {
		return nil, err
	}

	// Parse result
	result := &Result{
		TaskID:   info.TaskID,
		CostTime: info.CostTime,
	}

	if info.ResultJSON != "" {
		result.Raw = json.RawMessage(info.ResultJSON)
		var data struct {
			ResultURLs []string `json:"resultUrls"`
		}
		if err := json.Unmarshal([]byte(info.ResultJSON), &data); err == nil {
			result.URLs = data.ResultURLs
		}
	}

	return result, nil
}

// GenerateModel implements the model.Generator interface.
// This allows the Client to be used with RequestBuilder.Generate().
func (c *Client) GenerateModel(ctx context.Context, modelID string, input any, opts model.WaitOptions) (*Result, error) {
	// Create task
	task, err := c.CreateTask(ctx, &CreateTaskRequest{
		Model: modelID,
		Input: input,
	})
	if err != nil {
		return nil, err
	}

	// Wait for completion
	info, err := c.WaitForTask(ctx, task.TaskID,
		WithWaitTimeout(opts.Timeout),
		WithPollInterval(opts.PollInterval),
		WithMaxPollInterval(opts.MaxPollInterval),
	)
	if err != nil {
		return nil, err
	}

	// Parse result
	result := &Result{
		TaskID:   info.TaskID,
		CostTime: info.CostTime,
	}

	if info.ResultJSON != "" {
		result.Raw = json.RawMessage(info.ResultJSON)
		var data struct {
			ResultURLs []string `json:"resultUrls"`
		}
		if err := json.Unmarshal([]byte(info.ResultJSON), &data); err == nil {
			result.URLs = data.ResultURLs
		}
	}

	return result, nil
}
