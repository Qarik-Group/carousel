package bosh

type Variable struct {
	ID         string
	Name       string
	Deployment string
	Definition *VariableDefinition
}

type VariableDefinition struct {
	Name       string                 `yaml:"name" json:"name"`
	Type       string                 `yaml:"type" json:"type"`
	UpdateMode UpdateMode             `yaml:"update_mode,omitempty" json:"update_mode,omitempty"`
	Options    map[string]interface{} `yaml:"options,omitempty" json:"options,omitempty"`
}

type UpdateMode string

const (
	NoOverwrite, Overwrite, Converge UpdateMode = "no-overwrite", "overwrite", "converge"
)

func (v *VariableDefinition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// update_mode [String, optional]: Update mode to use when generating credentials.
	// Currently supported update modes are no-overwrite, overwrite, and converge. Defaults to no-overwrite
	// https://bosh.io/docs/manifest-v2/#variables

	type VariableDefinitionDefaulted VariableDefinition
	var defaults = VariableDefinitionDefaulted{
		UpdateMode: NoOverwrite,
	}

	out := defaults
	err := unmarshal(&out)
	*v = VariableDefinition(out)
	return err
}
