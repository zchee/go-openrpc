// Copyright 2019 The go-openrpc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openrpc

import (
	"encoding/json"
	"strconv"

	"github.com/zchee/go-openrpc/internal/jsonschema"
)

// Schema represents a root document object of the OpenRPC document.
type Schema struct {
	// This string MUST be the semantic version number of the OpenRPC Specification version that the OpenRPC document uses.
	// The `openrpc` field SHOULD be used by tooling specifications and clients to interpret the OpenRPC document. This is not related to the API `info.version` string.
	//
	// REQUIRED.
	OpenRPC string `json:"openrpc"`

	// Provides metadata about the API. The metadata MAY be used by tooling as required.
	//
	// REQUIRED.
	Info *Info `json:"info"`

	// An array of Server Objects, which provide connectivity information to a target server.
	// If the servers property is not provided, or is an empty array, the default value would be a Server with a url value of `localhost`.
	Servers []*Server `json:"servers,omitempty"`

	// The available methods for the API. While it is required, the array may be empty (to handle security filtering, for example).
	//
	// REQUIRED.
	Methods []*Method `json:"methods"`

	// An element to hold various schemas for the specification.
	Components *Components `json:"components,omitempty"`

	// Additional external documentation.
	ExternalDocs *ExternalDocumentation `json:"externaldocs,omitempty"`

	// Allows extensions to the OpenRPC Schema.
	Extensions []*Extension `json:"-"`
}

// Info provides metadata about the API.
//
// The metadata MAY be used by the clients if needed, and MAY be presented in editing or documentation generation tools for convenience.
type Info struct {
	// The title of the application.
	//
	// REQUIRED.
	Title string `json:"title"`

	// A verbose description of the application. GitHub Flavored Markdown syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// A URL to the Terms of Service for the API. MUST be in the format of a URL.
	TermsOfService string `json:"termsOfService,omitempty"`

	// The contact information for the exposed API.
	Contact *Contact `json:"contact,omitempty"`

	// A URL to the Terms of Service for the API. MUST be in the format of a URL.
	License *License `json:"license,omitempty"`

	// The version of the OpenRPC document (which is distinct from the OpenRPC Specification version or the API implementation version).
	//
	// REQUIRED.
	Version string `json:"version"`

	// Allows extensions to the OpenRPC Schema.
	Extensions []*Extension `json:"-"`
}

// Contact represents a information for the exposed API.
type Contact struct {
	// The identifying name of the contact person/organization.
	Name string `json:"name,omitempty"`

	// The URL pointing to the contact information. MUST be in the format of a URL.
	URL string `json:"url,omitempty"`

	// The email address of the contact person/organization. MUST be in the format of an email address.
	Email string `json:"email,omitempty"`

	// Allows extensions to the OpenRPC Schema.
	Extensions []*Extension `json:"-"`
}

// License represents a information for the exposed API.
type License struct {
	// The license name used for the API.
	//
	// REQUIRED.
	Name string `json:"name"`

	// A URL to the license used for the API. MUST be in the format of a URL.
	URL string `json:"url,omitempty"`

	// Allows extensions to the OpenRPC Schema.
	Extensions []*Extension `json:"-"`
}

// Server represents a Server.
type Server struct {
	// A name to be used as the canonical name for the server.
	//
	// REQUIRED.
	Name string `json:"name"`

	// A URL to the target host.
	// This URL supports Server Variables and MAY be relative, to indicate that the host location is relative to the location where the OpenRPC document is being served.
	// Server Variables are passed into the RuntimeExpression to produce a server URL.
	//
	// REQUIRED.
	URL string `json:"url"`

	// A short summary of what the server is.
	Summary string `json:"summary,omitempty"`

	// An optional string describing the host designated by the URL. GitHub Flavored Markdown syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// A map between a variable name and its value. The value is passed into the RuntimeExpression to produce a server URL.
	Variables map[string]*ServerVariables `json:"variables,omitempty"`

	// Allows extensions to the OpenRPC Schema.
	Extensions []*Extension `json:"-"`
}

