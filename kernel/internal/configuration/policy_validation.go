package configuration

// ValidateSettingPolicy validates one existing P02.09 scoped-setting policy
// against a governed configuration registry. It exposes the existing validation
// contract for trusted kernel integrations without creating a second policy
// authority or accepting untrusted tenant/organization identifiers.
func ValidateSettingPolicy(registry *Registry, policy SettingPolicy) error {
	return validateSettingPolicy(registry, policy)
}
