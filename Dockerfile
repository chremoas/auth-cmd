FROM --platform=$BUILDPLATFORM scratch AS build
MAINTAINER Brian Hechinger <wonko@4amlunch.net>
ARG BUILDOS
ARG BUILDARCH

FROM build
ADD auth-cmd-${BUILDOS}-${BUILDARCH} auth-cmd
VOLUME /etc/chremoas

ENTRYPOINT ["/auth-cmd", "--configuration_file", "/etc/chremoas/chremoas.yaml"]
