# BlackCat

* First seen: {{ .Fam.FirstSeen }}
* Aliases: {{$aliases := .Fam.Aliases}}
{{- range .Fam.Aliases }}
{{ . }}
{{- if ne . ($aliases | last) }}, {{- end }}
{{- end }}


## Basic Properties

| Property | Value |
| --- | --- |
| Size | {{ .File.Size }} bytes |
| CRC32 | {{ .File.Crc32 }} |
| MD5 | {{ .File.MD5 }} |
| SHA1 | {{ .File.SHA1 }} |
| SHA256 | {{ .File.SHA256 }} |
| SHA512 | {{ .File.SHA512 }} |
| Ssdeep | {{ .File.Ssdeep }} |
| Magic | {{ .File.Magic }}  |
| Packer | {{ range .File.Packer }}{{ . }}<br />{{ end }} |
| TrID | {{ range .File.TriD }}{{ . }}<br />{{ end }} |

## Antivirus Scan

```diff
{{- range $k, $v := .File.MultiAV.last_scan }}
{{- if $v.output }}
- {{ $k | title }}: {{ $v.output }}
{{- else }}
+ {{ $k | title }}: clean
{{- end }}
{{- end }}
```

## References

{{- range .Fam.References }}
- [{{ . }}]({{ . }})
{{- end }}
