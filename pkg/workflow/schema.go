package workflow

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

// WorkflowDefinition represents a complete workflow definition
type WorkflowDefinition struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description,omitempty"`
	Version     string            `yaml:"version"`
	Variables   map[string]string `yaml:"variables,omitempty"`
	Operations  []OperationDefinition `yaml:"operations"`
}

// OperationDefinition defines a single operation step
type OperationDefinition struct {
	ID          string            `yaml:"id"`
	Tool        string            `yaml:"tool"`
	Command     string            `yaml:"command"`
	Parameters  map[string]string `yaml:"parameters,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty"`
	Timeout     string            `yaml:"timeout,omitempty"`
	Retry       int               `yaml:"retry,omitempty"`
}

// WorkflowParser parses YAML workflow files
type WorkflowParser struct{}

func NewWorkflowParser() *WorkflowParser {
	return &WorkflowParser{}
}

func (p *WorkflowParser) Parse(data []byte) (*WorkflowDefinition, error) {
	var workflow WorkflowDefinition
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return nil, fmt.Errorf("failed to parse workflow YAML: %w", err)
	}
	
	// Validate workflow format
	if err := p.validate(&workflow); err != nil {
		return nil, err
	}
	
	return &workflow, nil
}

func (p *WorkflowParser) validate(workflow *WorkflowDefinition) error {
	if workflow.Name == "" {
		return fmt.Errorf("workflow name is required")
	}
	
	// Check operation dependencies
	operationIDs := make(map[string]bool)
	for _, op := range workflow.Operations {
		operationIDs[op.ID] = true
	}
	
	for _, op := range workflow.Operations {
		for _, depID := range op.DependsOn {
			if !operationIDs[depID] {
				return fmt.Errorf("operation %s depends on unknown operation %s", op.ID, depID)
			}
		}
	}
	
	return nil
}