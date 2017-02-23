# s3-fsnotify-go
Watches a directory for changes and uploads modified or created files to S3.  Tested on Darwin and Linux.

The best way to run this is with a failure-tolerant init script.  go2s3.service is a systemd exmaple.  Here are some installation notes:
```bash
vim /usr/lib/systemd/system/go2s3.service
sudo systemctl enable go2s3.service
service go2s3 status
service go2s3 start
```

Alternatively, you can set environment variables:
```bash
export syncdir=.
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
$ systemctl start go2s3
$ touch outputs.json
Feb 22 20:39:33  s3-fsnotify-go[1503]: 2017/02/22 20:39:33 modified file: ./outputs.json
Feb 22 20:39:35  s3-fsnotify-go[1503]: Successfully uploaded ./outputs.json to bucket-to-upload-to
```
