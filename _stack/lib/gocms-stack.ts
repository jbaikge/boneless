import { RemovalPolicy, Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as apigw from 'aws-cdk-lib/aws-apigateway';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import { join } from 'path';


export class GocmsStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    // DynamoDB table
    const table = new dynamodb.Table(this, 'GoCMSTable', {
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      // Not used when billingMode is PAY_PER_REQUEST
      // readCapacity: 1,
      // writeCapacity: 1,
      removalPolicy: RemovalPolicy.DESTROY,
      partitionKey: {
        name: 'PrimaryKey',
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: 'SortKey',
        type: dynamodb.AttributeType.STRING
      },
    });

    // Asset directory where all the lambda binaries come from
    const assetDir = join(__dirname, '..', '..', 'assets');

    // Create class lambda function
    const createClassLambda = new lambda.Function(this, 'CreateClassHandler', {
      environment: {
        'DYNAMODB_TABLE': table.tableName
      },
      runtime: lambda.Runtime.GO_1_X,
      code:    lambda.Code.fromAsset(join(assetDir, 'create-class')),
      handler: 'handler'
    });
    table.grantWriteData(createClassLambda);

    // Get class lambda function
    const getClassByIdLambda = new lambda.Function(this, 'GetClassByIdHandler', {
      environment: {
        'DYNAMODB_TABLE': table.tableName
      },
      runtime: lambda.Runtime.GO_1_X,
      code:    lambda.Code.fromAsset(join(assetDir, 'get-class-by-id')),
      handler: 'handler'
    });
    table.grantReadData(getClassByIdLambda);

    // List class lambda function
    const listClassesLambda = new lambda.Function(this, 'ListClassesHandler', {
      environment: {
        'DYNAMODB_TABLE': table.tableName
      },
      runtime: lambda.Runtime.GO_1_X,
      code:    lambda.Code.fromAsset(join(assetDir, 'list-classes')),
      handler: 'handler'
    });
    table.grantReadData(listClassesLambda);

    // Update class lambda function
    const updateClassLambda = new lambda.Function(this, 'UpdateClassHandler', {
      environment: {
        'DYNAMODB_TABLE': table.tableName
      },
      runtime: lambda.Runtime.GO_1_X,
      code:    lambda.Code.fromAsset(join(assetDir, 'update-class')),
      handler: 'handler'
    });
    table.grantWriteData(updateClassLambda);

    // REST API
    const api = new apigw.RestApi(this, 'GoCMS Endpoint', {
      defaultCorsPreflightOptions: {
        allowOrigins: apigw.Cors.ALL_ORIGINS,
        allowHeaders: [
          'Content-Type',
          'Range',
          'Authorization'
        ]
      }
    });

    // Class endpoints
    const classResource = api.root.addResource('classes');
    classResource.addMethod('GET', new apigw.LambdaIntegration(listClassesLambda));
    classResource.addMethod('POST', new apigw.LambdaIntegration(createClassLambda));

    const classIdResource = classResource.addResource('{id}')
    classIdResource.addMethod('GET', new apigw.LambdaIntegration(getClassByIdLambda));
    classIdResource.addMethod('PUT', new apigw.LambdaIntegration(updateClassLambda));
  }
}
