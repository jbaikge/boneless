import { RemovalPolicy, Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as apigw from 'aws-cdk-lib/aws-apigateway';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import { join } from 'path';
import { LambdaIntegration } from 'aws-cdk-lib/aws-apigateway';


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

    // REST API
    const api = new apigw.RestApi(this, 'GoCMS Endpoint', {});

    // Class endpoints
    const classResource = api.root.addResource('class');
    classResource.addMethod('POST', new LambdaIntegration(createClassLambda));

    const classIdResource = classResource.addResource('{id}')
    classIdResource.addMethod('GET', new LambdaIntegration(getClassByIdLambda));
  }
}
