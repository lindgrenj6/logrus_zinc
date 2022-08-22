# Logrus Zinc

This is a hook for use with the amazing [ZincSearch](https://zincsearch.com/) just in case anyone wanted to not have to run an entire elastic + kibana cluster locally.

### Adding to your Logrus logger
1.  `go get github.com/lindgrenj6/logrus_zinc`
2. All of the logic is contained in the `LocalZincHook` struct - so adding a call like this should be enough:
```go
logger.AddHook(&LocalZincHook{
    URL: "http://localhost:4080,
    Index: "myapp",
    Username: "admin",
    Password: "changeme",
})
```

_note: URL will default to localhost:4080 and index will default to...default. so they are not required_

3. After that just run the program like normal and your logs should be geting shipped to your local zinc instance. Sweet!

Just browse to http://localhost:4080 and you will see a lightweight verion of the Kibana UI we're all used to.


### Running Zinc
I usually run zinc with podman, but docker will work fine too. Run an ephemeral instance of zinc using this command:
```shell
mkdir -p /tmp/zincdata && chmod 777 $_

/usr/bin/podman run \
        -it --rm \
        -v /tmp/zincdata:/data \
        -e ZINC_DATA_PATH=/data \
        -e ZINC_FIRST_ADMIN_USER=admin \
        -e ZINC_FIRST_ADMIN_PASSWORD=changeme \
        -p 4080:4080 \
        --name zinc \
        -it public.ecr.aws/zinclabs/zinc:latest
```
and there will be a new zinc instance running on your localhost using the credentials provided.

if there are any issues, here are the amazing zinc docs for installation: https://docs.zincsearch.com/installation/
