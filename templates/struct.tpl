// Создание структуры
{{define "struct" }}
{{with .disableConstructor }}
{{else}}
func New{{toUpper .name}}() *{{.name}} {
    model := new({{.name}})
    {{ range $key, $field := .fields}}
    {{if hasPrefix $field.type "map"}}
    model.{{$field.name}} = make({{$field.type}})
    {{end}}  
    {{ end }}
    return model
}
{{end}}

{{- with .comment }}// {{.name}} {{.comment}}{{end}}
type {{.name}} struct {
    {{ range $key, $field := .fields}}
    {{with $field.comment}}// {{$field.comment}}{{end}}
    {{$field.name}} {{$field.type}} {{template "structtags" $field.tag}}  
    {{ end }}
}

{{- $structName := .name}}
{{range $key, $field := .fields}}
{{- template "setter" (map "structname" $structName "field" $field) -}}
{{- template "getter" (map "structname" $structName "field" $field) -}}  
{{ end }}

{{ if .transform}}
{{- template "transform" (map "transform" .transform "structname" $structName) -}}
{{end}}

{{ end }}

{{define "structtags" }}{{with .}}`{{.}}`{{end}}{{end}}

// Вспомогательные функции для структуры (в зависимости от типа поля)
{{define "setter"}}

{{- if and (eq (hasPrefix .field.type "map") false) (eq (hasPrefix .field.type "[]") false) }}


// Set{{.field.name}} set {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Set{{.field.name}}(v {{.field.type}}) {
    {{ .structname | firstLower }}.{{.field.name}} = v
}
{{ end }} 

{{- if (hasPrefix .field.type "[]") }}


{{/* Поддержка только элементарных типов */}}
{{if (intersection (substring .field.type 2) "string" "interface{}")}}
// Set{{.field.name}} set all elements {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Set{{.field.name}}(v {{.field.type}}) {
   
    for _, value := range v {
        {{ .structname | firstLower }}.Add{{.field.name}}(value)
    }
}

// Add{{.field.name}} add element {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Add{{.field.name}}(v {{substring .field.type 2}}) {
    if {{ .structname | firstLower }}.Include{{.field.name}}(v) {
        return
    }
    
    {{ .structname | firstLower }}.{{.field.name}} = append({{ .structname | firstLower }}.{{.field.name}}, v)
}

// Remove{{.field.name}} remove element {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Remove{{.field.name}}(v {{substring .field.type 2}}) {
    if !{{ .structname | firstLower }}.Include{{.field.name}}(v) {
        return
    }
    
    _i := {{ .structname | firstLower }}.Index{{.field.name}}(v)
    
    {{ .structname | firstLower }}.{{.field.name}} = append({{ .structname | firstLower }}.{{.field.name}}[:_i], {{ .structname | firstLower }}.{{.field.name}}[_i+1:]...)
}
{{ end }}

{{ end }}

{{- if (hasPrefix .field.type "map") }}
{{ $regexp := "map\\[(?P<key>[a-zA-Z0-9{}]+)\\](?P<item>[a-zA-Z0-9{}]+)" }}
{{ $keyType := (index (regexp .field.type $regexp) 1)}}
{{ $valueType := (index (regexp .field.type $regexp) 2)}}

// Set{{.field.name}} set all elements {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Set{{.field.name}}(v {{.field.type}}) {
    {{ .structname | firstLower }}.{{.field.name}} = make({{.field.type}})
    
    for key, value := range v {
        {{ .structname | firstLower }}.{{.field.name}}[key] = value
    }
}

// Add{{.field.name}} add element by key
func ({{ .structname | firstLower }} *{{.structname}}) SetOne{{.field.name}}(k {{ $keyType}}, v {{ $valueType }}) {
    {{ .structname | firstLower }}.{{.field.name}}[k] = v
}

// Remove{{.field.name}} remove element by key
func ({{ .structname | firstLower }} *{{.structname}}) Remove{{.field.name}}(k {{ $keyType}}) {
    if _, exist := {{ .structname | firstLower }}.{{.field.name}}[k]; exist {
        delete({{ .structname | firstLower }}.{{.field.name}}, k)  
    } 
}
{{ end }}

{{ end }}

// Вспомогательные функции для структуры (в зависимости от типа поля)
{{define "getter"}}

// Get{{.field.name}} get {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Get{{.field.name}}() {{.field.type}} {
    return {{ .structname | firstLower }}.{{.field.name}}
}

{{- if (hasPrefix .field.type "[]") }}

{{/* Поддержка только элементарных типов */}}
{{if (intersection (substring .field.type 2) "string" "interface{}")}}
// Index{{.field.name}} get index element {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Index{{.field.name}}(v {{substring .field.type 2}}) int {
    for _index, _v := range {{ .structname | firstLower }}.{{.field.name}} {
        if _v == v {
            return _index
        }
    }
    return -1
}

// Include{{.field.name}} has exist value {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) Include{{.field.name}}(v {{substring .field.type 2}}) bool {
    return {{ .structname | firstLower }}.Index{{.field.name}}(v) > -1
}
{{ end}}

{{ end }}

{{- if (hasPrefix .field.type "map") }}
{{ $regexp := "map\\[(?P<key>[a-zA-Z0-9{}]+)\\](?P<item>[a-zA-Z0-9{}]+)" }}
{{ $keyType := (index (regexp .field.type $regexp) 1)}}
{{ $valueType := (index (regexp .field.type $regexp) 2)}}

// Exist{{.field.name}} has exist key {{.field.name}}
func ({{ .structname | firstLower }} *{{.structname}}) ExistKey{{.field.name}}(k {{ $keyType}}) bool {
     _, exist := {{ .structname | firstLower }}.{{.field.name}}[k]
     
     return exist
}

func ({{ .structname | firstLower }} *{{.structname}}) GetOne{{.field.name}}(k {{ $keyType}}) {{ $valueType }} {
    return {{ .structname | firstLower }}.{{.field.name}}[k]
}
{{ end }}

{{ end }}

{{ define "transform" }}
{{ $self := (firstLower .structname )}}
func ({{ $self }} *{{.structname}}) TransformFrom (v interface{}) error {
    switch v.(type) {
    {{ range $index, $transform := .transform}}
    case *{{$transform.type}}:
        d := v.(*{{$transform.type}})
        {{ range $iindex, $map := $transform.map}}
        {{$self}}.{{$map.to}} = d.{{$map.from}}{{ end }}
        {{with $transform.custom}}{{$transform.custom}}{{end}}
        
    {{ end }}
	default:
		glog.Errorf("Not supported type %v", v)
		return fmt.Errorf("not_supported")
	}
    
	return nil
}
{{ end }}