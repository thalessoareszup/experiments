package workflow

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Workflow represents a YAML-defined workflow file.
// Example YAML:
//
// name: My workflow
// steps:
//   - id: setup
//     name: Setup project
//     description: ...
//
// Steps must have unique IDs and at least one entry.
type Workflow struct {
	Name  string         `yaml:"name"`
	Steps []WorkflowStep `yaml:"steps"`
}

// WorkflowStep is a single step within a workflow.
type WorkflowStep struct {
	ID                   string `yaml:"id"`
	Name                 string `yaml:"name"`
	Description          string `yaml:"description"`
	RequiresConfirmation bool   `yaml:"requires_confirmation"`
}

// Load loads and validates a workflow from the given YAML file path.
// It does not apply any default path logic; callers are responsible for
// deciding where the workflow file lives.
func Load(path string) (*Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read workflow: %w", err)
	}

	var wf Workflow
	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("parse workflow yaml: %w", err)
	}

	if err := wf.Validate(); err != nil {
		return nil, err
	}

	return &wf, nil
}

// Validate performs basic validation on the workflow definition.
func (w *Workflow) Validate() error {
	if len(w.Steps) == 0 {
		return fmt.Errorf("workflow must have at least one step")
	}

	seen := make(map[string]struct{}, len(w.Steps))
	for i, s := range w.Steps {
		if s.ID == "" {
			return fmt.Errorf("step %d is missing id", i)
		}
		if _, ok := seen[s.ID]; ok {
			return fmt.Errorf("duplicate step id: %s", s.ID)
		}
		seen[s.ID] = struct{}{}
	}

	return nil
}
