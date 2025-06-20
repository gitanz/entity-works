package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanInstantiateConfigurationBuilder(t *testing.T) {
	ymlParser := NewYmlParser()
	assert.NotNil(t, ymlParser)
}

func TestInvalidYmlMapInBuildThrowsAnException(t *testing.T) {
	ymlParser := NewYmlParser()
	_, err := ymlParser.Parse("Invalid YML")
	assert.NotNil(t, err)
}

func TestConfigurationCanBuild(t *testing.T) {
	ymlParser := NewYmlParser()
	ymlSchema, _ := ymlParser.Parse(`Name: Example`)
	assert.IsType(t, ymlSchema, YmlSchema{})
}

func TestConfigurationBuilderParsesYmlFile(t *testing.T) {
	ymlParser := NewYmlParser()
	ymlSchema, err := ymlParser.Parse(
		`
    Name: Example
    Description: Example YML configuration
  `)
	assert.Nil(t, err)
	assert.Equal(t, "Example", ymlSchema.Name)
	assert.Equal(t, "Example YML configuration", ymlSchema.Description)
}

func TestConfigurationBuilderParsesYmlFileWithMinimalResourcesSection(t *testing.T) {
	ymlParser := NewYmlParser()
	ymlSchema, err := ymlParser.Parse(
		`
    Name: Example
    Description: Example YML configuration
    Resources:
      MyTestResource:
        TableName: my_test_table
  `)
	assert.Nil(t, err)
	assert.Equal(t, "Example", ymlSchema.Name)
	assert.Equal(t, "Example YML configuration", ymlSchema.Description)
	assert.Equal(t, "my_test_table", ymlSchema.Resources["MyTestResource"].TableName)
}

func TestConfigurationBuilderParsesYmlWithCompleteResourcesSection(t *testing.T) {
	ymlParser := NewYmlParser()
	ymlSchema, err := ymlParser.Parse(
		`
    Name: Example
    Description: Example YML configuration
    Resources:
      MyTestResource:
        TableName: my_test_table
        PrimaryKey:
          - my_test_table.id
        AutoIncrement: true
        Index:
          IDX1:
            - my_test_table.idx1
          IDX2:
            - my_test_table.idx1
            - my_test_table.idx2

  `)
	assert.Nil(t, err)
	assert.Equal(t, "Example", ymlSchema.Name)
	assert.Equal(t, "Example YML configuration", ymlSchema.Description)
	assert.Equal(t, "my_test_table", ymlSchema.Resources["MyTestResource"].TableName)
	assert.ElementsMatch(t, []string{"my_test_table.id"}, ymlSchema.Resources["MyTestResource"].PrimaryKey)
	assert.True(t, ymlSchema.Resources["MyTestResource"].AutoIncrement)
	assert.ElementsMatch(t, []string{"my_test_table.idx1"}, ymlSchema.Resources["MyTestResource"].Index["IDX1"])
}

func TestConfigurationBuilderParsesYmlWithEntities(t *testing.T) {
	ymlParser := NewYmlParser()
	ymlSchema, err := ymlParser.Parse(
		`
    Name: Example
    Description: Example YML configuration
    Entities:
      MyTestEntity:
        Description: This is my test entity
  `)
	assert.Nil(t, err)
	assert.IsType(t, YmlEntity{}, ymlSchema.Entities["MyTestEntity"])
	assert.Equal(t, "This is my test entity", ymlSchema.Entities["MyTestEntity"].Description)
}

