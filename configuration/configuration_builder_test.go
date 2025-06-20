package configuration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getYmlSchema() YmlSchema {
	ymlContent, _ := os.ReadFile(Path + "/entities_test.yml")
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
	assert.IsType(t, map[string]Element{}, subjectComponent.elements)
	subjectElement := subjectComponent.elements["ElementA"]
	assert.IsType(t, Resource{}, subjectElement.resource)
	assert.Equal(t, "my_test_table", subjectElement.resource.tableName)
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
	assert.IsType(t, map[string]Element{}, subjectPhase.elements)
	subjectElement := subjectPhase.elements["ElementA"]

	assert.IsType(t, Resource{}, subjectElement.resource)
	assert.Equal(t, "my_test_table", subjectElement.resource.tableName)
}

func TestTasksInConfigurationBuiltUsingBuilderContainsResourceAndSelectionCriteria(t *testing.T) {
	configurationBuilder := NewConfigurationBuilderYml()
	configuration := configurationBuilder.Build(getYmlSchema())

	assert.IsType(t, map[string]Entity{}, configuration.entities)
	subjectEntity := configuration.entities["MyTestEntity"]
	subjectPhase := subjectEntity.components["MyTestComponent1"]
	elementA := subjectPhase.elements["ElementA"]

	assert.IsType(t, Resource{}, elementA.resource)
	assert.Nil(t, elementA.selectionCriteria)

	elementB := subjectPhase.elements["ElementB"]
	assert.IsType(t, Resource{}, elementB.resource)
	assert.NotNil(t, elementB.selectionCriteria)
	assert.IsType(t, &CustomSelectionCriteria{}, elementB.selectionCriteria)

	elementC := subjectPhase.elements["ElementC"]
	assert.IsType(t, Resource{}, elementC.resource)
	assert.NotNil(t, elementC.selectionCriteria)
	assert.IsType(t, &IndexedSelectionCriteria{}, elementC.selectionCriteria)

	elementD := subjectPhase.elements["ElementD"]
	assert.IsType(t, Resource{}, elementD.resource)
	assert.NotNil(t, elementD.selectionCriteria)
	assert.IsType(t, &RelatedSelectionCriteria{}, elementD.selectionCriteria)
}
