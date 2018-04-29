# s3-fsnotify-go

[![Build Status](https://travis-ci.org/andymotta/s3-fsnotify-go.svg?branch=master)](https://travis-ci.org/andymotta/s3-fsnotify-go)

Watches a directory for changes and uploads modified or created files to S3.  Tested on Darwin, Linux and Windows 386.

The best way to run this is with a failure-tolerant init script.  go2s3.service is a systemd example.  Here are some installation notes:
```bash
vim /usr/lib/systemd/system/go2s3.service
sudo systemctl enable go2s3.service
service go2s3 status
service go2s3 start
```

For testing, you can set environment variables:
```bash
export syncdir=/directory/to/sync/ # (optional, CWD otherwise)
export bucket=bucket-to-upload-to
go run s3_fsnotify.go &
```

For AWS auth you’ll need a default config (~/.aws/config) with at least:
```bash
[default]
region = us-west-2 # can be any region
```
and  ~/.aws/credentials with:
```
[default]
aws_access_key_id = <your_aws_access_key_id>
aws_secret_access_key = <you_aws_secret_access_key>
```

Make sure this ^ user can PutObject on the bucket you are specifying.

To use AWS profiles outside of default, simply set:
```bash
export AWS_PROFILE=altprofile
```

Test Run:
```bash
$ systemctl start go2s3 # or go run s3_fsnotify.go &
$ touch outputs.json
Feb 22 20:39:33  s3-fsnotify-go[1503]: 2017/02/22 20:39:33 modified file: ./outputs.json
Feb 22 20:39:35  s3-fsnotify-go[1503]: Successfully uploaded ./outputs.json to bucket-to-upload-to

$ vim README.md # write, quit
$ Successfully uploaded README.md to bucket-to-upload-to
2017/04/25 13:10:57 event: "README.md": CHMOD
2017/04/25 13:10:59 event: "README.md": CHMOD
```

### Roadmap
- STS Assume role in addition to shared credentials
