package main

import (
	pwl "github.com/justjanne/powerline-go/powerline"
	"fmt"
	"os"
	"strings"
)

func segmentPaasCF(p *powerline) []pwl.Segment {
	
	bblStateDir, _ := os.LookupEnv("BBL_STATE_DIR")
	if bblStateDir == "" {
		return []pwl.Segment{}
	}
	idx := strings.LastIndexByte(bblStateDir, '/');
	instance := bblStateDir;
	if ( idx >= 0 ) {
		instance = bblStateDir[idx+1:];
	}

	return []pwl.Segment{{
		Name:       "paascf",
		Content:    fmt.Sprintf("🎡 %s", instance),
		Foreground: p.theme.PaasCfFg,
		Background: p.theme.PaasCfBg,
	}}
}
