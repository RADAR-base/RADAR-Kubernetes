# Limit Access To Output Storage
In order to limit access to only a certain project in any S3 compatible object storage you can use this policy. Be sure to change `<project_name>` with name of your project in MP.
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetBucketLocation",
                "s3:ListAllMyBuckets"
            ],
            "Resource": [
                "arn:aws:s3:::radar-output-storage"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::radar-output-storage"
            ],
            "Condition": {
                "StringLike": {
                    "s3:prefix": [
                        "output/<project_name>/*"
                    ]
                }
            }
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetBucketLocation",
                "s3:GetObject"
            ],
            "Resource": [
                "arn:aws:s3:::radar-output-storage/output/<project_name>/*"
            ]
        }
    ]
}
```
