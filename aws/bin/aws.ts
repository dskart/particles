#!/usr/bin/env node
import * as cdk from "aws-cdk-lib";
import { ParticlesStack } from "../lib/particles-stack";
import { EcrStack } from "../lib/ecr-stack";

const app = new cdk.App();

const env = {
  account: process.env.CDK_DEFAULT_ACCOUNT,
  region: process.env.CDK_DEFAULT_REGION,
};

// ECR Stack
new EcrStack(app, "ParticlesEcrStack", {
  stackName: app.node.tryGetContext("stack-name"),
  repositoryName: "particles",
  env,
});

// Main application stack
new ParticlesStack(app, "ParticlesStack", {
  stackName: app.node.tryGetContext("stack-name"),
  env,
});
