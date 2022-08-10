## AWS Storage Service

-Amazon Elastic Block store : useful for storing multiple files as a block
can be used for storing unprocessed or raw data

-AMazon S3 Glacier: is a service and is focused on soting archieved files

-AMazon Elastic File System: used for storing of files

-AWS Storage Gateway : are mostly used when transferring files from an on premise server to aws. its main purpose is to make sure the data is transferred smoothly without breaks

-Amazon S3 Storage: object storage, where we can store multiple files as part of the s3 sotorage

## Amazon s3 (simple storage service)

Amazon S3 has a simple web service interface that we can
use to store and retrieve any amount and kind of data at any time anywhere on the web
s3 is:
-available
-cost friendly
-flexible
-durable
-secure
-scalable

# Storage Option for s3

- Standard: we store all the file we require on a frequent basis. there's no min requirement or file size

- Infrequent Access: we store files which we do not frequently but if needed can be fetched quickly. min store duration for 30days

-Glacier: historic data can be stored here.we store files that arent frequently read. like archieved files. files with larger volume can be here, files to be stored here take like 4hours. min store duration of 90days

## Buckets

a single object cant be > 5Tb

a container we use to store multiple objects with their meta data i.e keys, version, description ,etc
100 buckets per account
20 buckets per region

## Versioning

## ACL

public-read
