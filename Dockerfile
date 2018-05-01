FROM scratch
MAINTAINER Brian Hechinger <wonko@4amlunch.net>

ADD auth-cmd-linux-amd64 auth-cmd
VOLUME /etc/chremoas

ENTRYPOINT ["/auth-cmd", "--configuration_file", "/etc/chremoas/chremoas.yaml"]
