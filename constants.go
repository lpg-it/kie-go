package kie

// ================================================================================
// Common Parameter Constants
// ================================================================================
// These constants can be used with RequestBuilder.Set() or in Params
// to ensure correct parameter values.

// ================================================================================
// Aspect Ratio Constants
// ================================================================================

const (
	Ratio1x1  = "1:1"
	Ratio16x9 = "16:9"
	Ratio9x16 = "9:16"
	Ratio4x3  = "4:3"
	Ratio3x4  = "3:4"
	Ratio3x2  = "3:2"
	Ratio2x3  = "2:3"
	Ratio5x4  = "5:4"
	Ratio4x5  = "4:5"
	Ratio21x9 = "21:9"
	Ratio9x21 = "9:21"
	RatioAuto = "auto"
)

// ================================================================================
// Resolution Constants (for image models like NanoBananaPro)
// ================================================================================

const (
	Resolution1K = "1K"
	Resolution2K = "2K"
	Resolution4K = "4K"
)

// ================================================================================
// Output Format Constants
// ================================================================================

const (
	FormatPNG  = "png"
	FormatJPG  = "jpg"
	FormatJPEG = "jpeg"
)

// ================================================================================
// Video Duration Constants (with "s" suffix)
// ================================================================================

const (
	Duration4s  = "4s"
	Duration5s  = "5s"
	Duration8s  = "8s"
	Duration10s = "10s"
)

// ================================================================================
// Quality Constants
// ================================================================================

const (
	QualityBasic  = "basic"
	QualityMedium = "medium"
	QualityHigh   = "high"
)

// ================================================================================
// Image Size Constants (for models that use image_size instead of aspect_ratio)
// ================================================================================

const (
	Size1x1  = "1:1"
	Size9x16 = "9:16"
	Size16x9 = "16:9"
	Size3x4  = "3:4"
	Size4x3  = "4:3"
	Size3x2  = "3:2"
	Size2x3  = "2:3"
	Size5x4  = "5:4"
	Size4x5  = "4:5"
)

// ================================================================================
// Seedream Image Size Constants
// ================================================================================

const (
	SeedreamSquare       = "square"
	SeedreamSquareHD     = "square_hd"
	SeedreamPortrait43   = "portrait_4_3"
	SeedreamPortrait32   = "portrait_3_2"
	SeedreamPortrait169  = "portrait_16_9"
	SeedreamLandscape43  = "landscape_4_3"
	SeedreamLandscape32  = "landscape_3_2"
	SeedreamLandscape169 = "landscape_16_9"
	SeedreamLandscape219 = "landscape_21_9"
)

// ================================================================================
// Video Mode Constants
// ================================================================================

const (
	ModeFun    = "fun"
	ModeNormal = "normal"
	ModeSpicy  = "spicy"
)

// ================================================================================
// Upscale Factor Constants
// ================================================================================

const (
	Upscale1x = "1"
	Upscale2x = "2"
	Upscale4x = "4"
	Upscale8x = "8"
)

// ================================================================================
// Rendering Speed Constants (for Ideogram models)
// ================================================================================

const (
	SpeedTurbo    = "TURBO"
	SpeedBalanced = "BALANCED"
	SpeedQuality  = "QUALITY"
)

// ================================================================================
// Ideogram Style Constants
// ================================================================================

const (
	StyleAuto      = "AUTO"
	StyleGeneral   = "GENERAL"
	StyleRealistic = "REALISTIC"
	StyleDesign    = "DESIGN"
)

