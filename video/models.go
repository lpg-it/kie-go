// Package video provides video generation model definitions.
package video

import (
	"time"

	"github.com/lpg-it/kie-go/model"
)

// ================================================================================
// Video Models - All video generation models defined declaratively
// ================================================================================

// Seedance15Pro - High quality video generation
var Seedance15Pro = model.Define(
	"seedance/1.5-pro",
	"Seedance 1.5 Pro",
	model.CategoryTextToVideo,
	model.WithProvider("bytedance"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description for video generation"), model.MaxLen(5000)),
).Optional(
	model.Enum("aspect_ratio", []string{"16:9", "9:16", "1:1"}, model.Desc("Video aspect ratio"), model.Default("16:9")),
	model.Enum("duration", []string{"5s", "10s"}, model.Desc("Video duration"), model.Default("5s")),
	model.Str("negative_prompt", model.Desc("What to exclude from the video"), model.MaxLen(5000)),
)

// Seedance15ImageToVideo - Seedance 1.5 Image to Video
var Seedance15ImageToVideo = model.Define(
	"seedance/1.5-image-to-video",
	"Seedance 1.5 Image To Video",
	model.CategoryImageToVideo,
	model.WithProvider("bytedance"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description for video generation"), model.MaxLen(5000)),
	model.Strings("image_urls", model.Desc("Source images"), model.MaxItems(1)),
).Optional(
	model.Enum("duration", []string{"5s", "10s"}, model.Desc("Video duration"), model.Default("5s")),
	model.Str("negative_prompt", model.Desc("What to exclude from the video"), model.MaxLen(5000)),
)

// KlingVideo - Kling Video text to video
var KlingVideo = model.Define(
	"kling/video",
	"Kling Video",
	model.CategoryTextToVideo,
	model.WithProvider("kuaishou"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description for video generation"), model.MaxLen(5000)),
).Optional(
	model.Enum("aspect_ratio", []string{"16:9", "9:16", "1:1"}, model.Desc("Video aspect ratio"), model.Default("16:9")),
	model.Enum("duration", []string{"5s", "10s"}, model.Desc("Video duration"), model.Default("5s")),
)

// RunwayGen3 - Advanced video generation
var RunwayGen3 = model.Define(
	"runway/gen3",
	"Runway Gen3",
	model.CategoryTextToVideo,
	model.WithProvider("runway"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description for video generation"), model.MaxLen(5000)),
).Optional(
	model.Enum("aspect_ratio", []string{"16:9", "9:16", "1:1"}, model.Desc("Video aspect ratio"), model.Default("16:9")),
	model.Enum("duration", []string{"4s", "8s"}, model.Desc("Video duration"), model.Default("4s")),
)

// PikaVideo - Creative video generation
var PikaVideo = model.Define(
	"pika/video",
	"Pika Video",
	model.CategoryTextToVideo,
	model.WithProvider("pika"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description for video generation"), model.MaxLen(5000)),
).Optional(
	model.Enum("aspect_ratio", []string{"16:9", "9:16", "1:1"}, model.Desc("Video aspect ratio"), model.Default("16:9")),
)

// GrokImagineImageToVideo - Grok Imagine Image to Video
var GrokImagineImageToVideo = model.Define(
	"grok-imagine/image-to-video",
	"Grok Imagine Image To Video",
	model.CategoryImageToVideo,
	model.WithProvider("xai"),
	model.WithTimeout(20*time.Minute),
).Optional(
	model.Strings("image_urls", model.Desc("External image URL as reference for video generation (only one image supported)"), model.MaxItems(1)),
	model.Str("task_id", model.Desc("Task ID of a previously generated Grok image on Kie"), model.MaxLen(100)),
	model.Int("index", model.Desc("When using task_id, specify which image to use (0-5, Grok generates 6 images per task)"), model.Min(0), model.Max(5), model.Default(0)),
	model.Str("prompt", model.Desc("Text prompt describing the desired video motion"), model.MaxLen(5000)),
	model.Enum("mode", []string{"fun", "normal", "spicy"}, model.Desc("Video generation mode (Spicy mode not supported with external image inputs)"), model.Default("normal")),
)

// GrokImagineTextToVideo - Grok Imagine Text to Video
var GrokImagineTextToVideo = model.Define(
	"grok-imagine/text-to-video",
	"Grok Imagine Text To Video",
	model.CategoryTextToVideo,
	model.WithProvider("xai"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the desired video"), model.MaxLen(5000)),
).Optional(
	model.Enum("aspect_ratio", []string{"2:3", "3:2", "1:1", "9:16", "16:9"}, model.Desc("Aspect ratio of the generated video"), model.Default("2:3")),
	model.Enum("mode", []string{"fun", "normal", "spicy"}, model.Desc("Video generation mode"), model.Default("normal")),
)

// Kling26MotionControl - Kling 2.6 Motion Control video generation
var Kling26MotionControl = model.Define(
	"kling-2.6/motion-control",
	"Kling 2.6 Motion Control",
	model.CategoryImageToVideo,
	model.WithProvider("kuaishou"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Strings("input_urls", model.Desc("Image URLs showing the subject's head, shoulders, and torso"), model.MaxItems(1)),
	model.Strings("video_urls", model.Desc("Video URLs for motion reference (3-30 seconds, min 720p)"), model.MaxItems(1)),
	model.Enum("character_orientation", []string{"image", "video"}, model.Desc("Character orientation: 'image' for same as picture (max 10s), 'video' for same as video (max 30s)"), model.Default("video")),
	model.Enum("mode", []string{"720p", "1080p"}, model.Desc("Output resolution mode"), model.Default("720p")),
).Optional(
	model.Str("prompt", model.Desc("Text description of the desired output"), model.MaxLen(2500)),
)

// Kling30MotionControl - Kling 3.0 Motion Control video generation
var Kling30MotionControl = model.Define(
	"kling-3.0/motion-control",
	"Kling 3.0 Motion Control",
	model.CategoryImageToVideo,
	model.WithProvider("kling"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Strings("input_urls", model.Desc("Image URLs showing the subject's head, shoulders, and torso"), model.MaxItems(1)),
	model.Strings("video_urls", model.Desc("Video URLs for motion reference (3-30 seconds, max 100MB)"), model.MaxItems(1)),
).Optional(
	model.Str("prompt", model.Desc("Optional text prompt to guide the animation"), model.MaxLen(2500)),
	model.Enum("mode", []string{"720p", "1080p", "std", "pro"}, model.Desc("Video quality mode; docs mention std/pro while examples use 720p/1080p"), model.Default("720p")),
	model.Enum("character_orientation", []string{"video", "image"}, model.Desc("Reference source for character orientation"), model.Default("video")),
	model.Enum("background_source", []string{"input_video", "input_image"}, model.Desc("Background source for the generated video"), model.Default("input_video")),
)

// BytedanceSeedance15Pro - Bytedance Seedance 1.5 Pro video generation with audio support
var BytedanceSeedance15Pro = model.Define(
	"bytedance/seedance-1.5-pro",
	"Bytedance Seedance 1.5 Pro",
	model.CategoryTextToVideo,
	model.WithProvider("bytedance"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description for video generation (3-2500 characters)"), model.MaxLen(2500)),
	model.Enum("aspect_ratio", []string{"1:1", "21:9", "4:3", "3:4", "16:9", "9:16"}, model.Desc("Video aspect ratio"), model.Default("1:1")),
	model.Enum("duration", []string{"4", "8", "12"}, model.Desc("Video duration in seconds"), model.Default("8")),
).Optional(
	model.Strings("input_urls", model.Desc("Image URLs for reference (0-2 images, max 10MB each)"), model.MaxItems(2)),
	model.Enum("resolution", []string{"480p", "720p"}, model.Desc("Video resolution (Standard 480p or High 720p)"), model.Default("720p")),
	model.Bool("fixed_lens", model.Desc("Keep camera view static and stable")),
	model.Bool("generate_audio", model.Desc("Create sound effects for the video (additional cost)")),
)

// BytedanceSeedance20Fast - Bytedance Seedance 2 Fast video generation
var BytedanceSeedance20Fast = model.Define(
	"bytedance/seedance-2-fast",
	"Bytedance Seedance 2 Fast",
	model.CategoryTextToVideo,
	model.WithProvider("bytedance"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt used to generate the video"), model.MaxLen(20000)),
	model.Bool("web_search", model.Desc("Whether to use online search")),
).Optional(
	model.Str("first_frame_url", model.Desc("First frame image URL or asset://assetId"), model.MaxLen(2048)),
	model.Str("last_frame_url", model.Desc("Last frame image URL or asset://assetId"), model.MaxLen(2048)),
	model.Strings("reference_image_urls", model.Desc("Reference image URLs or asset://assetId values"), model.MaxItems(9)),
	model.Strings("reference_video_urls", model.Desc("Reference video URLs or asset://assetId values"), model.MaxItems(3)),
	model.Strings("reference_audio_urls", model.Desc("Reference audio URLs or asset://assetId values"), model.MaxItems(3)),
	model.Bool("return_last_frame", model.Desc("Whether to return the last frame as an image"), model.Default(false)),
	model.Bool("generate_audio", model.Desc("Whether to generate audio for the video"), model.Default(true)),
	model.Enum("resolution", []string{"480p", "720p"}, model.Desc("Video resolution"), model.Default("720p")),
	model.Enum("aspect_ratio", []string{"1:1", "4:3", "3:4", "16:9", "9:16", "21:9", "adaptive"}, model.Desc("Aspect ratio of the generated video"), model.Default("16:9")),
	model.Int("duration", model.Desc("Video duration in seconds"), model.Min(4), model.Max(15), model.Default(5)),
	model.Bool("nsfw_checker", model.Desc("Disable platform-side NSFW filtering when set to false"), model.Default(false)),
)

// BytedanceSeedance20 - Bytedance Seedance 2 video generation
var BytedanceSeedance20 = model.Define(
	"bytedance/seedance-2",
	"Bytedance Seedance 2",
	model.CategoryTextToVideo,
	model.WithProvider("bytedance"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt used to generate the video"), model.MaxLen(20000)),
	model.Bool("web_search", model.Desc("Whether to use online search")),
).Optional(
	model.Str("first_frame_url", model.Desc("First frame image URL or asset://assetId"), model.MaxLen(2048)),
	model.Str("last_frame_url", model.Desc("Last frame image URL or asset://assetId"), model.MaxLen(2048)),
	model.Strings("reference_image_urls", model.Desc("Reference image URLs or asset://assetId values"), model.MaxItems(9)),
	model.Strings("reference_video_urls", model.Desc("Reference video URLs or asset://assetId values"), model.MaxItems(3)),
	model.Strings("reference_audio_urls", model.Desc("Reference audio URLs or asset://assetId values"), model.MaxItems(3)),
	model.Bool("return_last_frame", model.Desc("Whether to return the last frame as an image"), model.Default(false)),
	model.Bool("generate_audio", model.Desc("Whether to generate audio for the video"), model.Default(true)),
	model.Enum("resolution", []string{"480p", "720p", "1080p"}, model.Desc("Video resolution"), model.Default("720p")),
	model.Enum("aspect_ratio", []string{"1:1", "4:3", "3:4", "16:9", "9:16", "21:9", "adaptive"}, model.Desc("Aspect ratio of the generated video"), model.Default("16:9")),
	model.Int("duration", model.Desc("Video duration in seconds"), model.Min(4), model.Max(15), model.Default(5)),
	model.Bool("nsfw_checker", model.Desc("Disable platform-side NSFW filtering when set to false"), model.Default(false)),
)

// Wan26TextToVideo - Wan 2.6 Text to Video generation with multi-shot support
var Wan26TextToVideo = model.Define(
	"wan/2-6-text-to-video",
	"Wan 2.6 Text To Video",
	model.CategoryTextToVideo,
	model.WithProvider("wan"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description for video generation (supports Chinese and English)"), model.MaxLen(5000)),
).Optional(
	model.Enum("duration", []string{"5", "10", "15"}, model.Desc("Video duration in seconds"), model.Default("5")),
	model.Enum("resolution", []string{"720p", "1080p"}, model.Desc("Video resolution"), model.Default("1080p")),
	model.Bool("multi_shots", model.Desc("Enable multi-shot composition with transitions instead of single continuous shot")),
)

// Wan26ImageToVideo - Wan 2.6 Image to Video generation with multi-shot support
var Wan26ImageToVideo = model.Define(
	"wan/2-6-image-to-video",
	"Wan 2.6 Image To Video",
	model.CategoryImageToVideo,
	model.WithProvider("wan"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description for video generation (supports Chinese and English)"), model.MaxLen(5000)),
	model.Strings("image_urls", model.Desc("Image URLs (min 256x256px, max 10MB each)"), model.MaxItems(1)),
).Optional(
	model.Enum("duration", []string{"5", "10", "15"}, model.Desc("Video duration in seconds"), model.Default("5")),
	model.Enum("resolution", []string{"720p", "1080p"}, model.Desc("Video resolution"), model.Default("1080p")),
	model.Bool("multi_shots", model.Desc("Enable multi-shot composition with transitions instead of single continuous shot")),
)

// Wan26VideoToVideo - Wan 2.6 Video to Video transformation with multi-shot support
var Wan26VideoToVideo = model.Define(
	"wan/2-6-video-to-video",
	"Wan 2.6 Video To Video",
	model.CategoryVideoToVideo,
	model.WithProvider("wan"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description for video generation (supports Chinese and English)"), model.MaxLen(5000)),
	model.Strings("video_urls", model.Desc("Video URLs for transformation (max 10MB, mp4/mov/mkv)"), model.MaxItems(1)),
).Optional(
	model.Enum("duration", []string{"5", "10"}, model.Desc("Video duration in seconds"), model.Default("5")),
	model.Enum("resolution", []string{"720p", "1080p"}, model.Desc("Video resolution"), model.Default("1080p")),
	model.Bool("multi_shots", model.Desc("Enable multi-shot composition with transitions instead of single continuous shot")),
)

// Wan27TextToVideo - Wan 2.7 Text To Video generation
var Wan27TextToVideo = model.Define(
	"wan/2-7-text-to-video",
	"Wan 2.7 Text To Video",
	model.CategoryTextToVideo,
	model.WithProvider("wan"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the desired video"), model.MaxLen(5000)),
).Optional(
	model.Str("negative_prompt", model.Desc("What to exclude from the video"), model.MaxLen(500)),
	model.Str("audio_url", model.Desc("Optional custom audio URL"), model.MaxLen(2048)),
	model.Enum("resolution", []string{"720p", "1080p"}, model.Desc("Video resolution"), model.Default("1080p")),
	model.Enum("ratio", []string{"16:9", "9:16", "1:1", "4:3", "3:4"}, model.Desc("Aspect ratio of the generated video"), model.Default("16:9")),
	model.Int("duration", model.Desc("Video duration in seconds"), model.Min(2), model.Max(15), model.Default(5)),
	model.Bool("prompt_extend", model.Desc("Enable intelligent prompt rewriting"), model.Default(true)),
	model.Bool("watermark", model.Desc("Whether to add an AI-generated watermark"), model.Default(false)),
	model.Int("seed", model.Desc("Random seed for reproducibility"), model.Min(0), model.Max(2147483647)),
	model.Bool("nsfw_checker", model.Desc("Disable platform-side NSFW filtering when set to false"), model.Default(false)),
)

// Wan27ImageToVideo - Wan 2.7 Image To Video generation
var Wan27ImageToVideo = model.Define(
	"wan/2-7-image-to-video",
	"Wan 2.7 Image To Video",
	model.CategoryImageToVideo,
	model.WithProvider("wan"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the desired video"), model.MaxLen(5000)),
).Optional(
	model.Str("negative_prompt", model.Desc("What to exclude from the video"), model.MaxLen(500)),
	model.Str("first_frame_url", model.Desc("First frame image URL"), model.MaxLen(2048)),
	model.Str("last_frame_url", model.Desc("Last frame image URL"), model.MaxLen(2048)),
	model.Str("first_clip_url", model.Desc("First clip video URL for continuation mode"), model.MaxLen(2048)),
	model.Str("driving_audio_url", model.Desc("Driving audio URL"), model.MaxLen(2048)),
	model.Enum("resolution", []string{"720p", "1080p"}, model.Desc("Video resolution"), model.Default("1080p")),
	model.Int("duration", model.Desc("Total output video duration in seconds"), model.Min(2), model.Max(15), model.Default(5)),
	model.Bool("prompt_extend", model.Desc("Enable intelligent prompt rewriting"), model.Default(true)),
	model.Bool("watermark", model.Desc("Whether to add an AI-generated watermark"), model.Default(false)),
	model.Int("seed", model.Desc("Random seed for reproducibility"), model.Min(0), model.Max(2147483647)),
	model.Bool("nsfw_checker", model.Desc("Disable platform-side NSFW filtering when set to false"), model.Default(false)),
)

// Wan27ReferenceToVideo - Wan 2.7 Reference To Video generation
var Wan27ReferenceToVideo = model.Define(
	"wan/2-7-r2v",
	"Wan 2.7 Reference To Video",
	model.CategoryTextToVideo,
	model.WithProvider("wan"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the desired video"), model.MaxLen(5000)),
).Optional(
	model.Str("negative_prompt", model.Desc("What to exclude from the video"), model.MaxLen(500)),
	model.Strings("reference_image", model.Desc("Reference image URLs; provide at least one reference image or video"), model.MaxItems(5)),
	model.Strings("reference_video", model.Desc("Reference video URLs; provide at least one reference image or video"), model.MaxItems(5)),
	model.Str("first_frame", model.Desc("First frame image URL"), model.MaxLen(2048)),
	model.Str("reference_voice", model.Desc("Reference audio URL used to control voice timbre"), model.MaxLen(2048)),
	model.Enum("resolution", []string{"720p", "1080p"}, model.Desc("Video resolution"), model.Default("1080p")),
	model.Enum("aspect_ratio", []string{"16:9", "9:16", "1:1", "4:3", "3:4"}, model.Desc("Aspect ratio of the generated video"), model.Default("16:9")),
	model.Int("duration", model.Desc("Output video duration in seconds"), model.Min(2), model.Max(10), model.Default(5)),
	model.Bool("prompt_extend", model.Desc("Enable intelligent prompt rewriting"), model.Default(true)),
	model.Bool("watermark", model.Desc("Whether to add an AI-generated watermark"), model.Default(false)),
	model.Int("seed", model.Desc("Random seed for reproducibility"), model.Min(0), model.Max(2147483647)),
	model.Bool("nsfw_checker", model.Desc("Disable platform-side NSFW filtering when set to false"), model.Default(false)),
)

// Wan27VideoEdit - Wan 2.7 Video Edit generation
var Wan27VideoEdit = model.Define(
	"wan/2-7-videoedit",
	"Wan 2.7 Video Edit",
	model.CategoryVideoToVideo,
	model.WithProvider("wan"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("video_url", model.Desc("Source video URL to edit"), model.MaxLen(2048)),
).Optional(
	model.Str("prompt", model.Desc("Text prompt describing the desired edit"), model.MaxLen(5000)),
	model.Str("negative_prompt", model.Desc("What to exclude from the edited video"), model.MaxLen(500)),
	model.Str("reference_image", model.Desc("Reference image URL for character, clothing, or style guidance"), model.MaxLen(2048)),
	model.Enum("resolution", []string{"720p", "1080p"}, model.Desc("Output video resolution"), model.Default("1080p")),
	model.Enum("aspect_ratio", []string{"16:9", "9:16", "1:1", "4:3", "3:4"}, model.Desc("Output video aspect ratio")),
	model.Int("duration", model.Desc("Output video duration in seconds; 0 keeps the full input length"), model.Min(0), model.Max(10), model.Default(0)),
	model.Enum("audio_setting", []string{"auto", "origin"}, model.Desc("Audio handling mode"), model.Default("auto")),
	model.Bool("prompt_extend", model.Desc("Enable intelligent prompt rewriting"), model.Default(true)),
	model.Bool("watermark", model.Desc("Whether to add an AI-generated watermark"), model.Default(false)),
	model.Int("seed", model.Desc("Random seed for reproducibility"), model.Min(0), model.Max(2147483647)),
	model.Bool("nsfw_checker", model.Desc("Disable platform-side NSFW filtering when set to false"), model.Default(false)),
)

// Kling26ImageToVideo - Kling 2.6 Image to Video generation with sound support
var Kling26ImageToVideo = model.Define(
	"kling-2.6/image-to-video",
	"Kling 2.6 Image To Video",
	model.CategoryImageToVideo,
	model.WithProvider("kuaishou"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt to generate the video"), model.MaxLen(2500)),
	model.Strings("image_urls", model.Desc("Image URLs to generate video from (max 10MB each)"), model.MaxItems(1)),
	model.Bool("sound", model.Desc("Whether the generated video contains sound")),
	model.Enum("duration", []string{"5", "10"}, model.Desc("Duration of the video in seconds"), model.Default("5")),
)

// Kling26TextToVideo - Kling 2.6 Text to Video generation with sound support
var Kling26TextToVideo = model.Define(
	"kling-2.6/text-to-video",
	"Kling 2.6 Text To Video",
	model.CategoryTextToVideo,
	model.WithProvider("kuaishou"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt to generate the video"), model.MaxLen(2500)),
	model.Bool("sound", model.Desc("Whether the generated video contains sound")),
	model.Enum("aspect_ratio", []string{"1:1", "16:9", "9:16"}, model.Desc("Aspect ratio of the video"), model.Default("1:1")),
	model.Enum("duration", []string{"5", "10"}, model.Desc("Duration of the video in seconds"), model.Default("5")),
)

// BytedanceV1ProFastImageToVideo - Bytedance V1 Pro Fast Image to Video generation
var BytedanceV1ProFastImageToVideo = model.Define(
	"bytedance/v1-pro-fast-image-to-video",
	"Bytedance V1 Pro Fast Image To Video",
	model.CategoryImageToVideo,
	model.WithProvider("bytedance"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt to generate the video"), model.MaxLen(10000)),
	model.Str("image_url", model.Desc("URL of the image to generate video from (max 10MB)"), model.MaxLen(2048)),
).Optional(
	model.Enum("resolution", []string{"720p", "1080p"}, model.Desc("Video resolution"), model.Default("720p")),
	model.Enum("duration", []string{"5", "10"}, model.Desc("Duration of the video in seconds"), model.Default("5")),
)

// ================================================================================
// Model Registry
// ================================================================================

// Models contains all video generation models.
var Models = []*model.Model{
	Seedance15Pro,
	Seedance15ImageToVideo,
	KlingVideo,
	RunwayGen3,
	PikaVideo,
	GrokImagineImageToVideo,
	GrokImagineTextToVideo,
	Kling26MotionControl,
	Kling30MotionControl,
	BytedanceSeedance15Pro,
	BytedanceSeedance20Fast,
	BytedanceSeedance20,
	Wan26TextToVideo,
	Wan26ImageToVideo,
	Wan26VideoToVideo,
	Wan27TextToVideo,
	Wan27ImageToVideo,
	Wan27ReferenceToVideo,
	Wan27VideoEdit,
	Kling26ImageToVideo,
	Kling26TextToVideo,
	BytedanceV1ProFastImageToVideo,
	Hailuo23ImageToVideoPro,
	Hailuo23ImageToVideoStandard,
	Sora2ProStoryboard,
	Sora2ProTextToVideo,
	Sora2ProImageToVideo,
	Sora2Characters,
	SoraWatermarkRemover,
	Kling25TurboTextToVideoPro,
	Kling25TurboImageToVideoPro,
	Wan25ImageToVideo,
	Wan25TextToVideo,
	Wan22AnimateMove,
	Wan22AnimateReplace,
	TopazVideoUpscale,
	InfinitalkFromAudio,
	Wan22A14bSpeechToVideoTurbo,
	KlingV1AvatarStandard,
	KlingAiAvatarV1Pro,
	Wan22A14bTextToVideoTurbo,
	Wan22A14bImageToVideoTurbo,
	KlingV21MasterImageToVideo,
	KlingV21Pro,
	KlingV21Standard,
	KlingV21MasterTextToVideo,
	BytedanceV1ProImageToVideo,
	BytedanceV1LiteImageToVideo,
	BytedanceV1ProTextToVideo,
	BytedanceV1LiteTextToVideo,
	Hailuo02TextToVideoStandard,
	Hailuo02ImageToVideoStandard,
	Hailuo02ImageToVideoPro,
	Hailuo02TextToVideoPro,
	Sora2ImageToVideo,
	Sora2TextToVideo,
}

// Hailuo23ImageToVideoPro - Hailuo 2.3 Image to Video Pro generation
var Hailuo23ImageToVideoPro = model.Define(
	"hailuo/2-3-image-to-video-pro",
	"Hailuo 2.3 Image To Video Pro",
	model.CategoryImageToVideo,
	model.WithProvider("hailuo"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the desired video animation"), model.MaxLen(5000)),
	model.Str("image_url", model.Desc("URL of the image to animate (max 10MB)"), model.MaxLen(2048)),
).Optional(
	model.Enum("duration", []string{"6", "10"}, model.Desc("Duration of the video in seconds (10s not supported for 1080P)"), model.Default("6")),
	model.Enum("resolution", []string{"768P", "1080P"}, model.Desc("Resolution of the generated video"), model.Default("768P")),
)

// Hailuo23ImageToVideoStandard - Hailuo 2.3 Image to Video Standard generation
var Hailuo23ImageToVideoStandard = model.Define(
	"hailuo/2-3-image-to-video-standard",
	"Hailuo 2.3 Image To Video Standard",
	model.CategoryImageToVideo,
	model.WithProvider("hailuo"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the desired video animation"), model.MaxLen(5000)),
	model.Str("image_url", model.Desc("URL of the image to animate (max 10MB)"), model.MaxLen(2048)),
).Optional(
	model.Enum("duration", []string{"6", "10"}, model.Desc("Duration of the video in seconds (10s not supported for 1080P)"), model.Default("6")),
	model.Enum("resolution", []string{"768P", "1080P"}, model.Desc("Resolution of the generated video"), model.Default("768P")),
)

// Hailuo02TextToVideoStandard - Hailuo 02 Text to Video Standard generation
var Hailuo02TextToVideoStandard = model.Define(
	"hailuo/02-text-to-video-standard",
	"Hailuo 02 Text To Video Standard",
	model.CategoryTextToVideo,
	model.WithProvider("hailuo"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description for video generation"), model.MaxLen(1500)),
).Optional(
	model.Enum("duration", []string{"6", "10"}, model.Desc("Duration of the video in seconds (10s not supported for 1080p)"), model.Default("6")),
	model.Bool("prompt_optimizer", model.Desc("Whether to use the model's prompt optimizer"), model.Default(true)),
)

// Hailuo02TextToVideoPro - Hailuo 02 Text to Video Pro generation
var Hailuo02TextToVideoPro = model.Define(
	"hailuo/02-text-to-video-pro",
	"Hailuo 02 Text To Video Pro",
	model.CategoryTextToVideo,
	model.WithProvider("hailuo"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt for video generation"), model.MaxLen(1500)),
).Optional(
	model.Bool("prompt_optimizer", model.Desc("Whether to use the model's prompt optimizer"), model.Default(true)),
)

// Hailuo02ImageToVideoStandard - Hailuo 02 Image to Video Standard generation
var Hailuo02ImageToVideoStandard = model.Define(
	"hailuo/02-image-to-video-standard",
	"Hailuo 02 Image To Video Standard",
	model.CategoryImageToVideo,
	model.WithProvider("hailuo"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the video to generate"), model.MaxLen(1500)),
	model.Str("image_url", model.Desc("URL of the image to use as the first frame"), model.MaxLen(2048)),
).Optional(
	model.Str("end_image_url", model.Desc("URL of the image to use as the last frame"), model.MaxLen(2048)),
	model.Enum("duration", []string{"6", "10"}, model.Desc("Duration of the video in seconds (10s not supported for 1080p)"), model.Default("10")),
	model.Enum("resolution", []string{"512P", "768P"}, model.Desc("Resolution of the generated video"), model.Default("768P")),
	model.Bool("prompt_optimizer", model.Desc("Whether to use the model's prompt optimizer"), model.Default(true)),
)

// Hailuo02ImageToVideoPro - Hailuo 02 Image to Video Pro generation
var Hailuo02ImageToVideoPro = model.Define(
	"hailuo/02-image-to-video-pro",
	"Hailuo 02 Image To Video Pro",
	model.CategoryImageToVideo,
	model.WithProvider("hailuo"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the desired video animation"), model.MaxLen(1500)),
	model.Str("image_url", model.Desc("URL of the image to animate"), model.MaxLen(2048)),
).Optional(
	model.Str("end_image_url", model.Desc("URL of the image to use as the last frame"), model.MaxLen(2048)),
	model.Bool("prompt_optimizer", model.Desc("Whether to use the model's prompt optimizer"), model.Default(true)),
)

// Sora2ProStoryboard - Sora 2 Pro Storyboard video generation
var Sora2ProStoryboard = model.Define(
	"sora-2-pro-storyboard",
	"Sora 2 Pro Storyboard",
	model.CategoryImageToVideo,
	model.WithProvider("openai"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Enum("n_frames", []string{"10", "15", "25"}, model.Desc("Total length of the video in seconds"), model.Default("15")),
).Optional(
	model.Strings("image_urls", model.Desc("Image URLs to use as input for the storyboard"), model.MaxItems(10)),
	model.Enum("aspect_ratio", []string{"portrait", "landscape"}, model.Desc("Aspect ratio of the video"), model.Default("landscape")),
)

// Sora2ProTextToVideo - Sora 2 Pro Text to Video generation
var Sora2ProTextToVideo = model.Define(
	"sora-2-pro-text-to-video",
	"Sora 2 Pro Text To Video",
	model.CategoryTextToVideo,
	model.WithProvider("openai"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the desired video motion"), model.MaxLen(10000)),
).Optional(
	model.Enum("aspect_ratio", []string{"portrait", "landscape"}, model.Desc("Aspect ratio of the video"), model.Default("landscape")),
	model.Enum("n_frames", []string{"10", "15"}, model.Desc("Length of the video in seconds"), model.Default("10")),
	model.Enum("size", []string{"standard", "high"}, model.Desc("Quality/size of the generated video"), model.Default("high")),
	model.Bool("remove_watermark", model.Desc("Remove watermarks from the generated video"), model.Default(true)),
)

// Sora2ProImageToVideo - Sora 2 Pro Image to Video generation
var Sora2ProImageToVideo = model.Define(
	"sora-2-pro-image-to-video",
	"Sora 2 Pro Image To Video",
	model.CategoryImageToVideo,
	model.WithProvider("openai"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the desired video motion"), model.MaxLen(10000)),
	model.Strings("image_urls", model.Desc("URLs of images to use as first frames (publicly accessible)"), model.MaxItems(10)),
).Optional(
	model.Enum("aspect_ratio", []string{"portrait", "landscape"}, model.Desc("Aspect ratio of the video"), model.Default("landscape")),
	model.Enum("n_frames", []string{"10", "15"}, model.Desc("Length of the video in seconds"), model.Default("10")),
	model.Enum("size", []string{"standard", "high"}, model.Desc("Quality/size of the generated video"), model.Default("standard")),
	model.Bool("remove_watermark", model.Desc("Remove watermarks from the generated video"), model.Default(true)),
)

// Sora2ImageToVideo - Sora 2 Image to Video generation
var Sora2ImageToVideo = model.Define(
	"sora-2-image-to-video",
	"Sora 2 Image To Video",
	model.CategoryImageToVideo,
	model.WithProvider("openai"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the desired video motion"), model.MaxLen(10000)),
	model.Strings("image_urls", model.Desc("URLs of images to use as first frames (publicly accessible)"), model.MaxItems(10)),
).Optional(
	model.Enum("aspect_ratio", []string{"portrait", "landscape"}, model.Desc("Aspect ratio of the video"), model.Default("landscape")),
	model.Enum("n_frames", []string{"10", "15"}, model.Desc("Length of the video in seconds"), model.Default("10")),
	model.Bool("remove_watermark", model.Desc("Remove watermarks from the generated video"), model.Default(true)),
)

// Sora2TextToVideo - Sora 2 Text to Video generation
var Sora2TextToVideo = model.Define(
	"sora-2-text-to-video",
	"Sora 2 Text To Video",
	model.CategoryTextToVideo,
	model.WithProvider("openai"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the desired video"), model.MaxLen(10000)),
).Optional(
	model.Enum("aspect_ratio", []string{"portrait", "landscape"}, model.Desc("Aspect ratio of the video"), model.Default("landscape")),
	model.Enum("n_frames", []string{"10", "15"}, model.Desc("Length of the video in seconds"), model.Default("10")),
	model.Bool("remove_watermark", model.Desc("Remove watermarks from the generated video"), model.Default(true)),
)

// Sora2Characters - Sora 2 character creation for video generation
var Sora2Characters = model.Define(
	"sora-2-characters",
	"Sora 2 Characters",
	model.CategoryTextToVideo,
	model.WithProvider("openai"),
	model.WithTimeout(10*time.Minute),
).Optional(
	model.Str("character_prompt", model.Desc("Stable character traits description (e.g., 'cheerful barista, green apron')"), model.MaxLen(5000)),
	model.Str("safety_instruction", model.Desc("Content boundaries and limits (e.g., 'no violence, PG-13 max')"), model.MaxLen(5000)),
)

// SoraWatermarkRemover - Sora video watermark removal
var SoraWatermarkRemover = model.Define(
	"sora-watermark-remover",
	"Sora Watermark Remover",
	model.CategoryVideoToVideo,
	model.WithProvider("openai"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("video_url", model.Desc("Sora 2 video URL from OpenAI (must start with sora.chatgpt.com)"), model.MaxLen(500)),
)

// Kling25TurboTextToVideoPro - Kling 2.5 Turbo Text to Video Pro generation
var Kling25TurboTextToVideoPro = model.Define(
	"kling/v2-5-turbo-text-to-video-pro",
	"Kling 2.5 Turbo Text To Video Pro",
	model.CategoryTextToVideo,
	model.WithProvider("kling"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description of the video to generate"), model.MaxLen(2500)),
).Optional(
	model.Enum("duration", []string{"5", "10"}, model.Desc("Duration of the video in seconds"), model.Default("5")),
	model.Enum("aspect_ratio", []string{"16:9", "9:16", "1:1"}, model.Desc("Aspect ratio of the video"), model.Default("16:9")),
	model.Str("negative_prompt", model.Desc("Things to avoid in the generated video"), model.MaxLen(2500)),
	model.Float("cfg_scale", model.Desc("CFG scale - how closely to follow the prompt (0-1)"), model.Min(0), model.Max(1), model.Default(0.5)),
)

// Kling25TurboImageToVideoPro - Kling 2.5 Turbo Image to Video Pro generation
var Kling25TurboImageToVideoPro = model.Define(
	"kling/v2-5-turbo-image-to-video-pro",
	"Kling 2.5 Turbo Image To Video Pro",
	model.CategoryImageToVideo,
	model.WithProvider("kling"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description for the video generation"), model.MaxLen(2500)),
	model.Str("image_url", model.Desc("URL of the image to be used for the video"), model.MaxLen(2048)),
).Optional(
	model.Str("tail_image_url", model.Desc("Tail frame image URL for the video"), model.MaxLen(2048)),
	model.Enum("duration", []string{"5", "10"}, model.Desc("Duration of the video in seconds"), model.Default("5")),
	model.Str("negative_prompt", model.Desc("Elements to avoid in the video"), model.MaxLen(2496)),
	model.Float("cfg_scale", model.Desc("CFG scale - how closely to follow the prompt (0-1)"), model.Min(0), model.Max(1), model.Default(0.5)),
)

// Wan25ImageToVideo - Wan 2.5 Image to Video generation
var Wan25ImageToVideo = model.Define(
	"wan/2-5-image-to-video",
	"Wan 2.5 Image To Video",
	model.CategoryImageToVideo,
	model.WithProvider("wan"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the desired video motion"), model.MaxLen(800)),
	model.Str("image_url", model.Desc("URL of the image to use as the first frame"), model.MaxLen(2048)),
).Optional(
	model.Enum("duration", []string{"5", "10"}, model.Desc("Duration of the video in seconds"), model.Default("5")),
	model.Enum("resolution", []string{"720p", "1080p"}, model.Desc("Video resolution"), model.Default("1080p")),
	model.Str("negative_prompt", model.Desc("Content to avoid"), model.MaxLen(500)),
	model.Bool("enable_prompt_expansion", model.Desc("Enable prompt rewriting using LLM"), model.Default(true)),
	model.Int("seed", model.Desc("Random seed for reproducibility")),
)

// Wan25TextToVideo - Wan 2.5 Text to Video generation
var Wan25TextToVideo = model.Define(
	"wan/2-5-text-to-video",
	"Wan 2.5 Text To Video",
	model.CategoryTextToVideo,
	model.WithProvider("wan"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt for video generation (max 800 characters)"), model.MaxLen(800)),
).Optional(
	model.Enum("duration", []string{"5", "10"}, model.Desc("Duration of the video in seconds"), model.Default("5")),
	model.Enum("aspect_ratio", []string{"16:9", "9:16", "1:1"}, model.Desc("Aspect ratio of the video"), model.Default("16:9")),
	model.Enum("resolution", []string{"720p", "1080p"}, model.Desc("Video resolution"), model.Default("1080p")),
	model.Str("negative_prompt", model.Desc("Content to avoid"), model.MaxLen(500)),
	model.Bool("enable_prompt_expansion", model.Desc("Enable prompt rewriting using LLM"), model.Default(true)),
	model.Int("seed", model.Desc("Random seed for reproducibility")),
)

// Wan22AnimateMove - Wan 2.2 Animate Move (image+video to video)
var Wan22AnimateMove = model.Define(
	"wan/2-2-animate-move",
	"Wan 2.2 Animate Move",
	model.CategoryVideoToVideo,
	model.WithProvider("wan"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("video_url", model.Desc("URL of the input video"), model.MaxLen(2048)),
	model.Str("image_url", model.Desc("URL of the input image"), model.MaxLen(2048)),
).Optional(
	model.Enum("resolution", []string{"480p", "580p", "720p"}, model.Desc("Resolution of the generated video"), model.Default("480p")),
)

// Wan22AnimateReplace - Wan 2.2 Animate Replace (image+video to video)
var Wan22AnimateReplace = model.Define(
	"wan/2-2-animate-replace",
	"Wan 2.2 Animate Replace",
	model.CategoryVideoToVideo,
	model.WithProvider("wan"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("video_url", model.Desc("URL of the input video"), model.MaxLen(2048)),
	model.Str("image_url", model.Desc("URL of the input image"), model.MaxLen(2048)),
).Optional(
	model.Enum("resolution", []string{"480p", "580p", "720p"}, model.Desc("Resolution of the generated video"), model.Default("480p")),
)

// TopazVideoUpscale - Topaz Video Upscale for video enhancement
var TopazVideoUpscale = model.Define(
	"topaz/video-upscale",
	"Topaz Video Upscale",
	model.CategoryVideoToVideo,
	model.WithProvider("topaz"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("video_url", model.Desc("URL of the video to upscale"), model.MaxLen(2048)),
).Optional(
	model.Enum("upscale_factor", []string{"1", "2", "4"}, model.Desc("Factor to upscale the video by"), model.Default("2")),
)

// InfinitalkFromAudio - Infinitalk From Audio (image + audio to talking video)
var InfinitalkFromAudio = model.Define(
	"infinitalk/from-audio",
	"Infinitalk From Audio",
	model.CategoryImageToVideo,
	model.WithProvider("infinitalk"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("image_url", model.Desc("URL of the input image"), model.MaxLen(2048)),
	model.Str("audio_url", model.Desc("URL of the audio file (max 15 seconds)"), model.MaxLen(2048)),
	model.Str("prompt", model.Desc("Text prompt to guide video generation"), model.MaxLen(5000)),
).Optional(
	model.Enum("resolution", []string{"480p", "720p"}, model.Desc("Resolution of the generated video"), model.Default("480p")),
	model.Int("seed", model.Desc("Random seed for reproducibility"), model.Min(10000), model.Max(1000000)),
)

// Wan22A14bSpeechToVideoTurbo - Wan 2.2 A14b Speech To Video Turbo (image + audio to talking video)
var Wan22A14bSpeechToVideoTurbo = model.Define(
	"wan/2-2-a14b-speech-to-video-turbo",
	"Wan 2.2 A14b Speech To Video Turbo",
	model.CategoryImageToVideo,
	model.WithProvider("wan"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt for video generation"), model.MaxLen(5000)),
	model.Str("image_url", model.Desc("URL of the input image"), model.MaxLen(2048)),
	model.Str("audio_url", model.Desc("URL of the audio file"), model.MaxLen(2048)),
).Optional(
	model.Int("num_frames", model.Desc("Number of frames to generate (40-120, multiple of 4)"), model.Min(40), model.Max(120), model.Default(80)),
	model.Int("frames_per_second", model.Desc("Frames per second of the video"), model.Min(4), model.Max(60), model.Default(16)),
	model.Enum("resolution", []string{"480p", "580p", "720p"}, model.Desc("Resolution of the generated video"), model.Default("480p")),
	model.Str("negative_prompt", model.Desc("What to exclude from the video"), model.MaxLen(500)),
	model.Int("seed", model.Desc("Random seed for reproducibility")),
	model.Int("num_inference_steps", model.Desc("Number of inference steps"), model.Min(2), model.Max(40), model.Default(27)),
	model.Float("guidance_scale", model.Desc("Guidance scale for prompt adherence"), model.Min(1), model.Max(10), model.Default(3.5)),
	model.Float("shift", model.Desc("Shift value for the video"), model.Min(1), model.Max(10), model.Default(5)),
	model.Bool("enable_safety_checker", model.Desc("Enable safety checking"), model.Default(true)),
)

// KlingV1AvatarStandard - Kling V1 Avatar Standard (image + audio to avatar video)
var KlingV1AvatarStandard = model.Define(
	"kling/v1-avatar-standard",
	"Kling V1 Avatar Standard",
	model.CategoryImageToVideo,
	model.WithProvider("kling"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("image_url", model.Desc("URL of the avatar image"), model.MaxLen(2048)),
	model.Str("audio_url", model.Desc("URL of the audio file"), model.MaxLen(2048)),
	model.Str("prompt", model.Desc("Text prompt for video generation"), model.MaxLen(5000)),
)

// KlingAiAvatarV1Pro - Kling AI Avatar V1 Pro (image + audio to avatar video)
var KlingAiAvatarV1Pro = model.Define(
	"kling/ai-avatar-v1-pro",
	"Kling AI Avatar V1 Pro",
	model.CategoryImageToVideo,
	model.WithProvider("kling"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("image_url", model.Desc("URL of the avatar image"), model.MaxLen(2048)),
	model.Str("audio_url", model.Desc("URL of the audio file"), model.MaxLen(2048)),
	model.Str("prompt", model.Desc("Text prompt for video generation"), model.MaxLen(5000)),
)

// Wan22A14bTextToVideoTurbo - Wan 2.2 A14b Text To Video Turbo
var Wan22A14bTextToVideoTurbo = model.Define(
	"wan/2-2-a14b-text-to-video-turbo",
	"Wan 2.2 A14b Text To Video Turbo",
	model.CategoryTextToVideo,
	model.WithProvider("wan"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt to guide video generation"), model.MaxLen(5000)),
).Optional(
	model.Enum("resolution", []string{"480p", "720p"}, model.Desc("Resolution of the generated video"), model.Default("720p")),
	model.Enum("aspect_ratio", []string{"16:9", "9:16"}, model.Desc("Aspect ratio of the generated video"), model.Default("16:9")),
	model.Bool("enable_prompt_expansion", model.Desc("Enable prompt expansion with LLM"), model.Default(false)),
	model.Int("seed", model.Desc("Random seed for reproducibility"), model.Min(0), model.Max(2147483647)),
	model.Enum("acceleration", []string{"none", "regular"}, model.Desc("Acceleration level"), model.Default("none")),
)

// Wan22A14bImageToVideoTurbo - Wan 2.2 A14b Image To Video Turbo
var Wan22A14bImageToVideoTurbo = model.Define(
	"wan/2-2-a14b-image-to-video-turbo",
	"Wan 2.2 A14b Image To Video Turbo",
	model.CategoryImageToVideo,
	model.WithProvider("wan"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("image_url", model.Desc("URL of the input image"), model.MaxLen(2048)),
	model.Str("prompt", model.Desc("Text prompt to guide video generation"), model.MaxLen(5000)),
).Optional(
	model.Enum("resolution", []string{"480p", "720p"}, model.Desc("Resolution of the generated video"), model.Default("720p")),
	model.Bool("enable_prompt_expansion", model.Desc("Enable prompt expansion with LLM"), model.Default(false)),
	model.Int("seed", model.Desc("Random seed for reproducibility"), model.Min(0), model.Max(2147483647)),
	model.Enum("acceleration", []string{"none", "regular"}, model.Desc("Acceleration level"), model.Default("none")),
)

// KlingV21MasterImageToVideo - Kling V2.1 Master Image To Video
var KlingV21MasterImageToVideo = model.Define(
	"kling/v2-1-master-image-to-video",
	"Kling V2.1 Master Image To Video",
	model.CategoryImageToVideo,
	model.WithProvider("kling"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the video to generate"), model.MaxLen(5000)),
	model.Str("image_url", model.Desc("URL of the input image"), model.MaxLen(2048)),
).Optional(
	model.Enum("duration", []string{"5", "10"}, model.Desc("Duration of the video in seconds"), model.Default("5")),
	model.Str("negative_prompt", model.Desc("Elements to exclude from the video"), model.MaxLen(500)),
	model.Float("cfg_scale", model.Desc("CFG scale for prompt adherence"), model.Min(0), model.Max(1), model.Default(0.5)),
)

// KlingV21Pro - Kling V2.1 Pro (image to video with optional tail image)
var KlingV21Pro = model.Define(
	"kling/v2-1-pro",
	"Kling V2.1 Pro",
	model.CategoryImageToVideo,
	model.WithProvider("kling"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the video to generate"), model.MaxLen(5000)),
	model.Str("image_url", model.Desc("URL of the input image"), model.MaxLen(2048)),
).Optional(
	model.Enum("duration", []string{"5", "10"}, model.Desc("Duration of the video in seconds"), model.Default("5")),
	model.Str("negative_prompt", model.Desc("Terms to avoid in the video"), model.MaxLen(500)),
	model.Float("cfg_scale", model.Desc("CFG scale for prompt adherence"), model.Min(0), model.Max(1), model.Default(0.5)),
	model.Str("tail_image_url", model.Desc("URL of the end frame image"), model.MaxLen(2048)),
)

// KlingV21Standard - Kling V2.1 Standard (image to video)
var KlingV21Standard = model.Define(
	"kling/v2-1-standard",
	"Kling V2.1 Standard",
	model.CategoryImageToVideo,
	model.WithProvider("kling"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the video to generate"), model.MaxLen(5000)),
	model.Str("image_url", model.Desc("URL of the input image"), model.MaxLen(2048)),
).Optional(
	model.Enum("duration", []string{"5", "10"}, model.Desc("Duration of the video in seconds"), model.Default("5")),
	model.Str("negative_prompt", model.Desc("Elements to avoid in the video"), model.MaxLen(500)),
	model.Float("cfg_scale", model.Desc("CFG scale for prompt adherence"), model.Min(0), model.Max(1), model.Default(0.5)),
)

// KlingV21MasterTextToVideo - Kling V2.1 Master Text To Video
var KlingV21MasterTextToVideo = model.Define(
	"kling/v2-1-master-text-to-video",
	"Kling V2.1 Master Text To Video",
	model.CategoryTextToVideo,
	model.WithProvider("kling"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt describing the video to generate"), model.MaxLen(5000)),
).Optional(
	model.Enum("duration", []string{"5", "10"}, model.Desc("Duration of the video in seconds"), model.Default("5")),
	model.Enum("aspect_ratio", []string{"16:9", "9:16", "1:1"}, model.Desc("Aspect ratio of the video"), model.Default("16:9")),
	model.Str("negative_prompt", model.Desc("Elements to avoid in the video"), model.MaxLen(500)),
	model.Float("cfg_scale", model.Desc("CFG scale for prompt adherence"), model.Min(0), model.Max(1), model.Default(0.5)),
)

// BytedanceV1ProImageToVideo - Bytedance V1 Pro Image To Video
var BytedanceV1ProImageToVideo = model.Define(
	"bytedance/v1-pro-image-to-video",
	"Bytedance V1 Pro Image To Video",
	model.CategoryImageToVideo,
	model.WithProvider("bytedance"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt for video generation"), model.MaxLen(10000)),
	model.Str("image_url", model.Desc("URL of the input image"), model.MaxLen(2048)),
).Optional(
	model.Enum("resolution", []string{"480p", "720p", "1080p"}, model.Desc("Video resolution"), model.Default("720p")),
	model.Enum("duration", []string{"5", "10"}, model.Desc("Duration of the video in seconds"), model.Default("5")),
	model.Bool("camera_fixed", model.Desc("Whether to fix the camera position"), model.Default(false)),
	model.Int("seed", model.Desc("Random seed (-1 for random)"), model.Min(-1), model.Max(2147483647), model.Default(-1)),
	model.Bool("enable_safety_checker", model.Desc("Enable safety checking"), model.Default(true)),
)

// BytedanceV1LiteImageToVideo - Bytedance V1 Lite Image To Video
var BytedanceV1LiteImageToVideo = model.Define(
	"bytedance/v1-lite-image-to-video",
	"Bytedance V1 Lite Image To Video",
	model.CategoryImageToVideo,
	model.WithProvider("bytedance"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt for video generation"), model.MaxLen(10000)),
	model.Str("image_url", model.Desc("URL of the input image"), model.MaxLen(2048)),
).Optional(
	model.Enum("resolution", []string{"480p", "720p", "1080p"}, model.Desc("Video resolution"), model.Default("720p")),
	model.Enum("duration", []string{"5", "10"}, model.Desc("Duration of the video in seconds"), model.Default("5")),
	model.Bool("camera_fixed", model.Desc("Whether to fix the camera position"), model.Default(false)),
	model.Int("seed", model.Desc("Random seed (-1 for random)"), model.Min(-1), model.Max(2147483647), model.Default(-1)),
	model.Bool("enable_safety_checker", model.Desc("Enable safety checking"), model.Default(true)),
	model.Str("end_image_url", model.Desc("URL of the image the video ends with"), model.MaxLen(2048)),
)

// BytedanceV1ProTextToVideo - Bytedance V1 Pro Text To Video
var BytedanceV1ProTextToVideo = model.Define(
	"bytedance/v1-pro-text-to-video",
	"Bytedance V1 Pro Text To Video",
	model.CategoryTextToVideo,
	model.WithProvider("bytedance"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt for video generation"), model.MaxLen(10000)),
).Optional(
	model.Enum("aspect_ratio", []string{"21:9", "16:9", "4:3", "1:1", "3:4", "9:16"}, model.Desc("Aspect ratio of the video"), model.Default("16:9")),
	model.Enum("resolution", []string{"480p", "720p", "1080p"}, model.Desc("Video resolution"), model.Default("720p")),
	model.Enum("duration", []string{"5", "10"}, model.Desc("Duration of the video in seconds"), model.Default("5")),
	model.Bool("camera_fixed", model.Desc("Whether to fix the camera position"), model.Default(false)),
	model.Int("seed", model.Desc("Random seed (-1 for random)"), model.Min(-1), model.Max(2147483647), model.Default(-1)),
	model.Bool("enable_safety_checker", model.Desc("Enable safety checking"), model.Default(true)),
)

// BytedanceV1LiteTextToVideo - Bytedance V1 Lite Text To Video
var BytedanceV1LiteTextToVideo = model.Define(
	"bytedance/v1-lite-text-to-video",
	"Bytedance V1 Lite Text To Video",
	model.CategoryTextToVideo,
	model.WithProvider("bytedance"),
	model.WithTimeout(20*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt for video generation"), model.MaxLen(10000)),
).Optional(
	model.Enum("aspect_ratio", []string{"16:9", "4:3", "1:1", "3:4", "9:16", "9:21"}, model.Desc("Aspect ratio of the video"), model.Default("16:9")),
	model.Enum("resolution", []string{"480p", "720p", "1080p"}, model.Desc("Video resolution"), model.Default("720p")),
	model.Enum("duration", []string{"5", "10"}, model.Desc("Duration of the video in seconds"), model.Default("5")),
	model.Bool("camera_fixed", model.Desc("Whether to fix the camera position"), model.Default(false)),
	model.Int("seed", model.Desc("Random seed (-1 for random)"), model.Min(-1), model.Max(2147483647), model.Default(-1)),
	model.Bool("enable_safety_checker", model.Desc("Enable safety checking"), model.Default(true)),
)

// Get returns a video model by identifier.
func Get(id string) *model.Model {
	for _, m := range Models {
		if m.Identifier == id {
			return m
		}
	}
	return nil
}
