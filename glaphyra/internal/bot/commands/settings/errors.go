package settings

type ValidationError struct{}

func (e *ValidationError) Error() string {
	return "settings validation failed"
}
