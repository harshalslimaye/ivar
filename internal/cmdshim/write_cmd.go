package cmdshim

import (
	"fmt"
	"regexp"
	"strings"
)

var cmdHeader string = `@ECHO off
GOTO start 
:find_dp0 
SET dp0=%~dp0 
EXIT /b 
:start 
SETLOCAL 
CALL :find_dp0
`

var longProgContent string = `%s%s
IF EXIST "%s" (
  SET "_prog=%s"
) ELSE (
  SET "_prog=%s"
  SET PATHEXT=%%PATHEXT:;.JS;=;%%
)

endLocal & goto #_undefined_# 2>NUL || title %%COMSPEC%% & "%%_prog%%" %s %s %%*
`

func WriteCmd(params *Params) string {
	progName := regexp.MustCompile(`(^")|("$)`).ReplaceAllString(params.Prog, "")

	var cmd string
	if params.LongProg != "" {
		args := strings.TrimSpace(params.Args)
		variablesBatch := ConvertToSetCommands(params.Variables)

		cmd = fmt.Sprintf(longProgContent,
			cmdHeader, variablesBatch, params.LongProg, params.LongProg,
			progName, args, params.Target)
	} else {
		cmd = fmt.Sprintf(`%s%s %s %s %%*`,
			cmdHeader, params.Prog, params.Args, params.Target)
	}

	return cmd
}
