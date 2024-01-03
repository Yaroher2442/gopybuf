# GOPYBUF

A Go project for calling Python code in runtime with cgo

## Table of Contents

- [Introduction](#introduction)
- [Using](#using)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

## Introduction

Imagine yourself as a Python!

This project aims to ensure smooth integration of Go and Python by allowing Python
functions to be called at runtime in a Go application.
Integration is achieved through the use of cgo, which allows direct calls to Python
functions via a custom rpc client-server based on the classic protobuf.

## Using

### Install

- install golang lib

  `` go get github.com/Yaroher2442/gopybuf ``
- install python lib

  `` poetry add gopybuf ``

  now we use "betterproto" library as main proto compiler, maybe changed in latest version

### Compile protobuf

- (version<1.0.0) compile your protobuf files

  `` python -m grpc_tools.protoc -I <.proto path> --python_betterproto_out=<out path>
  protoc -I <.proto path> --go_out=<out path> --go-grpc_out=<out path> ``

### Write some code

- In golang side:
  use custom client which implements grpc.ClientConnInterface

``` go
 client, err := gopybuf.NewClient("<your target python file>", "your target python funnc")
	if err != nil {
		panic(err)
		return
	}
  serviceClient := <your_proto_sdk>.<sour_service_constructor>(client)
  // example: serviceClient := gopybufSdk.NewTestServiceClient(client) 
  // and call:
  // test, err := serviceClient.Test(context.Background(), &gopybufSdk.TestMessage{
  //	Name:   "test",
  //	Age:    1,
  //	Scores: []int32{1, 2, 3},
  //})
 ```

- In python side: use custom server as service registrant

``` python
  import asyncio
  from typing import Tuple
  from gopybuf.server import IncomingBytes, OutgoingBytes, ErrorBytes, call_async, register_service
  
  from gopybuf_py_sdk import TestServiceBase, TestMessage
  
  
  class TestService(TestServiceBase): # this is rpc service inplementation fo test
    async def test(self, test_message: TestMessage) -> TestMessage:
      return test_message
  
  
  register_service(TestService()) # register service in custom server on import of file
  
  # this example of function to be called from cgo 
  # pastate in constructor in go: gopybuf.NewClient("main.py", "go_py_buf")
  def go_py_buf(method_name: str, arg: IncomingBytes) -> Tuple[OutgoingBytes, ErrorBytes]:
    return asyncio.run(call_async(method_name, arg))
```

### Let's compile it all

- first we need to have python3.X.a file with needed python version and we can find it if you already install python
  on your system

``` text
~$ find /usr -name libpython3.10.a
/usr/lib/x86_64-linux-gnu/libpython3.10.a
/usr/lib/python3.10/config-3.10-x86_64-linux-gnu/libpython3.10.a
```

this file needed "to be python", yes we compile python inside cgo and **this method will not be multiplatform**

- Now copy file in root application directory, example:

``` text
cp /usr/lib/x86_64-linux-gnu/libpython3.10.a <your_project_path>
```

- and now compile.
  on this step we need python headers files like /usr/include/python3.X.
  note: python3.X.a file in application root dir

``` text
CGO_CFLAGS=-I/usr/include/python3.10 CGO_LDFLAGS=-lpython3.10 go build
```

## Roadmap

- ~~Add grpc unary-unary requests~~
- Add grpc streaming stream-unary unary-stream stream-stream
- Add non-grpc handles and callers over bytes
    - support custom mappings with theirs registration
    - add streaming custom handlers
- Add function to call "full" python file with optional rpc
    - check if python file has function for rpc and then handle opportunity to call it

## Contributing

Feel free to contribute to this project by opening issues or submitting pull requests. Your feedback and contributions
are highly appreciated.

- Please search existing issues before creating a new one to avoid duplicates.
- Be descriptive and provide as much detail as possible.
- Follow the provided templates to ensure the necessary information is included.
- Engage in constructive discussions with the community.

### Feature Requests / Bug Report

If you have a feature | bug in mind that you'd like to see in "gopybuf", follow these steps when creating a request:

1. Click on the "Issues" tab.
2. Click the "New Issue" button.
3. Select the "Feature Request" or "Bug Report" template.
4. For "Feature Requests" - fill in the requested information, including a clear description of the problem you're
   trying to solve and your
   proposed solution.
   Assign relevant labels such as "enhancement," "new feature," and any applicable priority labels.
5. For "Bug Report" - Provide details about the bug, including steps to reproduce it and the expected vs. actual
   behavior.
   Assign relevant labels such as "bug" and any other descriptive labels.

Thank you for contributing to gopybuf!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.