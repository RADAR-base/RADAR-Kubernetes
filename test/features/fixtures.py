from behave import fixture
from minio import Minio

@fixture
def s3(context):
    minio_client = Minio(
        context.config.userdata["s3_connection_url"],
        context.config.userdata["s3_key"],
        context.config.userdata["s3_secret"],
    )
    # Remove all objects in intermediate and outout storage buckets.
    for bucket in ["radar-intermediate-storage", "radar-output-storage"]:
        objects = minio_client.list_objects(
            bucket_name=bucket,
            recursive=True
        )
        minio_client.remove_objects(
            bucket_name=bucket,
            delete_object_list=objects
        )
