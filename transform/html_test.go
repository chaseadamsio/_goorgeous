package transform

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/chaseadamsio/goorgeous/parse"
	"github.com/chaseadamsio/goorgeous/testdata"
	"github.com/google/go-cmp/cmp"
)

var update = flag.Bool("update", false, "update golden files")

func snapshotPath(filename string) string {
	return fmt.Sprintf("snapshots/%s.html", filename)
}

func TestTransform(t *testing.T) {
	for _, filename := range testdata.Tests {

		t.Run(filename, func(t *testing.T) {

			value := testdata.GetOrgStr(filename)
			ast := parse.Parse(value)

			snapshotPath := snapshotPath(filename)

			out := TransformToHTML(ast)
			fmt.Println(out)

			if *update {
				err := os.MkdirAll(filepath.Dir(snapshotPath), os.ModePerm)
				if err != nil {
					t.Errorf("failed to make directories for %s: %s", snapshotPath, err)
				}
				if err := ioutil.WriteFile(snapshotPath, []byte(out), os.ModePerm); err != nil {
					t.Errorf("failed to write %s file: %s", snapshotPath, err)
				}
				return
			}

			snapshot, err := ioutil.ReadFile(snapshotPath)
			if err != nil {
				t.Fatalf("failed to read %s file: %s", snapshotPath, err)
			}

			var expected map[string]interface{}
			err = json.Unmarshal([]byte(snapshot), &expected)
			if err != nil {
				t.Errorf("failed to unmarshal expected string: %s", err)
				// use an empty string to get a view of the world and uncomment this next line...
				// fmt.Printf("\nexpected:\n\t%v\nactual:\n\t%v", expected, ast)
			}

			if cmp.Equal(ast, expected) {
				t.Errorf("expected %s AST shape to match expected.", snapshotPath)
			}
		})
	}
}
