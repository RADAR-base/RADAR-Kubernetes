from pyspark.sql import SparkSession

# Initialize Spark with Iceberg and S3 configurations
spark = SparkSession.builder \
    .appName("IcebergUpsertProcessor") \
    # 1. Configure the Iceberg Catalog (SeaweedFS)
    .config("spark.sql.catalog.seaweed", "org.apache.iceberg.spark.SparkCatalog") \
    .config("spark.sql.catalog.seaweed.type", "hadoop") \
    .config("spark.sql.catalog.seaweed.warehouse", "s3a://radar-output-storage/warehouse") \
    # 2. S3 / SeaweedFS Connection
    .config("spark.hadoop.fs.s3a.endpoint", "http://radar-seaweedfs-s3:8333") \
    .config("spark.hadoop.fs.s3a.access.key", "your-access-key") \
    .config("spark.hadoop.fs.s3a.secret.key", "your-secret-key") \
    .config("spark.hadoop.fs.s3a.path.style.access", "true") \
    .getOrCreate()

# Tables to process
tables = ["topic_table_1", "topic_table_2"]

for table_name in tables:
    # Read from intermediate storage
    # Assuming input is also registered in the catalog or read via path
    new_data_df = spark.read.format("iceberg").load(f"s3a://radar-intermediate-storage/{table_name}")

    # Process: Extract keys from the 'key' column and deduplicate the batch
    processed_df = new_data_df.withColumn("projectId", new_data_df["key.projectId"]) \
                              .withColumn("subjectId", new_data_df["key.subjectId"]) \
                              .dropDuplicates(["projectId", "subjectId", "time"])

    # Register the batch as a temporary view to use Spark SQL
    processed_df.createOrReplaceTempView("batch_updates")

    # 3. Iceberg MERGE INTO Logic
    # This handles the "datum from the past" by updating if exists, inserting if not.
    spark.sql(f"""
        MERGE INTO seaweed.db.{table_name} AS target
        USING batch_updates AS source
        ON target.projectId = source.projectId 
           AND target.subjectId = source.subjectId 
           AND target.time = source.time
        WHEN MATCHED THEN
            UPDATE SET *
        WHEN NOT MATCHED THEN
            INSERT *
    """)

spark.stop()