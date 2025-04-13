package configuration

type ForeignKey struct {
	keyType         string
	key             string
	foreignResource string
	foreignKey      string
}

func NewForeignKey(yamlForeignKey YmlForeignKey) *ForeignKey {
	return &ForeignKey{
		keyType:         yamlForeignKey.Type,
		key:             yamlForeignKey.Key,
		foreignResource: yamlForeignKey.ResourceName,
		foreignKey:      yamlForeignKey.ForeignKey,
	}
}

type Resource struct {
	tableName     string
	primaryKey    []string
	autoIncrement bool
	index         map[string][]string
	foreignKeys   []ForeignKey
}

func NewResource(ymlResource YmlResource) *Resource {
	resource := &Resource{
		tableName:     ymlResource.TableName,
		primaryKey:    ymlResource.PrimaryKey,
		autoIncrement: ymlResource.AutoIncrement,
		index:         ymlResource.Index,
	}

	resource.foreignKeys = []ForeignKey{}
	for _, ymlForeignKey := range ymlResource.ForeignKeys {
		resource.foreignKeys = append(resource.foreignKeys, *NewForeignKey(ymlForeignKey))
	}

	return resource
}

type Relation struct {
	fromTable string
	fromKey   string
	toTable   string
	toKey     string
	keyType   string
}

func NewRelation(
	fromTable string,
	fromKey string,
	toTable string,
	toKey string,
	keyType string,
) *Relation {
	return &Relation{
		fromTable: fromTable,
		fromKey:   fromKey,
		toTable:   toTable,
		toKey:     toKey,
		keyType:   keyType,
	}
}

type Relations map[string]Relation

func NewRelations(relationMap map[string]Relation) *Relations {
	var relations Relations = relationMap
	return &relations
}

type FromRelationship struct {
	to Relations
}

func NewFromRelationship(relations Relations) *FromRelationship {
	fromRelationship := &FromRelationship{
		to: relations,
	}

	return fromRelationship
}

type ToRelationship struct {
	from Relations
}

func NewToRelationship(relations Relations) *ToRelationship {
	toRelationship := &ToRelationship{
		from: relations,
	}

	return toRelationship
}

type Relationships struct {
	from map[string]FromRelationship
	to   map[string]ToRelationship
}

func NewRelationships(from map[string]FromRelationship, to map[string]ToRelationship) *Relationships {
	relationships := &Relationships{
		from: from,
		to:   to,
	}

	return relationships
}

type SelectionCriteria interface {
}

type CustomSelectionCriteria struct {
	criteria string
}

func NewCustomSelectionCriteria() *CustomSelectionCriteria {
	return &CustomSelectionCriteria{}
}

type RelatedSelectionCritiera struct {
	tasks []Task
}

func NewRelatedSelectionCriteria() *RelatedSelectionCritiera {
	return &RelatedSelectionCritiera{}
}

type IndexedSelectionCriteria struct {
	tasks []Task
}

func NewIndexedSelectionCriteria() *IndexedSelectionCriteria {
	return &IndexedSelectionCriteria{}
}

type NewIndexedSelectionCritier struct {
}

type Task struct {
	resource          Resource
	selectionCriteria SelectionCriteria
}

func NewTask(resource Resource) *Task {
	return &Task{
		resource: resource,
	}
}

type Phase struct {
	description string
	tasks       map[string]Task
}

func NewPhase(description string) *Phase {
	phase := &Phase{}
	phase.description = description

	return phase
}

type Entity struct {
	description string
	phases      map[string]Phase
}

func NewEntity(description string) *Entity {
	entity := &Entity{}
	entity.description = description

	return entity
}

type Configuration struct {
	name          string
	description   string
	resources     map[string]Resource
	relationships Relationships
	entities      map[string]Entity
}

func NewConfiguration() *Configuration {
	configuration := &Configuration{}
	configuration.resources = make(map[string]Resource)

	return configuration
}
