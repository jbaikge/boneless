package gocms

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/zeebo/assert"
)

var dynamoTablePrefix = time.Now().Format("20060102-150405-")

func testDynamoConfig(tableName string) (cfg aws.Config, err error) {
	endpointResolverFunc := func(service string, region string, options ...interface{}) (endpoint aws.Endpoint, err error) {
		endpoint = aws.Endpoint{
			PartitionID:   "aws",
			URL:           "http://localhost:4566", // 4566 for LocalStack; 8000 for amazon/dynamodb-local
			SigningRegion: "local",
		}
		return
	}
	endpointResolver := aws.EndpointResolverWithOptionsFunc(endpointResolverFunc)

	cfg, err = config.LoadDefaultConfig(
		context.Background(),
		config.WithEndpointResolverWithOptions(endpointResolver),
	)
	if err != nil {
		return
	}

	// Build table before returning
	createTable := &dynamodb.CreateTableInput{
		TableName:   &tableName,
		BillingMode: types.BillingModePayPerRequest,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
	}

	client := dynamodb.NewFromConfig(cfg)
	if _, err = client.CreateTable(context.Background(), createTable); err != nil {
		return
	}

	return
}

func testDynamoEmptyTable(cfg aws.Config, table string) (err error) {
	client := dynamodb.NewFromConfig(cfg)
	params := &dynamodb.ScanInput{
		TableName: &table,
	}
	response, err := client.Scan(context.Background(), params)
	if err != nil {
		return
	}

	for _, item := range response.Items {
		delete := &dynamodb.DeleteItemInput{
			TableName: &table,
			Key: map[string]types.AttributeValue{
				"PK": item["PK"],
				"SK": item["SK"],
			},
		}
		if _, err = client.DeleteItem(context.Background(), delete); err != nil {
			return
		}
	}
	return
}

func TestDynamoFromClass(t *testing.T) {
	class := Class{
		Id:     "from_class",
		Name:   t.Name(),
		Fields: []Field{{Name: "field_1"}},
	}

	dc := new(dynamoClass)
	dc.FromClass(&class)

	assert.Equal(t, dynamoClassPrefix+class.Id, dc.PK)
	assert.Equal(t, "class#v0000", dc.SK)
	assert.Equal(t, t.Name(), dc.Name)
	assert.DeepEqual(t, class.Fields, dc.Fields)
}

func TestDynamoToClass(t *testing.T) {
	id := "to_class"
	dc := dynamoClass{
		PK:     dynamoClassPrefix + id,
		SK:     "class_v0",
		Name:   t.Name(),
		Fields: []Field{{Name: "field_1"}},
	}
	class := dc.ToClass()

	assert.Equal(t, id, class.Id)
	assert.Equal(t, t.Name(), class.Name)
	assert.DeepEqual(t, dc.Fields, class.Fields)
}

