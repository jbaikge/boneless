import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as apigateway from '@aws-cdk/aws-apigatewayv2-alpha';
import * as integration from '@aws-cdk/aws-apigatewayv2-integrations-alpha';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as s3 from 'aws-cdk-lib/aws-s3';
import { join } from 'path';
import { BucketDeployment, Source } from 'aws-cdk-lib/aws-s3-deployment';
import { Distribution, OriginAccessIdentity } from 'aws-cdk-lib/aws-cloudfront';
import { S3Origin } from 'aws-cdk-lib/aws-cloudfront-origins';
import { CfnOutput } from 'aws-cdk-lib';


export class GocmsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // DynamoDB table
    const repositoryTable = new dynamodb.Table(this, 'RepositoryTable', {
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      // Not used when billingMode is PAY_PER_REQUEST
      // readCapacity: 1,
      // writeCapacity: 1,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      partitionKey: {
        name: 'PK',
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: 'SK',
        type: dynamodb.AttributeType.STRING,
      },
    });

    // S3 Repository Bucket (document values and templates)
    const repositoryBucket = new s3.Bucket(this, 'RepositoryBucket', {
      accessControl: s3.BucketAccessControl.PRIVATE,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    // Environment variables
    const environment = {
      'REPOSITORY_BUCKET': repositoryBucket.bucketName,
      'REPOSITORY_TABLE': repositoryTable.tableName,
    }

    // Path back to repo root
    const rootDir = join(__dirname, '..', '..');

    // Asset directory where all the lambda binaries come from
    const assetDir = join(rootDir, 'assets');

    // Admin handler
    const adminLambda = new lambda.Function(this, 'AdminHandler', {
      environment: environment,
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset(join(assetDir, 'admin')),
      handler: 'handler',
    });
    repositoryBucket.grantReadWrite(adminLambda);
    repositoryTable.grantReadWriteData(adminLambda);

    /*
    // Create class lambda function
    const createClassLambda = new lambda.Function(this, 'CreateClassHandler', {
      environment: environment,
      runtime:     lambda.Runtime.GO_1_X,
      code:        lambda.Code.fromAsset(join(assetDir, 'create-class')),
      handler:     'handler'
    });
    repositoryTable.grantWriteData(createClassLambda);

    // Get class lambda function
    const getClassByIdLambda = new lambda.Function(this, 'GetClassByIdHandler', {
      environment: environment,
      runtime:     lambda.Runtime.GO_1_X,
      code:        lambda.Code.fromAsset(join(assetDir, 'get-class-by-id')),
      handler:     'handler'
    });
    repositoryTable.grantReadData(getClassByIdLambda);

    // List class lambda function
    const listClassesLambda = new lambda.Function(this, 'ListClassesHandler', {
      environment: environment,
      runtime:     lambda.Runtime.GO_1_X,
      code:        lambda.Code.fromAsset(join(assetDir, 'list-classes')),
      handler:     'handler'
    });
    repositoryTable.grantReadData(listClassesLambda);

    // Update class lambda function
    const updateClassLambda = new lambda.Function(this, 'UpdateClassHandler', {
      environment: environment,
      runtime:     lambda.Runtime.GO_1_X,
      code:        lambda.Code.fromAsset(join(assetDir, 'update-class')),
      handler:     'handler'
    });
    repositoryTable.grantWriteData(updateClassLambda);
    */

    // // REST API (v1)
    // const api = new apigw.RestApi(this, 'GoCMS Endpoint', {
    //   defaultCorsPreflightOptions: {
    //     allowOrigins: apigw.Cors.ALL_ORIGINS,
    //     allowHeaders: [
    //       'Content-Type',
    //       'Range',
    //       'Authorization'
    //     ],
    //     exposeHeaders: [
    //       'Content-Range',
    //       'X-Total-Count',
    //     ],
    //   }
    // });

    // Class endpoints
    // Allow the Range header with requests for pagination
    // 2022-07-22: After days of this not working, just leaving URLs here for
    // posterity in case another attempt is made.
    // Ref: https://rahullokurte.com/how-to-validate-requests-to-the-aws-api-gateway-using-cdk
    // Ref: https://stackoverflow.com/a/68305757
    // Ref: https://blog.kewah.com/2020/api-gateway-caching-with-aws-cdk/
    /*
    const classResource = api.root.addResource('classes');
    classResource.addMethod('GET', new apigw.LambdaIntegration(listClassesLambda));
    classResource.addMethod('POST', new apigw.LambdaIntegration(createClassLambda));

    const classIdResource = classResource.addResource('{id}')
    classIdResource.addMethod('GET', new apigw.LambdaIntegration(getClassByIdLambda));
    classIdResource.addMethod('PUT', new apigw.LambdaIntegration(updateClassLambda));
    */

    const api = new apigateway.HttpApi(this, 'Endpoint', {
      createDefaultStage: true,
      corsPreflight: {
        allowOrigins: [
          '*',
        ],
        allowHeaders: [
          'Content-Type',
          'Range',
          'Authorization',
        ],
        exposeHeaders: [
          'Content-Range',
          'X-Total-Count',
        ],
      },
    });

    new CfnOutput(this, 'EndpointUrl', { value: api.url! })

    const adminIntegration = new integration.HttpLambdaIntegration('AdminIntegration', adminLambda)
    api.addRoutes({
      path: '/classes',
      integration: adminIntegration,
      methods: [
        apigateway.HttpMethod.GET,
        apigateway.HttpMethod.POST,
      ],
    });

    api.addRoutes({
      path: '/classes/{class_id}',
      integration: adminIntegration,
      methods: [
        apigateway.HttpMethod.GET,
        apigateway.HttpMethod.PUT,
        apigateway.HttpMethod.DELETE,
      ],
    });

    api.addRoutes({
      path: '/classes/{class_id}/documents',
      integration: adminIntegration,
      methods: [
        apigateway.HttpMethod.GET,
        apigateway.HttpMethod.POST,
      ],
    });

    api.addRoutes({
      path: '/classes/{class_id}/documents/{doc_id}',
      integration: adminIntegration,
      methods: [
        apigateway.HttpMethod.GET,
        apigateway.HttpMethod.PUT,
        apigateway.HttpMethod.DELETE,
      ],
    });

    api.addRoutes({
      path: '/documents',
      integration: adminIntegration,
      methods: [
        apigateway.HttpMethod.GET,
        apigateway.HttpMethod.POST,
      ],
    });

    api.addRoutes({
      path: '/documents/{doc_id}',
      integration: adminIntegration,
      methods: [
        apigateway.HttpMethod.GET,
        apigateway.HttpMethod.PUT,
        apigateway.HttpMethod.DELETE,
      ],
    });

    /*
    const adminIntegration = new apigw.LambdaIntegration(adminLambda);
    const classResource = api.root.addResource('classes');
    classResource.addMethod('GET', adminIntegration);
    classResource.addMethod('POST', adminIntegration);

    const classItemResource = classResource.addResource('{class_id}');
    classItemResource.addMethod('GET', adminIntegration);
    classItemResource.addMethod('PUT', adminIntegration);
    classItemResource.addMethod('DELETE', adminIntegration);

    const classDocumentResource = classItemResource.addResource('documents');
    classDocumentResource.addMethod('GET', adminIntegration);
    classDocumentResource.addMethod('POST', adminIntegration);

    const classDocumentItemResource = classDocumentResource.addResource('{doc_id}')
    classDocumentItemResource.addMethod('GET', adminIntegration);
    classDocumentItemResource.addMethod('PUT', adminIntegration);
    classDocumentItemResource.addMethod('DELETE', adminIntegration);

    const documentResource = api.root.addResource('documents')
    documentResource.addMethod('GET', adminIntegration);
    documentResource.addMethod('POST', adminIntegration);

    const documentItemResource = documentResource.addResource('{doc_id}')
    documentItemResource.addMethod('GET', adminIntegration);
    documentItemResource.addMethod('PUT', adminIntegration);
    documentItemResource.addMethod('DELETE', adminIntegration);
    */

    // Admin frontend
    // https://aws-cdk.com/deploying-a-static-website-using-s3-and-cloudfront/
    /*
    const adminBucket = new s3.Bucket(this, 'FrontendAdminBucket', {
      accessControl: s3.BucketAccessControl.PRIVATE,
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
    */

    new cdk.CfnOutput(this, 'repository_bucket', {
      value: repositoryBucket.bucketName,
      description: 'Repository S3 bucket',
    });
  }
}
