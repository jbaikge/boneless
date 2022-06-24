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

    const table = new dynamodb.Table(this, 'GoCMSTable', {
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
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


    const assetDir = join(__dirname, '..', '..', 'assets');

    const pingPongFunction = new lambda.Function(this, 'PingPongHandler', {
      environment: {
        'DYNAMODB_TABLE': table.tableName
      },
      runtime: lambda.Runtime.GO_1_X,
      code:    lambda.Code.fromAsset(join(assetDir, 'ping-pong')),
      handler: 'handler'
    });
    table.grantReadWriteData(pingPongFunction);

    const api = new apigw.RestApi(this, 'GoCMS Endpoint', {});

    const pingPongResource = api.root.addResource('ping-pong');
    pingPongResource.addMethod('GET', new LambdaIntegration(pingPongFunction));
  }
}
