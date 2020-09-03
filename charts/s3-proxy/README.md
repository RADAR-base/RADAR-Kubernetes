# S3Proxy

Uses https://hub.docker.com/r/andrewgaul/s3proxy to proxy S3 API requests to any supported cloud provider.

Find some example configurations at https://github.com/gaul/s3proxy/wiki/Storage-backend-examples.


For example, set

```yaml
# Credentials used to access this proxy
s3:
  identity: MyUser
  credential: MySecret

# Where requests should be proxied to
target:
  provider: azureblob
  endpoint: http://MyCloud.com
  identity: MyCloudUser
  credentials: MyCloudSecret
```