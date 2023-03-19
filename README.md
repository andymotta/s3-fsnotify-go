# s3-fsnotify-go

Watches a directory for changes and uploads modified or created files to S3.  Tested on Darwin, Linux and Windows 386.

## Installation

1. Download the latest release for your platform from the [Releases page](https://github.com/andymotta/s3-fsnotify-go/releases).

2. Move the `s3-fsnotify-go` binary to a suitable location (e.g., `/usr/local/bin` or a location in your `$PATH`).

## Usage

1. Configure the `go2s3.service` file with the appropriate path to the `s3-fsnotify-go` binary.

2. Copy the `go2s3.service` file to `/etc/systemd/system/`.

3. Enable and start the service:

```bash
sudo systemctl enable go2s3.service
sudo systemctl start go2s3.service
```

4. To check the service status, use:
`sudo systemctl status go2s3.service`

5. To view the logs generated by the service, use:
`sudo journalctl -u go2s3.service`