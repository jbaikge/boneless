import * as cdk from 'aws-cdk-lib';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as cloudfront from 'aws-cdk-lib/aws-cloudfront';

// Ref: https://bobbyhadz.com/blog/aws-cdk-share-resources-between-stacks
export interface LambdaStackProps extends cdk.StackProps {
  dbTable: dynamodb.Table;
  dbBucket: s3.Bucket;
  staticBucket: s3.Bucket;
  staticDistribution: cloudfront.Distribution;
}
