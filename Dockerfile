FROM ghcr.io/jenkins-x/jx-boot:latest

ENTRYPOINT ["jx-test-collector"]

COPY ./build/linux/jx-test-collector /usr/bin/jx-test-collector