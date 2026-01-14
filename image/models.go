// Package image provides image generation model definitions.
package image

import (
	"time"

	"github.com/lpg-it/kie-go/model"
)

// ================================================================================
// Image Models - All image generation models defined declaratively
// ================================================================================

// GoogleImagen4 - High quality text-to-image generation
var GoogleImagen4 = model.Define(
	"google/imagen4",
	"Google Imagen4",
	model.CategoryTextToImage,
	model.WithProvider("google"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description for image generation"), model.MaxLen(5000)),
).Optional(
	model.Str("negative_prompt", model.Desc("What to exclude from the image"), model.MaxLen(5000)),
	model.Enum("aspect_ratio", []string{"1:1", "16:9", "9:16", "3:4", "4:3"}, model.Desc("Aspect ratio"), model.Default("1:1")),
	model.Str("seed", model.Desc("Random seed for reproducibility"), model.MaxLen(500)),
)

// GoogleImagen4Fast - Faster variant of Imagen4
var GoogleImagen4Fast = model.Define(
	"google/imagen4-fast",
	"Google Imagen4 Fast",
	model.CategoryTextToImage,
	model.WithProvider("google"),
	model.WithTimeout(5*time.Minute),
).Required(
	model.Str("prompt", model.Desc("The text prompt describing what you want to see"), model.MaxLen(5000)),
).Optional(
	model.Str("negative_prompt", model.Desc("A description of what to discourage in the generated images"), model.MaxLen(5000)),
	model.Enum("aspect_ratio", []string{"1:1", "16:9", "9:16", "3:4", "4:3"}, model.Desc("The aspect ratio of the generated image"), model.Default("16:9")),
	model.Enum("num_images", []string{"1", "2", "3", "4"}, model.Desc("Number of images to generate"), model.Default("1")),
	model.Int("seed", model.Desc("Random seed for reproducible generation")),
)

// GoogleImagen4Ultra - Highest quality Imagen4 variant
var GoogleImagen4Ultra = model.Define(
	"google/imagen4-ultra",
	"Google Imagen4 Ultra",
	model.CategoryTextToImage,
	model.WithProvider("google"),
	model.WithTimeout(15*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description for image generation"), model.MaxLen(5000)),
).Optional(
	model.Str("negative_prompt", model.Desc("What to exclude from the image"), model.MaxLen(5000)),
	model.Enum("aspect_ratio", []string{"1:1", "16:9", "9:16", "3:4", "4:3"}, model.Desc("Aspect ratio"), model.Default("1:1")),
	model.Str("seed", model.Desc("Random seed for reproducibility"), model.MaxLen(500)),
)

// GoogleNanoBanana - Lightweight image generation
var GoogleNanoBanana = model.Define(
	"google/nano-banana",
	"Google Nano Banana",
	model.CategoryTextToImage,
	model.WithProvider("google"),
	model.WithTimeout(5*time.Minute),
).Required(
	model.Str("prompt", model.Desc("The prompt for image generation"), model.MaxLen(20000)),
).Optional(
	model.Enum("output_format", []string{"png", "jpeg"}, model.Desc("Output format for the images"), model.Default("png")),
	model.Enum("image_size", []string{"1:1", "9:16", "16:9", "3:4", "4:3", "3:2", "2:3", "5:4", "4:5", "21:9", "auto"}, model.Desc("Image size ratio"), model.Default("1:1")),
)

// GoogleNanoBananaEdit - Image editing variant
var GoogleNanoBananaEdit = model.Define(
	"google/nano-banana-edit",
	"Google Nano Banana Edit",
	model.CategoryImageEdit,
	model.WithProvider("google"),
	model.WithTimeout(5*time.Minute),
).Required(
	model.Str("prompt", model.Desc("The prompt for image editing"), model.MaxLen(20000)),
	model.Strings("image_urls", model.Desc("List of URLs of input images for editing"), model.MaxItems(10)),
).Optional(
	model.Enum("output_format", []string{"png", "jpeg"}, model.Desc("Output format for the images"), model.Default("png")),
	model.Enum("image_size", []string{"1:1", "9:16", "16:9", "3:4", "4:3", "3:2", "2:3", "5:4", "4:5", "21:9", "auto"}, model.Desc("Image size ratio"), model.Default("1:1")),
)

// NanoBananaPro - Advanced image generation with reference support
var NanoBananaPro = model.Define(
	"nano-banana-pro",
	"Nano Banana Pro",
	model.CategoryTextToImage,
	model.WithProvider("google"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description of the image to generate"), model.MaxLen(20000)),
).Optional(
	model.Strings("image_input", model.Desc("Input images to transform or use as reference"), model.MaxItems(8)),
	model.Enum("aspect_ratio", []string{"1:1", "2:3", "3:2", "3:4", "4:3", "4:5", "5:4", "9:16", "16:9", "21:9", "auto"}, model.Desc("Aspect ratio of the generated image"), model.Default("1:1")),
	model.Enum("resolution", []string{"1K", "2K", "4K"}, model.Desc("Resolution of the generated image"), model.Default("1K")),
	model.Enum("output_format", []string{"png", "jpg"}, model.Desc("Format of the output image"), model.Default("png")),
)

// GrokImagineTextToImage - Text to image generation
var GrokImagineTextToImage = model.Define(
	"grok-imagine/text-to-image",
	"Grok Imagine Text To Image",
	model.CategoryTextToImage,
	model.WithProvider("xai"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description for image generation"), model.MaxLen(5000)),
).Optional(
	model.Enum("aspect_ratio", []string{"2:3", "3:2", "1:1", "9:16", "16:9"}, model.Desc("Aspect ratio of the generated image"), model.Default("3:2")),
)

// GrokImagineImageToImage - Image to image transformation
var GrokImagineImageToImage = model.Define(
	"grok-imagine/image-to-image",
	"Grok Imagine Image To Image",
	model.CategoryImageToImage,
	model.WithProvider("xai"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Strings("image_urls", model.Desc("Reference image URLs for transformation"), model.MaxItems(5)),
).Optional(
	model.Str("prompt", model.Desc("Text description specifying the desired content or style"), model.MaxLen(390000)),
)

// GrokImagineUpscale - Image upscaling using Kie AI task reference
var GrokImagineUpscale = model.Define(
	"grok-imagine/upscale",
	"Grok Imagine Upscale",
	model.CategoryUpscale,
	model.WithProvider("xai"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("task_id", model.Desc("Kie AI generated task ID to upscale"), model.MaxLen(100)),
)

// Seedream45TextToImage - Seedream 4.5 Text to Image
var Seedream45TextToImage = model.Define(
	"seedream/4.5-text-to-image",
	"Seedream 4.5 Text To Image",
	model.CategoryTextToImage,
	model.WithProvider("bytedance"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("A text description of the image you want to generate"), model.MaxLen(3000)),
	model.Enum("aspect_ratio", []string{"1:1", "4:3", "3:4", "16:9", "9:16", "2:3", "3:2", "21:9"}, model.Desc("Width-height ratio of the image"), model.Default("1:1")),
	model.Enum("quality", []string{"basic", "high"}, model.Desc("Basic outputs 2K images, High outputs 4K images"), model.Default("basic")),
)

// Seedream45Edit - Seedream 4.5 Edit
var Seedream45Edit = model.Define(
	"seedream/4.5-edit",
	"Seedream 4.5 Edit",
	model.CategoryImageEdit,
	model.WithProvider("bytedance"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("A text description of the image you want to generate"), model.MaxLen(3000)),
	model.Strings("image_urls", model.Desc("Input image URLs to edit"), model.MaxItems(10)),
	model.Enum("aspect_ratio", []string{"1:1", "4:3", "3:4", "16:9", "9:16", "2:3", "3:2", "21:9"}, model.Desc("Width-height ratio of the image"), model.Default("1:1")),
	model.Enum("quality", []string{"basic", "high"}, model.Desc("Basic outputs 2K images, High outputs 4K images"), model.Default("basic")),
)

// BytedanceSeedreamV4TextToImage - Seedream V4 Text to Image
var BytedanceSeedreamV4TextToImage = model.Define(
	"bytedance/seedream-v4-text-to-image",
	"Seedream V4 Text To Image",
	model.CategoryTextToImage,
	model.WithProvider("bytedance"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt to generate the image"), model.MaxLen(5000)),
).Optional(
	model.Enum("image_size", []string{"square", "square_hd", "portrait_4_3", "portrait_3_2", "portrait_16_9", "landscape_4_3", "landscape_3_2", "landscape_16_9", "landscape_21_9"}, model.Desc("Size of the generated image"), model.Default("square_hd")),
	model.Enum("image_resolution", []string{"1K", "2K", "4K"}, model.Desc("Resolution of the generated image"), model.Default("1K")),
	model.Int("max_images", model.Desc("Maximum number of images to generate (1-6)"), model.Min(1), model.Max(6), model.Default(1)),
	model.Int("seed", model.Desc("Random seed for reproducibility")),
)

// BytedanceSeedreamV4Edit - Seedream V4 Edit
var BytedanceSeedreamV4Edit = model.Define(
	"bytedance/seedream-v4-edit",
	"Seedream V4 Edit",
	model.CategoryImageEdit,
	model.WithProvider("bytedance"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text prompt to edit the image"), model.MaxLen(5000)),
	model.Strings("image_urls", model.Desc("Input image URLs to edit"), model.MaxItems(10)),
).Optional(
	model.Enum("image_size", []string{"square", "square_hd", "portrait_4_3", "portrait_3_2", "portrait_16_9", "landscape_4_3", "landscape_3_2", "landscape_16_9", "landscape_21_9"}, model.Desc("Size of the generated image"), model.Default("square_hd")),
	model.Enum("image_resolution", []string{"1K", "2K", "4K"}, model.Desc("Resolution of the generated image"), model.Default("1K")),
	model.Int("max_images", model.Desc("Maximum number of images to generate (1-6)"), model.Min(1), model.Max(6), model.Default(1)),
	model.Int("seed", model.Desc("Random seed for reproducibility")),
)

// BytedanceSeedream - Bytedance Seedream 3.0 Text to Image
var BytedanceSeedream = model.Define(
	"bytedance/seedream",
	"Bytedance Seedream",
	model.CategoryTextToImage,
	model.WithProvider("bytedance"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("The text prompt used to generate the image"), model.MaxLen(5000)),
).Optional(
	model.Enum("image_size", []string{"square", "square_hd", "portrait_4_3", "portrait_16_9", "landscape_4_3", "landscape_16_9"}, model.Desc("Size of the generated image"), model.Default("square_hd")),
	model.Float("guidance_scale", model.Desc("Controls how closely the output image aligns with the input prompt"), model.Min(1), model.Max(10), model.Default(2.5)),
	model.Int("seed", model.Desc("Random seed to control the stochasticity of image generation")),
	model.Bool("enable_safety_checker", model.Desc("If set to true, the safety checker will be enabled"), model.Default(true)),
)

// QwenImageToImage - Qwen Image to Image transformation
var QwenImageToImage = model.Define(
	"qwen/image-to-image",
	"Qwen Image To Image",
	model.CategoryImageToImage,
	model.WithProvider("qwen"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("The prompt to generate the image with"), model.MaxLen(5000)),
	model.Str("image_url", model.Desc("The reference image URL to guide the generation")),
).Optional(
	model.Float("strength", model.Desc("Denoising strength. 1.0 = fully remake; 0.0 = preserve original"), model.Min(0), model.Max(1), model.Default(0.8)),
	model.Enum("output_format", []string{"png", "jpeg"}, model.Desc("The format of the generated image"), model.Default("png")),
	model.Enum("acceleration", []string{"none", "regular", "high"}, model.Desc("Acceleration level for image generation"), model.Default("none")),
	model.Str("negative_prompt", model.Desc("The negative prompt for the generation"), model.MaxLen(500)),
	model.Int("seed", model.Desc("Random seed for reproducibility")),
	model.Int("num_inference_steps", model.Desc("The number of inference steps to perform"), model.Min(2), model.Max(250), model.Default(30)),
	model.Float("guidance_scale", model.Desc("The CFG scale for prompt adherence"), model.Min(0), model.Max(20), model.Default(2.5)),
	model.Bool("enable_safety_checker", model.Desc("Enable or disable safety checker"), model.Default(true)),
)

// QwenTextToImage - Qwen Text to Image generation
var QwenTextToImage = model.Define(
	"qwen/text-to-image",
	"Qwen Text To Image",
	model.CategoryTextToImage,
	model.WithProvider("qwen"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("The prompt to generate the image with"), model.MaxLen(5000)),
).Optional(
	model.Enum("image_size", []string{"square", "square_hd", "portrait_4_3", "portrait_16_9", "landscape_4_3", "landscape_16_9"}, model.Desc("The size of the generated image"), model.Default("square_hd")),
	model.Int("num_inference_steps", model.Desc("The number of inference steps to perform"), model.Min(2), model.Max(250), model.Default(30)),
	model.Int("seed", model.Desc("Random seed for reproducibility")),
	model.Float("guidance_scale", model.Desc("The CFG scale for prompt adherence"), model.Min(0), model.Max(20), model.Default(2.5)),
	model.Bool("enable_safety_checker", model.Desc("Enable or disable safety checker"), model.Default(true)),
	model.Enum("output_format", []string{"png", "jpeg"}, model.Desc("The format of the generated image"), model.Default("png")),
	model.Str("negative_prompt", model.Desc("The negative prompt for the generation"), model.MaxLen(500)),
	model.Enum("acceleration", []string{"none", "regular", "high"}, model.Desc("Acceleration level for image generation"), model.Default("none")),
)

// QwenImageEdit - Qwen Image Edit
var QwenImageEdit = model.Define(
	"qwen/image-edit",
	"Qwen Image Edit",
	model.CategoryImageEdit,
	model.WithProvider("qwen"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("The prompt to generate the image with"), model.MaxLen(2000)),
	model.Str("image_url", model.Desc("The URL of the image to edit")),
).Optional(
	model.Enum("acceleration", []string{"none", "regular", "high"}, model.Desc("Acceleration level for image generation"), model.Default("none")),
	model.Enum("image_size", []string{"square", "square_hd", "portrait_4_3", "portrait_16_9", "landscape_4_3", "landscape_16_9"}, model.Desc("The size of the generated image"), model.Default("landscape_4_3")),
	model.Int("num_inference_steps", model.Desc("The number of inference steps to perform"), model.Min(2), model.Max(49), model.Default(25)),
	model.Int("seed", model.Desc("Random seed for reproducibility")),
	model.Float("guidance_scale", model.Desc("The CFG scale for prompt adherence"), model.Min(0), model.Max(20), model.Default(4)),
	model.Bool("sync_mode", model.Desc("Wait for image generation before returning response"), model.Default(false)),
	model.Enum("num_images", []string{"1", "2", "3", "4"}, model.Desc("Number of images to generate")),
	model.Bool("enable_safety_checker", model.Desc("Enable or disable safety checker"), model.Default(true)),
	model.Enum("output_format", []string{"jpeg", "png"}, model.Desc("The format of the generated image"), model.Default("png")),
	model.Str("negative_prompt", model.Desc("The negative prompt for the generation"), model.MaxLen(500)),
)

// RecraftCrispUpscale - Recraft Crisp Upscale
var RecraftCrispUpscale = model.Define(
	"recraft/crisp-upscale",
	"Recraft Crisp Upscale",
	model.CategoryUpscale,
	model.WithProvider("recraft"),
	model.WithTimeout(5*time.Minute),
).Required(
	model.Str("image", model.Desc("URL of the image to upscale")),
)

// RecraftRemoveBackground - Recraft Remove Background
var RecraftRemoveBackground = model.Define(
	"recraft/remove-background",
	"Recraft Remove Background",
	model.CategoryImageEdit,
	model.WithProvider("recraft"),
	model.WithTimeout(5*time.Minute),
).Required(
	model.Str("image", model.Desc("URL of the image to remove background from")),
)

// TopazImageUpscale - Topaz Image Upscale
var TopazImageUpscale = model.Define(
	"topaz/image-upscale",
	"Topaz Image Upscale",
	model.CategoryUpscale,
	model.WithProvider("topaz"),
	model.WithTimeout(5*time.Minute),
).Required(
	model.Str("image_url", model.Desc("URL of the image to upscale")),
	model.Enum("upscale_factor", []string{"1", "2", "4", "8"}, model.Desc("Factor to upscale the image by"), model.Default("2")),
)

// GptImage15ImageToImage - GPT Image 1.5 Image To Image
var GptImage15ImageToImage = model.Define(
	"gpt-image/1.5-image-to-image",
	"GPT Image 1.5 Image To Image",
	model.CategoryImageToImage,
	model.WithProvider("gpt-image"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Strings("input_urls", model.Desc("Input image URLs to transform"), model.MaxItems(10)),
	model.Str("prompt", model.Desc("A text description of the image you want to generate"), model.MaxLen(3000)),
	model.Enum("aspect_ratio", []string{"1:1", "2:3", "3:2"}, model.Desc("Width-height ratio of the image"), model.Default("3:2")),
	model.Enum("quality", []string{"medium", "high"}, model.Desc("Quality: medium=balanced, high=slow/detailed"), model.Default("medium")),
)

// GptImage15TextToImage - GPT Image 1.5 Text To Image
var GptImage15TextToImage = model.Define(
	"gpt-image/1.5-text-to-image",
	"GPT Image 1.5 Text To Image",
	model.CategoryTextToImage,
	model.WithProvider("gpt-image"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("A text description of the image you want to generate"), model.MaxLen(3000)),
	model.Enum("aspect_ratio", []string{"1:1", "2:3", "3:2"}, model.Desc("Width-height ratio of the image"), model.Default("3:2")),
	model.Enum("quality", []string{"medium", "high"}, model.Desc("Quality: medium=balanced, high=slow/detailed"), model.Default("medium")),
)

// ZImage - Z Image text to image generation
var ZImage = model.Define(
	"z-image",
	"Z Image",
	model.CategoryTextToImage,
	model.WithProvider("z-image"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("A text description of the image you want to generate"), model.MaxLen(1000)),
	model.Enum("aspect_ratio", []string{"1:1", "4:3", "3:4", "16:9", "9:16"}, model.Desc("Aspect ratio for the generated image"), model.Default("1:1")),
)

// Flux2ProImageToImage - Flux 2 Pro Image To Image
var Flux2ProImageToImage = model.Define(
	"flux-2/pro-image-to-image",
	"Flux 2 Pro Image To Image",
	model.CategoryImageToImage,
	model.WithProvider("flux-2"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Strings("input_urls", model.Desc("Input reference images (1-8 images)"), model.MaxItems(8)),
	model.Str("prompt", model.Desc("A text description of the image you want to generate"), model.MaxLen(5000)),
	model.Enum("aspect_ratio", []string{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3", "auto"}, model.Desc("Aspect ratio for the generated image"), model.Default("4:3")),
	model.Enum("resolution", []string{"1K", "2K"}, model.Desc("Output image resolution"), model.Default("1K")),
)

// Flux2FlexImageToImage - Flux 2 Flex Image To Image
var Flux2FlexImageToImage = model.Define(
	"flux-2/flex-image-to-image",
	"Flux 2 Flex Image To Image",
	model.CategoryImageToImage,
	model.WithProvider("flux-2"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Strings("input_urls", model.Desc("Input reference images (1-8 images)"), model.MaxItems(8)),
	model.Str("prompt", model.Desc("Text description of the image to generate"), model.MaxLen(5000)),
	model.Enum("aspect_ratio", []string{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3", "auto"}, model.Desc("Aspect ratio for the generated image"), model.Default("1:1")),
	model.Enum("resolution", []string{"1K", "2K"}, model.Desc("Output image resolution"), model.Default("1K")),
)

// Flux2FlexTextToImage - Flux 2 Flex Text To Image
var Flux2FlexTextToImage = model.Define(
	"flux-2/flex-text-to-image",
	"Flux 2 Flex Text To Image",
	model.CategoryTextToImage,
	model.WithProvider("flux-2"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description of the image to generate"), model.MaxLen(5000)),
	model.Enum("aspect_ratio", []string{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3", "auto"}, model.Desc("Aspect ratio for the generated image"), model.Default("1:1")),
	model.Enum("resolution", []string{"1K", "2K"}, model.Desc("Output image resolution"), model.Default("1K")),
)

// Flux2ProTextToImage - Flux 2 Pro Text To Image
var Flux2ProTextToImage = model.Define(
	"flux-2/pro-text-to-image",
	"Flux 2 Pro Text To Image",
	model.CategoryTextToImage,
	model.WithProvider("flux-2"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Text description of the image to generate"), model.MaxLen(5000)),
	model.Enum("aspect_ratio", []string{"1:1", "4:3", "3:4", "16:9", "9:16", "3:2", "2:3", "auto"}, model.Desc("Aspect ratio for the generated image"), model.Default("1:1")),
	model.Enum("resolution", []string{"1K", "2K"}, model.Desc("Output image resolution"), model.Default("1K")),
)

// IdeogramV3Reframe - Ideogram V3 Reframe
var IdeogramV3Reframe = model.Define(
	"ideogram/v3-reframe",
	"Ideogram V3 Reframe",
	model.CategoryImageEdit,
	model.WithProvider("ideogram"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("image_url", model.Desc("URL of the image to reframe")),
	model.Enum("image_size", []string{"square", "square_hd", "portrait_4_3", "portrait_16_9", "landscape_4_3", "landscape_16_9"}, model.Desc("Resolution for the reframed output image"), model.Default("square_hd")),
).Optional(
	model.Enum("rendering_speed", []string{"TURBO", "BALANCED", "QUALITY"}, model.Desc("Rendering speed"), model.Default("BALANCED")),
	model.Enum("style", []string{"AUTO", "GENERAL", "REALISTIC", "DESIGN"}, model.Desc("Style type to generate with"), model.Default("AUTO")),
	model.Enum("num_images", []string{"1", "2", "3", "4"}, model.Desc("Number of images to generate"), model.Default("1")),
	model.Int("seed", model.Desc("Seed for the random number generator")),
)

// IdeogramV3TextToImage - Ideogram V3 Text to Image
var IdeogramV3TextToImage = model.Define(
	"ideogram/v3-text-to-image",
	"Ideogram V3 Text To Image",
	model.CategoryTextToImage,
	model.WithProvider("ideogram"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("Description of the image to generate"), model.MaxLen(5000)),
).Optional(
	model.Enum("rendering_speed", []string{"TURBO", "BALANCED", "QUALITY"}, model.Desc("The rendering speed to use"), model.Default("BALANCED")),
	model.Enum("style", []string{"AUTO", "GENERAL", "REALISTIC", "DESIGN"}, model.Desc("The style type to generate with"), model.Default("AUTO")),
	model.Bool("expand_prompt", model.Desc("Determine if MagicPrompt should be used"), model.Default(true)),
	model.Enum("image_size", []string{"square", "square_hd", "portrait_4_3", "portrait_16_9", "landscape_4_3", "landscape_16_9"}, model.Desc("The resolution of the generated image"), model.Default("square_hd")),
	model.Int("seed", model.Desc("Seed for the random number generator")),
	model.Str("negative_prompt", model.Desc("Description of what to exclude from an image"), model.MaxLen(5000)),
)

// IdeogramV3Edit - Ideogram V3 Edit (inpainting)
var IdeogramV3Edit = model.Define(
	"ideogram/v3-edit",
	"Ideogram V3 Edit",
	model.CategoryImageEdit,
	model.WithProvider("ideogram"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("The prompt to fill the masked part of the image"), model.MaxLen(5000)),
	model.Str("image_url", model.Desc("The image URL to generate an image from")),
	model.Str("mask_url", model.Desc("The mask URL to inpaint the image")),
).Optional(
	model.Enum("rendering_speed", []string{"TURBO", "BALANCED", "QUALITY"}, model.Desc("The rendering speed to use"), model.Default("BALANCED")),
	model.Bool("expand_prompt", model.Desc("Determine if MagicPrompt should be used"), model.Default(true)),
	model.Int("seed", model.Desc("Seed for the random number generator")),
)

// IdeogramV3Remix - Ideogram V3 Remix (image-to-image)
var IdeogramV3Remix = model.Define(
	"ideogram/v3-remix",
	"Ideogram V3 Remix",
	model.CategoryImageToImage,
	model.WithProvider("ideogram"),
	model.WithTimeout(10*time.Minute),
).Required(
	model.Str("prompt", model.Desc("The prompt to remix the image with"), model.MaxLen(5000)),
	model.Str("image_url", model.Desc("The image URL to remix")),
).Optional(
	model.Enum("rendering_speed", []string{"TURBO", "BALANCED", "QUALITY"}, model.Desc("The rendering speed to use"), model.Default("BALANCED")),
	model.Enum("style", []string{"AUTO", "GENERAL", "REALISTIC", "DESIGN"}, model.Desc("The style type to generate with"), model.Default("AUTO")),
	model.Bool("expand_prompt", model.Desc("Determine if MagicPrompt should be used"), model.Default(true)),
	model.Enum("image_size", []string{"square", "square_hd", "portrait_4_3", "portrait_16_9", "landscape_4_3", "landscape_16_9"}, model.Desc("The resolution of the generated image"), model.Default("square_hd")),
	model.Enum("num_images", []string{"1", "2", "3", "4"}, model.Desc("Number of images to generate"), model.Default("1")),
	model.Int("seed", model.Desc("Seed for the random number generator")),
	model.Float("strength", model.Desc("Strength of the input image in the remix"), model.Min(0.01), model.Max(1), model.Default(0.8)),
	model.Str("negative_prompt", model.Desc("Description of what to exclude from an image"), model.MaxLen(5000)),
)

// ================================================================================
// Model Registry
// ================================================================================

// Models contains all image generation models.
var Models = []*model.Model{
	GoogleImagen4,
	GoogleImagen4Fast,
	GoogleImagen4Ultra,
	GoogleNanoBanana,
	GoogleNanoBananaEdit,
	NanoBananaPro,
	GrokImagineTextToImage,
	GrokImagineImageToImage,
	GrokImagineUpscale,
	Seedream45TextToImage,
	Seedream45Edit,
	BytedanceSeedreamV4TextToImage,
	BytedanceSeedreamV4Edit,
	RecraftCrispUpscale,
	RecraftRemoveBackground,
	TopazImageUpscale,
	GptImage15ImageToImage,
	GptImage15TextToImage,
	ZImage,
	Flux2ProImageToImage,
	Flux2FlexImageToImage,
	Flux2FlexTextToImage,
	Flux2ProTextToImage,
	IdeogramV3Reframe,
	IdeogramV3TextToImage,
	IdeogramV3Edit,
	IdeogramV3Remix,
	BytedanceSeedream,
	QwenImageToImage,
	QwenTextToImage,
	QwenImageEdit,
}

// Get returns an image model by identifier.
func Get(id string) *model.Model {
	for _, m := range Models {
		if m.Identifier == id {
			return m
		}
	}
	return nil
}
