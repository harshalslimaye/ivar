package registry

import (
	"fmt"
	"io"
	"net/http"

	"github.com/harshalslimaye/ivar/internal/jsonparser"
	"github.com/logrusorgru/aurora"
	"github.com/valyala/fastjson"
)

func FetchDependencies(name, version string) (*jsonparser.JsonParser, error) {
	endpoint := fmt.Sprintf("https://registry.npmjs.org/%s/%s", name, version)

	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch dependency %s@%s: %s", aurora.Yellow(name), aurora.Yellow(version), aurora.Yellow(res.Status))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var p fastjson.Parser
	val, err := p.ParseBytes(body)
	if err != nil {
		return nil, err
	}

	parser := jsonparser.NewJsonParser(val)

	return parser, nil
}
