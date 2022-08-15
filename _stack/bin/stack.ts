#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { AdminStack } from '../lib/admin-stack';
import { RepositoryStack } from '../lib/repository-stack';

const app = new cdk.App();
new RepositoryStack(app, 'RepositoryStack', {});
new AdminStack(app, 'AdminStack', {});