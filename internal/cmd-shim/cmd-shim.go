package cmdShim

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/harshalslimaye/ivar/internal/logger"
)

var PwshHeader string = `#!/usr/bin/env pwsh
$basedir=Split-Path $MyInvocation.MyCommand.Definition -Parent

$exe=""
if ($PSVersionTable.PSVersion -lt "6.0" -or $IsWindows) {
  # Fix case when both the Windows and Linux builds of Node
  # are installed in the same directory
  $exe=".exe"
}
`
var PwshWithShLong string = `$ret=0
if (Test-Path %s) {
  # Support pipeline input
  if ($MyInvocation.ExpectingInput) {
    $input | & %s %s %s $args
  } else {
    & %s %s %s $args
  }
  $ret=$LASTEXITCODE
} else {
  # Support pipeline input
  if ($MyInvocation.ExpectingInput) {
    $input | & %s %s %s $args
  } else {
    & %s %s %s $args
  }
  $ret=$LASTEXITCODE
}
exit $ret
`

var PwshWithoutLong string = `# Support pipeline input
if ($MyInvocation.ExpectingInput) {
  $input | & %s %s %s $args
} else {
  & %s %s %s $args
}
exit $LASTEXITCODE
`

type Params struct {
	ShTarget     string
	Target       string
	LongProg     string
	ShProg       string
	ShLongProg   string
	PwshProg     string
	PwshLongProg string
	Prog         string
	Args         string
	Variables    string
	From         string
	To           string
}

func NewParams(from, to, prog, args, variables string) *Params {
	params := Params{
		ShTarget:     "",
		Target:       "",
		LongProg:     "",
		ShProg:       "",
		ShLongProg:   "",
		PwshProg:     "",
		PwshLongProg: "",
		Prog:         prog,
		Args:         args,
		Variables:    variables,
		From:         from,
		To:           to,
	}

	shTarget, err := filepath.Rel(filepath.Dir(to), from)
	if err != nil {
		fmt.Println("Error:", err)
	}

	params.Target = strings.ReplaceAll(shTarget, "/", "\\")
	params.ShTarget = strings.ReplaceAll(shTarget, "\\", "/")

	if params.Prog == "" {
		params.Prog = `"%dp0%\\${target}"`
		params.ShProg = "\"$basedir/" + shTarget + "\""
		params.PwshProg = params.ShProg
		params.Args = ""
		params.Target = ""
		params.ShTarget = ""
	} else {
		params.LongProg = "%dp0%\\" + params.Prog + ".exe"
		params.ShLongProg = "\"$basedir/" + params.Prog + "\""
		params.PwshLongProg = "\"" + "$basedir/" + params.Prog + "$exe" + "\""
		params.Target = "\"" + "%dp0%\\" + params.Target + "\""
		params.ShProg = strings.ReplaceAll(params.Prog, "\\", "/")
		params.PwshProg = "\"" + params.ShProg + "$exe\""
		params.ShTarget = "\"$basedir/" + params.ShTarget + "\""
	}

	return &params
}

var shebangExpr *regexp.Regexp = regexp.MustCompile(`^#!\s*(?:\/usr\/bin\/env\s+(?:-S\s+)?((?:[^ \t=]+=[^ \t=]+\s+)*))?([^ \t]+)(.*)$`)

func CmdShim(from, to string) {
	rm(to)
	rm(to + ".cmd")
	rm(to + ".ps1")

	writeShim(from, to)
}

func rm(to string) {
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		if err := os.Remove(to); err != nil {
			logger.Error("error while deleting existing links: ", err)
		}
	}
}

func writeShim(from, to string) {
	binPath := filepath.Join("node_modules", ".bin")
	if _, binExists := os.Stat(binPath); os.IsNotExist(binExists) {
		if err := os.MkdirAll(binPath, 0755); err != nil {
			logger.Error("error while creating .bin folder", err)
		}
	}

	file, err := os.Open(from)
	if err != nil {
		logger.Error(fmt.Sprintf("error while reading file %s: ", from), err)
	}

	scanner := bufio.NewScanner(file)
	var firstLine string

	if scanner.Scan() {
		firstLine = scanner.Text()
	} else {
		logger.Error("empty executable file: ", errors.New(fmt.Sprintf("File %s is empty", to)))
	}

	shebang := shebangExpr.FindStringSubmatch(firstLine)

	if shebang == nil {
		writeShim_(from, to, "", "", "")
		return
	}

	vars := shebang[1]
	prog := shebang[2]
	args := shebang[3]

	writeShim_(from, to, prog, vars, args)
}

func writeShim_(from, to, prog, args, variables string) {
	params := NewParams(from, to, prog, args, variables)
	writeCmd(params)
	writeSh(params)
	writePwsh(params)
}

