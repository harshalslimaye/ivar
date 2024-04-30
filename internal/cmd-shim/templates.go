package cmdShim

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
