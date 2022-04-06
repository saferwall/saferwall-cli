# BlackCat

* First seen: {{ .Fam.FirstSeen }}
* Aliases: {{ range .Fam.Aliases }}{{ . }}, {{ end }}


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

| Vendors     | status | Vendors | status |
| ----------- | ------ | ------- | ------ |
| Avast       | {{ .File.MultiAV.last_scan.avast.output }} | FSecure | {{ .File.MultiAV.last_scan.fsecure.output }} |
| Avira       | {{ .File.MultiAV.last_scan.avira.output }} | Kaspersky | {{ .File.MultiAV.last_scan.kaspersky.output }} |
| Bitdefender | {{ .File.MultiAV.last_scan.bitdefender.output }} | McAfee | {{ .File.MultiAV.last_scan.mcafee.output }} |
| ClamAV      | {{ .File.MultiAV.last_scan.clamav.output }} | Sophos | {{ .File.MultiAV.last_scan.sophos.output }} |
| Comodo      | {{ .File.MultiAV.last_scan.comodo.output }} | Symantec | {{ .File.MultiAV.last_scan.symantec.output }} |
| ESET        | {{ .File.MultiAV.last_scan.eset.output }} | Windows Defender | {{ .File.MultiAV.windefender.avast.output }} |
| TrendMicro  | {{ .File.MultiAV.last_scan.trendmicro.output }} | DrWeb | {{ .File.MultiAV.last_scan.drweb.output }} |

