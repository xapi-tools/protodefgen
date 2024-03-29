package protodefgen

import (
	"os"
	"path"
	"testing"
)

func TestToStringBuilder(t *testing.T) {
	tests := []struct {
		name    string
		pw      ProtoDefWriter
		wantErr bool
	}{
		{
			name: "example",
			pw: NewProtoDefWriter(
				&Proto{
					Description: "This is an example file\nAutogenerated by protodefgen",
					Package:     "example",
					Imports:     []string{"google/protobuf/empty.proto"},
					Enums: []Enum{
						{
							Name:        "GlobalEnum",
							Description: "This is an Enum with global scope",
							Constants: []EnumConstant{
								{
									Name:        "VARIANT_Z",
									Value:       0,
									Description: "This is variant Z",
								},
								{
									Name:        "VARIANT_A",
									Value:       1,
									Description: "This is variant A",
								},
								{
									Name:        "VARIANT_B",
									Value:       2,
									Description: "This is variant B",
								},
							},
						},
					},
					Messages: []Message{
						{
							Name:        "BasicType",
							Description: "This message contains basic types",
							Fields: []MessageField{
								{
									Description: "This is required string field",
									Name:        "name",
									Type:        "string",
									Id:          1,
								},
								{
									Description: "This is optional numeric field",
									Name:        "num",
									Type:        "uint32",
									Optional:    true,
									Id:          2,
								},
								{
									Description: "This is optional enum field",
									Name:        "global_enum",
									Type:        "GlobalEnum",
									Optional:    true,
									Id:          3,
								},
							},
						},
						{
							Name:        "ListType",
							Description: "This message contains list types",
							Fields: []MessageField{
								{
									Description: "This is list of strings",
									Name:        "names",
									Type:        "string",
									Repeated:    true,
									Id:          1,
								},
								{
									Description: "This is list of numbers",
									Name:        "num",
									Type:        "uint32",
									Repeated:    true,
									Id:          2,
								},
								{
									Description: "This is list of enums",
									Name:        "global_enums",
									Type:        "GlobalEnum",
									Repeated:    true,
									Id:          3,
								},
							},
						},
						{
							Name:        "ComplexType",
							Description: "This message contains complex types",
							Enums: []Enum{
								{
									Name:        "NestedEnum",
									Description: "This is an Enum with message scope",
									Constants: []EnumConstant{
										{
											Name:        "VARIANT_Z",
											Value:       0,
											Description: "This is variant Z",
										},
										{
											Name:        "VARIANT_A",
											Value:       1,
											Description: "This is variant A",
										},
										{
											Name:        "VARIANT_B",
											Value:       2,
											Description: "This is variant B",
										},
									},
								},
							},
							Messages: []Message{
								{
									Name:        "NestedMessage",
									Description: "This is a nested message",
									Enums: []Enum{
										{
											Name:        "DoublyNestedEnum",
											Description: "This is a double nested Enum with message scope",
											Constants: []EnumConstant{
												{
													Name:        "VARIANT_Z",
													Value:       0,
													Description: "This is variant Z",
												},
												{
													Name:        "VARIANT_A",
													Value:       1,
													Description: "This is variant A",
												},
												{
													Name:        "VARIANT_B",
													Value:       2,
													Description: "This is variant B",
												},
											},
										},
									},
									Fields: []MessageField{
										{
											Description: "This is required string field",
											Name:        "name",
											Type:        "string",
											Id:          1,
										},
										{
											Description: "This is optional numeric field",
											Name:        "num",
											Type:        "uint32",
											Optional:    true,
											Id:          2,
										},
										{
											Description: "This is optional enum field",
											Name:        "global_enum",
											Type:        "GlobalEnum",
											Optional:    true,
											Id:          3,
										},
									},
								},
							},
							Fields: []MessageField{
								{
									Description: "This holds BasicType",
									Name:        "basic",
									Type:        "BasicType",
									Id:          1,
								},
								{
									Description: "This holds ListType",
									Name:        "list",
									Type:        "ListType",
									Id:          2,
								},
								{
									Description: "This holds NestedMessage",
									Name:        "nested_message",
									Type:        "NestedMessage",
									Id:          3,
								},
								{
									Description: "This holds NestedEnum",
									Name:        "nested_enum",
									Type:        "NestedEnum",
									Id:          4,
								},
								{
									Description: "This holds DoublyNestedEnum",
									Name:        "doubly_nested_enum",
									Type:        "NestedMessage.DoublyNestedEnum",
									Id:          5,
								},
							},
						},
					},
					Services: []Service{
						{
							Name:        "BasicService",
							Description: "This is a service exercising BasicType",
							Methods: []ServiceMethod{
								{
									Name:        "GetBasic",
									Description: "Get BasicType",
									Request:     "google.protobuf.Empty",
									Response:    "BasicType",
								},
								{
									Name:        "PushBasic",
									Description: "Push BasicType",
									Request:     "BasicType",
									Response:    "google.protobuf.Empty",
								},
							},
						},
						{
							Name:        "ComplexService",
							Description: "This is a service exercising ComplexType",
							Methods: []ServiceMethod{
								{
									Name:          "StreamRequest",
									Description:   "RPC with streaming request",
									Request:       "BasicType",
									StreamRequest: true,
									Response:      "google.protobuf.Empty",
								},
								{
									Name:           "StreamResponse",
									Description:    "RPC with streaming response",
									Request:        "google.protobuf.Empty",
									Response:       "BasicType",
									StreamResponse: true,
								},
								{
									Name:           "StreamRequestResponse",
									Description:    "RPC with bidirectional streaming",
									Request:        "BasicType",
									StreamRequest:  true,
									Response:       "ComplexType",
									StreamResponse: true,
								},
							},
						},
					},
				},
				&ProtoWriterOpts{
					IndentWidth: 4,
				},
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pw.ToStringBuilder()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestToStringBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got.String())

			if b, err := os.ReadFile(path.Join("testdata", tt.name+".proto")); err != nil {
				t.Errorf("TestToStringBuilder() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				want := string(b)
				if got.String() != want {
					t.Errorf("ProtoDefWriter.ToStringBuilder() = %v, want %v", got, want)
				}
			}
		})
	}
}
