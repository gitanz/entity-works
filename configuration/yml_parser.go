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
	Parts    []string `yaml:"Parts,omitempty"`
	Index    string   `yaml:"Index,omitempty"`
}

type YmlPart struct {
	Resource          string               `yaml:"Resource"`
	Shares            bool                 `yaml:"Shares,omitempty"`
	SelectionCriteria YmlSelectionCriteria `yaml:"SelectionCriteria"`
}

type YmlComponent struct {
	Description string             `yaml:"Description"`
	Parts       map[string]YmlPart `yaml:"Parts,omitempty"`
}

type YmlEntity struct {
	Description string                  `yaml:"Description"`
	Components  map[string]YmlComponent `yaml:"Components,omitempty"`
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
