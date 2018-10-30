# Influxdata-Operator

The Influxdata Operator creates, configures and manages Influxdb OSS running on Kubernetes.

InfluxDB
An Open-Source Time Series Database
InfluxDB is an open source time series database built by the folks over at InfluxData with no external dependencies. It's useful for recording metrics, events, and performing analytics.

QuickStart

kubectl create -f deploy/crds/influxdata_v1alpha1_influxdb_cr.yaml

Prerequisites:

dep,
git,
go,
docker,
kubectl,
And Access to a kubernetes cluster.



Installing the Operator-SDK

$ mkdir -p $GOPATH/src/github.com/operator-framework

$ cd $GOPATH/src/github.com/operator-framework

$ git clone https://github.com/operator-framework/operator-sdk

$ cd operator-sdk

$ make dep

$ make install

# Create and deploy an influxdata-operator using the SDK CLI:
$ operator-sdk new influxdata-operator --api-version=dev9-labs.bitbucket.org/v1alpha1 --kind=Influxdb

# Add a new controller that watches for Influxdb
$ operator-sdk add controller  --api-version=dev9-labs.bitbucket.org/v1alpha1 --kind=Influxdb 




# Deploy the Influxdata Operator && Custom Resource for Influxdata Installation
$ kubectl create -f deploy/storageclass-gcp.yaml

$ kubectl create -f deploy/crds/influxdata_v1alpha1_influxdb_cr.yaml

Note: storageclass-gcp.yaml this yaml file will create PD in GCP , storageclass-aws.yaml this yaml file will create EBS in AWS
and storageclass-nfs.yaml will create PVC based on nfs .

# Cleanup
$ kubectl create -f deploy/crds/influxdata_v1alpha1_influxdb_cr.yaml

$ kubectl create -f deploy/storageclass-gcp.yaml


# Persistence
The InfluxDB image stores data in the /var/lib/influxdb directory in the container.
The Operator mounts a Persistent Volume volume at this location. The volume is created using dynamic volume provisioning.
