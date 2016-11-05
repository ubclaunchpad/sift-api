FROM golang:1.6.3-wheezy
MAINTAINER Jordan Schalm <jordan.schalm@gmail.com>

ENV SIFT_API_PATH /go/src/github.com/ubclaunchpad/sift-api
EXPOSE 9090

COPY . $SIFT_API_PATH
RUN cd $SIFT_API_PATH && go get && go install

CMD /go/bin/sift-api
