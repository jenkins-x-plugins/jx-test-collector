FROM gcr.io/jenkinsxio-labs-private/jx-cli-base:0.0.3

ENTRYPOINT ["jx-test-collector"]

COPY ./build/linux/jx-test-collector /usr/bin/jx-test-collector