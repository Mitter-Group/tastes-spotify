## DynamoDB go pomelo client

### Install

```sh
go get github.com/pomelo-la/fraud-libs/aws/dynamodbv2
```

### Import
```go
import "github.com/pomelo-la/fraud-libs/aws/dynamodbv2"
```


### How to use DynamoDB client

```json
"aws": 
{
    "table1": {
      "table_name": "use1-dev-table2",
      "primary_key_field": "pk_name",
      "sort_key_field": "sk_name"
      },
    "table2": {
      "table_name": "use1-dev-table2",
      "primary_key_field": "pk_name",
      "sort_key_field": "sk_name"
      }
  }
```
```go
    //constructor
	dynamoV2 := dynamov2.NewDynamoClientv2(awsConfig,
		dynamov2.WithTable(cfg.AWS.table1),
		dynamov2.WithTable(cfg.AWS.table2),
	)

    //mehtod  example
    err := c.dynamov2.Save("tablename", body)
	    if err != nil {
	    	return err
	    }
	    return nil
```

### How to work with the library locally?
You can use localstack or the bundled local feature.

First intialize the client:

```go
    client := dynamodb.NewLocalClient() 
// add tables to the client
    client.WithTable(dynamodb.DynamoTable{
    TableName:         "tableName",
    PartitionKeyField: "id",
    SortKeyField:      "name",
    })
```

You can also preload a list of json items in the client:

```go
dynamo.WithPreloadedItems("tableName", "./preloaded.json")
```

The local clients expects a json file with the following format:
```json
[
  {
    "id": "1",
    "name": "John",
    "age": "30",
    "city": "New York"
  },
  {
    "id": "2",
    "name": "Samuel",
    "age": "23",
    "city": "Montevideo"
  }
]
```



### How to mock DynamoDB client

We use `testify/mock` to mock DynamoDB client.

```go
    dynamoClient := dynamodb.DynamoMock{}
    dynamoClient.On("GetOneWithSort", "entity", "test_entity", mock.Anything).Return(nil).Run(
    func(args mock.Arguments) {
    entity := args.Get(2).(*model.Entity)
    entity.Id = "test_entity"
    entity.Type = "Test Event"
    })
```