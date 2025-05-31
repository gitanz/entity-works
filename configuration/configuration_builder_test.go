package configuration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getYmlSchema() YmlSchema {
	ymlContent, _ := os.ReadFile(CONFIGURATION_PATH + "/entities_test.yml")
	ymlSchema, _ := NewYmlParser().Parse(string(ymlContent))
	return ymlSchema
}

func TestConfigurationBuilderCanBuildConfigurationFromFile(t *testing.T) {
	configurationBuilder := NewConfigurationBuilderYml()
	configuration := configurationBuilder.Build(getYmlSchema())

	assert.IsType(t, &Configuration{}, configuration)
}

func TestConfigurationBuilderCanBuildConfigurationWithResources(t *testing.T) {
	configurationBuilder := NewConfigurationBuilderYml()
	configuration := configurationBuilder.Build(getYmlSchema())

	assert.IsType(t, &Configuration{}, configuration)
	resources := configuration.resources
	assert.Contains(t, resources, "MyTestResource")
	assert.Contains(t, resources, "MyTestResource2")

	for _, resource := range resources {
		assert.IsType(t, Resource{}, resource)
	}
}

func TestConfigurationBuilderCanExtractRelationshipsFromGivenTable(t *testing.T) {
	configurationBuilder := NewConfigurationBuilderYml()
	configuration := configurationBuilder.Build(getYmlSchema())

	assert.IsType(t, Relationships{}, configuration.relationships)

	assert.Contains(t, configuration.relationships.from, "MyTestResource")
	assert.IsType(t, FromRelationship{}, configuration.relationships.from["MyTestResource"])
	assert.IsType(t, Relations{}, configuration.relationships.from["MyTestResource"].to)
	assert.Contains(t, configuration.relationships.from["MyTestResource"].to, "MyTestResource2")
	assert.IsType(t, Relation{}, configuration.relationships.from["MyTestResource"].to["MyTestResource2"])
	assert.Equal(t, "my_test_table", configuration.relationships.from["MyTestResource"].to["MyTestResource2"].fromTable)
	assert.Equal(t, "my_test_table.fk1", configuration.relationships.from["MyTestResource"].to["MyTestResource2"].fromKey)
	assert.Equal(t, "my_test_table2", configuration.relationships.from["MyTestResource"].to["MyTestResource2"].toTable)
	assert.Equal(t, "my_test_table2.id", configuration.relationships.from["MyTestResource"].to["MyTestResource2"].toKey)
	assert.Equal(t, "NORMAL", configuration.relationships.from["MyTestResource"].to["MyTestResource2"].keyType)
}

func TestConfigurationBuilderCanExtractRelationshipsToGivenTable(t *testing.T) {

	configurationBuilder := NewConfigurationBuilderYml()
	configuration := configurationBuilder.Build(getYmlSchema())

	assert.IsType(t, Relationships{}, configuration.relationships)

	assert.Contains(t, configuration.relationships.to, "MyTestResource2")
	assert.IsType(t, ToRelationship{}, configuration.relationships.to["MyTestResource2"])
	assert.IsType(t, Relations{}, configuration.relationships.to["MyTestResource2"].from)
	assert.Contains(t, configuration.relationships.to["MyTestResource2"].from, "MyTestResource")
	assert.Equal(t, "my_test_table", configuration.relationships.to["MyTestResource2"].from["MyTestResource"].fromTable)
	assert.Equal(t, "my_test_table.fk1", configuration.relationships.to["MyTestResource2"].from["MyTestResource"].fromKey)
	assert.Equal(t, "my_test_table2", configuration.relationships.to["MyTestResource2"].from["MyTestResource"].toTable)
	assert.Equal(t, "my_test_table2.id", configuration.relationships.to["MyTestResource2"].from["MyTestResource"].toKey)
	assert.Equal(t, "NORMAL", configuration.relationships.to["MyTestResource2"].from["MyTestResource"].keyType)
}

func TestConfigurationBuilderCanBuildConfigurationWithEntities(t *testing.T) {
	configurationBuilder := NewConfigurationBuilderYml()
	configuration := configurationBuilder.Build(getYmlSchema())

	assert.IsType(t, map[string]Entity{}, configuration.entities)
	subjectEntity := configuration.entities["MyTestEntity"]
	assert.Equal(t, "This is my test entity", subjectEntity.description)
	assert.IsType(t, map[string]Component{}, subjectEntity.components)
	subjectComponent := subjectEntity.components["MyTestComponent1"]
	assert.Equal(t, "This is a test component", subjectComponent.description)
	assert.IsType(t, map[string]Part{}, subjectComponent.parts)
	subjectPart := subjectComponent.parts["PartA"]

	assert.IsType(t, Resource{}, subjectPart.resource)
	assert.Equal(t, "my_test_table", subjectPart.resource.tableName)
}

func TestEntityFromConfigurationBuiltUsingBuilderContainsPhasesAndTasks(t *testing.T) {
	configurationBuilder := NewConfigurationBuilderYml()
	configuration := configurationBuilder.Build(getYmlSchema())

	assert.IsType(t, map[string]Entity{}, configuration.entities)
	subjectEntity := configuration.entities["MyTestEntity"]
	assert.Equal(t, "This is my test entity", subjectEntity.description)
	assert.IsType(t, map[string]Component{}, subjectEntity.components)
	subjectPhase := subjectEntity.components["MyTestComponent1"]
	assert.Equal(t, "This is a test component", subjectPhase.description)
	assert.IsType(t, map[string]Part{}, subjectPhase.parts)
	subjectPart := subjectPhase.parts["PartA"]

	assert.IsType(t, Resource{}, subjectPart.resource)
	assert.Equal(t, "my_test_table", subjectPart.resource.tableName)
}

func TestTasksInConfigurationBuiltUsingBuilderContainsResourceAndSelectionCriteria(t *testing.T) {
	configurationBuilder := NewConfigurationBuilderYml()
	configuration := configurationBuilder.Build(getYmlSchema())

	assert.IsType(t, map[string]Entity{}, configuration.entities)
	subjectEntity := configuration.entities["MyTestEntity"]
	subjectPhase := subjectEntity.components["MyTestComponent1"]
	partA := subjectPhase.parts["PartA"]

	assert.IsType(t, Resource{}, partA.resource)
	assert.Nil(t, partA.selectionCriteria)

	partB := subjectPhase.parts["PartB"]
	assert.IsType(t, Resource{}, partB.resource)
	assert.NotNil(t, partB.selectionCriteria)
	assert.IsType(t, &CustomSelectionCriteria{}, partB.selectionCriteria)

	partC := subjectPhase.parts["PartC"]
	assert.IsType(t, Resource{}, partC.resource)
	assert.NotNil(t, partC.selectionCriteria)
	assert.IsType(t, &IndexedSelectionCriteria{}, partC.selectionCriteria)

	partD := subjectPhase.parts["PartD"]
	assert.IsType(t, Resource{}, partD.resource)
	assert.NotNil(t, partD.selectionCriteria)
	assert.IsType(t, &RelatedSelectionCriteria{}, partD.selectionCriteria)
}
