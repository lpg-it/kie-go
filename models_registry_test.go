package kie

import "testing"

func TestNewModelsAreRegistered(t *testing.T) {
	cases := []*Model{
		GoogleNanoBanana2,
		GptImage2ImageToImage,
		GptImage2TextToImage,
		Wan27Image,
		Wan27ImagePro,
		Qwen2ImageEdit,
		Kling30MotionControl,
		BytedanceSeedance20Fast,
		BytedanceSeedance20,
		Wan27TextToVideo,
		Wan27ImageToVideo,
		Wan27ReferenceToVideo,
		Wan27VideoEdit,
	}

	for _, m := range cases {
		if m == nil {
			t.Fatalf("expected model to be initialized")
		}
		if got := GetModel(m.Identifier); got == nil {
			t.Fatalf("expected GetModel(%q) to return a model", m.Identifier)
		}
	}
}

func TestNewModelsValidateRepresentativeInputs(t *testing.T) {
	tests := []struct {
		name   string
		model  *Model
		params Params
	}{
		{
			name:  "google nano banana 2",
			model: GoogleNanoBanana2,
			params: Params{
				"prompt": "A comic poster with a banana hero.",
			},
		},
		{
			name:  "gpt image 2 text to image",
			model: GptImage2TextToImage,
			params: Params{
				"prompt": "A cinematic night city poster with neon reflections on a rainy street.",
			},
		},
		{
			name:  "gpt image 2 image to image",
			model: GptImage2ImageToImage,
			params: Params{
				"prompt":     "Transform this product image into a premium e-commerce poster style.",
				"input_urls": []string{"https://example.com/input.png"},
			},
		},
		{
			name:  "wan 2.7 image",
			model: Wan27Image,
			params: Params{
				"prompt": "Generate a cozy tea shop scene.",
				"color_palette": []map[string]string{
					{"hex": "#C2D1E6", "ratio": "23.51%"},
					{"hex": "#F5E6CA", "ratio": "41.49%"},
					{"hex": "#5F4B32", "ratio": "35.00%"},
				},
			},
		},
		{
			name:  "wan 2.7 image pro",
			model: Wan27ImagePro,
			params: Params{
				"prompt":     "Edit the uploaded photo into a clean product shot.",
				"input_urls": []string{"https://example.com/input.jpg"},
				"bbox_list": []any{
					[]any{},
				},
			},
		},
		{
			name:  "qwen2 image edit",
			model: Qwen2ImageEdit,
			params: Params{
				"prompt":    "Replace the background with a minimalist studio.",
				"image_url": "https://example.com/input.png",
			},
		},
		{
			name:  "kling 3.0 motion control",
			model: Kling30MotionControl,
			params: Params{
				"input_urls": []string{"https://example.com/subject.png"},
				"video_urls": []string{"https://example.com/motion.mp4"},
			},
		},
		{
			name:  "seedance 2 fast",
			model: BytedanceSeedance20Fast,
			params: Params{
				"prompt":     "A cinematic beach sunset.",
				"web_search": false,
			},
		},
		{
			name:  "seedance 2",
			model: BytedanceSeedance20,
			params: Params{
				"prompt":     "A cinematic beach sunset.",
				"web_search": false,
			},
		},
		{
			name:  "wan 2.7 text to video",
			model: Wan27TextToVideo,
			params: Params{
				"prompt": "A futuristic city street at night.",
			},
		},
		{
			name:  "wan 2.7 image to video",
			model: Wan27ImageToVideo,
			params: Params{
				"prompt":          "Animate the character waving.",
				"first_frame_url": "https://example.com/first.png",
			},
		},
		{
			name:  "wan 2.7 reference to video",
			model: Wan27ReferenceToVideo,
			params: Params{
				"prompt":          "Have the singer perform on stage.",
				"reference_image": []string{"https://example.com/ref.png"},
			},
		},
		{
			name:  "wan 2.7 video edit",
			model: Wan27VideoEdit,
			params: Params{
				"video_url": "https://example.com/source.mp4",
			},
		},
		{
			name:  "qwen text to image with nsfw checker",
			model: QwenTextToImage,
			params: Params{
				"prompt":       "A photorealistic mountain landscape.",
				"nsfw_checker": false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.model.Validate(tt.params); err != nil {
				t.Fatalf("validate failed: %v", err)
			}
		})
	}
}
