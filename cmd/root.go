package cmd

import (
	"github.com/CiscoAMPBeats/beater"
	"github.com/elastic/beats/libbeat/cmd"
)

// Name of this beat
var Name = "CiscoAMPBeats"

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmd(Name, "", beater.New)
