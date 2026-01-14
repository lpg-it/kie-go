// Package model provides the core model framework for defining AI models.
package model

import (
	"fmt"
	"strings"
)

// FieldType represents the type of a field.
type FieldType int

const (
	TypeString FieldType = iota
	TypeInt
	TypeFloat
	TypeBool
	TypeStringSlice
	TypeEnum
)

func (t FieldType) String() string {
	switch t {
	case TypeString:
		return "string"
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeBool:
		return "bool"
	case TypeStringSlice:
		return "[]string"
	case TypeEnum:
		return "enum"
	default:
		return "unknown"
	}
}

// Field represents a model field with its constraints.
type Field struct {
	Name        string
	Type        FieldType
	Required    bool
	Description string

	// Constraints
	MaxLength int
	MaxItems  int
	Min       *float64
	Max       *float64
	Default   any
	EnumVals  []string
}

// FieldOption configures a Field.
type FieldOption func(*Field)

// Desc sets the field description.
func Desc(d string) FieldOption {
	return func(f *Field) { f.Description = d }
}

// MaxLen sets maximum string length.
func MaxLen(n int) FieldOption {
	return func(f *Field) { f.MaxLength = n }
}

// MaxItems sets maximum slice items.
func MaxItems(n int) FieldOption {
	return func(f *Field) { f.MaxItems = n }
}

// Min sets minimum numeric value.
func Min(n float64) FieldOption {
	return func(f *Field) { f.Min = &n }
}

// Max sets maximum numeric value.
func Max(n float64) FieldOption {
	return func(f *Field) { f.Max = &n }
}

// Default sets the default value.
func Default(v any) FieldOption {
	return func(f *Field) { f.Default = v }
}

// Str creates a string field.
func Str(name string, opts ...FieldOption) Field {
	f := Field{Name: name, Type: TypeString}
	for _, opt := range opts {
		opt(&f)
	}
	return f
}

// Int creates an integer field.
func Int(name string, opts ...FieldOption) Field {
	f := Field{Name: name, Type: TypeInt}
	for _, opt := range opts {
		opt(&f)
	}
	return f
}

// Strings creates a string slice field.
func Strings(name string, opts ...FieldOption) Field {
	f := Field{Name: name, Type: TypeStringSlice}
	for _, opt := range opts {
		opt(&f)
	}
	return f
}

// Enum creates an enum field with allowed values.
func Enum(name string, values []string, opts ...FieldOption) Field {
	f := Field{Name: name, Type: TypeEnum, EnumVals: values}
	for _, opt := range opts {
		opt(&f)
	}
	return f
}

// Float creates a float field.
func Float(name string, opts ...FieldOption) Field {
	f := Field{Name: name, Type: TypeFloat}
	for _, opt := range opts {
		opt(&f)
	}
	return f
}

// Bool creates a boolean field.
func Bool(name string, opts ...FieldOption) Field {
	f := Field{Name: name, Type: TypeBool}
	for _, opt := range opts {
		opt(&f)
	}
	return f
}

// Validate validates a value against field constraints.
func (f *Field) Validate(value any) error {
	// Check required
	if f.Required {
		if value == nil {
			return fmt.Errorf("%s is required", f.Name)
		}
		switch f.Type {
		case TypeString, TypeEnum:
			if v, ok := value.(string); ok && v == "" {
				return fmt.Errorf("%s is required", f.Name)
			}
		case TypeStringSlice:
			if v, ok := value.([]string); ok && len(v) == 0 {
				return fmt.Errorf("%s is required", f.Name)
			}
		}
	}

	// Skip validation if nil/empty and not required
	if value == nil {
		return nil
	}

	switch f.Type {
	case TypeString:
		return f.validateString(value)
	case TypeInt:
		return f.validateInt(value)
	case TypeFloat:
		return f.validateFloat(value)
	case TypeBool:
		return f.validateBool(value)
	case TypeStringSlice:
		return f.validateStringSlice(value)
	case TypeEnum:
		return f.validateEnum(value)
	}

	return nil
}

func (f *Field) validateString(value any) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("%s must be a string", f.Name)
	}
	if v == "" {
		return nil
	}
	if f.MaxLength > 0 && len(v) > f.MaxLength {
		return fmt.Errorf("%s exceeds maximum length of %d", f.Name, f.MaxLength)
	}
	return nil
}

func (f *Field) validateInt(value any) error {
	var v float64
	switch n := value.(type) {
	case int:
		v = float64(n)
	case int64:
		v = float64(n)
	case float64:
		v = n
	default:
		return fmt.Errorf("%s must be a number", f.Name)
	}
	if f.Min != nil && v < *f.Min {
		return fmt.Errorf("%s must be at least %v", f.Name, *f.Min)
	}
	if f.Max != nil && v > *f.Max {
		return fmt.Errorf("%s must be at most %v", f.Name, *f.Max)
	}
	return nil
}

func (f *Field) validateStringSlice(value any) error {
	v, ok := value.([]string)
	if !ok {
		return fmt.Errorf("%s must be a string array", f.Name)
	}
	if f.MaxItems > 0 && len(v) > f.MaxItems {
		return fmt.Errorf("%s exceeds maximum of %d items", f.Name, f.MaxItems)
	}
	return nil
}

func (f *Field) validateEnum(value any) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("%s must be a string", f.Name)
	}
	if v == "" {
		return nil
	}
	for _, allowed := range f.EnumVals {
		if v == allowed {
			return nil
		}
	}
	return fmt.Errorf("%s must be one of: %s", f.Name, strings.Join(f.EnumVals, ", "))
}

func (f *Field) validateFloat(value any) error {
	var v float64
	switch n := value.(type) {
	case int:
		v = float64(n)
	case int64:
		v = float64(n)
	case float64:
		v = n
	case float32:
		v = float64(n)
	default:
		return fmt.Errorf("%s must be a number", f.Name)
	}
	if f.Min != nil && v < *f.Min {
		return fmt.Errorf("%s must be at least %v", f.Name, *f.Min)
	}
	if f.Max != nil && v > *f.Max {
		return fmt.Errorf("%s must be at most %v", f.Name, *f.Max)
	}
	return nil
}

func (f *Field) validateBool(value any) error {
	_, ok := value.(bool)
	if !ok {
		return fmt.Errorf("%s must be a boolean", f.Name)
	}
	return nil
}