func TestDynamoDBRepository(t *testing.T) {
	resources := DynamoDBResources{
		Table: dynamoTablePrefix + "Test",
	}
	cfg, err := testDynamoConfig(resources.Table)
	assert.NoError(t, err)

	repo := NewDynamoDBRepository(cfg, resources)
	ctx := context.Background()

	t.Run("CreateClass", func(t *testing.T) {
		class := Class{
			Id:   "class_create",
			Name: t.Name(),
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))
	})

	t.Run("GetClassByIdSuccess", func(t *testing.T) {
		class := Class{
			Id:   "get_success",
			Name: t.Name(),
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))

		check, err := repo.GetClassById(ctx, class.Id)
		assert.NoError(t, err)
		assert.Equal(t, class.Id, check.Id)
		assert.Equal(t, class.Name, check.Name)
	})

	t.Run("GetClassByIdFail", func(t *testing.T) {
		_, err := repo.GetClassById(ctx, "get_fail")
		assert.Equal(t, ErrNotExist, err)
	})

	t.Run("UpdateClass", func(t *testing.T) {
		class := Class{
			Id:          "update_class",
			Name:        t.Name(),
			TableLabels: "Field 1; Field 2",
			TableFields: "field_1; field_2",
			Created:     time.Now(),
			Updated:     time.Now(),
			Fields:      []Field{{Name: "field_2"}, {Name: "field_1"}},
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))

		class.Name = t.Name() + "-Updated"
		class.TableLabels += "; Field 3"
		class.TableFields += "; field_3"
		class.Updated = time.UnixMicro(time.Now().UnixMicro())
		class.Fields = append(class.Fields, Field{Name: "field_3"})
		assert.NoError(t, repo.UpdateClass(ctx, &class))

		check, err := repo.GetClassById(ctx, class.Id)
		assert.NoError(t, err)
		assert.Equal(t, class.Name, check.Name)
		assert.Equal(t, class.TableFields, check.TableFields)
		assert.Equal(t, class.TableLabels, check.TableLabels)
		assert.Equal(t, class.Updated, check.Updated)
		assert.DeepEqual(t, class.Fields, check.Fields)
	})

	t.Run("DeleteClass", func(t *testing.T) {
		class := Class{
			Id:   "delete_class",
			Name: t.Name(),
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))
		assert.NoError(t, repo.DeleteClass(ctx, class.Id))
		_, err := repo.GetClassById(ctx, class.Id)
		assert.Equal(t, ErrNotExist, err)
	})

	t.Run("ClassList", func(t *testing.T) {
		assert.NoError(t, testDynamoEmptyTable(cfg, resources.Table))

		t.Run("Empty", func(t *testing.T) {
			filter := ClassFilter{Range: Range{End: 9}}
			classes, r, err := repo.GetClassList(ctx, filter)
			assert.NoError(t, err)
			assert.DeepEqual(t, Range{}, r)
			assert.Equal(t, 0, len(classes))
		})

		for i := 0; i < 10; i++ {
			class := Class{
				Id:   fmt.Sprintf("class_list_%02d", i),
				Name: fmt.Sprintf("Class List (%02d)", i+1),
			}
			assert.NoError(t, repo.CreateClass(ctx, &class))
		}

		t.Run("All", func(t *testing.T) {
			filter := ClassFilter{Range: Range{End: 9}}
			classes, r, err := repo.GetClassList(ctx, filter)
			assert.NoError(t, err)
			assert.DeepEqual(t, Range{End: 9, Size: 10}, r)
			assert.Equal(t, 10, len(classes))
		})

		t.Run("LargeWindow", func(t *testing.T) {
			filter := ClassFilter{Range: Range{End: 99}}
			classes, r, err := repo.GetClassList(ctx, filter)
			assert.NoError(t, err)
			assert.DeepEqual(t, Range{End: 9, Size: 10}, r)
			assert.Equal(t, 10, len(classes))
		})

		t.Run("InvalidRange", func(t *testing.T) {
			filter := ClassFilter{Range: Range{Start: 90, End: 99}}
			classes, r, err := repo.GetClassList(ctx, filter)
			assert.Equal(t, ErrBadRange, err)
			assert.DeepEqual(t, Range{Size: 10}, r)
			assert.Equal(t, 0, len(classes))
		})

		t.Run("Beginning", func(t *testing.T) {
			filter := ClassFilter{Range: Range{Start: 0, End: 4}}
			classes, r, err := repo.GetClassList(ctx, filter)
			assert.NoError(t, err)
			assert.DeepEqual(t, Range{End: 4, Size: 10}, r)
			assert.Equal(t, 5, len(classes))
		})

		t.Run("End", func(t *testing.T) {
			filter := ClassFilter{Range: Range{Start: 5, End: 9}}
			classes, r, err := repo.GetClassList(ctx, filter)
			assert.NoError(t, err)
			assert.DeepEqual(t, Range{Start: 5, End: 9, Size: 10}, r)
			assert.Equal(t, 5, len(classes))
		})

		t.Run("Middle", func(t *testing.T) {
			filter := ClassFilter{Range: Range{Start: 3, End: 6}}
			classes, r, err := repo.GetClassList(ctx, filter)
			assert.NoError(t, err)
			assert.DeepEqual(t, Range{Start: 3, End: 6, Size: 10}, r)
			assert.Equal(t, 4, len(classes))
		})
	})

	t.Run("CreateDocument", func(t *testing.T) {
		class := Class{
			Id:   "create_document_class",
			Name: "Create Document Class",
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))

		doc := Document{
			Id:      "create_document",
			ClassId: class.Id,
			Name:    t.Name(),
		}
		assert.NoError(t, repo.CreateDocument(ctx, &doc))
	})

	t.Run("CreateDocumentWithPath", func(t *testing.T) {
		class := Class{
			Id:   "create_document_path_class",
			Name: "Create Document with Path Class",
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))

		doc := Document{
			Id:      "document_with_path",
			ClassId: class.Id,
			Name:    t.Name(),
			Path:    "/doc/with/path",
		}
		assert.NoError(t, repo.CreateDocument(ctx, &doc))
	})

	t.Run("CreateDocumentWithSort", func(t *testing.T) {
		assert.NoError(t, testDynamoEmptyTable(cfg, resources.Table))

		class := Class{
			Id:   "sort_class",
			Name: "Sort Class",
			Fields: []Field{
				{
					Name: "existing_field",
					Sort: true,
				},
				{
					Name: "nonexistant_field",
					Sort: true,
				},
				{
					Name: "ignore_field",
					Sort: false,
				},
				{
					Name: "time_field",
					Sort: true,
				},
			},
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))

		doc := Document{
			Id:      "sort_doc",
			ClassId: class.Id,
			Name:    t.Name(),
			Values: map[string]interface{}{
				"existing_field": "My Value",
				"ignore_field":   "Ignore Me",
				"extra_field":    "Extra Field",
				"time_field":     time.Now(),
			},
		}
		assert.NoError(t, repo.CreateDocument(ctx, &doc))
	})

	t.Run("GetDocumentByIdSuccess", func(t *testing.T) {
		class := Class{
			Id:   "get_document_success_class",
			Name: "Get Document Success Class",
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))

		doc := Document{
			Id:      "get_document_by_id_success",
			ClassId: class.Id,
			Name:    t.Name(),
		}
		assert.NoError(t, repo.CreateDocument(ctx, &doc))

		check, err := repo.GetDocumentById(ctx, doc.Id)
		assert.NoError(t, err)
		assert.Equal(t, doc.Id, check.Id)
		assert.Equal(t, doc.ClassId, check.ClassId)
		assert.Equal(t, doc.Name, check.Name)
	})

	t.Run("GetDocumentByIdFail", func(t *testing.T) {
		_, err := repo.GetDocumentById(ctx, "bad_doc_id")
		assert.Equal(t, ErrNotExist, err)
	})

	t.Run("GetDocumentByPath", func(t *testing.T) {
		class := Class{
			Id:   "document_by_path_class",
			Name: "Get Document By Path Class",
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))

		doc := Document{
			Id:      "document_by_path",
			Name:    "Document By Path",
			ClassId: class.Id,
			Path:    "/document/by/path",
		}
		assert.NoError(t, repo.CreateDocument(ctx, &doc))

		check, err := repo.GetDocumentByPath(ctx, doc.Path)
		assert.NoError(t, err)
		assert.Equal(t, doc.Path, check.Path)

		_, err = repo.GetDocumentByPath(ctx, "/invalid/path")
		assert.Equal(t, ErrNotExist, err)
	})

	t.Run("UpdateDocument", func(t *testing.T) {
		class := Class{
			Id:   "update_document_class",
			Name: "Update Document Class",
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))

		t.Run("NoPathNoPath", func(t *testing.T) {
			doc := Document{
				Id:      "no_path_no_path",
				ClassId: class.Id,
				Name:    t.Name(),
				Path:    "",
			}
			assert.NoError(t, repo.CreateDocument(ctx, &doc))

			doc.Name += "-Updated"
			assert.NoError(t, repo.UpdateDocument(ctx, &doc))
		})

		t.Run("NoPathYesPath", func(t *testing.T) {
			doc := Document{
				Id:      "no_path_yes_path",
				ClassId: class.Id,
				Name:    t.Name(),
				Path:    "",
			}
			assert.NoError(t, repo.CreateDocument(ctx, &doc))

			doc.Name += "-Updated"
			doc.Path = "/no/path/yes/path"
			assert.NoError(t, repo.UpdateDocument(ctx, &doc))
		})

		t.Run("YesPathNoPath", func(t *testing.T) {
			doc := Document{
				Id:      "yes_path_no_path",
				ClassId: class.Id,
				Name:    t.Name(),
				Path:    "/yes/path/no/path",
			}
			assert.NoError(t, repo.CreateDocument(ctx, &doc))

			doc.Name += "-Updated"
			doc.Path = ""
			assert.NoError(t, repo.UpdateDocument(ctx, &doc))
		})

		t.Run("YesPathYesPath", func(t *testing.T) {
			doc := Document{
				Id:      "yes_path_yes_path",
				ClassId: class.Id,
				Name:    t.Name(),
				Path:    "/yes/path/yes/path",
			}
			assert.NoError(t, repo.CreateDocument(ctx, &doc))

			doc.Name += "-Updated"
			doc.Path = "/yes/path/yes/path/updated"
			assert.NoError(t, repo.UpdateDocument(ctx, &doc))
		})

		t.Run("ForceTableScan", func(t *testing.T) {
			doc := Document{
				Id:      "force_table_scan",
				ClassId: class.Id,
				Name:    t.Name(),
				Path:    "/force/scan/original",
			}
			assert.NoError(t, repo.CreateDocument(ctx, &doc))

			// Manually override the path
			pk, _ := attributevalue.Marshal(dynamoDocPrefix + doc.Id)
			sk, _ := attributevalue.Marshal(fmt.Sprintf(dynamoDocSortF, 0))
			path, _ := attributevalue.Marshal("/force/scan/override")
			client := dynamodb.NewFromConfig(cfg)
			_, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
				TableName: &resources.Table,
				Key: map[string]types.AttributeValue{
					"PK": pk,
					"SK": sk,
				},
				UpdateExpression: aws.String("SET #path = :path"),
				ExpressionAttributeNames: map[string]string{
					"#path": "Path",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":path": path,
				},
			})
			assert.NoError(t, err)

			// Force table scan during update
			doc.Name += "-Updated"
			doc.Path = "/force/scan/updated"
			assert.NoError(t, repo.UpdateDocument(ctx, &doc))
		})
	})

	t.Run("DeleteDocument", func(t *testing.T) {
		class := Class{
			Id:   "delete_document_class",
			Name: "Delete Document Class",
			Fields: []Field{
				{Name: "field_1", Sort: true},
				{Name: "field_2", Sort: true},
				{Name: "field_3", Sort: true},
			},
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))

		doc := Document{
			Id:      "delete_me",
			Name:    "Delete Me",
			ClassId: class.Id,
			Path:    "/delete/me",
			Values: map[string]interface{}{
				"field_1": "My first value",
				"field_2": "My second value",
				"field_3": []int{1, 2, 3},
			},
		}
		assert.NoError(t, repo.CreateDocument(ctx, &doc))

		assert.NoError(t, repo.DeleteDocument(ctx, doc.Id))
	})
}

