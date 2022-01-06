package xxl

import "strings"

const (
	DefaultExecutorPort = "9999"
	DefaultRegistryKey  = "golang-jobs"
)

const (
	GlueTypeBEAN       = "BEAN"
	GlueTypeGROOVY     = "GLUE(Java)"
	GlueTypeSHELL      = "GLUE(Shell)"
	GlueTypePYTHON     = "GLUE(Python)"
	GlueTypePHP        = "GLUE(PHP)"
	GlueTypeNODEJS     = "GLUE(Nodejs)"
	GlueTypePOWERSHELL = "GLUE(PowerShell)"
)

var script = []string{
	GlueTypeSHELL,
	GlueTypePYTHON,
	GlueTypePHP,
	GlueTypeNODEJS,
	GlueTypePOWERSHELL,
}

func IsScript(glueType string) bool {
	for _, s := range script {
		if strings.EqualFold(s, glueType) {
			return true
		}
	}
	return false
}

const (
	SerialExecution = "SERIAL_EXECUTION"
	DiscardLater    = "DISCARD_LATER"
	CoverEarly      = "COVER_EARLY"
)
