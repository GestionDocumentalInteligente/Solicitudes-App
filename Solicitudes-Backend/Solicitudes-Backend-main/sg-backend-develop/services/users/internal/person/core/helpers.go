package person

import "regexp"

// Helper function to clean Cuil (remove any non-numeric characters)
func keepOnlyNumbers(cuil string) string {
	// Regular expression to remove all non-numeric characters
	re := regexp.MustCompile("[^0-9]")
	return re.ReplaceAllString(cuil, "")
}
