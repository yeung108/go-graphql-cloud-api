# How to run
- `GOOS=linux GOARCH=amd64 go build main.go` in mac
- git clone the repository onto the server
- create .env file and public.pem 
- type `./main` to run the program

# Run as service
- `sudo nano /lib/systemd/system/mclient-v2.service`
- 
```
[Unit]
Description=mclient v2 service
ConditionPathExists=/home/ubuntu/go-graphql-cloud-api
After=network.target

[Service]
Type=simple
User=ubuntu
Group=ubuntu
LimitNOFILE=1024

Restart=on-failure
RestartSec=10
StartLimitIntervalSec=60

WorkingDirectory=/home/ubuntu/go-graphql-cloud-api
ExecStart=/home/ubuntu/go-graphql-cloud-api/main

# make sure log directory exists and owned by syslog
PermissionsStartOnly=true
ExecStartPre=/bin/mkdir -p /var/log/mclient-v2-service
ExecStartPre=/bin/chown ubuntu:ubuntu /var/log/mclient-v2-service
ExecStartPre=/bin/chmod 755 /var/log/mclient-v2-service
StandardOutput=/var/log/mclient-v2-service/mclient-v2.log
StandardError=/var/log/mclient-v2-service/error.log
SyslogIdentifier=mclient-v2-service

[Install]
WantedBy=multi-user.target
```
- `sudo systemctl enable mclient-v2.service`
- `sudo systemctl start mclient-v2.service`

