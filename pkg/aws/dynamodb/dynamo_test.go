package dynamodb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type person struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Age     string `json:"age"`
	City    string `json:"city"`
	Country string `json:"country"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

func TestDynamoLocalDevelopmentGetOne(t *testing.T) {
	client := NewLocalClient().
		WithTable(DynamoTable{
			TableName:         "person",
			PartitionKeyField: "id",
		}).
		WithPreloadedItems("person", "/test.json")
	var p person
	err := client.GetOne("person", "1", &p)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "1", p.Id)
	assert.Equal(t, "John", p.Name)
	assert.Equal(t, "30", p.Age)

	err = client.GetOne("person", "2", &p)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "2", p.Id)
	assert.Equal(t, "Peter", p.Name)
	assert.Equal(t, "40", p.Age)
}

func TestDynamoLocalDevelopmentSave(t *testing.T) {
	client := NewLocalClient().
		WithTable(DynamoTable{
			TableName:         "person",
			PartitionKeyField: "id",
		}).
		WithPreloadedItems("person", "/test.json")

	p := person{
		Id:      "4",
		Name:    "James",
		Age:     "60",
		City:    "London",
		Country: "UK",
		Email:   "]}"}
	err := client.Save("person", p)
	if err != nil {
		t.Error(err)
	}

	var p2 person
	err = client.GetOne("person", "4", &p2)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "4", p2.Id)
	assert.Equal(t, "James", p2.Name)
	assert.Equal(t, "60", p2.Age)
	assert.Equal(t, "London", p2.City)
	assert.Equal(t, "UK", p2.Country)
}

func TestDynamoLocalDevelopmentGetOneWithSort(t *testing.T) {
	client := NewLocalClient().
		WithTable(DynamoTable{
			TableName:         "person",
			PartitionKeyField: "id",
			SortKeyField:      "name",
		}).
		WithPreloadedItems("person", "/test.json")
	var p person
	err := client.GetOneWithSort("person", "1", "John", &p)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "1", p.Id)
	assert.Equal(t, "John", p.Name)
	assert.Equal(t, "30", p.Age)

	err = client.GetOneWithSort("person", "2", "Peter", &p)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "2", p.Id)
	assert.Equal(t, "Peter", p.Name)
	assert.Equal(t, "40", p.Age)

	err = client.GetOneWithSort("person", "2", "John", &p)
	assert.Error(t, err)
	assert.Equal(t, "item not found", err.Error())
}

func TestDynamoLocalDevelopmentQueryOne(t *testing.T) {
	client := NewLocalClient().
		WithTable(DynamoTable{
			TableName:         "person",
			PartitionKeyField: "id",
		}).
		WithPreloadedItems("person", "/test.json")

	var p person
	err := client.QueryOne("person", "3", 1, &p)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "3", p.Id)
	assert.Equal(t, "Suzanne", p.Name)
	assert.Equal(t, "50", p.Age)
}
