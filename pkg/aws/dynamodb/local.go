package dynamodb

import (
	"encoding/json"
	"fmt"
	"os"
)

type LocalClient struct {
	data          map[string][]interface{}
	preloadedFile string
	tables        map[string]*DynamoTable
}

type LocalTableConfig struct {
	PartitionKey string
	SortKey      string
}

func NewLocalClient() *LocalClient {
	return &LocalClient{
		data:   make(map[string][]interface{}),
		tables: make(map[string]*DynamoTable),
	}
}

// WithTable initializes the given table with the given config.
func (l *LocalClient) WithTable(table DynamoTable) *LocalClient {
	l.tables[table.TableName] = &table
	return l
}

// WithPreloadedItems loads the given file into the local client.
// The file should be a JSON file with the following format:
// [{...}, {...}]
func (l *LocalClient) WithPreloadedItems(table string, filePath string) *LocalClient {
	//load file located in filePath
	//add to l.data
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	fileName := mydir + filePath
	configFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("opening config file error: ", err.Error())
		return nil
	}

	var items []map[string]interface{}
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&items); err != nil {
		fmt.Println("parsing preloaded file", err.Error())
	}

	l.data[table] = []interface{}{}
	for _, item := range items {
		l.data[table] = append(l.data[table], item)
	}

	return l
}

func (l *LocalClient) Save(table string, values interface{}) error {
	l.data[table] = append(l.data[table], values)
	return nil
}

func (l *LocalClient) GetOne(table string, partitionKey string, bindTo interface{}) error {
	if l.data[table] == nil {
		return fmt.Errorf("table %s not found", table)
	}
	if l.tables[table] == nil {
		return fmt.Errorf("table %s not initialized", table)
	}
	partitionField := l.tables[table].PartitionKeyField
	for _, item := range l.data[table] {
		var itemMap map[string]interface{}
		inrec, _ := json.Marshal(item)
		json.Unmarshal(inrec, &itemMap)
		if itemMap[partitionField] == partitionKey {
			b, err := json.Marshal(itemMap)
			if err != nil {
				return err
			}
			return json.Unmarshal(b, bindTo)
		}
	}

	return fmt.Errorf("item not found")
}

func (l *LocalClient) GetOneWithSort(table string, partitionKey string, sortKey string, bindTo interface{}) error {
	if l.data[table] == nil {
		return fmt.Errorf("table %s not found", table)
	}
	if l.tables[table] == nil {
		return fmt.Errorf("table %s not initialized", table)
	}
	partitionField := l.tables[table].PartitionKeyField
	sortField := l.tables[table].SortKeyField
	for _, item := range l.data[table] {
		var itemMap map[string]interface{}
		inrec, _ := json.Marshal(item)
		json.Unmarshal(inrec, &itemMap)
		if itemMap[partitionField] == partitionKey && itemMap[sortField] == sortKey {
			b, err := json.Marshal(itemMap)
			if err != nil {
				return err
			}
			return json.Unmarshal(b, bindTo)
		}
	}
	return fmt.Errorf("item not found")
}

func (l *LocalClient) QueryOne(table string, partitionKey string, limit int32, bindTo interface{}) error {
	if l.data[table] == nil {
		return fmt.Errorf("table %s not found", table)
	}
	if l.tables[table] == nil {
		return fmt.Errorf("table %s not initialized", table)
	}
	partitionField := l.tables[table].PartitionKeyField
	for _, item := range l.data[table] {
		var itemMap map[string]interface{}
		inrec, _ := json.Marshal(item)
		json.Unmarshal(inrec, &itemMap)
		if itemMap[partitionField] == partitionKey {
			b, err := json.Marshal(itemMap)
			if err != nil {
				return err
			}
			return json.Unmarshal(b, bindTo)
		}
	}
	return fmt.Errorf("item not found")
}

func (l *LocalClient) QueryMultiple(table string, partitionKey string, limit int32, bindTo interface{}) error {
	if l.data[table] == nil {
		return fmt.Errorf("table %s not found", table)
	}
	if l.tables[table] == nil {
		return fmt.Errorf("table %s not initialized", table)
	}
	partitionField := l.tables[table].PartitionKeyField
	var items []interface{}
	for _, item := range l.data[table] {
		var itemMap map[string]interface{}
		inrec, _ := json.Marshal(item)
		json.Unmarshal(inrec, &itemMap)
		if itemMap[partitionField] == partitionKey {
			items = append(items, item)
		}
	}
	b, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, bindTo)

}

func (l *LocalClient) BatchGetWithSort(values map[string]interface{}) error {
	return nil
}
