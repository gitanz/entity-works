package configuration

import (
	"os"
)

var CONFIGURATION_PATH string

func init() {
	dir, _ := os.Getwd()
	CONFIGURATION_PATH = dir + "/configurations"
}

type Builder interface {
	Build(ymlSchema YmlSchema) []byte
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

func (configurationBuilder *ConfigurationBuilderYml) setEntities(entities map[string]Entity) *ConfigurationBuilderYml {
	configurationBuilder.configuration.entities = entities

	return configurationBuilder
}

func (configurationBuilder *ConfigurationBuilderYml) get() *Configuration {
	return configurationBuilder.configuration
}

func (configurationBuilder *ConfigurationBuilderYml) Build(ymlSchema YmlSchema) *Configuration {
	configuration := configurationBuilder.
		setName(ymlSchema.Name).
		setDescription(ymlSchema.Description)

	resources := configurationBuilder.buildResources(ymlSchema.Resources)
	configuration.setResources(resources)

	relationships := configurationBuilder.buildRelationships(resources)
	configuration.setRelationships(relationships)

	entities := configurationBuilder.buildEntities(ymlSchema.Entities)
	configuration.setEntities(entities)

	return configuration.get()
}

func (configurationBuilder *ConfigurationBuilderYml) buildResources(ymlResources map[string]YmlResource) map[string]Resource {
	resources := make(map[string]Resource)
	for resourceName, ymlResource := range ymlResources {
		resource := NewResource(ymlResource)
		resources[resourceName] = *resource
	}

	return resources
}

func (configurationBuilder *ConfigurationBuilderYml) buildRelationships(resources map[string]Resource) Relationships {
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

	return *NewRelationships(fromRelationshipMap, toRelationshipMap)
}

func (configurationBuilder *ConfigurationBuilderYml) buildEntities(ymlEntities map[string]YmlEntity) map[string]Entity {
	entities := make(map[string]Entity)
	for entityName, ymlEntity := range ymlEntities {
		entity := *NewEntity(ymlEntity.Description)
		if entity.components == nil {
			entity.components = make(map[string]Component)
		}
		entity.components = configurationBuilder.buildComponents(ymlEntity.Components)
		entities[entityName] = entity
	}

	return entities
}

func (configurationBuilder *ConfigurationBuilderYml) buildComponents(ymlComponents map[string]YmlComponent) map[string]Component {
	components := make(map[string]Component)
	for componentName, ymlComponent := range ymlComponents {
		component := *NewComponent(ymlComponent.Description)
		if component.parts == nil {
			component.parts = make(map[string]Part)
		}

		component.parts = configurationBuilder.buildParts(ymlComponent.Parts)
		components[componentName] = component
	}

	return components
}

func (configurationBuilder *ConfigurationBuilderYml) buildParts(ymlParts map[string]YmlPart) map[string]Part {
	parts := make(map[string]Part)

	for partName, ymlPart := range ymlParts {
		part := *NewPart(configurationBuilder.configuration.resources[ymlPart.Resource])

		if part.selectionCriteria == "Related" {
			relatedSelectionCriteria := NewRelatedSelectionCriteria()
			part.selectionCriteria = relatedSelectionCriteria
		}

		parts[partName] = part
	}

	for partName, ymlPart := range ymlParts {
		ymlSelectionCriteria := ymlPart.SelectionCriteria
		part := parts[partName]
		switch ymlSelectionCriteria.Type {

		case "Custom":
			customSelectionCriteria := NewCustomSelectionCriteria()
			customSelectionCriteria.criteria = ymlSelectionCriteria.Criteria
			part.selectionCriteria = customSelectionCriteria

		case "Index":
			indexedSelectionCriteria := NewIndexedSelectionCriteria()
			relatedParts := []Part{}
			for _, relatedPartName := range ymlSelectionCriteria.Parts {
				relatedParts = append(relatedParts, parts[relatedPartName])
			}

			indexedSelectionCriteria.parts = relatedParts
			part.selectionCriteria = indexedSelectionCriteria

		case "Related":
			relatedSelectionCriteria := NewRelatedSelectionCriteria()
			relatedParts := []Part{}
			for _, relatedPartName := range ymlSelectionCriteria.Parts {
				relatedParts = append(relatedParts, parts[relatedPartName])
			}

			relatedSelectionCriteria.parts = relatedParts
			part.selectionCriteria = relatedSelectionCriteria

		}

		parts[partName] = part
	}

	return parts
}
