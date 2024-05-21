package cmdshim

import (
	"fmt"
	"path/filepath"
	"strings"
)

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
