import {
  aws_autoscaling as autoscaling,
  aws_ec2 as ec2,
  aws_iam as iam,
  CfnParameter,
  aws_ecs as ecs,
} from 'aws-cdk-lib';
import { Construct } from 'constructs';

interface ClusterProps {
  readonly vpc: ec2.IVpc;
  readonly clusterName: string;
  readonly instanceType: CfnParameter;
  readonly securityGroup: ec2.SecurityGroup;
}

export class Cluster extends Construct {
  readonly ecsCluster: ecs.Cluster;
  readonly autoScalingGroup: autoscaling.AutoScalingGroup;
  readonly capacityProvider: ecs.AsgCapacityProvider;

  constructor(scope: Construct, id: string, props: ClusterProps) {
    super(scope, id);

    this.ecsCluster = new ecs.Cluster(this, 'Cluster', {
      vpc: props.vpc,
      clusterName: props.clusterName,
    });

    this.autoScalingGroup = new autoscaling.AutoScalingGroup(this, 'AutoscalingGroup', {
      vpc: props.vpc,
      minCapacity: 1,
      maxCapacity: 1,
      desiredCapacity: 1,
      machineImage: ecs.EcsOptimizedImage.amazonLinux2(),
      instanceType: new ec2.InstanceType(props.instanceType.valueAsString),
      securityGroup: props.securityGroup,
      role: new iam.Role(this, 'EcsInstanceRole', {
        assumedBy: new iam.ServicePrincipal('ec2.amazonaws.com'),
        managedPolicies: [
          iam.ManagedPolicy.fromAwsManagedPolicyName('service-role/AmazonEC2ContainerServiceforEC2Role'),
        ],
      }),
    });

    this.capacityProvider = new ecs.AsgCapacityProvider(this, 'AsgCapacityProvider', {
      autoScalingGroup: this.autoScalingGroup,
      enableManagedScaling: false,
      enableManagedTerminationProtection: false,
    });

    this.ecsCluster.addAsgCapacityProvider(this.capacityProvider);
  }
}