package postman

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// UnmarshalJSON implements json.Unmarshaler.
func (j *Auth) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["type"]; !ok || v == nil {
		return fmt.Errorf("field type: required")
	}
	type Plain Auth
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Auth(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AuthAttribute) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["key"]; !ok || v == nil {
		return fmt.Errorf("field key: required")
	}
	type Plain AuthAttribute
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = AuthAttribute(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ProxyConfig) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	type Plain ProxyConfig
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["disabled"]; !ok || v == nil {
		plain.Disabled = false
	}
	if v, ok := raw["match"]; !ok || v == nil {
		plain.Match = "http+https://*/*"
	}
	if v, ok := raw["port"]; !ok || v == nil {
		plain.Port = 8080
	}
	if v, ok := raw["tunnel"]; !ok || v == nil {
		plain.Tunnel = false
	}
	*j = ProxyConfig(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ItemGroup) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["item"]; !ok || v == nil {
		return fmt.Errorf("field item: required")
	}
	type Plain ItemGroup
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ItemGroup(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AuthType) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_AuthType {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_AuthType, v)
	}
	*j = AuthType(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Item) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["request"]; !ok || v == nil {
		return fmt.Errorf("field request: required")
	}
	type Plain Item
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Item(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Variable) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	type Plain Variable
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["disabled"]; !ok || v == nil {
		plain.Disabled = false
	}
	if v, ok := raw["system"]; !ok || v == nil {
		plain.System = false
	}
	*j = Variable(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *VariableType) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_VariableType {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_VariableType, v)
	}
	*j = VariableType(v)
	return nil
}

type Request struct {
	Auth   Auth     `json:"auth,omitempty"`
	Method string   `json:"method"`
	Header []Header `json:"header"`
	Body   Body     `json:"body,omitempty"`
	URL    URL      `json:"url"`
}

// Represents authentication helpers provided by Postman
type Auth struct {
	// The attributes for API Key Authentication.
	Apikey []AuthAttribute `json:"apikey,omitempty"`

	// The attributes for [AWS
	// Auth](http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html).
	Awsv4 []AuthAttribute `json:"awsv4,omitempty"`

	// The attributes for [Basic
	// Authentication](https://en.wikipedia.org/wiki/Basic_access_authentication).
	Basic []AuthAttribute `json:"basic,omitempty"`

	// The helper attributes for [Bearer Token
	// Authentication](https://tools.ietf.org/html/rfc6750)
	Bearer []AuthAttribute `json:"bearer,omitempty"`

	// The attributes for [Digest
	// Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication).
	Digest []AuthAttribute `json:"digest,omitempty"`

	// The attributes for [Akamai EdgeGrid
	// Authentication](https://developer.akamai.com/legacy/introduction/Client_Auth.html).
	Edgegrid []AuthAttribute `json:"edgegrid,omitempty"`

	// The attributes for [Hawk Authentication](https://github.com/hueniverse/hawk)
	Hawk []AuthAttribute `json:"hawk,omitempty"`

	// Noauth corresponds to the JSON schema field "noauth".
	Noauth interface{} `json:"noauth,omitempty"`

	// The attributes for [NTLM
	// Authentication](https://msdn.microsoft.com/en-us/library/cc237488.aspx)
	Ntlm []AuthAttribute `json:"ntlm,omitempty"`

	// The attributes for [OAuth2](https://oauth.net/1/)
	Oauth1 []AuthAttribute `json:"oauth1,omitempty"`

	// Helper attributes for [OAuth2](https://oauth.net/2/)
	Oauth2 []AuthAttribute `json:"oauth2,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type AuthType `json:"type"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Info) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name: required")
	}
	if v, ok := raw["schema"]; !ok || v == nil {
		return fmt.Errorf("field schema: required")
	}
	type Plain Info
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Info(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Header) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["key"]; !ok || v == nil {
		return fmt.Errorf("field key: required")
	}
	if v, ok := raw["value"]; !ok || v == nil {
		return fmt.Errorf("field value: required")
	}
	type Plain Header
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["disabled"]; !ok || v == nil {
		plain.Disabled = false
	}
	*j = Header(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Event) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["listen"]; !ok || v == nil {
		return fmt.Errorf("field listen: required")
	}
	type Plain Event
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["disabled"]; !ok || v == nil {
		plain.Disabled = false
	}
	*j = Event(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Cookie) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["domain"]; !ok || v == nil {
		return fmt.Errorf("field domain: required")
	}
	if v, ok := raw["path"]; !ok || v == nil {
		return fmt.Errorf("field path: required")
	}
	type Plain Cookie
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Cookie(plain)
	return nil
}

