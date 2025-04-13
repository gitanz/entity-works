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

func (configurationBuilder *ConfigurationBuilderYml) buildEntities(ymlEntitites map[string]YmlEntity) map[string]Entity {
	entities := make(map[string]Entity)
	for entityName, ymlEntity := range ymlEntitites {
		entity := *NewEntity(ymlEntity.Description)
		if entity.phases == nil {
			entity.phases = make(map[string]Phase)
		}
		entity.phases = configurationBuilder.buildPhases(ymlEntity.Phases)
		entities[entityName] = entity
	}

	return entities
}

func (configurationBuilder *ConfigurationBuilderYml) buildPhases(ymlPhases map[string]YmlPhase) map[string]Phase {
	phases := make(map[string]Phase)
	for phaseName, ymlPhase := range ymlPhases {
		phase := *NewPhase(ymlPhase.Description)
		if phase.tasks == nil {
			phase.tasks = make(map[string]Task)
		}

		phase.tasks = configurationBuilder.buildTasks(ymlPhase.Tasks)
		phases[phaseName] = phase
	}

	return phases
}

func (configurationBuilder *ConfigurationBuilderYml) buildTasks(ymlTasks map[string]YmlTask) map[string]Task {
	tasks := make(map[string]Task)

	for taskName, ymlTask := range ymlTasks {
		task := *NewTask(configurationBuilder.configuration.resources[ymlTask.Resource])

		if task.selectionCriteria == "Related" {
			relatedSelectionCriteria := NewRelatedSelectionCriteria()
			task.selectionCriteria = relatedSelectionCriteria
		}

		tasks[taskName] = task
	}

	for taskName, ymlTask := range ymlTasks {
		ymlSelectionCriteria := ymlTask.SelectionCriteria
		task := tasks[taskName]
		switch ymlSelectionCriteria.Type {

		case "Custom":
			customSelectionCriteria := NewCustomSelectionCriteria()
			customSelectionCriteria.criteria = ymlSelectionCriteria.Criteria
			task.selectionCriteria = customSelectionCriteria

		case "Index":
			indexedSelectionCriteria := NewIndexedSelectionCriteria()
			relatedTasks := []Task{}
			for _, relatedTaskName := range ymlSelectionCriteria.Tasks {
				relatedTasks = append(relatedTasks, tasks[relatedTaskName])
			}

			indexedSelectionCriteria.tasks = relatedTasks
			task.selectionCriteria = indexedSelectionCriteria

		case "Related":
			relatedSelectionCriteria := NewRelatedSelectionCriteria()
			relatedTasks := []Task{}
			for _, relatedTaskName := range ymlSelectionCriteria.Tasks {
				relatedTasks = append(relatedTasks, tasks[relatedTaskName])
			}

			relatedSelectionCriteria.tasks = relatedTasks
			task.selectionCriteria = relatedSelectionCriteria

		}

		tasks[taskName] = task
	}

	return tasks
}
