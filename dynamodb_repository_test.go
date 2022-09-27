package boneless

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/zeebo/assert"
)

var (
	dynamoPrefix  = time.Now().Format("20060102-150405-")
	useLocalStack = false
)

func init() {
	flag.BoolVar(&useLocalStack, "localstack", useLocalStack, "Force path style URLs for LocalStack compatibility")
}

func testDynamoNewRepository(resources DynamoDBResources) (repo *DynamoDBRepository, err error) {
	endpointResolverFunc := func(service string, region string, options ...interface{}) (endpoint aws.Endpoint, err error) {
		endpoint = aws.Endpoint{
			PartitionID:   "aws",
			URL:           "http://localhost:4566", // 4566 for LocalStack; 8000 for amazon/dynamodb-local
			SigningRegion: "us-east-1",             // Must be a legitimate region for LocalStack S3 to work
		}
		return
	}
	endpointResolver := aws.EndpointResolverWithOptionsFunc(endpointResolverFunc)

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithEndpointResolverWithOptions(endpointResolver),
	)
	if err != nil {
		return
	}

	// Build table before returning
	createTable := &dynamodb.CreateTableInput{
		TableName:   &resources.Table,
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

	db := dynamodb.NewFromConfig(cfg)
	if _, err = db.CreateTable(context.Background(), createTable); err != nil {
		return
	}

	// Create S3 bucket
	createBucket := &s3.CreateBucketInput{
		Bucket: &resources.Bucket,
	}

	// UsePathStyle is required to prevent host lookup exception
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = useLocalStack })
	if _, err = s3Client.CreateBucket(context.Background(), createBucket); err != nil {
		return
	}

	return &DynamoDBRepository{
		db:        db,
		s3:        s3Client,
		resources: resources,
	}, nil
}