// Represents an attribute for any authorization method provided by Postman. For
// example `username` and `password` are set as auth attributes for Basic
// Authentication method.
type AuthAttribute struct {
	// Key corresponds to the JSON schema field "key".
	Key string `json:"key"`

	// Type corresponds to the JSON schema field "type".
	Type *string `json:"type,omitempty"`

	// Value corresponds to the JSON schema field "value".
	Value interface{} `json:"value,omitempty"`
}

type AuthType string

const AuthTypeApikey AuthType = "apikey"
const AuthTypeAwsv4 AuthType = "awsv4"
const AuthTypeBasic AuthType = "basic"
const AuthTypeBearer AuthType = "bearer"
const AuthTypeDigest AuthType = "digest"
const AuthTypeEdgegrid AuthType = "edgegrid"
const AuthTypeHawk AuthType = "hawk"
const AuthTypeNoauth AuthType = "noauth"
const AuthTypeNtlm AuthType = "ntlm"
const AuthTypeOauth1 AuthType = "oauth1"
const AuthTypeOauth2 AuthType = "oauth2"

// A representation of an ssl certificate
type Certificate struct {
	// An object containing path to file certificate, on the file system
	Cert *CertificateCert `json:"cert,omitempty"`

	// An object containing path to file containing private key, on the file system
	Key *CertificateKey `json:"key,omitempty"`

	// A list of Url match pattern strings, to identify Urls this certificate can be
	// used for.
	Matches []string `json:"matches,omitempty"`

	// A name for the certificate for user reference
	Name *string `json:"name,omitempty"`

	// The passphrase for the certificate
	Passphrase *string `json:"passphrase,omitempty"`
}

// An object containing path to file certificate, on the file system
type CertificateCert struct {
	// The path to file containing key for certificate, on the file system
	Src interface{} `json:"src,omitempty"`
}

// An object containing path to file containing private key, on the file system
type CertificateKey struct {
	// The path to file containing key for certificate, on the file system
	Src interface{} `json:"src,omitempty"`
}

// A representation of a list of ssl certificates
type CertificateList []Certificate

// A Cookie, that follows the [Google Chrome
// format](https://developer.chrome.com/extensions/cookies)
type Cookie struct {
	// The domain for which this cookie is valid.
	Domain string `json:"domain"`

	// When the cookie expires.
	Expires interface{} `json:"expires,omitempty"`

	// Custom attributes for a cookie go here, such as the [Priority
	// Field](https://code.google.com/p/chromium/issues/detail?id=232693)
	Extensions []interface{} `json:"extensions,omitempty"`

	// True if the cookie is a host-only cookie. (i.e. a request's URL domain must
	// exactly match the domain of the cookie).
	HostOnly *bool `json:"hostOnly,omitempty"`

	// Indicates if this cookie is HTTP Only. (if True, the cookie is inaccessible to
	// client-side scripts)
	HttpOnly *bool `json:"httpOnly,omitempty"`

	// MaxAge corresponds to the JSON schema field "maxAge".
	MaxAge *string `json:"maxAge,omitempty"`

	// This is the name of the Cookie.
	Name *string `json:"name,omitempty"`

	// The path associated with the Cookie.
	Path string `json:"path"`

	// Indicates if the 'secure' flag is set on the Cookie, meaning that it is
	// transmitted over secure connections only. (typically HTTPS)
	Secure *bool `json:"secure,omitempty"`

	// True if the cookie is a session cookie.
	Session *bool `json:"session,omitempty"`

	// The value of the Cookie.
	Value *string `json:"value,omitempty"`
}

// A representation of a list of cookies
type CookieList []Cookie

// A Description can be a raw text, or be an object, which holds the description
// along with its format.
type Description interface{}

// Defines a script associated with an associated event name
type Event struct {
	// Indicates whether the event is disabled. If absent, the event is assumed to be
	// enabled.
	Disabled bool `json:"disabled,omitempty"`

	// A unique identifier for the enclosing event.
	Id *string `json:"id,omitempty"`

	// Can be set to `test` or `prerequest` for test scripts or pre-request scripts
	// respectively.
	Listen string `json:"listen"`

	// Script corresponds to the JSON schema field "script".
	Script *Script `json:"script,omitempty"`
}

// Postman allows you to configure scripts to run when specific events occur. These
// scripts are stored here, and can be referenced in the collection by their ID.
type EventList []Event