func writePwsh(params *Params) {
	fileContent := PwshHeader

	if params.ShLongProg != "" {
		fileContent += fmt.Sprintf(
			PwshWithShLong,
			params.PwshLongProg,
			params.PwshLongProg,
			params.Args,
			params.ShTarget,
			params.PwshLongProg,
			params.Args,
			params.ShTarget,
			params.PwshProg,
			params.Args,
			params.ShTarget,
			params.PwshProg,
			params.Args,
			params.ShTarget,
		)
	} else {
		fileContent += fmt.Sprintf(
			PwshWithoutLong,
			params.PwshProg,
			params.Args,
			params.ShTarget,
			params.PwshProg,
			params.Args,
			params.ShTarget,
		)
	}

	CreateLink(params.To+".ps1", fileContent)
}

func writeSh(params *Params) {
	sh := "#!/bin/sh\n"

	sh = sh + "basedir=$(dirname \"$(echo \"$0\" | sed -e 's,\\\\,/,g')\")\n"
	sh = sh + "\n"
	sh = sh + "case `uname` in\n"
	sh = sh + "    *CYGWIN*|*MINGW*|*MSYS*)\n"
	sh = sh + "        if command -v cygpath > /dev/null 2>&1; then\n"
	sh = sh + "            basedir=`cygpath -w \"$basedir\"`\n"
	sh = sh + "        fi\n"
	sh = sh + "    ;;\n"
	sh = sh + "esac\n"
	sh = sh + "\n"

	if params.ShLongProg != "" {
		sh = sh + "if [ -x " + params.ShLongProg + " ]; then\n"
		sh = sh + "  exec " + params.Variables + params.ShLongProg + " " + params.Args + " " + params.ShTarget + " \"$@\"\n"
		sh = sh + "else \n"
		sh = sh + "  exec " + params.Variables + params.ShProg + " " + params.Args + " " + params.ShTarget + " \"$@\"\n"
		sh = sh + "fi\n"
	} else {
		sh = sh + "exec " + params.ShProg + " " + params.Args + " " + params.ShTarget + " \"$@\"\n"
	}
	CreateLink(params.To, sh)
}

func writeCmd(params *Params) {
	progName := regexp.MustCompile(`(^")|("$)`).ReplaceAllString(params.Prog, "")

	head := "@ECHO off\r\n" +
		"GOTO start\r\n" +
		":find_dp0\r\n" +
		"SET dp0=%~dp0\r\n" +
		"EXIT /b\r\n" +
		":start\r\n" +
		"SETLOCAL\r\n" +
		"CALL :find_dp0\r\n"

	var cmd string
	if params.LongProg != "" {
		args := strings.TrimSpace(params.Args)
		variablesBatch := ConvertToSetCommands(params.Variables)

		cmd = head + variablesBatch
		cmd = cmd + "\r\n"
		cmd = cmd + "IF EXIST " + "\"" + params.LongProg + "\"" + " (\r\n"
		cmd = cmd + "\tSET \"_prog=" + params.LongProg + "\"\r\n"
		cmd = cmd + ") ELSE (\r\n"
		cmd = cmd + "\tSET \"_prog=" + progName + "\"\r\n"
		cmd = cmd + "\tSET PATHEXT=%PATHEXT:;.JS;=;%\r\n"
		cmd = cmd + ")\r\n"
		cmd = cmd + "\r\n"
		cmd = cmd + "endLocal & goto #_undefined_# 2>NUL || title %COMSPEC% & "
		cmd = cmd + "\"%_prog%\" " + args + " " + params.Target + " " + "%*\r\n"
	} else {
		cmd = head + params.Prog + " " + params.Args + " " + params.Target + " %*\r\n"
	}
	CreateLink(params.To+".cmd", cmd)
}

func CreateLink(to string, content string) {
	file, err := os.Create(to)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	file.Write([]byte(content))
}

func ConvertToSetCommand(key, value string) string {
	var line string
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	if key != "" && value != "" && len(value) > 0 {
		line = "@SET " + key + "=" + ReplaceDollarWithPercentPair(value) + "\r\n"
	}
	return line
}

func ExtractVariableValuePairs(declarations []string) map[string]string {
	pairs := make(map[string]string)
	for _, declaration := range declarations {
		split := strings.Split(declaration, "=")
		if len(split) == 2 {
			pairs[strings.TrimSpace(split[0])] = strings.TrimSpace(split[1])
		}
	}
	return pairs
}

func ConvertToSetCommands(variableString string) string {
	variableValuePairs := ExtractVariableValuePairs(strings.Split(variableString, " "))
	var variableDeclarationsAsBatch string
	for key, value := range variableValuePairs {
		variableDeclarationsAsBatch += ConvertToSetCommand(key, value)
	}
	return variableDeclarationsAsBatch
}

func ReplaceDollarWithPercentPair(value string) string {
	dollarExpressions := regexp.MustCompile(`\$\{?([^$@#?\- \t{}:]+)\}?`)
	result := ""
	startIndex := 0
	for _, match := range dollarExpressions.FindAllStringSubmatchIndex(value, -1) {
		betweenMatches := value[startIndex:match[0]]
		result += betweenMatches + "%" + value[match[2]:match[3]] + "%"
		startIndex = match[1]
	}
	result += value[startIndex:]
	return result
}
