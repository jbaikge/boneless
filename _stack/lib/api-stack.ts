import * as cdk from 'aws-cdk-lib';
import * as constructs from 'constructs';
import * as lsp from './lambda-stack-props';
import * as apigateway from '@aws-cdk/aws-apigatewayv2-alpha';
import * as integration from '@aws-cdk/aws-apigatewayv2-integrations-alpha';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as path from 'path';

export class ApiStack extends cdk.Stack {
  public readonly api: apigateway.HttpApi;

  constructor(scope: constructs.Construct, id: string, props: lsp.LambdaStackProps) {
    super(scope, id, props);

    const handlerDir = path.resolve(__dirname, '..', '..', 'lambda', 'api');
    const apiLambda = new lambda.Function(this, 'ApiHandler', {
      code: lambda.Code.fromAsset(handlerDir),
      runtime: lambda.Runtime.GO_1_X,
      handler: 'handler',
      environment: {
        'REPOSITORY_BUCKET': props.bucket.bucketName,
        'REPOSITORY_TABLE': props.db.tableName,
        'STATIC_BUCKET': props.static.bucketName,
      },
    });
    props.bucket.grantReadWrite(apiLambda);
    props.db.grantReadWriteData(apiLambda);
    props.static.grantReadWrite(apiLambda);

    this.api = new apigateway.HttpApi(this, 'API', {
      createDefaultStage: true,
      corsPreflight: {
        allowOrigins: [
          '*',
        ],
        allowHeaders: [
          'Authorization',
          'Content-Type',
          'Range',
        ],
        exposeHeaders: [
          'Content-Range',
          'X-Total-Count',
        ],
        allowMethods: [
          apigateway.CorsHttpMethod.DELETE,
          apigateway.CorsHttpMethod.GET,
          apigateway.CorsHttpMethod.OPTIONS,
          apigateway.CorsHttpMethod.POST,
          apigateway.CorsHttpMethod.PUT,
        ],
      },
    });

    const apiIntegration = new integration.HttpLambdaIntegration('ApiIntegration', apiLambda);
    const routes: apigateway.AddRoutesOptions[] = [
      {
        path: '/classes',
        integration: apiIntegration,
        methods: [
          apigateway.HttpMethod.GET,
          apigateway.HttpMethod.POST,
        ],
      },
      {
        path: '/classes/{class_id}',
        integration: apiIntegration,
        methods: [
          apigateway.HttpMethod.GET,
          apigateway.HttpMethod.PUT,
          apigateway.HttpMethod.DELETE,
        ],
      },
      {
        path: '/classes/{class_id}/documents',
        integration: apiIntegration,
        methods: [
          apigateway.HttpMethod.GET,
          apigateway.HttpMethod.POST,
        ],
      },
      {
        path: '/classes/{class_id}/documents/{doc_id}',
        integration: apiIntegration,
        methods: [
          apigateway.HttpMethod.GET,
          apigateway.HttpMethod.PUT,
          apigateway.HttpMethod.DELETE,
        ],
      },
      {
        path: '/documents',
        integration: apiIntegration,
        methods: [
          apigateway.HttpMethod.GET,
          apigateway.HttpMethod.POST,
        ],
      },
      {
        path: '/documents/{doc_id}',
        integration: apiIntegration,
        methods: [
          apigateway.HttpMethod.GET,
          apigateway.HttpMethod.PUT,
          apigateway.HttpMethod.DELETE,
        ],
      },
      {
        path: '/templates',
        integration: apiIntegration,
        methods: [
          apigateway.HttpMethod.GET,
          apigateway.HttpMethod.POST,
        ],
      },
      {
        path: '/templates/{template_id}',
        integration: apiIntegration,
        methods: [
          apigateway.HttpMethod.GET,
          apigateway.HttpMethod.PUT,
          apigateway.HttpMethod.DELETE,
        ],
      },
    ];
    routes.forEach((route) => this.api.addRoutes(route));

    new cdk.CfnOutput(this, 'ApiUrl', {
      value: this.api.url!,
      description: 'API URL',
    });
  }
}
