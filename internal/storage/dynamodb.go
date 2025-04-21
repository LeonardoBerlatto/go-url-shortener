package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/leonardoberlatto/go-url-shortener/internal/logger"
	"github.com/leonardoberlatto/go-url-shortener/internal/models"
)

const (
	tableName    = "UrlMappings"
	partitionKey = "ShortID"
)

type DynamoDBStorage struct {
	client *dynamodb.Client
}

func NewDynamoDBStorage(endpoint, region, accessKeyID, secretAccessKey string) (*DynamoDBStorage, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		logger.Fatal(err)
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if ep := endpoint; ep != "" {
			o.BaseEndpoint = aws.String(ep)
		}
	})

	return &DynamoDBStorage{
		client: client,
	}, nil
}

func (d *DynamoDBStorage) Store(ctx context.Context, mapping models.URLMapping) error {
	item, err := attributevalue.MarshalMap(mapping)
	if err != nil {
		return err
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(ShortID)"),
	})

	if err != nil {
		if _, ok := err.(*types.ConditionalCheckFailedException); ok {
			return ErrorConflict
		}
		return err
	}

	return nil
}

func (d *DynamoDBStorage) Get(ctx context.Context, shortID string) (models.URLMapping, error) {
	result, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       mapPartitionKey(shortID),
	})

	if err != nil {
		return models.URLMapping{}, err
	}

	if result.Item == nil {
		return models.URLMapping{}, ErrorNotFound
	}

	var mapping models.URLMapping
	err = attributevalue.UnmarshalMap(result.Item, &mapping)
	if err != nil {
		return models.URLMapping{}, err
	}

	return mapping, nil
}

func (d *DynamoDBStorage) Delete(ctx context.Context, shortID string) error {
	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       mapPartitionKey(shortID),
	})
	return err
}

func (d *DynamoDBStorage) CheckExists(ctx context.Context, shortID string) (bool, error) {
	result, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:            aws.String(tableName),
		Key:                  mapPartitionKey(shortID),
		ProjectionExpression: aws.String(partitionKey),
	})

	if err != nil {
		return false, err
	}

	return result.Item != nil, nil
}

func (d *DynamoDBStorage) IncrementHits(ctx context.Context, shortID string) error {
	_, err := d.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        aws.String(tableName),
		Key:              mapPartitionKey(shortID),
		UpdateExpression: aws.String("ADD Hits :inc"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc": &types.AttributeValueMemberN{Value: "1"},
		},
	})
	return err
}

func mapPartitionKey(shortID string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"ShortID": &types.AttributeValueMemberS{Value: shortID},
	}
}

func (d *DynamoDBStorage) ListURLs(ctx context.Context, pageNumber, pageSize int) ([]models.URLMapping, int64, error) {
	offset := (pageNumber - 1) * pageSize

	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
		Limit:     aws.Int32(int32(pageSize)),
	}

	totalCount := int64(0)
	countResult, err := d.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(tableName),
		Select:    types.SelectCount,
	})
	if err != nil {
		return nil, 0, err
	}
	totalCount = int64(countResult.Count)

	var scannedCount int32 = 0
	var items []map[string]types.AttributeValue
	var lastEvaluatedKey map[string]types.AttributeValue

	for scannedCount < int32(offset) || len(items) < pageSize {
		if lastEvaluatedKey != nil {
			scanInput.ExclusiveStartKey = lastEvaluatedKey
		}

		result, err := d.client.Scan(ctx, scanInput)
		if err != nil {
			return nil, 0, err
		}

		if scannedCount < int32(offset) {
			if scannedCount+result.ScannedCount <= int32(offset) {
				scannedCount += result.ScannedCount
			} else {
				remainingOffset := int32(offset) - scannedCount
				itemsNeeded := int(result.ScannedCount) - int(remainingOffset)
				if itemsNeeded > pageSize-len(items) {
					itemsNeeded = pageSize - len(items)
				}
				items = append(items, result.Items[remainingOffset:remainingOffset+int32(itemsNeeded)]...)
				scannedCount += result.ScannedCount
			}
		} else {
			itemsToAdd := len(result.Items)
			if len(items)+itemsToAdd > pageSize {
				itemsToAdd = pageSize - len(items)
			}
			items = append(items, result.Items[:itemsToAdd]...)
		}

		if len(items) >= pageSize || result.LastEvaluatedKey == nil {
			break
		}

		lastEvaluatedKey = result.LastEvaluatedKey
	}

	var mappings []models.URLMapping
	for _, item := range items {
		var mapping models.URLMapping
		if err := attributevalue.UnmarshalMap(item, &mapping); err != nil {
			return nil, 0, err
		}
		mappings = append(mappings, mapping)
	}

	return mappings, totalCount, nil
}
