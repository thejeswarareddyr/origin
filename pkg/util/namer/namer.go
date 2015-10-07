package namer

import (
	"fmt"
	"hash/fnv"

	"k8s.io/kubernetes/pkg/util"
)

// GetName returns a name given a base ("deployment-5") and a suffix ("deploy")
// It will first attempt to join them with a dash. If the resulting name is longer
// than maxLength: if the suffix is too long, it will truncate the base name and add
// an 8-character hash of the [base]-[suffix] string.  If the suffix is not too long,
// it will truncate the base, add the hash of the base and return [base]-[hash]-[suffix]
func GetName(base, suffix string, maxLength int) string {
	name := fmt.Sprintf("%s-%s", base, suffix)
	if len(name) > maxLength {
		baseLength := maxLength - 10 /*length of -hash-*/ - len(suffix)

		// if the suffix is too long, ignore it
		if baseLength < 0 {
			prefix := base[0:min(len(base), maxLength-9)]
			// Calculate hash on initial base-suffix string
			return fmt.Sprintf("%s-%s", prefix, hash(name))
		}

		prefix := base[0:baseLength]
		// Calculate hash on initial base-suffix string
		return fmt.Sprintf("%s-%s-%s", prefix, hash(base), suffix)
	}

	return name
}

// GetPodName calls GetName with the length restriction for pods
func GetPodName(base, suffix string) string {
	return GetName(base, suffix, util.DNS1123SubdomainMaxLength)
}

// LimitLength returns a string of at most maxLength characters based on name.
// It returns name unchanged if it is short enough, otherwise it returns a
// prefix of name plus a suffix computed from a hash of name.
func LimitLength(name string, maxLength int) string {
	if maxLength <= 0 {
		return ""
	}
	if len(name) <= maxLength {
		return name
	}
	suffix := hash(name)
	prefix := name[:min(max(0, maxLength-len(suffix)), len(name))]
	shortName := fmt.Sprintf("%s%s", prefix, suffix)[:maxLength]
	return shortName
}

// max returns the greater of its 2 inputs
func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

// min returns the lesser of its 2 inputs
func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

// hash calculates the hexadecimal representation (8-chars)
// of the hash of the passed in string using the FNV-a algorithm
func hash(s string) string {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	intHash := hash.Sum32()
	result := fmt.Sprintf("%08x", intHash)
	return result
}
