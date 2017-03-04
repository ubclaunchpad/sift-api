FROM golang:1.7.4-wheezy
MAINTAINER Eric Stroczynski <ericstroczynski@gmail.com>

ENV SIFT_API_PATH /go/src/github.com/ubclaunchpad/sift-api
EXPOSE 9090

COPY . $SIFT_API_PATH
RUN apt-get update \
    && apt-get clean \
    && cd $SIFT_API_PATH \
    && go get \
    && go install
    
CMD ["/go/bin/sift-api","--dbhost=postgres"]
