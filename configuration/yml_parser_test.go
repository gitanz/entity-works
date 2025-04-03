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
        Phases:
          MyTestPhase1:
            Description: This is a test phase
          MyTestPhase2:
            Description: This is a test phase 2
  `)

	assert.Nil(t, err)
	assert.IsType(t, YmlPhase{}, ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"])
	assert.Equal(t, "This is a test phase", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Description)
	assert.Equal(t, "This is a test phase 2", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase2"].Description)
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
        Phases:
          MyTestPhase1:
            Description: This is a test phase
            Tasks:
              TaskA:
                Resource: MyTestResource
                Shares: true
              TaskB:
                Resource: MyTestResource2
                SelectionCriteria:
                  Type: Custom
                  Criteria: |
                    1 = 1
              TaskC:
                Resource: MyTestResource2
                SelectionCriteria:
                  Type: Index
                  FromTask: TaskA
                  Index: IDX1
              TaskD:
                Resource: MyTestResource3
                SelectionCriteria:
                  Type: Related
                  Tasks: 
                    - TaskB
  `)

	assert.Nil(t, err)
	assert.IsType(t, YmlTask{}, ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskA"])
	assert.Equal(t, "MyTestResource", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskA"].Resource)
	assert.True(t, ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskA"].Shares)

	assert.IsType(t, YmlSelectionCriteria{}, ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskB"].SelectionCriteria)
	assert.Equal(t, "Custom", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskB"].SelectionCriteria.Type)
	assert.Equal(t, "1 = 1\n", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskB"].SelectionCriteria.Criteria)

	assert.Equal(t, "Index", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskC"].SelectionCriteria.Type)
	assert.Equal(t, "TaskA", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskC"].SelectionCriteria.FromTask)
	assert.Equal(t, "IDX1", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskC"].SelectionCriteria.Index)

	assert.Equal(t, "Related", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskD"].SelectionCriteria.Type)
	assert.ElementsMatch(t, []string{"TaskB"}, ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskD"].SelectionCriteria.Tasks)
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
        Phases:
          MyTestPhase1:
            Description: This is a test phase
            Tasks:
              TaskA:
                Resource: MyTestResource
                Shares: true
              TaskB:
                Resource: MyTestResource2
                SelectionCriteria:
                  Type: Custom
                  Criteria: |
                    1 = 1
              TaskC:
                Resource: MyTestResource2
                SelectionCriteria:
                  Type: Index
                  FromTask: TaskA
                  Index: IDX1
              TaskD:
                Resource: MyTestResource3
                SelectionCriteria:
                  Type: Related
                  Tasks: 
                    - TaskB
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

	assert.IsType(t, YmlTask{}, ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskA"])
	assert.Equal(t, "MyTestResource", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskA"].Resource)
	assert.True(t, ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskA"].Shares)

	assert.IsType(t, YmlSelectionCriteria{}, ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskB"].SelectionCriteria)
	assert.Equal(t, "Custom", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskB"].SelectionCriteria.Type)
	assert.Equal(t, "1 = 1\n", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskB"].SelectionCriteria.Criteria)

	assert.Equal(t, "Index", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskC"].SelectionCriteria.Type)
	assert.Equal(t, "TaskA", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskC"].SelectionCriteria.FromTask)
	assert.Equal(t, "IDX1", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskC"].SelectionCriteria.Index)

	assert.Equal(t, "Related", ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskD"].SelectionCriteria.Type)
	assert.ElementsMatch(t, []string{"TaskB"}, ymlSchema.Entities["MyTestEntity"].Phases["MyTestPhase1"].Tasks["TaskD"].SelectionCriteria.Tasks)
}
