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
{{range $index, $name := . -}}
var {{$name | toUpper}} = "{{$name | toLower "_"}}"
{{ end }} 
{{ end }}

{{ define "errors" }}
{{range $error := .errors -}}
{{- $message := (or $error.message (toLower $error.name  "_"))}}
{{- with $error.comment}}// Err{{$error.name | toUpper}} {{ $error.comment }}.{{ end }}
var Err{{$error.name | toUpper}} = fmt.Errorf("{{$message}}")
{{ end }}
{{ end }}