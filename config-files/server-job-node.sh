#cloud-config
runcmd:
  - apt-get install -y wget
  - wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
  - tar -C /usr/local -xzf go1.4.2.linux-amd64.tar.gz
  - echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.bashrc
  - mkdir /root/go
  - export HOME=/root
  - echo 'export GOROOT=$HOME/go' >> /root/.bashrc
  - echo 'export PATH=$PATH:$GOROOT/bin' >> /root/.bashrc
  - export GOPATH=/root/go
  - /usr/local/go/bin/go get github.com/MohamedBassem/getaredis/...
  - apt-get install -y supervisor
write_files:
  - path: /etc/supervisor/conf.d/go_jobs.conf
    content: |
        [program:go_jobs]
        command=$GOROOT/bin/getaredis-jobs -f /root/config.yml
        autostart=true
        autorestart=true
        stderr_logfile=/var/log/go_jobs.err.log
        stdout_logfile=/var/log/go_jobs.out.log

        [program:go_server]
        command=$GOROOT/bin/getaredis-server -f /root/config.yml
        autostart=true
        autorestart=true
        stderr_logfile=/var/log/getaredis-server.err.log
        stdout_logfile=/var/log/getaredis-server.out.log
  - path: /root/config.yml
    content: |
        # Please replace me!