func TestDynamoDBRepositoryDocumentList(t *testing.T) {
	resources := DynamoDBResources{
		Table: dynamoTablePrefix + "List",
	}
	cfg, err := testDynamoConfig(resources.Table)
	assert.NoError(t, err)

	repo := NewDynamoDBRepository(cfg, resources)
	ctx := context.Background()

	for _, class := range testClasses() {
		assert.NoError(t, repo.CreateClass(ctx, &class))
	}

	for _, document := range testDocuments() {
		assert.NoError(t, repo.CreateDocument(ctx, &document))
	}

	t.Run("ListPagesByTitle", func(t *testing.T) {
		filter := DocumentFilter{
			ClassId: "page",
			Field:   "title",
			Range:   Range{End: 9},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, Range{End: 1, Size: 2}, r)
		assert.Equal(t, 2, len(docs))
		assert.Equal(t, "page-2", docs[0].Id)
		assert.Equal(t, "page-1", docs[1].Id)
	})

	t.Run("ListSessionsByStart", func(t *testing.T) {
		filter := DocumentFilter{
			ClassId: "session",
			Field:   "start",
			Range:   Range{End: 9},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, Range{End: 4, Size: 5}, r)
		assert.Equal(t, 5, len(docs))
	})

	t.Run("ListSessionsByEvent", func(t *testing.T) {
		filter := DocumentFilter{
			ClassId:  "session",
			ParentId: "event-1",
			Field:    "start",
			Range:    Range{End: 9},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, Range{End: 2, Size: 3}, r)
		assert.Equal(t, 3, len(docs))
	})
}