// Represents a single HTTP Header
type Header struct {
	// Description corresponds to the JSON schema field "description".
	Description HeaderDescription `json:"description,omitempty"`

	// If set to true, the current header will not be sent with requests.
	Disabled bool `json:"disabled,omitempty"`

	// This holds the LHS of the HTTP Header, e.g ``Content-Type`` or
	// ``X-Custom-Header``
	Key string `json:"key"`

	// The value (or the RHS) of the Header is stored in this field.
	Value string `json:"value"`
}

type HeaderDescription interface{}

// A representation for a list of headers
type HeaderList []Header

// Detailed description of the info block
type Info struct {
	// Every collection is identified by the unique value of this field. The value of
	// this field is usually easiest to generate using a UID generator function. If
	// you already have a collection, it is recommended that you maintain the same id
	// since changing the id usually implies that is a different collection than it
	// was originally.
	//  *Note: This field exists for compatibility reasons with Collection Format V1.*
	PostmanId *string `json:"_postman_id,omitempty"`

	// Description corresponds to the JSON schema field "description".
	Description InfoDescription `json:"description,omitempty"`

	// A collection's friendly name is defined by this field. You would want to set
	// this field to a value that would allow you to easily identify this collection
	// among a bunch of other collections, as such outlining its usage or content.
	Name string `json:"name"`

	// This should ideally hold a link to the Postman schema that is used to validate
	// this collection. E.g: https://schema.getpostman.com/collection/v1
	Schema string `json:"schema"`

	// Version corresponds to the JSON schema field "version".
	Version InfoVersion `json:"version,omitempty"`
}

type InfoDescription interface{}

type InfoVersion interface{}

// Items are entities which contain an actual HTTP request, and sample responses
// attached to it.
type Item struct {
	// Description corresponds to the JSON schema field "description".
	Description ItemDescription `json:"description,omitempty"`

	// Event corresponds to the JSON schema field "event".
	Event EventList `json:"event,omitempty"`

	// A unique ID that is used to identify collections internally
	Id *string `json:"id,omitempty"`

	// A human readable identifier for the current item.
	Name *string `json:"name,omitempty"`

	// ProtocolProfileBehavior corresponds to the JSON schema field
	// "protocolProfileBehavior".
	ProtocolProfileBehavior ProtocolProfileBehavior `json:"protocolProfileBehavior,omitempty"`

	// Request corresponds to the JSON schema field "request".
	Request Request `json:"request"`

	// Response corresponds to the JSON schema field "response".
	Response []Response `json:"response,omitempty"`

	// Variable corresponds to the JSON schema field "variable".
	Variable VariableList `json:"variable,omitempty"`
}

type ItemDescription interface{}

// One of the primary goals of Postman is to organize the development of APIs. To
// this end, it is necessary to be able to group requests together. This can be
// achived using 'Folders'. A folder just is an ordered set of requests.
type ItemGroup struct {
	// Auth corresponds to the JSON schema field "auth".
	Auth interface{} `json:"auth,omitempty"`

	// Description corresponds to the JSON schema field "description".
	Description ItemGroupDescription `json:"description,omitempty"`

	// Event corresponds to the JSON schema field "event".
	Event EventList `json:"event,omitempty"`

	// Items are entities which contain an actual HTTP request, and sample responses
	// attached to it. Folders may contain many items.
	Items []Item `json:"item"`

	// A folder's friendly name is defined by this field. You would want to set this
	// field to a value that would allow you to easily identify this folder.
	Name *string `json:"name,omitempty"`

	// ProtocolProfileBehavior corresponds to the JSON schema field
	// "protocolProfileBehavior".
	ProtocolProfileBehavior ProtocolProfileBehavior `json:"protocolProfileBehavior,omitempty"`

	// Variable corresponds to the JSON schema field "variable".
	Variable VariableList `json:"variable,omitempty"`
}

type ItemGroupDescription interface{}

type Postman struct {
	// Auth corresponds to the JSON schema field "auth".
	Auth interface{} `json:"auth,omitempty"`

	// Event corresponds to the JSON schema field "event".
	Event EventList `json:"event,omitempty"`

	// Info corresponds to the JSON schema field "info".
	Info Info `json:"info"`

	// Items are the basic unit for a Postman collection. You can think of them as
	// corresponding to a single API endpoint. Each Item has one request and may have
	// multiple API responses associated with it.
	Items []Item `json:"item"`

	// ProtocolProfileBehavior corresponds to the JSON schema field
	// "protocolProfileBehavior".
	ProtocolProfileBehavior ProtocolProfileBehavior `json:"protocolProfileBehavior,omitempty"`

	// Variable corresponds to the JSON schema field "variable".
	Variable VariableList `json:"variable,omitempty"`
}

// Set of configurations used to alter the usual behavior of sending the request
type ProtocolProfileBehavior map[string]interface{}

