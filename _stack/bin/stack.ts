#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { DatabaseStack } from '../lib/database-stack';
import { ApiStack } from '../lib/api-stack';
import { FrontendStack } from '../lib/frontend-stack';
import { AdminStack } from '../lib/admin-stack';
import { StaticStack } from '../lib/static-stack';

const app = new cdk.App();

const databaseStack = new DatabaseStack(app, 'DatabaseStack', {});

const staticStack = new StaticStack(app, 'StaticStack', {});

const apiStack = new ApiStack(app, 'ApiStack', {
  dbTable:     databaseStack.table,
  dbBucket: databaseStack.bucket,
  staticBucket: staticStack.bucket,
  staticDistribution: staticStack.distribution,
});

new AdminStack(app, 'AdminStack', {
  api: apiStack.api,
});

new FrontendStack(app, 'FrontendStack', {
  dbTable:     databaseStack.table,
  dbBucket: databaseStack.bucket,
  staticBucket: staticStack.bucket,
  staticDistribution: staticStack.distribution,
});
