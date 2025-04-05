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
		relations := make(map[string]Relation)
		for _, foreignKey := range resource.foreignKeys {
			relations[foreignKey.foreignResource] = *NewRelation(
				resources[resourceName].tableName,
				foreignKey.key,
				resources[foreignKey.foreignResource].tableName,
				foreignKey.foreignKey,
				foreignKey.keyType,
			)
		}

		toRelations := *NewRelations(relations)
		fromRelationshipMap[resourceName] = *NewFromRelationship(toRelations)
	}

	toRelationshipMap := make(map[string]ToRelationship)
	relations := make(map[string]map[string]Relation)
	for resourceName, resource := range resources {
		for _, foreignKey := range resource.foreignKeys {
			if _, exists := relations[foreignKey.foreignResource]; !exists {
				relations[foreignKey.foreignResource] = make(map[string]Relation)
			}
			relations[foreignKey.foreignResource][resourceName] = *NewRelation(
				resources[resourceName].tableName,
				foreignKey.key,
				resources[foreignKey.foreignResource].tableName,
				foreignKey.foreignKey,
				foreignKey.keyType,
			)
		}
	}
	for resourceName, relations := range relations {
		fromRelations := *NewRelations(relations)
		toRelationshipMap[resourceName] = *NewToRelationship(fromRelations)
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
