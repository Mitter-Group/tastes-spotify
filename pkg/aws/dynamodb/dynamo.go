package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/chunnior/spotify/internal/util/log"
)

var (
	// ErrNotFound is returned when no items could be found in Get or OldValue and similar operations.
	ErrNotFound = errors.New("dynamo: no item found")
	Tagkey      = "dynamo"
	tables      = make(map[string]DynamoTable)
)

type Implementation struct {
	client            *dynamodb.Client
	partitionKeyField string
	sortKeyField      *string
	DynamoTables      map[string]DynamoTable
}

type DynamoTable struct {
	TableName         string `json:"table_name"`
	PartitionKeyField string `json:"primary_key_field"`
	SortKeyField      string `json:"sort_key_field"`
}

type funcTable func(i *Implementation)

func WithTable(arg DynamoTable) funcTable {
	return func(i *Implementation) {
		tables[arg.TableName] = DynamoTable{
			TableName:         arg.TableName,
			PartitionKeyField: arg.PartitionKeyField,
			SortKeyField:      arg.SortKeyField,
		}
		i.DynamoTables = tables
	}
}

func NewDynamoClient(awsConfig aws.Config, funcTableArray ...funcTable) Client {
	var i Implementation
	db := dynamodb.NewFromConfig(awsConfig)
	for _, ft := range funcTableArray {
		ft(&i)
	}
	i.client = db
	return &i
}

func (i *Implementation) Save(table string, values interface{}) error {
	log.Debugf("[DynamoDB] executing put query")

	item, err := attributevalue.MarshalMapWithOptions(values, func(h *attributevalue.EncoderOptions) {
		h.TagKey = Tagkey
	})

	log.Debugf("[DynamoDB] item to save: %s", item)
	if err != nil {
		panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
	}

	_, err = i.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(i.DynamoTables[table].TableName),
		Item:      item,
	})
	return err
}

func (i *Implementation) getItem(table string, key map[string]types.AttributeValue, bindTo interface{}) error {
	log.Debugf("[DynamoDB] executing get query")

	out, err := i.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(i.DynamoTables[table].TableName),
		Key:       key,
	})
	if err != nil {
		return err
	}
	if out.Item == nil {
		return ErrNotFound
	}
	err = attributevalue.UnmarshalMapWithOptions(out.Item, &bindTo, func(options *attributevalue.DecoderOptions) {
		options.TagKey = Tagkey
	})

	return err
}

func (i *Implementation) getItemQuery(table string, key string, limit int32, bindTo interface{}) error {
	log.Debugf("[DynamoDB] executing get query")
	keyEx := expression.Key(i.DynamoTables[table].PartitionKeyField).Equal(expression.Value(key))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return err
	}
	out, err := i.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String(i.DynamoTables[table].TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		Limit:                     aws.Int32(limit),
	})
	if out != nil {
		if len(out.Items) == 0 {
			return ErrNotFound
		}
	}
	if err != nil {
		return err
	}

	err = attributevalue.UnmarshalListOfMapsWithOptions(out.Items, &bindTo, func(options *attributevalue.DecoderOptions) {
		options.TagKey = Tagkey
	})
	return err
}

func (i *Implementation) ItemQueryExpression(table string, query expression.Expression, limit int32, bindTo interface{}) error {
	log.Debugf("[DynamoDB] executing get query")

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(i.DynamoTables[table].TableName),
		ExpressionAttributeNames:  query.Names(),
		ExpressionAttributeValues: query.Values(),
		KeyConditionExpression:    query.KeyCondition(),
	}
	if limit > -1 {
		queryInput.Limit = aws.Int32(limit)
	}
	out, err := i.client.Query(context.TODO(), queryInput)
	if out != nil {
		if len(out.Items) == 0 {
			return ErrNotFound
		}
	}
	if err != nil {
		return err
	}

	err = attributevalue.UnmarshalListOfMapsWithOptions(out.Items, &bindTo, func(options *attributevalue.DecoderOptions) {
		options.TagKey = Tagkey
	})
	return err
}

func (i *Implementation) BatchGetItem(key map[string]types.KeysAndAttributes, values map[string]interface{}) error {
	log.Debugf("[DynamoDB] batch get query")
	out, err := i.client.BatchGetItem(context.TODO(), &dynamodb.BatchGetItemInput{
		RequestItems: key,
	})
	if err != nil {
		return err
	}
	var bindList []interface{}
	for i, o := range out.Responses {
		for t, v := range values {
			if i == t {
				v, _ := v.([]interface{})
				bindList := append(bindList, v[2])
				err = attributevalue.UnmarshalListOfMapsWithOptions(o, &bindList, func(options *attributevalue.DecoderOptions) {
					options.TagKey = Tagkey
				})
			}
		}
	}
	return err
}

func (i *Implementation) GetOne(table string, partitionKey string, bindTo interface{}) error {
	log.Debugf("[DynamoDB] executing get query")
	return i.getItem(table,
		map[string]types.AttributeValue{
			i.DynamoTables[table].PartitionKeyField: &types.AttributeValueMemberS{Value: partitionKey},
		},
		bindTo)
}

func (i *Implementation) GetOneWithSort(table string, partitionKey string, sortKey string, bindTo interface{}) error {
	log.Debugf("[DynamoDB] executing get query with sortkey [pk:%s][sk:%s]", partitionKey, sortKey)
	tData := i.DynamoTables[table]

	return i.getItem(table,
		map[string]types.AttributeValue{
			tData.PartitionKeyField: &types.AttributeValueMemberS{Value: partitionKey},
			tData.SortKeyField:      &types.AttributeValueMemberS{Value: sortKey},
		},
		bindTo)
}

func (i *Implementation) QueryOne(table string, partitionKey string, limit int32, bindTo interface{}) error {
	log.Debugf("[DynamoDB] executing get query with [pk:%s]", partitionKey)

	return i.getItemQuery(table, partitionKey, limit, bindTo)

}

func (i *Implementation) BatchGetWithSort(values map[string]interface{}) error {
	log.Debugf("[DynamoDB] executing BatchGet query")
	batchkeys := make(map[string]types.KeysAndAttributes)
	for t, v := range values {
		v, _ := v.([]interface{})
		batchkeys[t] = types.KeysAndAttributes{
			Keys: []map[string]types.AttributeValue{
				{
					i.DynamoTables[t].PartitionKeyField: &types.AttributeValueMemberS{Value: v[0].(string)},
					i.DynamoTables[t].SortKeyField:      &types.AttributeValueMemberS{Value: v[1].(string)},
				},
			},
		}

	}

	return i.BatchGetItem(batchkeys, values)
}

// QueryExpression  returns multiple items by using a query expression
// -1 for limit means no limit
func (i *Implementation) QueryExpression(table string, query expression.Expression, limit int32, bindTo interface{}) error {
	return i.ItemQueryExpression(table, query, limit, bindTo)
}
