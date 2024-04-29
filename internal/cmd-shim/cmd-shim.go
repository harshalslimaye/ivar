package cmdShim

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var shebangExpr *regexp.Regexp = regexp.MustCompile(`^#!\s*(?:\/usr\/bin\/env\s+(?:-S\s+)?((?:[^ \t=]+=[^ \t=]+\s+)*))?([^ \t]+)(.*)$`)

func CmdShim(from, to string) {
	// rm(to)
	// rm(to + ".cmd")
	// rm(to + ".ps1")

	writeShim(from, to)
}

func rm(to string) {
	if err := os.Remove(to); err != nil {
		fmt.Println("Error while deleting existing links: ", err)
	}
}

func writeShim(from, to string) {
	if err := os.MkdirAll(filepath.Join("node_modules", ".bin"), 0755); err != nil {
		fmt.Println("Error while creating .bin folder")
	}

	file, err := os.Open(from)
	if err != nil {
		fmt.Println("Error while reading file: ", err)
	}

	scanner := bufio.NewScanner(file)
	var firstLine string

	if scanner.Scan() {
		firstLine = scanner.Text()
	} else {
		fmt.Print("Error file is empty " + to)
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
	fmt.Println(prog)
	writeCmd(from, to, prog, args, variables)
	writeSh(from, to, prog, args, variables)
	writePwSh(from, to, prog, args, variables)
}

func writePwSh(from, to, prog, args, variables string) {
	shTarget, err := filepath.Rel(filepath.Dir(to), from)
	if err != nil {
		// Handle error if any
		fmt.Println("Error:", err)
		return
	}
	// longProg := ""
	target := strings.ReplaceAll(shTarget, "/", "\\")
	pwshLongProg := ""
	shLongProg := ""
	pwshProg := ""
	shProg := ""

	if prog == "" {
		prog = `"%dp0%\\${target}"`
		args = ""
		target = ""
		shProg = "\"$basedir/" + shTarget + "\""
		pwshProg = shProg
	} else {
		pwshLongProg = "\"" + "$basedir/" + prog + "$exe" + "\""
		target = "\"" + "%dp0%\\" + target + "\""
		shLongProg = "\"$basedir/" + prog + "\""
		shProg = strings.ReplaceAll(prog, "\\", "/")
		shTarget = "\"$basedir/" + shTarget + "\""
	}

	pwsh := "#!/usr/bin/env pwsh\n"
	pwsh = pwsh + "$basedir=Split-Path $MyInvocation.MyCommand.Definition -Parent\n"
	pwsh = pwsh + "\n"
	pwsh = pwsh + "$exe=\"\"\n"
	pwsh = pwsh + "if ($PSVersionTable.PSVersion -lt \"6.0\" -or $IsWindows) {\n"
	pwsh = pwsh + "  # Fix case when both the Windows and Linux builds of Node\n"
	pwsh = pwsh + "  # are installed in the same directory\n"
	pwsh = pwsh + "  $exe=\".exe\"\n"
	pwsh = pwsh + "}\n"

	if shLongProg != "" {
		pwsh = pwsh + "$ret=0\n"
		pwsh = pwsh + "if (Test-Path " + pwshLongProg + ") {\n"
		pwsh = pwsh + "  # Support pipeline input\n"
		pwsh = pwsh + "  if ($MyInvocation.ExpectingInput) {\n"
		pwsh = pwsh + "    $input | & " + pwshLongProg + " " + args + " " + shTarget + " $args\n"
		pwsh = pwsh + "  } else {\n"
		pwsh = pwsh + "    & " + pwshLongProg + " " + args + " " + shTarget + " $args\n"
		pwsh = pwsh + "  }\n"
		pwsh = pwsh + "  $ret=$LASTEXITCODE\n"
		pwsh = pwsh + "} else {\n"
		pwsh = pwsh + "  # Support pipeline input\n"
		pwsh = pwsh + "  if ($MyInvocation.ExpectingInput) {\n"
		pwsh = pwsh + "    $input | & " + pwshProg + " " + args + " " + shTarget + " $args\n"
		pwsh = pwsh + "  } else {\n"
		pwsh = pwsh + "    & " + pwshProg + " " + args + " " + shTarget + " $args\n"
		pwsh = pwsh + "  }\n"
		pwsh = pwsh + "  $ret=$LASTEXITCODE\n"
		pwsh = pwsh + "}\n"
		pwsh = pwsh + "exit $ret\n"
	} else {
		pwsh = pwsh + "# Support pipeline input\n"
		pwsh = pwsh + "if ($MyInvocation.ExpectingInput) {\n"
		pwsh = pwsh + "  $input | & " + pwshProg + " " + args + " " + shTarget + " $args\n"
		pwsh = pwsh + "} else {\n"
		pwsh = pwsh + "  & " + pwshProg + " " + args + " " + shTarget + " $args\n"
		pwsh = pwsh + "}\n"
		pwsh = pwsh + "exit $LASTEXITCODE\n"
	}
	CreateLink(to+".ps1", pwsh)
}

func writeSh(from, to, prog, args, variables string) {
	shTarget, err := filepath.Rel(filepath.Dir(to), from)
	if err != nil {
		// Handle error if any
		fmt.Println("Error:", err)
		return
	}
	shLongProg := ""
	shProg := ""

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if prog == "" {
		shTarget = ""
		shProg = "\"$basedir/" + shTarget + "\""
	} else {
		shLongProg = "\"$basedir/" + prog + "\""
		shTarget = "\"$basedir/" + shTarget + "\""
		shProg = strings.ReplaceAll(prog, "\\", "/")
	}

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

	if shLongProg != "" {
		sh = sh + "if [ -x " + shLongProg + " ]; then\n"
		sh = sh + "  exec " + variables + shLongProg + " " + args + " " + shTarget + "\"$@\"\n"
		sh = sh + "else \n"
		sh = sh + "  exec " + variables + shProg + " " + args + " " + shTarget + "\"$@\"\n"
		sh = sh + "fi\n"
	} else {
		sh = sh + "exec " + shProg + " " + args + " " + shTarget + " \"$@\"\n"
	}
	CreateLink(to, sh)
}

func writeCmd(from, to, prog, args, variables string) {
	shTarget, err := filepath.Rel(filepath.Dir(to), from)
	if err != nil {
		// Handle error if any
		fmt.Println("Error:", err)
		return
	}
	longProg := ""
	target := strings.ReplaceAll(shTarget, "/", "\\")
	progName := regexp.MustCompile(`(^")|("$)`).ReplaceAllString(prog, "")

	if prog == "" {
		prog = `"%dp0%\\${target}"`
		args = ""
		target = ""
	} else {
		longProg = "%dp0%\\" + prog + ".exe"
		target = "\"" + "%dp0%\\" + target + "\""
	}

	head := "@ECHO off\r\n" +
		"GOTO start\r\n" +
		":find_dp0\r\n" +
		"SET dp0=%~dp0\r\n" +
		"EXIT /b\r\n" +
		":start\r\n" +
		"SETLOCAL\r\n" +
		"CALL :find_dp0\r\n"

	var cmd string
	if longProg != "" {
		args = strings.TrimSpace(args)
		variablesBatch := ConvertToSetCommands(variables)

		cmd = head + variablesBatch
		cmd = cmd + "\r\n"
		cmd = cmd + "IF EXIST " + "\"" + longProg + "\"" + " (\r\n"
		cmd = cmd + "\tSET \"_prog=" + longProg + "\"\r\n"
		cmd = cmd + ") ELSE (\r\n"
		cmd = cmd + "\tSET \"_prog=" + progName + "\"\r\n"
		cmd = cmd + "\tSET PATHEXT=%PATHEXT:;.JS;=;%\r\n"
		cmd = cmd + ")\r\n"
		cmd = cmd + "\r\n"
		cmd = cmd + "endLocal & goto #_undefined_# 2>NUL || title %COMSPEC% & "
		cmd = cmd + "\"%_prog%\" " + args + " " + target + " " + "%*\r\n"
	} else {
		cmd = head + prog + " " + args + " " + target + " %*\r\n"
	}
	CreateLink(to+".cmd", cmd)
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