// ServerVariables a Server Variable for server URL template substitution.
type ServerVariables struct {
	// An enumeration of string values to be used if the substitution options are from a limited set.
	Enum []string `json:"enum,omitempty"`

	// The default value to use for substitution, which SHALL be sent if an alternate value is not supplied.
	// Note this behavior is different than the Schema Object’s treatment of default values, because in those cases parameter values are optional.
	//
	// REQUIRED.
	Default string `json:"default"`

	// An optional description for the server variable. GitHub Flavored Markdown syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`
}

// ParamStructure is the expected format of the parameters.
//
// As per the JSON-RPC 2.0 specification, params may be either an array, an object, or either. Defaults to "by-position".
type ParamStructure int

const (
	ByPosition ParamStructure = iota
	ByName
	Either
)

// String implements fmt.Stringer.
func (p ParamStructure) String() string {
	switch p {
	case ByPosition:
		return "by-position"
	case ByName:
		return "by-name"
	case Either:
		return "either"
	default:
		return strconv.FormatInt(int64(p), 10)
	}
}

// Method describes the interface for the given method name.
//
// The method name is used as the method field of the JSON-RPC body. It therefore MUST be unique.
type Method struct {
	// Name is the canonical name for the method. The name MUST be unique within the methods array.
	//
	// REQUIRED.
	Name string `json:"name"`

	// A list of tags for API documentation control. Tags can be used for logical grouping of methods by resources or any other qualifier.
	Tags []*Tag `json:"tags,omitempty"`

	// A short summary of what the method does.
	Summary string `json:"summary,omitempty"`

	// A verbose explanation of the method behavior. GitHub Flavored Markdown syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Additional external documentation for this method.
	ExternalDocs *ExternalDocumentation `json:"externaldocs,omitempty"`

	// A list of parameters that are applicable for this method. The list MUST NOT include duplicated parameters and therefore require name to be unique.
	// The list can use the Reference Object to link to parameters that are defined by the ContentDescriptor.
	// It may also nest the content descriptor or reference object inside of a OneOf Object.
	//
	// All optional params (content descriptor objects with “required”: false) MUST be positioned after all required params in the list.
	//
	// REQUIRED.
	Params []*ContentDescriptor `json:"params"`

	// The description of the result returned by the method. It MUST be a Content Descriptor.
	//
	// REQUIRED.
	Result *ContentDescriptor `json:"result"`

	// Declares this method to be deprecated. Consumers SHOULD refrain from usage of the declared method. Default value is `false`.
	Deprecated bool `json:"deprecated,omitempty"`

	// An alternative `servers` array to service this method. If an alternative `servers` array is specified at the Root level, it will be overridden by this value.
	Servers []*Server `json:"servers,omitempty"`

	// A list of custom application defined errors that MAY be returned. The Errors MUST have unique error codes.
	Errors []*Error `json:"errors,omitempty"`

	// A list of possible links from this method call.
	Links []*Link `json:"links,omitempty"`

	// The expected format of the parameters. As per the JSON-RPC 2.0 specification, params may be either an array, an object, or either. Defaults to "by-position".
	ParamStructure ParamStructure `json:"paramStructure,omitempty"`

	// Array of Example Pairing Object where each example includes a valid params-to-result Content Descriptor pairing.
	Examples []*ExamplePairing `json:"examples,omitempty"`

	// Allows extensions to the OpenRPC Schema.
	Extensions []*Extension `json:"-"`
}

// ContentDescriptor descriptors are objects that do just as they suggest - describe content.
//
// They are reusable ways of describing either parameters or result. They MUST have a schema.
type ContentDescriptor struct {
	// Array of Example Pairing Object where each example includes a valid params-to-result Content Descriptor pairing.
	Examples []*ExamplePairing `json:"examples,omitempty"`

	// Name of the content that is being described.
	//
	// REQUIRED.
	Name string `json:"name"`

	// A short summary of the content that is being described.
	Summary string `json:"summary,omitempty"`

	// A verbose explanation of the content descriptor behavior. GitHub Flavored Markdown syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Schema that describes the content.
	//
	// REQUIRED.
	Schema *JSONSchema `json:"schema"`

	// Determines if the content is a required field. Default value is `false`.
	Required bool `json:"required,omitempty"`

	// Specifies that the content is deprecated and SHOULD be transitioned out of usage. Default value is `false`.
	Deprecated bool `json:"deprecated,omitempty"`
}

