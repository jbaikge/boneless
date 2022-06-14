import { Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigw from 'aws-cdk-lib/aws-apigateway';
import { join } from 'path';
import { LambdaIntegration } from 'aws-cdk-lib/aws-apigateway';

export class GocmsStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const assetDir = join(__dirname, '..', '..', 'assets');

    const pingPongFunction = new lambda.Function(this, 'PingPongHandler', {
      runtime: lambda.Runtime.GO_1_X,
      code:    lambda.Code.fromAsset(join(assetDir, 'ping-pong')),
      handler: 'ping-pong'
    });

    const api = new apigw.RestApi(this, 'Endpoint', {});

    const pingPongResource = api.root.addResource('ping-pong');
    pingPongResource.addMethod('GET', new LambdaIntegration(pingPongFunction));
  }
}
