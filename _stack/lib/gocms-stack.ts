import { RemovalPolicy, Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as apigw from 'aws-cdk-lib/aws-apigateway';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import { join } from 'path';
import { Bucket, BucketAccessControl } from 'aws-cdk-lib/aws-s3';
import { BucketDeployment, Source } from 'aws-cdk-lib/aws-s3-deployment';
import { Distribution, OriginAccessIdentity } from 'aws-cdk-lib/aws-cloudfront';
import { S3Origin } from 'aws-cdk-lib/aws-cloudfront-origins';


export class GocmsStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    // DynamoDB tables
    const classTable = new dynamodb.Table(this, 'Classes', {
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      // Not used when billingMode is PAY_PER_REQUEST
      // readCapacity: 1,
      // writeCapacity: 1,
      removalPolicy: RemovalPolicy.DESTROY,
      partitionKey: {
        name: 'ClassId',
        type: dynamodb.AttributeType.STRING,
      },
    });

    const docTable = new dynamodb.Table(this, 'Documents', {
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: RemovalPolicy.DESTROY,
      partitionKey: {
        name: 'DocumentId',
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: 'Version',
        type: dynamodb.AttributeType.NUMBER,
      },
    });

    docTable.addGlobalSecondaryIndex({
      indexName: 'GSI-Class',
      partitionKey: {
        name: 'ClassId',
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: 'Version',
        type: dynamodb.AttributeType.NUMBER,
      },
      projectionType: dynamodb.ProjectionType.ALL,
    });

    docTable.addGlobalSecondaryIndex({
      indexName: 'GSI-Parent',
      partitionKey: {
        name: 'ParentId',
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: 'Version',
        type: dynamodb.AttributeType.NUMBER,
      },
      projectionType: dynamodb.ProjectionType.ALL,
    });

    const sortTable = new dynamodb.Table(this, 'Sort', {
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: RemovalPolicy.DESTROY,
      partitionKey: {
        name: 'ClassField',
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: 'Value',
        type: dynamodb.AttributeType.STRING,
      },
    });

    sortTable.addGlobalSecondaryIndex({
      indexName: 'GSI-Documents',
      partitionKey: {
        name: 'DocumentId',
        type: dynamodb.AttributeType.STRING,
      },
    });

    // Path back to repo root
    const rootDir = join(__dirname, '..', '..');

    // Asset directory where all the lambda binaries come from
    const assetDir = join(rootDir, 'assets');

    // Create class lambda function
    const createClassLambda = new lambda.Function(this, 'CreateClassHandler', {
      environment: {
        'DYNAMODB_CLASS_TABLE': classTable.tableName,
      },
      runtime: lambda.Runtime.GO_1_X,
      code:    lambda.Code.fromAsset(join(assetDir, 'create-class')),
      handler: 'handler'
    });
    classTable.grantWriteData(createClassLambda);

    // Get class lambda function
    const getClassByIdLambda = new lambda.Function(this, 'GetClassByIdHandler', {
      environment: {
        'DYNAMODB_CLASS_TABLE': classTable.tableName,
      },
      runtime: lambda.Runtime.GO_1_X,
      code:    lambda.Code.fromAsset(join(assetDir, 'get-class-by-id')),
      handler: 'handler'
    });
    classTable.grantReadData(getClassByIdLambda);

    // List class lambda function
    const listClassesLambda = new lambda.Function(this, 'ListClassesHandler', {
      environment: {
        'DYNAMODB_CLASS_TABLE': classTable.tableName,
      },
      runtime: lambda.Runtime.GO_1_X,
      code:    lambda.Code.fromAsset(join(assetDir, 'list-classes')),
      handler: 'handler'
    });
    classTable.grantReadData(listClassesLambda);

    // Update class lambda function
    const updateClassLambda = new lambda.Function(this, 'UpdateClassHandler', {
      environment: {
        'DYNAMODB_CLASS_TABLE': classTable.tableName,
      },
      runtime: lambda.Runtime.GO_1_X,
      code:    lambda.Code.fromAsset(join(assetDir, 'update-class')),
      handler: 'handler'
    });
    classTable.grantWriteData(updateClassLambda);

    // REST API
    const api = new apigw.RestApi(this, 'GoCMS Endpoint', {
      defaultCorsPreflightOptions: {
        allowOrigins: apigw.Cors.ALL_ORIGINS,
        allowHeaders: [
          'Content-Type',
          'Range',
          'Authorization'
        ],
        exposeHeaders: [
          'Content-Range',
          'X-Total-Count',
        ],
      }
    });

    // Class endpoints
    // Allow the Range header with requests for pagination
    // 2022-07-22: After days of this not working, just leaving URLs here for
    // posterity in case another attempt is made.
    // Ref: https://rahullokurte.com/how-to-validate-requests-to-the-aws-api-gateway-using-cdk
    // Ref: https://stackoverflow.com/a/68305757
    // Ref: https://blog.kewah.com/2020/api-gateway-caching-with-aws-cdk/
    const classResource = api.root.addResource('classes');
    classResource.addMethod('GET', new apigw.LambdaIntegration(listClassesLambda));
    classResource.addMethod('POST', new apigw.LambdaIntegration(createClassLambda));

    const classIdResource = classResource.addResource('{id}')
    classIdResource.addMethod('GET', new apigw.LambdaIntegration(getClassByIdLambda));
    classIdResource.addMethod('PUT', new apigw.LambdaIntegration(updateClassLambda));

    // Admin frontend
    // https://aws-cdk.com/deploying-a-static-website-using-s3-and-cloudfront/
    const adminBucket = new Bucket(this, 'FrontendAdminBucket', {
      accessControl: BucketAccessControl.PRIVATE,
    });

    new BucketDeployment(this, 'FrontendAdminBucketDeployment', {
      destinationBucket: adminBucket,
      sources: [Source.asset(join(rootDir, '_frontend-admin', 'build'))],
    });

    const adminOriginAccessIdentity = new OriginAccessIdentity(this, 'FrontendAdminOAI');
    adminBucket.grantRead(adminOriginAccessIdentity);

    new Distribution(this, 'FrontendAdminDistribution', {
      defaultRootObject: 'index.html',
      defaultBehavior: {
        origin: new S3Origin(adminBucket, {originAccessIdentity: adminOriginAccessIdentity}),
      }
    });
  }
}
