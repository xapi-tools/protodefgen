# protofilegen - Protocol Buffer File Generator

This package provides Go structures for constructs used in Protocol Buffer file and a helper to write those constructs to `.proto` file with valid syntax.

### Usage

1. Get the package

    ```bash
    go get github.com/xapi-tools/protofilegen@latest
    ```

2. Create Go structures for Protocol Buffers and generate file

    ```go
    package main

    import pfg "github.com/xapi-tools/protofilegen"

    func main() {
        pw := pfg.NewProtoFileWriter(
            &pfg.Proto{
                Package: "test",
                Imports: []string{"google/protobuf/empty.proto"},
                Messages: []pfg.Message{
                    {
                        Name:        "BasicType",
                        Description: "This message contains basic types",
                        Fields: []pfg.MessageField{
                            {
                                Description: "This is required string field",
                                Name:        "name",
                                Type:        "string",
                                Id:          0,
                            },
                        },
                    },
                },
                Services: []pfg.Service{
                    {
                        Name:        "BasicService",
                        Description: "This is a service exercising BasicType",
                        Methods: []pfg.ServiceMethod{
                            {
                                Name:        "GetBasic",
                                Description: "Get BasicType",
                                Request:     "google.protobuf.Empty",
                                Response:    "BasicType",
                            },
                        },
                    },
                },
            },
            &pfg.ProtoWriterOpts{
                IndentWidth: 4,
            },
        )

        pw.ToFile("./test.proto")
    }
    ```


    Check following for more comprehensive examples:
    - [Unit Test](gen_test.go)
    - [Output](testdata/example.proto)