func testDynamoEmptyTable(repo *DynamoDBRepository, table string) (err error) {
	params := &dynamodb.ScanInput{
		TableName: &table,
	}
	response, err := repo.db.Scan(context.Background(), params)
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
		if _, err = repo.db.DeleteItem(context.Background(), delete); err != nil {
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
		Bucket: dynamoPrefix + "test",
		Table:  dynamoPrefix + "Test",
	}
	repo, err := testDynamoNewRepository(resources)
	assert.NoError(t, err)

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
			Id:      "update_class",
			Name:    t.Name(),
			Created: time.Now(),
			Updated: time.Now(),
			Fields:  []Field{{Name: "field_2"}, {Name: "field_1"}},
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))

		class.Name = t.Name() + "-Updated"
		class.Updated = time.UnixMicro(time.Now().UnixMicro())
		class.Fields = append(class.Fields, Field{Name: "field_3"})
		assert.NoError(t, repo.UpdateClass(ctx, &class))

		check, err := repo.GetClassById(ctx, class.Id)
		assert.NoError(t, err)
		assert.Equal(t, class.Name, check.Name)
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
		assert.NoError(t, testDynamoEmptyTable(repo, resources.Table))

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
			Path:    "/doc/with/path",
		}
		assert.NoError(t, repo.CreateDocument(ctx, &doc))
	})

	t.Run("CreateDocumentWithSort", func(t *testing.T) {
		assert.NoError(t, testDynamoEmptyTable(repo, resources.Table))

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
		}
		assert.NoError(t, repo.CreateDocument(ctx, &doc))

		check, err := repo.GetDocumentById(ctx, doc.Id)
		assert.NoError(t, err)
		assert.Equal(t, doc.Id, check.Id)
		assert.Equal(t, doc.ClassId, check.ClassId)
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
			Fields: []Field{
				{Name: "field1", Sort: true},
				{Name: "field2"},
				{Name: "field3", Sort: true},
			},
		}
		assert.NoError(t, repo.CreateClass(ctx, &class))

		t.Run("NoPathNoPath", func(t *testing.T) {
			doc := Document{
				Id:      "no_path_no_path",
				ClassId: class.Id,
				Path:    "",
			}
			assert.NoError(t, repo.CreateDocument(ctx, &doc))

			assert.NoError(t, repo.UpdateDocument(ctx, &doc))
		})

		t.Run("NoPathYesPath", func(t *testing.T) {
			doc := Document{
				Id:      "no_path_yes_path",
				ClassId: class.Id,
				Path:    "",
			}
			assert.NoError(t, repo.CreateDocument(ctx, &doc))

			doc.Path = "/no/path/yes/path"
			assert.NoError(t, repo.UpdateDocument(ctx, &doc))
		})

		t.Run("YesPathNoPath", func(t *testing.T) {
			doc := Document{
				Id:      "yes_path_no_path",
				ClassId: class.Id,
				Path:    "/yes/path/no/path",
			}
			assert.NoError(t, repo.CreateDocument(ctx, &doc))

			doc.Path = ""
			assert.NoError(t, repo.UpdateDocument(ctx, &doc))
		})

		t.Run("YesPathYesPath", func(t *testing.T) {
			doc := Document{
				Id:      "yes_path_yes_path",
				ClassId: class.Id,
				Path:    "/yes/path/yes/path",
			}
			assert.NoError(t, repo.CreateDocument(ctx, &doc))

			doc.Path = "/yes/path/yes/path/updated"
			assert.NoError(t, repo.UpdateDocument(ctx, &doc))
		})

		t.Run("ForceTableScan", func(t *testing.T) {
			doc := Document{
				Id:      "force_table_scan",
				ClassId: class.Id,
				Path:    "/force/scan/original",
			}
			assert.NoError(t, repo.CreateDocument(ctx, &doc))

			// Manually override the path - Not possible with the repo API
			pk, _ := attributevalue.Marshal(dynamoDocPrefix + doc.Id)
			sk, _ := attributevalue.Marshal(fmt.Sprintf(dynamoDocSortF, 0))
			path, _ := attributevalue.Marshal("/force/scan/override")
			_, err := repo.db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
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
			doc.Path = "/force/scan/updated"
			assert.NoError(t, repo.UpdateDocument(ctx, &doc))
		})

		t.Run("SortUpdates", func(t *testing.T) {
			doc := Document{
				Id:      "sort_update",
				ClassId: class.Id,
				Values: map[string]interface{}{
					"field1": "v1",
					"field2": "v1",
					"field3": "v1",
				},
			}
			assert.NoError(t, repo.CreateDocument(ctx, &doc))

			doc.Values["field1"] = "v2"
			doc.Values["field2"] = "v2"
			doc.Values["field3"] = "v2"
			assert.NoError(t, repo.UpdateDocument(ctx, &doc))

			check, err := repo.GetDocumentById(ctx, doc.Id)
			assert.NoError(t, err)
			assert.DeepEqual(t, doc.Values, check.Values)
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
		Bucket: dynamoPrefix + "list",
		Table:  dynamoPrefix + "List",
	}
	repo, err := testDynamoNewRepository(resources)
	assert.NoError(t, err)

	ctx := context.Background()

	for _, class := range ClassData() {
		assert.NoError(t, repo.CreateClass(ctx, &class))
	}

	for _, document := range DocumentData() {
		assert.NoError(t, repo.CreateDocument(ctx, &document))
	}

	t.Run("ListPagesByTitle", func(t *testing.T) {
		filter := DocumentFilter{
			ClassId: "page",
			Sort:    DocumentFilterSort{Field: "title"},
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
			Sort:    DocumentFilterSort{Field: "start"},
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
			Sort:     DocumentFilterSort{Field: "start"},
			Range:    Range{End: 9},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, Range{End: 2, Size: 3}, r)
		assert.Equal(t, 3, len(docs))
	})

	t.Run("EmptyFilter", func(t *testing.T) {
		// Should list all documents, sorted by descending creation date
		filter := DocumentFilter{
			Range: Range{End: 99},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, Range{End: 19, Size: 20}, r)
		assert.Equal(t, 20, len(docs))
		assert.Equal(t, "speaker-6", docs[0].Id)
		assert.Equal(t, "event-1", docs[19].Id)
	})

	t.Run("AllChildren", func(t *testing.T) {
		// Should give all children of a parent document, regardless of class
		// then sorted by descending creation date
		filter := DocumentFilter{
			ParentId: "event-1",
			Range:    Range{End: 99},
		}
		docs, r, err := repo.GetDocumentList(ctx, filter)
		assert.NoError(t, err)
		assert.DeepEqual(t, Range{End: 2, Size: 3}, r)
		assert.Equal(t, 3, len(docs))
		assert.Equal(t, "session-3", docs[0].Id)
		assert.Equal(t, "session-1", docs[2].Id)
	})
}

