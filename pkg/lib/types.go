package lib

// Package example provides types for managing input and output data.
//
// This package defines default types for two main data structures: MInput and MOutput.

type InputOpts func(*MInput)

func DescriptionInput(description string) InputOpts {
	return func(input *MInput) {
		input.Description = description
	}
}

func DefaultInput(defaultValue interface{}) InputOpts {
	return func(input *MInput) {
		input.Default = defaultValue
	}
}

func RequiredInput() InputOpts {
	return func(input *MInput) {
		input.Required = true
	}
}

func ConstraintsInput(constraints ...Constraint) InputOpts {
	return func(input *MInput) {
		input.Constraints = constraints
	}
}

func InputString(name string, opts ...InputOpts) MInput {
	minput := MInput{
		Name:     name,
		MimeType: "text/plain",
	}

	for _, opt := range opts {
		opt(&minput)
	}

	return minput
}

func InputImage(name string, opts ...InputOpts) MInput {
	minput := MInput{
		Name:     name,
		MimeType: "image/png",
	}

	for _, opt := range opts {
		opt(&minput)
	}

	return minput
}

func InputInt(name string, opts ...InputOpts) MInput {
	minput := MInput{
		Name:     name,
		MimeType: "text/plain",
	}

	for _, opt := range opts {
		opt(&minput)
	}

	return minput
}

type OutputOpts func(*MOutput)

func DescriptionOutput(description string) OutputOpts {
	return func(output *MOutput) {
		output.Description = description
	}
}

func MimeTypeOutput(mimeType string) OutputOpts {
	return func(output *MOutput) {
		output.MimeType = mimeType
	}
}

func OutputString(name string, opts ...OutputOpts) MOutput {
	output := MOutput{
		Name:     name,
		MimeType: "text/plain",
	}

	for _, opt := range opts {
		opt(&output)
	}

	return output
}

func OutputImage(name string, opts ...OutputOpts) MOutput {
	output := MOutput{
		Name:     name,
		MimeType: "image/png",
	}

	for _, opt := range opts {
		opt(&output)
	}

	return output
}

func OutputInt(name string, opts ...OutputOpts) MOutput {
	output := MOutput{
		Name:     name,
		MimeType: "text/plain",
	}

	for _, opt := range opts {
		opt(&output)
	}

	return output
}
