// Example: Model and Parameter Discovery using KIE SDK
//
// This example demonstrates how to discover model parameters programmatically:
// 1. List all available models
// 2. Get required and optional fields for a model
// 3. Understand field types and constraints
// 4. Build a simple model explorer
//
// This is useful for:
// - Building dynamic UIs that adapt to model parameters
// - Generating documentation
// - Validating user input before API calls
// - Understanding what parameters are available
//
// Usage:
//
//	go run main.go
package main

import (
	"fmt"
	"strings"

	kie "github.com/lpg-it/kie-go"
)

func main() {
	// Example 1: List all available models
	fmt.Println("=== Example 1: List All Models ===")
	listAllModels()

	// Example 2: Explore a specific model
	fmt.Println("\n=== Example 2: Explore Model Parameters ===")
	exploreModel(kie.NanoBananaPro)

	// Example 3: Compare models
	fmt.Println("\n=== Example 3: Compare Models ===")
	compareModels([]*kie.Model{
		kie.GoogleImagen4,
		kie.NanoBananaPro,
		kie.Seedream45TextToImage,
	})

	// Example 4: Find models by category
	fmt.Println("\n=== Example 4: Models by Category ===")
	listModelsByCategory()

	// Example 5: Get model dynamically by ID
	fmt.Println("\n=== Example 5: Dynamic Model Lookup ===")
	dynamicModelLookup("nano-banana-pro")
	dynamicModelLookup("seedance/1.5-pro")
}

func listAllModels() {
	imageModels := kie.AllImageModels()
	videoModels := kie.AllVideoModels()

	fmt.Printf("Available Models:\n")
	fmt.Printf("  - Image Models: %d\n", len(imageModels))
	fmt.Printf("  - Video Models: %d\n", len(videoModels))
	fmt.Printf("  - Total: %d\n", len(imageModels)+len(videoModels))

	fmt.Println("\nImage Models:")
	for _, m := range imageModels[:5] { // Show first 5
		fmt.Printf("  - %s (%s)\n", m.Name, m.Identifier)
	}
	if len(imageModels) > 5 {
		fmt.Printf("  ... and %d more\n", len(imageModels)-5)
	}

	fmt.Println("\nVideo Models:")
	for _, m := range videoModels[:5] { // Show first 5
		fmt.Printf("  - %s (%s)\n", m.Name, m.Identifier)
	}
	if len(videoModels) > 5 {
		fmt.Printf("  ... and %d more\n", len(videoModels)-5)
	}
}

func exploreModel(model *kie.Model) {
	fmt.Printf("Model: %s\n", model.Name)
	fmt.Printf("Identifier: %s\n", model.Identifier)
	fmt.Printf("Category: %s\n", model.Category)
	fmt.Printf("Provider: %s\n", model.Provider)
	fmt.Printf("Timeout: %s\n", model.Timeout)

	// Required fields
	requiredFields := model.RequiredFields()
	fmt.Printf("\nRequired Fields (%d):\n", len(requiredFields))
	for _, f := range requiredFields {
		printField(f)
	}

	// Optional fields
	optionalFields := model.OptionalFields()
	fmt.Printf("\nOptional Fields (%d):\n", len(optionalFields))
	for _, f := range optionalFields {
		printField(f)
	}
}

func printField(f kie.Field) {
	fmt.Printf("  • %s (%s)\n", f.Name, f.Type)

	if f.Description != "" {
		fmt.Printf("    Description: %s\n", f.Description)
	}

	if len(f.EnumVals) > 0 {
		fmt.Printf("    Allowed values: %s\n", strings.Join(f.EnumVals, ", "))
	}

	if f.Default != nil {
		fmt.Printf("    Default: %v\n", f.Default)
	}

	if f.MaxLength > 0 {
		fmt.Printf("    Max length: %d\n", f.MaxLength)
	}

	if f.MaxItems > 0 {
		fmt.Printf("    Max items: %d\n", f.MaxItems)
	}

	if f.Min != nil {
		fmt.Printf("    Min: %v\n", *f.Min)
	}

	if f.Max != nil {
		fmt.Printf("    Max: %v\n", *f.Max)
	}
}

func compareModels(models []*kie.Model) {
	fmt.Printf("%-25s | %-12s | Required | Optional\n", "Model", "Category")
	fmt.Println(strings.Repeat("-", 60))

	for _, m := range models {
		req := len(m.RequiredFields())
		opt := len(m.OptionalFields())
		fmt.Printf("%-25s | %-12s | %8d | %8d\n",
			truncate(m.Name, 25),
			m.Category,
			req,
			opt,
		)
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func listModelsByCategory() {
	categories := make(map[kie.Category][]*kie.Model)

	for _, m := range kie.AllModels() {
		categories[m.Category] = append(categories[m.Category], m)
	}

	for category, models := range categories {
		fmt.Printf("%s (%d models):\n", category, len(models))
		for _, m := range models {
			fmt.Printf("  - %s\n", m.Name)
		}
		fmt.Println()
	}
}

func dynamicModelLookup(modelID string) {
	fmt.Printf("\nLooking up model: %s\n", modelID)

	model := kie.GetModel(modelID)
	if model == nil {
		fmt.Printf("  ❌ Model not found\n")
		return
	}

	fmt.Printf("  ✓ Found: %s (%s)\n", model.Name, model.Category)
	fmt.Printf("  Required fields: %v\n", fieldNames(model.RequiredFields()))
	fmt.Printf("  Optional fields: %v\n", fieldNames(model.OptionalFields()))
}

func fieldNames(fields []kie.Field) []string {
	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = f.Name
	}
	return names
}
