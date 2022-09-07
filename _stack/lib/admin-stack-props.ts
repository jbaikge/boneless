import * as cdk from 'aws-cdk-lib';
import * as apigateway from '@aws-cdk/aws-apigatewayv2-alpha';

// Ref: https://bobbyhadz.com/blog/aws-cdk-share-resources-between-stacks
export interface AdminStackProps extends cdk.StackProps {
  api: apigateway.HttpApi;
}
