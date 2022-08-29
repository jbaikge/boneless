import * as cdk from 'aws-cdk-lib';
import * as constructs from 'constructs';
import * as lsp from './lambda-stack-props';
import * as apigateway from '@aws-cdk/aws-apigatewayv2-alpha';
import * as integration from '@aws-cdk/aws-apigatewayv2-integrations-alpha';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as path from 'path';

export class FrontendStack extends cdk.Stack {
  constructor(scope: constructs.Construct, id: string, props: lsp.LambdaStackProps) {
    super(scope, id, props);

    const handlerDir = path.resolve(__dirname, '..', '..', 'lambda', 'frontend');
    const frontendLambda = new lambda.Function(this, 'FrontendHandler', {
      code: lambda.Code.fromAsset(handlerDir),
      runtime: lambda.Runtime.GO_1_X,
      handler: 'handler',
      environment: {
        'REPOSITORY_BUCKET': props.bucket.bucketName,
        'REPOSITORY_TABLE': props.db.tableName,
      },
    });
    props.bucket.grantRead(frontendLambda);
    props.db.grantReadData(frontendLambda);

    const frontendIntegration = new integration.HttpLambdaIntegration('FrontendIntegration', frontendLambda);

    // https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-develop-routes.html
    // https://docs.aws.amazon.com/cdk/api/v1/docs/aws-apigateway-readme.html#default-integration-and-method-options
    const api = new apigateway.HttpApi(this, 'Frontend', {
      defaultIntegration: frontendIntegration,
      createDefaultStage: true,
      corsPreflight: {
        allowOrigins: [
          '*',
        ],
        allowHeaders: [
        ],
        exposeHeaders: [
        ],
        allowMethods: [
          apigateway.CorsHttpMethod.GET,
        ],
      },
    });

    new cdk.CfnOutput(this, 'FrontendUrl', {
      value: api.url!,
      description: 'Frontend URL',
    });
  }
}
