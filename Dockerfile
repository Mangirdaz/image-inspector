FROM openshift/origin-base
MAINTAINER Federico Simoncelli <fsimonce@redhat.com>

RUN yum install -y golang openscap-scanner openscap-containers openscap-utils && yum clean all

ADD ./cmd/cmd /usr/local/bin/cmd

EXPOSE 8080

WORKDIR /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/cmd"]
