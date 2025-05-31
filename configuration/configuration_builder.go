package configuration

import (
	"os"
)

var Path string

func init() {
	dir, _ := os.Getwd()
	Path = dir + "/configurations"
}

type Builder interface {
	Build(ymlSchema YmlSchema) []byte
}

type BuilderYml struct {
	configuration *Configuration
}

func NewConfigurationBuilderYml() *BuilderYml {
	return &BuilderYml{
		configuration: NewConfiguration(),
	}
}

func (configurationBuilder *BuilderYml) setName(name string) *BuilderYml {
	configurationBuilder.configuration.name = name

	return configurationBuilder
}

func (configurationBuilder *BuilderYml) setDescription(description string) *BuilderYml {
	configurationBuilder.configuration.description = description

	return configurationBuilder
}

func (configurationBuilder *BuilderYml) setResources(resources map[string]Resource) *BuilderYml {
	configurationBuilder.configuration.resources = resources

	return configurationBuilder
}

func (configurationBuilder *BuilderYml) setRelationships(relationships Relationships) *BuilderYml {
	configurationBuilder.configuration.relationships = relationships

	return configurationBuilder
}

func (configurationBuilder *BuilderYml) setEntities(entities map[string]Entity) *BuilderYml {
	configurationBuilder.configuration.entities = entities

	return configurationBuilder
}

func (configurationBuilder *BuilderYml) get() *Configuration {
	return configurationBuilder.configuration
}

func (configurationBuilder *BuilderYml) Build(ymlSchema YmlSchema) *Configuration {
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

func (configurationBuilder *BuilderYml) buildResources(ymlResources map[string]YmlResource) map[string]Resource {
	resources := make(map[string]Resource)
	for resourceName, ymlResource := range ymlResources {
		resource := NewResource(ymlResource)
		resources[resourceName] = *resource
	}

	return resources
}

func (configurationBuilder *BuilderYml) buildRelationships(resources map[string]Resource) Relationships {
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

func (configurationBuilder *BuilderYml) buildEntities(ymlEntities map[string]YmlEntity) map[string]Entity {
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

func (configurationBuilder *BuilderYml) buildComponents(ymlComponents map[string]YmlComponent) map[string]Component {
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

func (configurationBuilder *BuilderYml) buildParts(ymlParts map[string]YmlPart) map[string]Part {
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
			var relatedParts []Part
			for _, relatedPartName := range ymlSelectionCriteria.Parts {
				relatedParts = append(relatedParts, parts[relatedPartName])
			}

			indexedSelectionCriteria.parts = relatedParts
			part.selectionCriteria = indexedSelectionCriteria

		case "Related":
			relatedSelectionCriteria := NewRelatedSelectionCriteria()
			var relatedParts []Part
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
