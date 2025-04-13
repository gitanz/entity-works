package configuration

import (
	"github.com/goccy/go-yaml"
)

type YmlForeignKey struct {
	Type         string `yaml:"Type"`
	Key          string `yaml:"Key"`
	ResourceName string `yaml:"ResourceName"`
	ForeignKey   string `yaml:"ForeignKey"`
}

type YmlResource struct {
	TableName     string              `yaml:"TableName"`
	PrimaryKey    []string            `yaml:"PrimaryKey,omitempty"`
	AutoIncrement bool                `yaml:"AutoIncrement,omitempty"`
	Index         map[string][]string `yaml:"Index,omitempty"`
	ForeignKeys   []YmlForeignKey     `yaml:"ForeignKeys,omitempty"`
}

type YmlSelectionCriteria struct {
	Type     string   `yaml:"Type"`
	Criteria string   `yaml:"Criteria,omitempty"`
	Tasks    []string `yaml:"Tasks,omitempty"`
	Index    string   `yaml:"Index,omitempty"`
}

type YmlTask struct {
	Resource          string               `yaml:"Resource"`
	Shares            bool                 `yaml:"Shares,omitempty"`
	SelectionCriteria YmlSelectionCriteria `yaml:"SelectionCriteria"`
}

type YmlPhase struct {
	Description string             `yaml:"Description"`
	Tasks       map[string]YmlTask `yaml:"Tasks"`
}

type YmlEntity struct {
	Description string              `yaml:"Description"`
	Phases      map[string]YmlPhase `yaml:"Phases,omitempty"`
}

type YmlSchema struct {
	Name        string                 `yaml:"Name"`
	Description string                 `yaml:"Description"`
	Resources   map[string]YmlResource `yaml:"Resources"`
	Entities    map[string]YmlEntity   `yaml:"Entities,omitempty"`
}

func NewYmlSchema() *YmlSchema {
	return &YmlSchema{}
}

type Parser interface {
	Parse(ymlConfiguration string) YmlSchema
}

type YmlParser struct {
}

func (ymlParser YmlParser) Parse(ymlConfiguration string) (YmlSchema, error) {
	ymlDefinition := NewYmlSchema()
	err := yaml.Unmarshal([]byte(ymlConfiguration), ymlDefinition)
	return *ymlDefinition, err
}

func NewYmlParser() *YmlParser {
	return &YmlParser{}
}
