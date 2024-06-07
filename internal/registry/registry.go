package registry

import (
	"fmt"
	"io"
	"net/http"

	"github.com/harshalslimaye/ivar/internal/jsonparser"
	cache "github.com/harshalslimaye/ivar/internal/store"
	"github.com/logrusorgru/aurora"
	"github.com/valyala/fastjson"
)

type Registry struct {
	store *cache.Store
}

func FetchDependencies(name, version string) ([]byte, error) {
	endpoint := fmt.Sprintf("https://registry.npmjs.org/%s/%s", name, version)
	fmt.Println(endpoint)

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

	return body, nil
}

func (r *Registry) Fetch(name, version string) (*jsonparser.JsonParser, error) {
	value := r.store.Get(name)

	data, okay := value.([]byte)

	if okay && data != nil {
		return getParser(data)
	}

	response, err := FetchDependencies(name, version)
	if err != nil {
		return nil, err
	}
	r.store.Set(name, response)

	return getParser(response)
}

func getParser(response []byte) (*jsonparser.JsonParser, error) {
	var p fastjson.Parser
	value, err := p.ParseBytes(response)
	if err != nil {
		return nil, err
	}

	return jsonparser.NewJsonParser(value), nil
}

func NewRegistry() *Registry {
	return &Registry{
		store: cache.GetStore(),
	}
}
