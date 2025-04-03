package configuration

import (
	"fmt"
	"os"
)

var CONFIGURATION_PATH string

func init() {
	dir, _ := os.Getwd()
	CONFIGURATION_PATH = dir + "/configurations"
}

type ConfigurationBuilder interface {
	Build(filename string) ([]byte, error)
}

type ConfigurationBuilderYml struct {
	configuration *Configuration
}

func NewConfigurationBuilderYml() *ConfigurationBuilderYml {
	return &ConfigurationBuilderYml{
		configuration: NewConfiguration(),
	}
}

func (configurationBuilder *ConfigurationBuilderYml) setName(name string) *ConfigurationBuilderYml {
	configurationBuilder.configuration.name = name

	return configurationBuilder
}

func (configurationBuilder *ConfigurationBuilderYml) setDescription(description string) *ConfigurationBuilderYml {
	configurationBuilder.configuration.description = description

	return configurationBuilder
}

func (configurationBuilder *ConfigurationBuilderYml) setResources(resources map[string]Resource) *ConfigurationBuilderYml {
	configurationBuilder.configuration.resources = resources

	return configurationBuilder
}

func (configurationBuilder *ConfigurationBuilderYml) setRelationships(relationships Relationships) *ConfigurationBuilderYml {
	configurationBuilder.configuration.relationships = relationships

	return configurationBuilder
}

func (configurationBuilder *ConfigurationBuilderYml) get() *Configuration {
	return configurationBuilder.configuration
}

func (configurationBuilder *ConfigurationBuilderYml) Build(filepath string) (*Configuration, error) {
	ymlContent, err := os.ReadFile(CONFIGURATION_PATH + "/" + filepath)

	if err != nil {
		fmt.Println("Error reading file:", err)
		return &Configuration{}, err
	}

	ymlSchema, err := NewYmlParser().Parse(string(ymlContent))

	resources := make(map[string]Resource)
	for resourceName, ymlResource := range ymlSchema.Resources {
		resource := NewResource(ymlResource)
		resources[resourceName] = *resource
	}

	fromRelationshipMap := make(map[string]FromRelationship)
	for resourceName, resource := range resources {
		relatedResources := make(map[string]Resource)
		for _, foreignKey := range resource.foreignKeys {
			relatedResources[foreignKey.foreignResource] = resources[foreignKey.foreignResource]
		}

		toRelatedResources := *NewRelatedResources(relatedResources)
		fromRelationshipMap[resourceName] = *NewFromRelationship(toRelatedResources)
	}

	toRelationshipMap := make(map[string]ToRelationship)
	relatedResources := make(map[string]map[string]Resource)
	for resourceName, resource := range resources {
		for _, foreignKey := range resource.foreignKeys {
			if _, exists := relatedResources[foreignKey.foreignResource]; !exists {
				relatedResources[foreignKey.foreignResource] = make(map[string]Resource)
			}
			relatedResources[foreignKey.foreignResource][resourceName] = resource
		}
	}
	for resourceName, relatedResources := range relatedResources {
		fromRelatedResources := *NewRelatedResources(relatedResources)
		toRelationshipMap[resourceName] = *NewToRelationship(fromRelatedResources)
	}

	relationships := *NewRelationships(fromRelationshipMap, toRelationshipMap)
	configuration := configurationBuilder.
		setName(ymlSchema.Name).
		setDescription(ymlSchema.Description).
		setResources(resources).
		setRelationships(relationships).
		get()

	return configuration, err
}
