package modules

import "strings"

const (
	// MaxDependencyConstraintBytes bounds untrusted schema-v2 module dependency constraints.
	MaxDependencyConstraintBytes = 256
	// MaxDependencyComparators bounds parser work for one schema-v2 module dependency constraint.
	MaxDependencyComparators = 16
)

type semverIdentifier struct {
	value   string
	numeric bool
}

type strictSemVer struct {
	major      string
	minor      string
	patch      string
	prerelease []semverIdentifier
}

type versionComparator struct {
	operator string
	version  strictSemVer
}

type constraintParseFailure string

const (
	constraintInvalid             constraintParseFailure = "invalid"
	constraintTooLong             constraintParseFailure = "too_long"
	constraintTooManyComparators  constraintParseFailure = "too_many_comparators"
)

func parseStrictSemVer(value string) (strictSemVer, bool) {
	if value == "" || len(value) > maxIdentifierLength {
		return strictSemVer{}, false
	}

	coreAndPre := value
	if plus := strings.IndexByte(coreAndPre, '+'); plus >= 0 {
		if strings.IndexByte(coreAndPre[plus+1:], '+') >= 0 || !validSemVerIdentifiers(coreAndPre[plus+1:]) {
			return strictSemVer{}, false
		}
		coreAndPre = coreAndPre[:plus]
	}

	core := coreAndPre
	var prerelease []semverIdentifier
	if dash := strings.IndexByte(coreAndPre, '-'); dash >= 0 {
		pre := coreAndPre[dash+1:]
		if pre == "" {
			return strictSemVer{}, false
		}
		core = coreAndPre[:dash]
		parts := strings.Split(pre, ".")
		prerelease = make([]semverIdentifier, 0, len(parts))
		for _, part := range parts {
			if !validSemVerIdentifier(part) {
				return strictSemVer{}, false
			}
			numeric := isASCIIUnsigned(part)
			if numeric && len(part) > 1 && part[0] == '0' {
				return strictSemVer{}, false
			}
			prerelease = append(prerelease, semverIdentifier{value: part, numeric: numeric})
		}
	}

	parts := strings.Split(core, ".")
	if len(parts) != 3 {
		return strictSemVer{}, false
	}
	for _, part := range parts {
		if !isASCIIUnsigned(part) || (len(part) > 1 && part[0] == '0') {
			return strictSemVer{}, false
		}
	}

	return strictSemVer{
		major:      parts[0],
		minor:      parts[1],
		patch:      parts[2],
		prerelease: prerelease,
	}, true
}

func validSemVerIdentifiers(value string) bool {
	if value == "" {
		return false
	}
	for _, part := range strings.Split(value, ".") {
		if !validSemVerIdentifier(part) {
			return false
		}
	}
	return true
}

func validSemVerIdentifier(value string) bool {
	if value == "" {
		return false
	}
	for index := 0; index < len(value); index++ {
		char := value[index]
		if (char >= '0' && char <= '9') || (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || char == '-' {
			continue
		}
		return false
	}
	return true
}

func isASCIIUnsigned(value string) bool {
	if value == "" {
		return false
	}
	for index := 0; index < len(value); index++ {
		if value[index] < '0' || value[index] > '9' {
			return false
		}
	}
	return true
}

func compareStrictSemVer(left, right strictSemVer) int {
	if comparison := compareNumericString(left.major, right.major); comparison != 0 {
		return comparison
	}
	if comparison := compareNumericString(left.minor, right.minor); comparison != 0 {
		return comparison
	}
	if comparison := compareNumericString(left.patch, right.patch); comparison != 0 {
		return comparison
	}

	if len(left.prerelease) == 0 && len(right.prerelease) == 0 {
		return 0
	}
	if len(left.prerelease) == 0 {
		return 1
	}
	if len(right.prerelease) == 0 {
		return -1
	}

	limit := len(left.prerelease)
	if len(right.prerelease) < limit {
		limit = len(right.prerelease)
	}
	for index := 0; index < limit; index++ {
		leftID := left.prerelease[index]
		rightID := right.prerelease[index]
		if leftID.numeric && rightID.numeric {
			if comparison := compareNumericString(leftID.value, rightID.value); comparison != 0 {
				return comparison
			}
			continue
		}
		if leftID.numeric != rightID.numeric {
			if leftID.numeric {
				return -1
			}
			return 1
		}
		if leftID.value < rightID.value {
			return -1
		}
		if leftID.value > rightID.value {
			return 1
		}
	}

	if len(left.prerelease) < len(right.prerelease) {
		return -1
	}
	if len(left.prerelease) > len(right.prerelease) {
		return 1
	}
	return 0
}

func compareNumericString(left, right string) int {
	if len(left) < len(right) {
		return -1
	}
	if len(left) > len(right) {
		return 1
	}
	if left < right {
		return -1
	}
	if left > right {
		return 1
	}
	return 0
}

func parseDependencyConstraint(value string) ([]versionComparator, constraintParseFailure) {
	if len(value) == 0 || strings.TrimSpace(value) != value || strings.ContainsAny(value, "\t\r\n") {
		return nil, constraintInvalid
	}
	if len(value) > MaxDependencyConstraintBytes {
		return nil, constraintTooLong
	}

	parts := strings.Split(value, " ")
	if len(parts) > MaxDependencyComparators {
		return nil, constraintTooManyComparators
	}
	if len(parts) == 0 {
		return nil, constraintInvalid
	}

	comparators := make([]versionComparator, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			return nil, constraintInvalid
		}
		operator, operand := splitComparator(part)
		if operator == "" || operand == "" {
			return nil, constraintInvalid
		}
		version, ok := parseStrictSemVer(operand)
		if !ok {
			return nil, constraintInvalid
		}
		comparators = append(comparators, versionComparator{operator: operator, version: version})
	}
	return comparators, ""
}

func splitComparator(value string) (string, string) {
	for _, operator := range []string{">=", "<=", "=", ">", "<"} {
		if strings.HasPrefix(value, operator) {
			return operator, value[len(operator):]
		}
	}
	return "", ""
}

func matchesDependencyConstraint(version string, constraint string) (bool, bool) {
	parsedVersion, ok := parseStrictSemVer(version)
	if !ok {
		return false, false
	}
	comparators, failure := parseDependencyConstraint(constraint)
	if failure != "" {
		return false, false
	}
	for _, comparator := range comparators {
		comparison := compareStrictSemVer(parsedVersion, comparator.version)
		switch comparator.operator {
		case "=":
			if comparison != 0 {
				return false, true
			}
		case ">":
			if comparison <= 0 {
				return false, true
			}
		case ">=":
			if comparison < 0 {
				return false, true
			}
		case "<":
			if comparison >= 0 {
				return false, true
			}
		case "<=":
			if comparison > 0 {
				return false, true
			}
		default:
			return false, false
		}
	}
	return true, true
}
