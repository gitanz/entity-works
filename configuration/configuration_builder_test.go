package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurationBuilderCanBuildConfigurationFromFile(t *testing.T) {
	configurationBuilder := NewConfigurationBuilderYml()
	configuration, _ := configurationBuilder.Build("entities_test.yml")

	assert.IsType(t, &Configuration{}, configuration)
}

func TestConfigurationBuilderCanBuildConfigurationWithResources(t *testing.T) {
	configurationBuilder := NewConfigurationBuilderYml()
	configuration, _ := configurationBuilder.Build("entities_test.yml")

	assert.IsType(t, &Configuration{}, configuration)
	resources := configuration.resources
	assert.Contains(t, resources, "MyTestResource")
	assert.Contains(t, resources, "MyTestResource2")

	for _, resource := range resources {
		assert.IsType(t, Resource{}, resource)
	}
}

func TestConfigurationBuilderCanExtractRelationshipsBetweenResources(t *testing.T) {
	configurationBuilder := NewConfigurationBuilderYml()
	configuration, _ := configurationBuilder.Build("entities_test.yml")

	assert.IsType(t, Relationships{}, configuration.relationships)

	assert.Contains(t, configuration.relationships.from, "MyTestResource")
	assert.IsType(t, FromRelationship{}, configuration.relationships.from["MyTestResource"])
	assert.IsType(t, RelatedResources{}, configuration.relationships.from["MyTestResource"].to)
	assert.Contains(t, configuration.relationships.from["MyTestResource"].to, "MyTestResource2")

	assert.Contains(t, configuration.relationships.to, "MyTestResource2")
	assert.IsType(t, ToRelationship{}, configuration.relationships.to["MyTestResource2"])
	assert.IsType(t, RelatedResources{}, configuration.relationships.to["MyTestResource2"].from)
	assert.Contains(t, configuration.relationships.to["MyTestResource2"].from, "MyTestResource")
}
