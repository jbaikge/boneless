import * as cdk from 'aws-cdk-lib';
import * as constructs from 'constructs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as s3Deployment from 'aws-cdk-lib/aws-s3-deployment';
import * as path from 'path';
import * as fs from 'fs-extra';
import * as cloudfront from 'aws-cdk-lib/aws-cloudfront';
import * as cloudfrontOrigins from 'aws-cdk-lib/aws-cloudfront-origins';
import exec from './exec';


export class AdminStack extends cdk.Stack {
  constructor(scope: constructs.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    // Path back to repo root
    const rootDir = path.resolve(__dirname, '..', '..');

    // Admin frontend
    // Ref: https://aws-cdk.com/deploying-a-static-website-using-s3-and-cloudfront/
    // Ref: https://github.com/aws-samples/cdk-build-bundle-deploy-example/blob/main/cdk-bundle-static-site-example/lib/static-site-stack.ts
    const adminBucket = new s3.Bucket(this, 'FrontendAdminBucket', {
      bucketName: cdk.PhysicalName.GENERATE_IF_NEEDED,
      encryption: s3.BucketEncryption.S3_MANAGED,
      accessControl: s3.BucketAccessControl.PRIVATE,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
    });

    const adminOriginAccessIdentity = new cloudfront.OriginAccessIdentity(this, 'FrontendAdminOAI');
    adminBucket.grantRead(adminOriginAccessIdentity);

    const distribution = new cloudfront.Distribution(this, 'FrontendAdminDistribution', {
      defaultRootObject: 'index.html',
      defaultBehavior: {
        origin: new cloudfrontOrigins.S3Origin(adminBucket, {originAccessIdentity: adminOriginAccessIdentity}),
      }
    });

    const adminFrontendDir = path.join(rootDir, '_frontend-admin');
    const adminDeployment = new s3Deployment.BucketDeployment(this, 'FrontendAdminDeployment', {
      destinationBucket: adminBucket,
      distribution: distribution,
      memoryLimit: 128,
      sources: [
        s3Deployment.Source.asset(adminFrontendDir, {
          bundling: {
            local: {
              tryBundle(outputDir: string) {
                try {
                  exec('npm run build', {
                    cwd: adminFrontendDir,
                    env: process.env,
                    stdio: [
                      'ignore',
                      process.stderr,
                      'inherit',
                    ],
                  });
                  fs.copySync(path.join(adminFrontendDir, 'build'), outputDir);
                } catch (error) {
                  console.error(error)
                  return false;
                }
                return true;
              }
            },
            image: cdk.DockerImage.fromRegistry('node:lts'),
            command: [
              'bash', '-c', [
                'cd /asset-input',
                'npm install',
                'npm run build',
                'cp -r /asset/input/build/* /asset-output/',
              ].join(' && '),
            ],
            environment: {},
          }
        }),
      ],
    });

    new cdk.CfnOutput(this, 'CloudFrontUrl', {
      value: distribution.distributionDomainName,
      description: 'CloudFront distribution domain',
    });
  }
}