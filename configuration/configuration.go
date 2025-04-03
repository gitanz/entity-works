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

type RelatedResources map[string]Resource

func NewRelatedResources(resourceMap map[string]Resource) *RelatedResources {
	var relatedResources RelatedResources = resourceMap
	return &relatedResources
}

type FromRelationship struct {
	to RelatedResources
}

func NewFromRelationship(relatedResources RelatedResources) *FromRelationship {
	fromRelationship := &FromRelationship{
		to: relatedResources,
	}

	return fromRelationship
}

type ToRelationship struct {
	from RelatedResources
}

func NewToRelationship(relatedResources RelatedResources) *ToRelationship {
	toRelationship := &ToRelationship{
		from: relatedResources,
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

type Configuration struct {
	name          string
	description   string
	resources     map[string]Resource
	relationships Relationships
}

func NewConfiguration() *Configuration {
	configuration := &Configuration{}
	configuration.resources = make(map[string]Resource)

	return configuration
}