func TestDynamoDBRepositoryValues(t *testing.T) {
	resources := DynamoDBResources{
		Bucket: dynamoPrefix + "values",
		Table:  dynamoPrefix + "Values",
	}
	repo, err := testDynamoNewRepository(resources)
	assert.NoError(t, err)

	ctx := context.Background()

	class := Class{
		Id:      "class",
		Name:    "Class",
		Created: time.Now(),
		Updated: time.Now(),
		Fields: []Field{
			{Name: "field1", Sort: true},
			{Name: "field2"},
			{Name: "field3", Sort: true},
			{Name: "field4"},
		},
	}
	assert.NoError(t, repo.CreateClass(ctx, &class))

	doc := Document{
		Id:      "doc",
		ClassId: "class",
		Values: map[string]interface{}{
			"field1": "abc",
			"field2": "def",
			"field3": "ghi",
			"field4": "jkl",
		},
	}
	assert.NoError(t, repo.CreateDocument(ctx, &doc))

	docCheck, err := repo.GetDocumentById(ctx, doc.Id)
	assert.NoError(t, err)

	assert.DeepEqual(t, doc.Values, docCheck.Values)
	// Can't use DeepEqual for ... reasons?
	// Need to cast everything to string for assert.Equal
	for key, expect := range doc.Values {
		assert.Equal(t, fmt.Sprint(expect), fmt.Sprint(docCheck.Values[key]))
	}

	extra := Document{
		Id:      "extra",
		ClassId: "class",
		Values: map[string]interface{}{
			"field0": "abc", // extra
			"field1": "def",
			// field2 intentionally left out
			"field3": "ghi",
			"field4": "jkl",
			"field5": "mno", // extra
		},
	}
	assert.NoError(t, repo.CreateDocument(ctx, &extra))

	extraCheck, err := repo.GetDocumentById(ctx, extra.Id)
	assert.NoError(t, err)

	// Should contain all the submitted values and match what was sent into the
	// repo.
	assert.DeepEqual(t, extra.Values, extraCheck.Values)
}

func TestDynamoDBRepositoryTemplates(t *testing.T) {
	resources := DynamoDBResources{
		Bucket: dynamoPrefix + "templates",
		Table:  dynamoPrefix + "Templates",
	}
	repo, err := testDynamoNewRepository(resources)
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("BasicInOut", func(t *testing.T) {
		template := Template{
			Id:      "template-1",
			Name:    "Test Template",
			Created: time.Now(),
			Updated: time.Now(),
			Body:    "This is a test template\n\nThis is another line in the test template",
		}
		assert.NoError(t, repo.CreateTemplate(ctx, &template))

		check, err := repo.GetTemplateById(ctx, template.Id)
		assert.NoError(t, err)
		assert.Equal(t, template.Body, check.Body)
	})

	t.Run("Version2", func(t *testing.T) {
		template := Template{
			Id:      "two-versions",
			Name:    "Version 1",
			Created: time.Now(),
			Updated: time.Now(),
			Body:    "This is version one content",
		}
		assert.NoError(t, repo.CreateTemplate(ctx, &template))

		template.Name = "Version 2"
		template.Body = "This is version two content"
		assert.NoError(t, repo.UpdateTemplate(ctx, &template))

		check, err := repo.GetTemplateById(ctx, template.Id)
		assert.NoError(t, err)
		assert.Equal(t, template.Body, check.Body)
	})

	t.Run("DeleteMultiVersion", func(t *testing.T) {
		template := Template{
			Id:      "delete-me",
			Name:    "Delete Me v1",
			Created: time.Now(),
			Updated: time.Now(),
			Body:    "Delete this content v1",
		}
		assert.NoError(t, repo.CreateTemplate(ctx, &template))

		template.Name = "Delete Me v2"
		template.Body = template.Body[len(template.Body)-2:] + "2"
		assert.NoError(t, repo.UpdateTemplate(ctx, &template))

		assert.NoError(t, repo.DeleteTemplate(ctx, template.Id))

		_, err := repo.GetTemplateById(ctx, template.Id)
		assert.Equal(t, ErrNotExist, err)
	})
}
