package audit

type OpenAPIDoc struct {
	OpenAPI    string                     `json:"openapi"`
	Info       OpenAPIInfo                `json:"info"`
	Servers    []OpenAPIServer            `json:"servers,omitempty"`
	Paths      map[string]OpenAPIPathItem `json:"paths"`
	Components *OpenAPIComponents         `json:"components,omitempty"`
}

type OpenAPIInfo struct {
	Title          string `json:"title"`
	Version        string `json:"version"`
	Description    string `json:"description,omitempty"`
	TermsOfService string `json:"termsOfService,omitempty"`
}

type OpenAPIServer struct {
	URL         string                    `json:"url"`
	Description string                    `json:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty"`
}

type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default"`
	Description string   `json:"description,omitempty"`
}

type OpenAPIPathItem struct {
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`

	GET     *OpenAPIOperation `json:"get,omitempty"`
	POST    *OpenAPIOperation `json:"post,omitempty"`
	PUT     *OpenAPIOperation `json:"put,omitempty"`
	PATCH   *OpenAPIOperation `json:"patch,omitempty"`
	DELETE  *OpenAPIOperation `json:"delete,omitempty"`
	HEAD    *OpenAPIOperation `json:"head,omitempty"`
	OPTIONS *OpenAPIOperation `json:"options,omitempty"`
	TRACE   *OpenAPIOperation `json:"trace,omitempty"`

	Parameters []OpenAPIParameter `json:"parameters,omitempty"`
}

type OpenAPIOperation struct {
	Tags        []string `json:"tags,omitempty"`
	OperationID string   `json:"operationId,omitempty"`
	Summary     string   `json:"summary,omitempty"`
	Description string   `json:"description,omitempty"`
	Deprecated  bool     `json:"deprecated,omitempty"`

	Parameters  []OpenAPIParameter         `json:"parameters,omitempty"`
	RequestBody *OpenAPIRequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]OpenAPIResponse `json:"responses"`

	Security []map[string][]string `json:"security,omitempty"`
}

type OpenAPIParameter struct {
	Name        string         `json:"name"`
	In          string         `json:"in"` // query, header, path, cookie
	Description string         `json:"description,omitempty"`
	Required    bool           `json:"required,omitempty"`
	Deprecated  bool           `json:"deprecated,omitempty"`
	Schema      *OpenAPISchema `json:"schema,omitempty"`
	Example     any            `json:"example,omitempty"`
}

type OpenAPIRequestBody struct {
	Description string                      `json:"description,omitempty"`
	Required    bool                        `json:"required,omitempty"`
	Content     map[string]OpenAPIMediaType `json:"content,omitempty"`
}

type OpenAPIResponse struct {
	Description string                      `json:"description"`
	Headers     map[string]OpenAPIHeader    `json:"headers,omitempty"`
	Content     map[string]OpenAPIMediaType `json:"content,omitempty"`
}

type OpenAPIHeader struct {
	Description string         `json:"description,omitempty"`
	Required    bool           `json:"required,omitempty"`
	Deprecated  bool           `json:"deprecated,omitempty"`
	Schema      *OpenAPISchema `json:"schema,omitempty"`
	Example     any            `json:"example,omitempty"`
}

type OpenAPIMediaType struct {
	Schema   *OpenAPISchema `json:"schema,omitempty"`
	Example  any            `json:"example,omitempty"`
	Examples map[string]any `json:"examples,omitempty"`
	Encoding map[string]any `json:"encoding,omitempty"`
}

type OpenAPIComponents struct {
	Schemas         map[string]OpenAPISchema         `json:"schemas,omitempty"`
	Responses       map[string]OpenAPIResponse       `json:"responses,omitempty"`
	Parameters      map[string]OpenAPIParameter      `json:"parameters,omitempty"`
	RequestBodies   map[string]OpenAPIRequestBody    `json:"requestBodies,omitempty"`
	Headers         map[string]OpenAPIHeader         `json:"headers,omitempty"`
	SecuritySchemes map[string]OpenAPISecurityScheme `json:"securitySchemes,omitempty"`
}

type OpenAPISecurityScheme struct {
	Type             string         `json:"type"` // apiKey, http, mutualTLS, oauth2, openIdConnect
	Description      string         `json:"description,omitempty"`
	Name             string         `json:"name,omitempty"`
	In               string         `json:"in,omitempty"`     // query, header, cookie
	Scheme           string         `json:"scheme,omitempty"` // bearer, basic
	BearerFormat     string         `json:"bearerFormat,omitempty"`
	Flows            map[string]any `json:"flows,omitempty"`
	OpenIDConnectURL string         `json:"openIdConnectUrl,omitempty"`
}

type OpenAPISchema struct {
	Ref string `json:"$ref,omitempty"`

	Type        string `json:"type,omitempty"`
	Format      string `json:"format,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Default     any    `json:"default,omitempty"`
	Example     any    `json:"example,omitempty"`

	Properties map[string]OpenAPISchema `json:"properties,omitempty"`
	Required   []string                 `json:"required,omitempty"`
	Items      *OpenAPISchema           `json:"items,omitempty"`

	Enum []any `json:"enum,omitempty"`

	Nullable  bool `json:"nullable,omitempty"`
	ReadOnly  bool `json:"readOnly,omitempty"`
	WriteOnly bool `json:"writeOnly,omitempty"`

	MinLength *int     `json:"minLength,omitempty"`
	MaxLength *int     `json:"maxLength,omitempty"`
	Minimum   *float64 `json:"minimum,omitempty"`
	Maximum   *float64 `json:"maximum,omitempty"`

	Pattern string `json:"pattern,omitempty"`

	AdditionalProperties any `json:"additionalProperties,omitempty"`

	AllOf []OpenAPISchema `json:"allOf,omitempty"`
	OneOf []OpenAPISchema `json:"oneOf,omitempty"`
	AnyOf []OpenAPISchema `json:"anyOf,omitempty"`
	Not   *OpenAPISchema  `json:"not,omitempty"`
}
