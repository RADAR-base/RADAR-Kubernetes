package wizard

import "strings"

// uncommentFirebaseBlock removes the surrounding YAML/Go-template comment wrapping
// from the radar_appserver firebase-adminsdk block in etc/production.yaml.gotmpl.
//
// Upstream ships it as:
//
//	#{{/*
//	#radar_appserver:
//	#  google_application_credentials: {{ readFile "../etc/radar-appserver/firebase-adminsdk.json" | quote }}
//	#*/}}
//
// We strip the `#{{/*` and `#*/}}` lines and the leading `#` on the two YAML lines
// between them. If the block isn't present (or already uncommented) the input is
// returned unchanged.
func uncommentFirebaseBlock(src string) string {
	lines := strings.Split(src, "\n")
	startIdx, endIdx := -1, -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if startIdx == -1 && trimmed == "#{{/*" {
			// Look ahead to confirm this is the firebase block (not some other commented gotmpl).
			for j := i + 1; j < len(lines) && j < i+5; j++ {
				if strings.Contains(lines[j], "google_application_credentials") {
					startIdx = i
					break
				}
				if strings.TrimSpace(lines[j]) == "#*/}}" {
					break
				}
			}
		} else if startIdx != -1 && strings.TrimSpace(line) == "#*/}}" {
			endIdx = i
			break
		}
	}
	if startIdx == -1 || endIdx == -1 {
		return src
	}

	out := make([]string, 0, len(lines))
	out = append(out, lines[:startIdx]...)
	for i := startIdx + 1; i < endIdx; i++ {
		out = append(out, strings.TrimPrefix(lines[i], "#"))
	}
	out = append(out, lines[endIdx+1:]...)
	return strings.Join(out, "\n")
}
