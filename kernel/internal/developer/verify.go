package developer

import (
	"context"
	"errors"
	"fmt"
	"io"
)

var errUnknownVerifyTarget = errors.New("unknown verification target")

type verificationStep struct {
	executable string
	arguments  []string
}

func runVerification(ctx context.Context, root string, environment []string, stdout io.Writer, stderr io.Writer, runner CommandRunner, target string) int {
	if target == "module-lifecycle" {
		_, _ = fmt.Fprintln(stdout, "omnexa verify module-lifecycle: N/A")
		return exitOK
	}

	steps, err := verificationSteps(target)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, "omnexa verify: unknown target")
		return exitUsage
	}
	for _, step := range steps {
		if err := runner.Run(ctx, root, environment, stdout, stderr, step.executable, step.arguments...); err != nil {
			_, _ = fmt.Fprintf(stderr, "omnexa verify %s: FAIL\n", target)
			return exitOperationFailed
		}
	}
	_, _ = fmt.Fprintf(stdout, "omnexa verify %s: PASS\n", target)
	return exitOK
}

func verificationSteps(target string) ([]verificationStep, error) {
	switch target {
	case "governance":
		return governanceSteps(), nil
	case "format", "lint", "static":
		return goQualitySteps(), nil
	case "unit":
		return []verificationStep{goStep("test", "./kernel/...")}, nil
	case "contracts":
		return packageVerifierSteps(3, 5, 6, 7, 8, 9, 10, 11), nil
	case "integration":
		return packageVerifierSteps(4, 5, 6), nil
	case "migrations":
		return packageVerifierSteps(4), nil
	case "security":
		steps := goQualitySteps()
		steps = append(steps, packageVerifierSteps(2, 3, 6, 7, 8, 10, 11)...)
		return steps, nil
	case "build":
		return []verificationStep{goStep("build", "./kernel/...")}, nil
	case "release":
		steps := goQualitySteps()
		steps = append(steps, goStep("mod", "verify"), goStep("build", "./kernel/..."))
		return steps, nil
	case "all":
		steps := governanceSteps()
		steps = append(steps, goQualitySteps()...)
		steps = append(steps, packageVerifierSteps(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)...)
		steps = append(steps, goStep("mod", "verify"), goStep("build", "./kernel/..."))
		return steps, nil
	default:
		return nil, errUnknownVerifyTarget
	}
}

func governanceSteps() []verificationStep {
	return []verificationStep{
		pythonStep("scripts/validate_governance.py"),
		pythonStep("scripts/validate_development_spec.py"),
		pythonStep("scripts/validate_operations_spec.py"),
		pythonStep("scripts/validate_freeze_review.py"),
		pythonStep("scripts/validate_p01_preparation.py"),
		pythonStep("scripts/validate_p01_package_specs.py"),
	}
}

func goQualitySteps() []verificationStep {
	return []verificationStep{bashStep("scripts/verify_go_quality.sh")}
}

func packageVerifierSteps(numbers ...int) []verificationStep {
	steps := make([]verificationStep, 0, len(numbers))
	for _, number := range numbers {
		steps = append(steps, bashStep(fmt.Sprintf("scripts/verify_p01_%02d.sh", number)))
	}
	return steps
}

func bashStep(script string) verificationStep {
	return verificationStep{executable: "bash", arguments: []string{script}}
}

func pythonStep(script string) verificationStep {
	return verificationStep{executable: "python", arguments: []string{script}}
}

func goStep(arguments ...string) verificationStep {
	return verificationStep{executable: "go", arguments: append([]string(nil), arguments...)}
}
