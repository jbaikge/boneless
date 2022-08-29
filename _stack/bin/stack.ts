#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { DatabaseStack } from '../lib/database-stack';
import { ApiStack } from '../lib/api-stack';
import { FrontendStack } from '../lib/frontend-stack';
import { AdminStack } from '../lib/admin-stack';
// import { AdminStack } from '../lib/admin-stack';
// import { RepositoryStack } from '../lib/repository-stack';

const app = new cdk.App();

const databaseStack = new DatabaseStack(app, 'DatabaseStack', {});

new ApiStack(app, 'ApiStack', {
    db:     databaseStack.db,
    bucket: databaseStack.bucket,
});

new FrontendStack(app, 'FrontendStack', {
    db:     databaseStack.db,
    bucket: databaseStack.bucket,
});

new AdminStack(app, 'AdminStack', {});