// Using the Proxy, you can configure your custom proxy into the postman for
// particular url match
type ProxyConfig struct {
	// When set to true, ignores this proxy configuration entity
	Disabled bool `json:"disabled,omitempty"`

	// The proxy server host
	Host *string `json:"host,omitempty"`

	// The Url match for which the proxy config is defined
	Match string `json:"match,omitempty"`

	// The proxy server port
	Port int `json:"port,omitempty"`

	// The tunneling details for the proxy config
	Tunnel bool `json:"tunnel,omitempty"`
}

// A request represents an HTTP request. If a string, the string is assumed to be
// the request URL and the method is assumed to be 'GET'.
// type Request interface{}

// A response represents an HTTP response.
type Response interface{}

// A script is a snippet of Javascript code that can be used to to perform setup or
// teardown operations on a particular response.
type Script struct {
	// Exec corresponds to the JSON schema field "exec".
	Exec interface{} `json:"exec,omitempty"`

	// A unique, user defined identifier that can  be used to refer to this script
	// from requests.
	Id *string `json:"id,omitempty"`

	// Script name
	Name *string `json:"name,omitempty"`

	// Src corresponds to the JSON schema field "src".
	Src ScriptSrc `json:"src,omitempty"`

	// Type of the script. E.g: 'text/javascript'
	Type *string `json:"type,omitempty"`
}

type ScriptSrc interface{}

// If object, contains the complete broken-down URL for this request. If string,
// contains the literal request URL.
// type Url interface{}
type URL struct {
	Raw  string   `json:"raw"`
	Host []string `json:"host"`
	Path []string `json:"path"`
}

type Urlencoded struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Formdata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Body struct {
	Mode       string        `json:"mode,omitempty"`
	Raw        string        `json:"raw,omitempty"`
	Urlencoded []*Urlencoded `json:"urlencoded,omitempty"`
	Formdata   []*Formdata   `json:"formdata,omitempty"`

	// --- NOT IMPLEMENTED HERE YET ---
	// File       File       `json:"file,omitempty"`
	// Graphql    Graphql    `json:"graphql,omitempty"`
	// Options    Options    `json:"options,omitempty"`
	// Disabled   Disabled   `json:"disabled,omitempty"`
}

// Using variables in your Postman requests eliminates the need to duplicate
// requests, which can save a lot of time. Variables can be defined, and referenced
// to from any part of a request.
type Variable struct {
	// Description corresponds to the JSON schema field "description".
	Description VariableDescription `json:"description,omitempty"`

	// Disabled corresponds to the JSON schema field "disabled".
	Disabled bool `json:"disabled,omitempty"`

	// A variable ID is a unique user-defined value that identifies the variable
	// within a collection. In traditional terms, this would be a variable name.
	Id *string `json:"id,omitempty"`

	// A variable key is a human friendly value that identifies the variable within a
	// collection. In traditional terms, this would be a variable name.
	Key *string `json:"key,omitempty"`

	// Variable name
	Name *string `json:"name,omitempty"`

	// When set to true, indicates that this variable has been set by Postman
	System bool `json:"system,omitempty"`

	// A variable may have multiple types. This field specifies the type of the
	// variable.
	Type *VariableType `json:"type,omitempty"`

	// The value that a variable holds in this collection. Ultimately, the variables
	// will be replaced by this value, when say running a set of requests from a
	// collection
	Value interface{} `json:"value,omitempty"`
}

type VariableDescription interface{}

// Collection variables allow you to define a set of variables, that are a *part of
// the collection*, as opposed to environments, which are separate entities.
// *Note: Collection variables must not contain any sensitive information.*
type VariableList []Variable

type VariableType string

const VariableTypeAny VariableType = "any"
const VariableTypeBoolean VariableType = "boolean"
const VariableTypeNumber VariableType = "number"
const VariableTypeString VariableType = "string"

// Postman allows you to version your collections as they grow, and this field
// holds the version number. While optional, it is recommended that you use this
// field to its fullest extent!
type Version interface{}

var enumValues_AuthType = []interface{}{
	"apikey",
	"awsv4",
	"basic",
	"bearer",
	"digest",
	"edgegrid",
	"hawk",
	"noauth",
	"oauth1",
	"oauth2",
	"ntlm",
}
var enumValues_VariableType = []interface{}{
	"string",
	"boolean",
	"any",
	"number",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Postman) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["info"]; !ok || v == nil {
		return fmt.Errorf("field info: required")
	}
	if v, ok := raw["item"]; !ok || v == nil {
		return fmt.Errorf("field item: required")
	}
	type Plain Postman
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Postman(plain)
	return nil
}
