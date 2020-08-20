### Linux

```shell
curl -L https://github.com/jenkins-x/jx-test-collector/releases/download/v{{.Version}}/jx-test-collector-linux-amd64.tar.gz | tar xzv 
sudo mv jx-test-collector /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x/jx-test-collector/releases/download/v{{.Version}}/jx-test-collector-darwin-amd64.tar.gz | tar xzv
sudo mv jx-test-collector /usr/local/bin
```

