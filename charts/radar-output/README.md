

# radar-output

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.2.1](https://img.shields.io/badge/AppVersion-1.2.1-informational?style=flat-square)

A Helm chart for RADAR-base radar-output service

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl |  |
| Joris Borgdorff | joris@thehyve.nl |  |
| Nivethika Mahasivam | nivethika@thehyve.nl |  |

## Source Code

* <https://github.com/RADAR-base/radar-output-restructure>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+
* PV provisioner support in the underlying infrastructure

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `1` |  |
| image.repository | string | `"radarbase/radar-output-restructure"` |  |
| image.tag | string | `"1.2.1"` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| imagePullSecrets | list | `[]` |  |
| nameOverride | string | `""` |  |
| fullnameOverride | string | `""` |  |
| resources.limits.cpu | string | `"1000m"` |  |
| resources.requests.cpu | string | `"100m"` |  |
| resources.requests.memory | string | `"400Mi"` |  |
| nodeSelector | object | `{}` |  |
| tolerations | list | `[]` |  |
| affinity | object | `{}` |  |
| persistence.enabled | bool | `false` |  |
| source.type | string | `"s3"` | Type of the intermediate storage of the RADAR-base pipeline e.g. s3, hdfs |
| source.s3.endpoint | string | `"http://minio:9000"` | s3 endpoint of the intermediate storage |
| source.s3.accessToken | string | `"access_key"` | s3 access-key of the intermediate storage |
| source.s3.secretKey | string | `"secret"` | s3 secret-key of the intermediate storage |
| source.s3.bucket | string | `"radar-intermediate-storage"` | s3 bucket name of the intermediate storage |
| source.azure.endpoint | string | `""` | Azure endpoint of the intermediate storage |
| source.azure.username | string | `""` | Azure username to access the s3 endpoint when using personal login |
| source.azure.password | string | `""` | Azure password when using personal login |
| source.azure.accountName | string | `""` | Azure account name when using shared access tokens |
| source.azure.accountKey | string | `""` | Azure account key when using shared access tokens |
| source.azure.sasToken | string | `""` | Azure SAS(shared access signature) token when using shared access tokens |
| source.azure.container | string | `""` | Azure blob container name |
| target.type | string | `"s3"` | Type of the output storage of the RADAR-base pipeline e.g. s3, local |
| target.s3.endpoint | string | `"http://minio:9000"` | s3 endpoint of the output storage |
| target.s3.accessToken | string | `"access_key"` | s3 access-key of the output storage |
| target.s3.secretKey | string | `"secret"` | s3 secret-key of the output storage |
| target.s3.bucket | string | `"radar-output-storage"` | s3 bucket name of the output storage |
| target.azure.endpoint | string | `""` | Azure endpoint of the output storage |
| target.azure.username | string | `""` | Azure username to access the s3 endpoint when using personal login |
| target.azure.password | string | `""` | Azure password when using personal login |
| target.azure.accountName | string | `""` | Azure account name when using shared access tokens |
| target.azure.accountKey | string | `""` | Azure account key when using shared access tokens |
| target.azure.sasToken | string | `""` | Azure SAS(shared access signature) token when using shared access tokens |
| target.azure.container | string | `""` | Azure blob container name |
| redis.uri | string | `"redis://redis-master:6379"` | URL of the redis database |
| worker.minimumFileAge | int | `900` | Minimum amount of time in seconds since a file was last modified for it to be considered for processing. |
| worker.numThreads | int | `2` | Number of threads to do processing on |
| cleaner.age | int | `7` | Number of days after which a source file is considered old |
| paths.input | string | `"topics"` | Relative path to intermediate storage root to browse for data |
| paths.output | string | `"output"` | Relative path to output storage to write data |
| paths.factory | string | `"org.radarbase.output.path.ObservationKeyPathFactory"` | # Output path construction factory |
| paths.properties | object | `{}` | Additional properties. For details see https://github.com/RADAR-base/radar-output-restructure/blob/master/restructure.yml |
