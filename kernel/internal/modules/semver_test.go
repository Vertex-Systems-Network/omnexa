package modules

import "testing"

func TestStrictSemVerAcceptsAndRejectsSemVer20Forms(t *testing.T) {
	valid := []string{
		"0.0.0",
		"1.2.3",
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-0.3.7",
		"1.0.0-x.7.z.92",
		"1.0.0+build.01",
		"999999999999999999999999.0.0",
	}
	for _, value := range valid {
		if _, ok := parseStrictSemVer(value); !ok {
			t.Fatalf("expected strict SemVer to accept %q", value)
		}
	}

	invalid := []string{
		"1.0",
		"01.0.0",
		"1.01.0",
		"1.0.01",
		"1.0.0-01",
		"1.0.0-alpha..1",
		"1.0.0-",
		"1.0.0+",
		"1.0.0+a_b",
		"v1.0.0",
	}
	for _, value := range invalid {
		if _, ok := parseStrictSemVer(value); ok {
			t.Fatalf("expected strict SemVer to reject %q", value)
		}
	}
}

func TestStrictSemVerPrecedenceAndBuildMetadataEquality(t *testing.T) {
	ordered := []string{
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-alpha.beta",
		"1.0.0-beta",
		"1.0.0-beta.2",
		"1.0.0-beta.11",
		"1.0.0-rc.1",
		"1.0.0",
	}
	for index := 0; index < len(ordered)-1; index++ {
		left, leftOK := parseStrictSemVer(ordered[index])
		right, rightOK := parseStrictSemVer(ordered[index+1])
		if !leftOK || !rightOK || compareStrictSemVer(left, right) >= 0 {
			t.Fatalf("expected %q < %q", ordered[index], ordered[index+1])
		}
	}

	left, _ := parseStrictSemVer("1.2.3+build.1")
	right, _ := parseStrictSemVer("1.2.3+build.999")
	if compareStrictSemVer(left, right) != 0 {
		t.Fatal("build metadata must not affect precedence equality")
	}
}

func TestDependencyConstraintGrammarIsBoundedAndExplicit(t *testing.T) {
	for _, value := range []string{
		"=1.2.3",
		">=1.2.3 <2.0.0",
		">1.0.0 <=1.9.9",
	} {
		if _, failure := parseDependencyConstraint(value); failure != "" {
			t.Fatalf("expected valid constraint %q, got %q", value, failure)
		}
	}

	invalid := []string{
		"1.2.3",
		"^1.2.3",
		"~1.2.3",
		">=1.2.3  <2.0.0",
		" >=1.2.3",
		">=1.2.3 ",
		">=1.2.3 || <2.0.0",
		">=1.2.3, <2.0.0",
		">=1.0.0-01",
	}
	for _, value := range invalid {
		if _, failure := parseDependencyConstraint(value); failure == "" {
			t.Fatalf("expected invalid constraint %q", value)
		}
	}

	tooMany := "=1.0.0"
	for index := 0; index < MaxDependencyComparators; index++ {
		tooMany += " =1.0.0"
	}
	if _, failure := parseDependencyConstraint(tooMany); failure != constraintTooManyComparators {
		t.Fatalf("expected comparator bound failure, got %q", failure)
	}

	tooLong := ">=" + string(make([]byte, MaxDependencyConstraintBytes))
	if _, failure := parseDependencyConstraint(tooLong); failure != constraintTooLong {
		t.Fatalf("expected length bound failure, got %q", failure)
	}
}

func TestDependencyConstraintComparisonUsesSemVerPrecedence(t *testing.T) {
	if matched, valid := matchesDependencyConstraint("1.5.0", ">=1.2.0 <2.0.0"); !valid || !matched {
		t.Fatal("expected 1.5.0 to satisfy bounded range")
	}
	if matched, valid := matchesDependencyConstraint("2.0.0", ">=1.2.0 <2.0.0"); !valid || matched {
		t.Fatal("expected 2.0.0 to fail upper bound")
	}
	if matched, valid := matchesDependencyConstraint("1.0.0+one", "=1.0.0+two"); !valid || !matched {
		t.Fatal("build metadata must not affect equality")
	}
}
