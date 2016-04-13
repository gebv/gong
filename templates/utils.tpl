{{define "helpfull" }}
// FromJson extract object from data (io.Reader OR []byte)
func FromJson(obj interface{}, data interface{}) error {
	switch data.(type) {
	case io.Reader:
		decoder := json.NewDecoder(data.(io.Reader))
		return decoder.Decode(obj)
	case []byte:
		return json.Unmarshal(data.([]byte), obj)
	}

	return ErrNotSupported
}
{{end}}

{{ define "stringer" }}
{{ $prefixField := .fieldprefix }}
{{range $field := .strings -}}
var {{toUpper $prefixField}}{{toUpper $field.name}} = "{{toLower $field.name "_"}}"
{{ end }} 
{{ end }}

{{ define "errors" }}
{{range $error := .errors -}}
{{- $message := (or $error.message (toLower $error.name  "_"))}}
{{- with $error.comment}}// Err{{$error.name | toUpper}} {{ $error.comment }}.{{ end }}
var Err{{$error.name | toUpper}} = fmt.Errorf("{{$message}}")
{{ end }}
{{ end }}