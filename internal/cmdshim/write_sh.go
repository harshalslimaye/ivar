package cmdshim

import "fmt"

var shHeader string = `#!/bin/sh
basedir=$(dirname "$(echo "$0" | sed -e 's,\\,/,g')")

case $(uname) in
    *CYGWIN*|*MINGW*|*MSYS*)
        if command -v cygpath > /dev/null 2>&1; then
            basedir=$(cygpath -w "$basedir")
        fi
    ;;
esac

`

var shContent string = `%sif [ -x %s ]; then
exec %s %s %s %s "$@"
else
exec %s %s %s "$@"
fi
`

func WriteSh(params *Params) string {

	var sh string
	if params.ShLongProg != "" {
		sh = fmt.Sprintf(shContent,
			shHeader, params.ShLongProg, params.Variables, params.ShLongProg,
			params.Args, params.ShTarget, params.Variables, params.ShProg,
			params.Args, params.ShTarget)
	} else {
		sh = fmt.Sprintf(`%sexec %s %s %s "$@"
`,
			shHeader, params.ShProg, params.Args, params.ShTarget)
	}

	return sh
}
