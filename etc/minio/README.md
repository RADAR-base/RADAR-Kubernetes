# minio

Minio object storage used for storing RADAR-Base intermediate and output data, based on the bitnami/minio chart. Some values have changed to adapt needs of RADAR-Base. For a complete list of values visit [the original chart](https://github.com/RADAR-base/radar-helm-charts/tree/main/external/minio). The Minio console is accessed via `https://s3.<radar-base-hostname>`.

The default root credentials should only be used for creating actual users and attaching policies to them. This can be done by adapting the `provisioning` values, specifically the `provisioning.users` and `provisioning.policies` lists.

By default, bucket `radar-output-storage` is provisioned to store output data and bucket `radar-intermediate-storage` is used to store raw output data from Kafka. In general, users should only get read-only access to the `radar-output-storage` bucket, by assigning the `read-output` policy to them.

Please be aware that this chart is optional: if the infrastructure has an (S3) object storage or Azure Blob Storage available, then it does not need to be added here.

## Minio client

To create users and define access policies use [the guide on Minio website](https://docs.min.io/docs/minio-client-quickstart-guide.html) to install and configure Minio client. The URL that minio client should use is `https://api.s3.<radar-base-hostname>`.