func TestConfigurationBuilderrParsesYmlWithEntitiesAndPhases(t *testing.T) {
	ymlParser := NewYmlParser()
	ymlSchema, err := ymlParser.Parse(
		`
    Name: Example
    Description: Example YML configuration
    Entities:
      MyTestEntity:
        Description: This is my test entity
        Components:
          MyTestComponent1:
            Description: This is component 1
          MyTestComponent2:
            Description: This is component 2
  `)

	assert.Nil(t, err)
	assert.IsType(t, YmlComponent{}, ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"])
	assert.Equal(t, "This is component 1", ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Description)
	assert.Equal(t, "This is component 2", ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent2"].Description)
}

func TestConfigurationBuilderParsesYmlWithEntitiesAndPhasesAndTasks(t *testing.T) {
	ymlParser := NewYmlParser()
	ymlSchema, err := ymlParser.Parse(
		`
   Name: Example
   Description: Example YML configuration
   Entities:
     MyTestEntity:
       Description: This is my test entity
       Components:
         MyTestComponent1:
           Description: This is a test component
           Elements:
             ElementA:
               Resource: MyTestResource
               Shares: true
             ElementB:
               Resource: MyTestResource2
               SelectionCriteria:
                 Type: Custom
                 Criteria: |
                   1 = 1
             ElementC:
               Resource: MyTestResource2
               SelectionCriteria:
                 Type: Index
                 Elements:
                   - ElementA
                 Index: IDX1
             ElementD:
               Resource: MyTestResource3
               SelectionCriteria:
                 Type: Related
                 Elements:
                   - ElementB
 `)

	assert.Nil(t, err)
	assert.IsType(t, YmlElement{}, ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementA"])
	assert.Equal(t, "MyTestResource", ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementA"].Resource)

	assert.IsType(t, YmlSelectionCriteria{}, ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementB"].SelectionCriteria)
	assert.Equal(t, "Custom", ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementB"].SelectionCriteria.Type)
	assert.Equal(t, "1 = 1\n", ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementB"].SelectionCriteria.Criteria)

	assert.Equal(t, "Index", ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementC"].SelectionCriteria.Type)
	assert.ElementsMatch(t, []string{"ElementA"}, ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementC"].SelectionCriteria.Elements)
	assert.Equal(t, "IDX1", ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementC"].SelectionCriteria.Index)

	assert.Equal(t, "Related", ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementD"].SelectionCriteria.Type)
	assert.ElementsMatch(t, []string{"ElementB"}, ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementD"].SelectionCriteria.Elements)
}

func TestConfigurationBuilderParsesYmlWithResourcesEntitiesPhasesAndTasks(t *testing.T) {
	ymlParser := NewYmlParser()
	ymlSchema, err := ymlParser.Parse(
		`
   Name: Example
   Description: Example YML configuration
   Resources:
     MyTestResource:
       TableName: my_test_table
       PrimaryKey:
         - my_test_table.id
       AutoIncrement: true
       Index:
         IDX1:
           - my_test_table.idx1
         IDX2:
           - my_test_table.idx1
           - my_test_table.idx2
       ForeignKeys:
         - Type: NORMAL
           Key: my_test_table.fk1
           ResourceName: MyTestResource2
           ForeignKey: my_test_table2.id

     MyTestResource2:
       TableName: my_test_table2
       PrimaryKey:
         - my_test_table2.id
       AutoIncrement: true
       Index:
         IDX1:
           - my_test_table2.idx1
         IDX2:
           - my_test_table2.idx1
           - my_test_table2.idx2
   Entities:
     MyTestEntity:
       Description: This is my test entity
       Components:
         MyTestComponent1:
           Description: This is a test component
           Elements:
             ElementA:
               Resource: MyTestResource
               Shares: true
             ElementB:
               Resource: MyTestResource2
               SelectionCriteria:
                 Type: Custom
                 Criteria: |
                   1 = 1
             ElementC:
               Resource: MyTestResource2
               SelectionCriteria:
                 Type: Index
                 Elements:
                   - ElementA
                 Index: IDX1
             ElementD:
               Resource: MyTestResource3
               SelectionCriteria:
                 Type: Related
                 Elements:
                   - ElementB
 `)

	assert.Nil(t, err)
	assert.Equal(t, "Example", ymlSchema.Name)
	assert.Equal(t, "Example YML configuration", ymlSchema.Description)

	assert.Equal(t, "my_test_table", ymlSchema.Resources["MyTestResource"].TableName)
	assert.ElementsMatch(t, []string{"my_test_table.id"}, ymlSchema.Resources["MyTestResource"].PrimaryKey)
	assert.True(t, ymlSchema.Resources["MyTestResource"].AutoIncrement)
	assert.ElementsMatch(t, []string{"my_test_table.idx1"}, ymlSchema.Resources["MyTestResource"].Index["IDX1"])

	assert.Equal(t, "my_test_table2", ymlSchema.Resources["MyTestResource2"].TableName)
	assert.ElementsMatch(t, []string{"my_test_table2.id"}, ymlSchema.Resources["MyTestResource2"].PrimaryKey)
	assert.True(t, ymlSchema.Resources["MyTestResource2"].AutoIncrement)
	assert.ElementsMatch(t, []string{"my_test_table2.idx1"}, ymlSchema.Resources["MyTestResource2"].Index["IDX1"])

	assert.IsType(t, YmlElement{}, ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementA"])
	assert.Equal(t, "MyTestResource", ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementA"].Resource)

	assert.IsType(t, YmlSelectionCriteria{}, ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementB"].SelectionCriteria)
	assert.Equal(t, "Custom", ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementB"].SelectionCriteria.Type)
	assert.Equal(t, "1 = 1\n", ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementB"].SelectionCriteria.Criteria)

	assert.Equal(t, "Index", ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementC"].SelectionCriteria.Type)
	assert.ElementsMatch(t, []string{"ElementA"}, ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementC"].SelectionCriteria.Elements)
	assert.Equal(t, "IDX1", ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementC"].SelectionCriteria.Index)

	assert.Equal(t, "Related", ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementD"].SelectionCriteria.Type)
	assert.ElementsMatch(t, []string{"ElementB"}, ymlSchema.Entities["MyTestEntity"].Components["MyTestComponent1"].Elements["ElementD"].SelectionCriteria.Elements)
}
