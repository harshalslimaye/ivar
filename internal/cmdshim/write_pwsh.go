package cmdshim

import "fmt"

var pwshHeader string = `#!/usr/bin/env pwsh
$basedir=Split-Path $MyInvocation.MyCommand.Definition -Parent

$exe=""
if ($PSVersionTable.PSVersion -lt "6.0" -or $IsWindows) {
  # Fix case when both the Windows and Linux builds of Node
  # are installed in the same directory
  $exe=".exe"
}
`
var pwshWithShLong string = `$ret=0
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

var pwshWithoutLong string = `# Support pipeline input
if ($MyInvocation.ExpectingInput) {
  $input | & %s %s %s $args
} else {
  & %s %s %s $args
}
exit $LASTEXITCODE
`

func WritePwsh(params *Params) string {
	fileContent := pwshHeader

	if params.ShLongProg != "" {
		fileContent += fmt.Sprintf(
			pwshWithShLong,
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
			pwshWithoutLong,
			params.PwshProg,
			params.Args,
			params.ShTarget,
			params.PwshProg,
			params.Args,
			params.ShTarget,
		)
	}

	return fileContent
}