// JSONSchema is the Schema Object allows the definition of input and output data types.
// The JSONSchema MUST follow the specifications outline in the JSON Schema Specification 7 Alternatively, any time a JSONSchema can be used, a Reference Object can be used in its place.
// This allows referencing definitions instead of defining them inline.
type JSONSchema struct {
	*jsonschema.Schema

	// Allows extensions to the OpenRPC Schema.
	Extensions []*Extension `json:"-"`
}

// ExamplePairing represents a pairing object consists of a set of example params and result.
//
// The result is what you can expect from the JSON-RPC service given the exact params.
type ExamplePairing struct {
	// Name for the example pairing.
	Name string `json:"name,omitempty"`

	// A verbose explanation of the example pairing.
	Description string `json:"description,omitempty"`

	// Short description for the example pairing.
	Summary string `json:"summary,omitempty"`

	// Example parameters.
	Params []*Example `json:"params,omitempty"`

	// Example result.
	Result *Example `json:"result,omitempty"`

	// Allows extensions to the OpenRPC Schema.
	Extensions []*Extension `json:"-"`
}

// Example defines an example that is intended to match a given Content Descriptor Schema.
//
// If the Content Descriptor Schema includes examples, the value from this Example Object supercedes the value of the schema example.
type Example struct {
	// Canonical name of the example.
	Name string `json:"name,omitempty"`

	// Short description for the example.
	Summary string `json:"tags,omitempty"`

	// A verbose explanation of the example. GitHub Flavored Markdown syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Embedded literal example.
	// The `value` field and 8externalValue8 field are mutually exclusive.
	//
	// To represent examples of media types that cannot naturally represented in JSON, use a string value to contain the example, escaping where necessary.
	Value interface{} `json:"value,omitempty"`

	// A URL that points to the literal example. This provides the capability to reference examples that cannot easily be included in JSON documents.
	// The `value` field and `externalValue` field are mutually exclusive.
	ExternalValue string `json:"externalValue,omitempty"`

	// Allows extensions to the OpenRPC Schema.
	Extensions []*Extension `json:"-"`
}

// Link represents a possible design-time link for a result.
//
// The presence of a link does not guarantee the caller’s ability to successfully invoke it, rather it provides a known relationship and traversal mechanism between results and other methods.
type Link struct {
	// Canonical name of the link.
	//
	// REQUIRED.
	Name string `json:"name"`

	// A description of the link. GitHub Flavored Markdown syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Short description for the link.
	Summary string `json:"summary,omitempty"`

	// The name of an existing, resolvable OpenRPC method, as defined with a unique `method`.
	//
	// This field MUST resolve to a unique Method. As opposed to Open Api, Relative `method` values ARE NOT permitted.
	Method string `json:"method,omitempty"`

	// A map representing parameters to pass to a method as specified with `method`.
	//
	// The key is the parameter name to be used, whereas the value can be a constant or a RuntimeExpression to be evaluated and passed to the linked method.
	Params map[interface{}]RuntimeExpressions `json:"params,omitempty"`

	// A server object to be used by the target method.
	Server *Server `json:"server,omitempty"`

	// Allows extensions to the OpenRPC Schema.
	Extensions []*Extension `json:"-"`
}

// RuntimeExpressions allow the user to define an expression which will evaluate to a string once the desired value(s) are known.
//
// They are used when the desired value of a link or server can only be constructed at run time. This mechanism is used by Link Objects and Server Variables.
type RuntimeExpressions string

// ErrorCode is a number that indicates the error type that occurred.
type ErrorCode int64

const (
	// ParseError is the invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.
	ParseError = ErrorCode(-32700)

	// InvalidRequest is the JSON sent is not a valid Request object.
	InvalidRequest = ErrorCode(-32600)

	// MethodNotFound is the method does not exist / is not available.
	MethodNotFound = ErrorCode(-32601)

	// InvalidParams is the invalid method parameter(s).
	InvalidParams = ErrorCode(-32602)

	// InternalError is the internal JSON-RPC error.
	InternalError = ErrorCode(-32603)

	// Reserved for implementation-defined server-errors.
	_ = ErrorCode(-32099) // codeServerErrorStart
	// Reserved for implementation-defined server-errors.
	_ = ErrorCode(-32000) // codeServerErrorEnd
)

