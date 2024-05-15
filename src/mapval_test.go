package main

import "testing"

func TestGetMapVal(t *testing.T) {
	ep := newTestEndpoint()

	content := readFileContent(
		fromTestFolder("dump/markdown/1.md"), ep,
	)
	validateGetMapVal(
		"front_matter.title", content, []string{"title1"}, t,
	)
	validateGetMapVal(
		"front_matter.tags", content, []string{"tag1", "tag2"}, t,
	)

	content = readFileContent(
		fromTestFolder("dump/yaml/cpx/data_aip.yaml"), ep,
	)
	validateGetMapVal(
		"title", content, []string{"Data Services @ AIP"}, t,
	)
	validateGetMapVal(
		"metadata.access", content, []string{"open"}, t,
	)
	validateGetMapVal(
		"metadata.tags", content, []string{"vo", "IVOA", "Daiquiri"}, t,
	)
	validateGetMapVal(
		"metadata.url", content, []string{"https://data.aip.de"}, t,
	)

}

func validateGetMapVal(key string, mp map[string]interface{}, exp []string, t *testing.T) {
	b := false
	res := getMapVal(key, mp)
	if len(exp) == len(res) {
		for i, x := range res {
			if x != exp[i] {
				b = true
			}
		}
	} else {
		b = true
	}
	if b {
		t.Errorf(
			"error get map val,\nkey: %+v,\nmap: %+v,\nexp: %v,\nres: %v\n\n",
			key, mp, exp, res,
		)
	}
}
