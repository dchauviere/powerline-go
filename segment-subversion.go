package main

import (
	"errors"
	"fmt"
	pwl "github.com/justjanne/powerline-go/powerline"
	"os/exec"
	"strings"
)

var otherModified int

func (r repoStats) SvnStatusColors(p *powerline) (uint8, uint8) {
	if r.conflicted > 0 {
		return p.theme.RepoConflictFg, p.theme.RepoConflictBg
	}

	if (r.staged + r.notStaged > 0) || otherModified > 0 {
		return p.theme.RepoDirtyFg, p.theme.RepoDirtyBg
	}

	return p.theme.RepoCleanFg, p.theme.RepoCleanBg
}

func svnAddRepoStats(nChanges int, symbol string) string {
	if nChanges > 0 {
		return fmt.Sprintf(" %d%s", nChanges, symbol)
	}
	return ""
}

func (r repoStats) SvnStats(p *powerline) string {
	stats := svnAddRepoStats(r.ahead, p.symbolTemplates.RepoAhead)
	stats = stats + svnAddRepoStats(r.behind, p.symbolTemplates.RepoBehind)
	stats = stats + svnAddRepoStats(r.staged, p.symbolTemplates.RepoStaged)
	stats = stats + svnAddRepoStats(r.notStaged, p.symbolTemplates.RepoNotStaged)
	stats = stats + svnAddRepoStats(r.untracked, p.symbolTemplates.RepoUntracked)
	stats = stats + svnAddRepoStats(r.conflicted, p.symbolTemplates.RepoConflicted)
	stats = stats + svnAddRepoStats(r.stashed, p.symbolTemplates.RepoStashed)
	return stats
}

func runSvnCommand(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	out, err := command.Output()
	return string(out), err
}

func parseSvnURL() (map[string]string, error) {
	info, err := runSvnCommand("svn", "info")
	if err != nil {
		return nil, errors.New("not a working copy")
	}

	svnInfo := make(map[string]string, 0)
	infos := strings.Split(info, "\n")
	if len(infos) > 1 {
		for _, line := range infos[:] {
			items := strings.Split(line, ": ")
			if len(items) >= 2 {
				svnInfo[items[0]] = items[1]
			}
		}
	}

	return svnInfo, nil
}

func ensureUnmodified(code string, stats repoStats) {
	if code != " " {
		otherModified++
	}
}

func parseSvnStatus() repoStats {
	stats := repoStats{}
	info, err := runSvnCommand("svn", "status", "-u")
	if err != nil {
		return stats
	}
	infos := strings.Split(info, "\n")
	if len(infos) > 1 {
		for _, line := range infos[:] {
			if len(line) >= 9 {
				code := line[0:1]
				switch code {
				case "?":
					stats.untracked++
				case "C":
					stats.conflicted++
				case "A", "D", "M":
					stats.notStaged++
				default:
					ensureUnmodified(code, stats)
				}
				code = line[1:2]
				switch code {
				case "C":
					stats.conflicted++
				case "M":
					stats.notStaged++
				default:
					ensureUnmodified(code, stats)
				}
				ensureUnmodified(line[2:3], stats)
				ensureUnmodified(line[3:4], stats)
				ensureUnmodified(line[4:5], stats)
				ensureUnmodified(line[5:6], stats)
				ensureUnmodified(line[6:7], stats)
				ensureUnmodified(line[7:8], stats)
				code = line[8:9]
				switch code {
				case "*":
					stats.behind++
				default:
					ensureUnmodified(code, stats)
				}
			}
		}
	}

	return stats
}

func segmentSubversion(p *powerline) []pwl.Segment {

	svnInfo, err := parseSvnURL()
	if err != nil {
		return []pwl.Segment{}
	}

	if len(p.ignoreRepos) > 0 {
		if p.ignoreRepos[svnInfo["URL"]] || p.ignoreRepos[svnInfo["Relative URL"]] {
			return []pwl.Segment{}
		}
	}

	svnStats := parseSvnStatus()

	var foreground, background uint8
	foreground, background = svnStats.SvnStatusColors(p)

	segments := []pwl.Segment{{
		Name:       "svn-branch",
		Content:    fmt.Sprintf("%s %s", svnInfo["Relative URL"], svnStats.SvnStats(p)),
		Foreground: foreground,
		Background: background,
	}}

	return segments
}
