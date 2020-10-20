# Runtimes

Contains a series of dockerfiles that can be used to build images that can be used a runtimes to run certain pieces of softwares.

## protobuild.dockerfile
It's the docker image used to build the .proto files, this is because we need to leverage gogoproto which is not currently
 supported by latest protobuf versions of protobuf. And also it exists so you do not have to download the whole protobuf,
protoc-gen-go executable. The image will build everything for you.