// ================================================================================
// Parameter Name Constants
// ================================================================================
// Use these constants as keys in Params or RequestBuilder.Set()

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
	ParamDuration = "duration"
	ParamMode     = "mode"
	ParamTaskID   = "task_id"
	ParamIndex    = "index"

	// Count parameters
	ParamN     = "n"
	ParamScale = "scale"

	// Upscale parameters
	ParamUpscaleFactor = "upscale_factor"

	// Seedream parameters
	ParamImageResolution = "image_resolution"
	ParamMaxImages       = "max_images"

	// Ideogram parameters
	ParamRenderingSpeed = "rendering_speed"
	ParamStyle          = "style"
	ParamNumImages      = "num_images"

	// Quality parameters
	ParamQuality = "quality"

	// Seedream 3.0 parameters
	ParamGuidanceScale       = "guidance_scale"
	ParamEnableSafetyChecker = "enable_safety_checker"

	// Qwen Image-to-Image parameters
	ParamStrength          = "strength"
	ParamAcceleration      = "acceleration"
	ParamNumInferenceSteps = "num_inference_steps"

	// Qwen Image Edit parameters
	ParamSyncMode = "sync_mode"

	// Ideogram V3 parameters
	ParamExpandPrompt = "expand_prompt"
	ParamMaskURL      = "mask_url"

	// Kling 2.6 Motion Control parameters
	ParamVideoURLs            = "video_urls"
	ParamCharacterOrientation = "character_orientation"
)

// ================================================================================
// Acceleration Constants
// ================================================================================

const (
	AccelerationNone    = "none"
	AccelerationRegular = "regular"
	AccelerationHigh    = "high"
)

// ================================================================================
// Video Resolution Constants
// ================================================================================

// Video resolution constant values (lowercase "p" suffix)
const (
	VideoResolution480p  = "480p"
	VideoResolution580p  = "580p"
	VideoResolution720p  = "720p"
	VideoResolution1080p = "1080p"
)

// Video resolution constant values (uppercase "P" suffix, for Hailuo models)
const (
	VideoResolution512P  = "512P"
	VideoResolution768P  = "768P"
	VideoResolution1080P = "1080P"
)

// ================================================================================
// Character Orientation Constants (for Kling 2.6 Motion Control)
// ================================================================================

const (
	OrientationImage = "image"
	OrientationVideo = "video"
)

// ================================================================================
// Video Duration Seconds Constants (without "s" suffix)
// ================================================================================

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

// ================================================================================
// Bytedance Seedance Parameter Constants
// ================================================================================

const (
	ParamFixedLens     = "fixed_lens"
	ParamGenerateAudio = "generate_audio"
)

// ================================================================================
// Wan Video Parameter Constants
// ================================================================================

const (
	ParamMultiShots = "multi_shots"
	ParamSound      = "sound"
)

// ================================================================================
// Sora Storyboard Constants
// ================================================================================

const (
	StoryboardPortrait  = "portrait"
	StoryboardLandscape = "landscape"
)

const (
	ParamNFrames = "n_frames"
)

// ================================================================================
// Sora 2 Pro Video Constants
// ================================================================================

const (
	SizeStandard = "standard"
	SizeHigh     = "high"
)

const (
	ParamSize            = "size"
	ParamRemoveWatermark = "remove_watermark"
)

// ================================================================================
// Sora 2 Characters Constants
// ================================================================================

const (
	ParamCharacterPrompt   = "character_prompt"
	ParamSafetyInstruction = "safety_instruction"
)

// ================================================================================
// Sora Watermark Remover Constants
// ================================================================================

const (
	ParamVideoURL = "video_url"
)

// ================================================================================
// Kling 2.5 Turbo Constants
// ================================================================================

const (
	ParamCfgScale     = "cfg_scale"
	ParamTailImageURL = "tail_image_url"
)

// ================================================================================
// Wan 2.5 Constants
// ================================================================================

const (
	ParamEnablePromptExpansion = "enable_prompt_expansion"
)

// ================================================================================
// Infinitalk Constants
// ================================================================================

const (
	ParamAudioURL = "audio_url"
)

// ================================================================================
// Wan Speech To Video Constants
// ================================================================================

const (
	ParamNumFrames       = "num_frames"
	ParamFramesPerSecond = "frames_per_second"
	ParamShift           = "shift"
)

// ================================================================================
// Bytedance Video Constants
// ================================================================================

const (
	ParamCameraFixed = "camera_fixed"
	ParamEndImageURL = "end_image_url"
)

// ================================================================================
// Hailuo Video Constants
// ================================================================================

const (
	ParamPromptOptimizer = "prompt_optimizer"
)
