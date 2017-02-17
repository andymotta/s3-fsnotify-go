# s3-fsnotify-go
Watches a file for changes and to S3 on change.  Tested on Darwin and Linux.

You will need to set the following environment variables for this to work:
```bash
export filename=/path/to/filetowatch.json
export bucket=bucket-to-upload-to
```

For AWS auth you’ll need a default config (~/.aws/config) with at least:
```
[default]
region = us-west-2
```
and  ~/.aws/credentials with:
```
[default]
aws_access_key_id = <your_aws_access_key_id>
aws_secret_access_key = <you_aws_secret_access_key>
```

Make sure this ^ user can PutObject on the bucket you are specifying.

Test Run:
```bash
$ export filename=file
$ export bucket=andy-cloudfront-dev
$ s3-fsnotify-go &
[1] 11446

$ echo "Hello World" > file
$ Successfully uploaded%!(EXTRA string=file, string=andy-cloudfront-dev)2017/02/16 23:03:59 event: "file": WRITE
2017/02/16 23:03:59 modified file: file
Successfully uploaded%!(EXTRA string=file, string=andy-cloudfront-dev)2017/02/16 23:04:00 event: "file": CHMOD
$ aws s3 ls s3://andy-cloudfront-dev/
2017-02-16 23:04:00         14 file

$ echo "" > file
$ Successfully uploaded%!(EXTRA string=file, string=andy-cloudfront-dev)2017/02/16 23:04:13 event: "file": CHMOD
$ aws s3 ls s3://andy-cloudfront-dev/
2017-02-16 23:04:13          0 file
```
