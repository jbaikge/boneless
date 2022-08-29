import * as cdk from 'aws-cdk-lib';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as s3 from 'aws-cdk-lib/aws-s3';

// Ref: https://bobbyhadz.com/blog/aws-cdk-share-resources-between-stacks
export interface LambdaStackProps extends cdk.StackProps {
    bucket: s3.Bucket;
    db: dynamodb.Table;
}