// Error defines an application level error.
type Error struct {
	// A Number that indicates the error type that occurred. This MUST be an integer.
	//
	// The error codes from and including -32768 to -32000 are reserved for pre-defined errors.
	// These pre-defined errors SHOULD be assumed to be returned from any JSON-RPC API.
	//
	// REQUIRED.
	Code ErrorCode `json:"code"`

	// A String providing a short description of the error. The message SHOULD be limited to a concise single sentence.
	//
	// REQUIRED.
	Message string `json:"message"`

	// A Primitive or Structured value that contains additional information about the error. This may be omitted.
	//
	// The value of this member is defined by the Server (e.g. detailed error information, nested errors etc.).
	Data json.RawMessage `json:"data,omitempty"`
}

// Components holds a set of reusable objects for different aspects of the OpenRPC.
//
// All objects defined within the components object will have no effect on the API unless they are explicitly referenced from properties outside the components object.
type Components struct {
	// An object to hold reusable ContentDescriptor.
	ContentDescriptors map[string]*ContentDescriptor `json:"contentDescriptors,omitempty"`

	// An object to hold reusable JSONSchema.
	Schemas map[string]*JSONSchema `json:"schemas,omitempty"`

	// An object to hold reusable Example.
	Examples map[string]*Example `json:"examples,omitempty"`

	// An object to hold reusable Link.
	Links map[string]*Link `json:"links,omitempty"`

	// An object to hold reusable Error.
	Errors map[string]*Error `json:"errors,omitempty"`

	// An object to hold reusable ExamplePairing.
	ExamplePairingObjects map[string]*ExamplePairing `json:"examplePairingObjects,omitempty"`

	// An object to hold reusable Tag.
	Tags map[string]*Tag `json:"tags,omitempty"`

	// Allows extensions to the OpenRPC Schema.
	Extensions []*Extension `json:"-"`
}

// Tag adds metadata to a single tag that is used by the Method Object.
//
// It is not mandatory to have a Tag Object per tag defined in the Method Object instances.
type Tag struct {
	// The name of the tag.
	//
	// REQUIRED.
	Name string `json:"name"`

	// A short summary of the tag.
	Summary string `json:"summary,omitempty"`

	// A verbose explanation for the tag. GitHub Flavored Markdown syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// Additional external documentation for this tag.
	ExternalDocs *ExternalDocumentation `json:"externaldocs,omitempty"`

	// Allows extensions to the OpenRPC Schema.
	Extensions []*Extension `json:"-"`
}

// ExternalDocumentation allows referencing an external resource for extended documentation.
type ExternalDocumentation struct {
	// A verbose explanation of the target documentation. GitHub Flavored Markdown syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`

	// The URL for the target documentation. Value MUST be in the format of a URL.
	//
	// REQUIRED.
	URL string `json:"url"`

	// Allows extensions to the OpenRPC Schema.
	Extensions []*Extension `json:"-"`
}

// Reference is a simple object to allow referencing other components in the specification, internally and externally.
//
// The Reference Object is defined by JSON Schema and follows the same structure, behavior and rules.
type Reference struct {
	// The reference string.
	//
	// REQUIRED.
	Ref string `json:"$ref"`
}

// OneOf a simple object allowing for conditional content descriptors.
//
// It MUST only be used in place of a content descriptor.
// It specifies that the content descriptor in question must match one of the listed content descriptors.
//
// This allows you to define content descriptors more granularly, without having to rely so heavily on json schemas.
type OneOf struct {
	// The reference string.
	//
	// REQUIRED.
	OneOf *ContentDescriptor `json:"oneOf"`
}

// Extension while the OpenRPC Specification tries to accommodate most use cases, additional data can be added to extend the specification at certain points.
//
// The extensions properties are implemented as patterned fields that are always prefixed by `"x-"`.
type Extension struct {
	// Allows extensions to the OpenRPC Schema.
	// The field name MUST begin with `x-`, for example, `x-internal-id`.
	//
	// The value can be `null`, a primitive, an array or an object. Can have any valid JSON format value.
	Pattern []interface{} `json:"-"` // ^x-
}
