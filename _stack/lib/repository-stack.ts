import * as cdk from 'aws-cdk-lib';
import * as constructs from 'constructs';
import * as apigateway from '@aws-cdk/aws-apigatewayv2-alpha';
import * as integration from '@aws-cdk/aws-apigatewayv2-integrations-alpha';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as path from 'path';
import exec from './exec';


export class RepositoryStack extends cdk.Stack {
  constructor(scope: constructs.Construct, id: string, props?: cdk.StackProps) {
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

    // Path back to repo root
    const rootDir = path.resolve(__dirname, '..', '..');

    // Admin handler
    // Ref: https://github.com/aws-samples/cdk-build-bundle-deploy-example/blob/main/cdk-bundle-go-lambda-example/lib/api-stack.ts
    const goEnvironment = {
      CGO_ENABLED: '0',
      GOOS:        'linux',
      GOARCH:      'amd64',
    };
    const adminLambdaDir = path.join(rootDir, 'lambda', 'admin');
    const adminLambda = new lambda.Function(this, 'AdminHandler', {
      code: lambda.Code.fromAsset(adminLambdaDir, {
        bundling: {
          local: {
            tryBundle(outputDir: string) {
              try {
                // build the binary
                exec(`go build -o ${path.join(outputDir, 'handler')}`, {
                  cwd: adminLambdaDir,
                  env: { ...process.env, ...goEnvironment },
                  stdio: [
                    'ignore',       // ignore stdio
                    process.stderr, // redirect stdout to stderr
                    'inherit',      // inherit stderr
                  ],
                })
              } catch (error) {
                // if we don't have go installed return false which
                // tells the CDK to try Docker bundling
                return false;
              }

              return true;
            },
          },
          image: lambda.Runtime.GO_1_X.bundlingImage, // lambci/lambda:build-go1.x
          command: [
            'bash', '-c', [
              'cd /asset-input',
              'go build -o /asset-output/handler',
            ].join(' && '),
          ],
          environment: goEnvironment,
        }
      }),
      runtime: lambda.Runtime.GO_1_X,
      handler: 'handler',
      environment: {
        'REPOSITORY_BUCKET': repositoryBucket.bucketName,
        'REPOSITORY_TABLE': repositoryTable.tableName,
      },
    });
    repositoryBucket.grantReadWrite(adminLambda);
    repositoryTable.grantReadWriteData(adminLambda);

    // REST API (v2)
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

    new cdk.CfnOutput(this, 'EndpointUrl', { value: api.url!, description: 'API URL' })

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
  }
}