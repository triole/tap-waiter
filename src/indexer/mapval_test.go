package indexer

import (
	"testing"
)

func TestGetContentVal(t *testing.T) {
	tc := InitTests(false)
	ep := newTestEndpoint()
	specs := tc.readSpecs("specs/mapval/spec.yaml")

	for _, specItf := range specs {
		spec := specItf.(map[string]interface{})
		spec["content"] = tc.ind.readFileContent(
			spec["content_file"].(string),
			ep,
		)

		exp := tc.ind.Util.ListInterfaceToListString(
			spec["exp"].([]interface{}),
		)

		tc.validateGetContentVal(
			spec["key"].(string),
			spec["content"].(FileContent),
			exp,
		)
	}
}

func (tc testContext) validateGetContentVal(key string, mp FileContent, exp []string) {
	b := false
	res := tc.ind.TapIndex.getContentVal(key, mp)
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
		tc.t.Errorf(
			"test fail get content val,\nkey: %+v,\nmap: %+v,\nexp: %v,\nres: %v\n\n",
			key, mp, exp, res,
		)
	}
}
