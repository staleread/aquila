// Command hypotheses validates hypotheses.yaml, the only claim registry of this
// repository, against the evidence it cites (see AGENTS.md).
//
//	go run ./tools/hypotheses check          static: registry, refs, test names
//	go run ./tools/hypotheses check --run    also run every tagged test
//	go run ./tools/hypotheses test H03 [...] run the tests of these ids and of their rejectors
//	go run ./tools/hypotheses next           print the next free id
//
// A Go test carries an id in its top-level name, delimited by an underscore or
// the name boundary after the Test prefix: TestH03_Bijective, TestSLE_H07.
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	registryFile = "hypotheses.yaml"
	bibFile      = "article/references.bib"
)

var (
	areas     = []string{"research", "model", "impl", "structure", "diffusion", "inversion", "feasibility", "novelty"}
	statuses  = []string{"open", "proven", "rejected"}
	evidences = []string{"derivation", "test", "experiment", "literature"}

	idPattern  = regexp.MustCompile(`^H(\d{2,3})$`)
	tagPattern = regexp.MustCompile(`(?:^|_)(H\d{2,3})(?:_|$)`)
)

type Hypothesis struct {
	Area       string   `yaml:"area"`
	Status     string   `yaml:"status"`
	Evidence   string   `yaml:"evidence"`
	Ref        string   `yaml:"ref,omitempty"`
	Parent     string   `yaml:"parent,omitempty"`
	RejectedBy []string `yaml:"rejectedBy,omitempty"`
	Statement  string   `yaml:"statement"`
}

type testRef struct{ Pkg, Name string }

type registry struct {
	ids   []string // file order
	items map[string]Hypothesis
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	var err error
	switch os.Args[1] {
	case "check":
		run := len(os.Args) > 2 && os.Args[2] == "--run"
		err = check(run)
	case "test":
		if len(os.Args) < 3 {
			usage()
		}
		err = testIDs(os.Args[2:])
	case "next":
		err = next()
	default:
		usage()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: hypotheses check [--run] | test ID... | next")
	os.Exit(2)
}

func load() (*registry, error) {
	data, err := os.ReadFile(registryFile)
	if err != nil {
		return nil, err
	}
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil { // also rejects duplicate keys
		return nil, fmt.Errorf("%s: %w", registryFile, err)
	}
	reg := &registry{items: map[string]Hypothesis{}}
	if len(root.Content) == 0 {
		return reg, nil
	}
	m := root.Content[0]
	if m.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%s: top level must be a mapping keyed by id", registryFile)
	}
	for i := 0; i < len(m.Content); i += 2 {
		id := m.Content[i].Value
		var h Hypothesis
		dec := m.Content[i+1]
		if err := strictDecode(dec, &h); err != nil {
			return nil, fmt.Errorf("%s: %s: %w", registryFile, id, err)
		}
		reg.ids = append(reg.ids, id)
		reg.items[id] = h
	}
	return reg, nil
}

// strictDecode rejects fields the schema does not know.
func strictDecode(n *yaml.Node, h *Hypothesis) error {
	known := []string{"area", "status", "evidence", "ref", "parent", "rejectedBy", "statement"}
	if n.Kind != yaml.MappingNode {
		return errors.New("entry must be a mapping")
	}
	for i := 0; i < len(n.Content); i += 2 {
		if k := n.Content[i].Value; !slices.Contains(known, k) {
			return fmt.Errorf("unknown field %q", k)
		}
	}
	return n.Decode(h)
}

// scanTests maps each id to the top-level Go tests whose names carry it.
func scanTests() (map[string][]testRef, error) {
	tags := map[string][]testRef{}
	fset := token.NewFileSet()
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != "." && (strings.HasPrefix(d.Name(), ".") || d.Name() == "references" || d.Name() == "article") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		for _, decl := range f.Decls {
			name := funcName(decl)
			if !strings.HasPrefix(name, "Test") {
				continue
			}
			for _, m := range tagPattern.FindAllStringSubmatch(strings.TrimPrefix(name, "Test"), -1) {
				tags[m[1]] = append(tags[m[1]], testRef{"./" + filepath.Dir(path), name})
			}
		}
		return nil
	})
	return tags, err
}

func validate(reg *registry, tags map[string][]testRef) []string {
	var errs []string
	bad := func(id, format string, a ...any) {
		errs = append(errs, id+": "+fmt.Sprintf(format, a...))
	}
	for _, id := range reg.ids {
		h := reg.items[id]
		if !idPattern.MatchString(id) {
			bad(id, "id must match %s", idPattern)
		}
		if !slices.Contains(areas, h.Area) {
			bad(id, "area %q not in %v", h.Area, areas)
		}
		if !slices.Contains(statuses, h.Status) {
			bad(id, "status %q not in %v", h.Status, statuses)
		}
		if !slices.Contains(evidences, h.Evidence) {
			bad(id, "evidence %q not in %v", h.Evidence, evidences)
		}
		if strings.TrimSpace(h.Statement) == "" {
			bad(id, "empty statement")
		}
		if h.Parent != "" {
			if _, ok := reg.items[h.Parent]; !ok || h.Parent == id {
				bad(id, "parent %q is not another registry id", h.Parent)
			}
		}
		provenRejector := false
		for _, r := range h.RejectedBy {
			rh, ok := reg.items[r]
			if !ok || r == id {
				bad(id, "rejectedBy %q is not another registry id", r)
				continue
			}
			provenRejector = provenRejector || rh.Status == "proven"
		}
		switch h.Status {
		case "rejected":
			if !provenRejector {
				bad(id, "rejected needs a proven hypothesis in rejectedBy")
			}
		case "open":
			if provenRejector {
				bad(id, "a rejector is proven; set status to rejected")
			}
		case "proven":
			if provenRejector {
				bad(id, "proven but a rejector is proven too")
			}
			if msg := checkEvidence(h, tags[id]); msg != "" {
				bad(id, "%s", msg)
			}
		}
	}
	for id, refs := range tags {
		if _, ok := reg.items[id]; !ok {
			for _, t := range refs {
				errs = append(errs, fmt.Sprintf("%s: test %s in %s carries an id missing from the registry", id, t.Name, t.Pkg))
			}
		}
	}
	sort.Strings(errs)
	return errs
}

func checkEvidence(h Hypothesis, tests []testRef) string {
	switch h.Evidence {
	case "test":
		if len(tests) == 0 {
			return "proven by test, but no test carries the id"
		}
	case "derivation":
		path, label, ok := strings.Cut(h.Ref, "#")
		if !ok || label == "" {
			return "derivation ref must be <file.typ>#<label>"
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return "derivation ref: " + err.Error()
		}
		if !bytes.Contains(data, []byte("<"+label+">")) {
			return fmt.Sprintf("label <%s> not found in %s", label, path)
		}
	case "literature":
		data, err := os.ReadFile(bibFile)
		if err != nil {
			return "literature ref: " + err.Error()
		}
		if h.Ref == "" || !regexp.MustCompile(`@\w+\{`+regexp.QuoteMeta(h.Ref)+`,`).Match(data) {
			return fmt.Sprintf("bib key %q not found in %s", h.Ref, bibFile)
		}
	case "experiment":
		if h.Ref == "" {
			return "experiment ref must name the protocol/results path"
		}
		if _, err := os.Stat(h.Ref); err != nil {
			return "experiment ref: " + err.Error()
		}
	}
	return ""
}

// runTests runs the given tests and returns pass/fail counts per test name.
func runTests(refs []testRef) (map[testRef]string, error) {
	byPkg := map[string][]string{}
	for _, t := range refs {
		byPkg[t.Pkg] = append(byPkg[t.Pkg], t.Name)
	}
	results := map[testRef]string{}
	pkgs := make([]string, 0, len(byPkg))
	for p := range byPkg {
		pkgs = append(pkgs, p)
	}
	sort.Strings(pkgs)
	for _, pkg := range pkgs {
		names := slices.Compact(slices.Sorted(slices.Values(byPkg[pkg])))
		cmd := exec.Command("go", "test", "-count=1", "-json", "-run", "^("+strings.Join(names, "|")+")$", pkg)
		cmd.Stderr = os.Stderr
		out, _ := cmd.Output() // failures are read from the JSON stream
		sc := bufio.NewScanner(bytes.NewReader(out))
		sc.Buffer(make([]byte, 1<<20), 1<<24)
		for sc.Scan() {
			var ev struct{ Action, Test string }
			if json.Unmarshal(sc.Bytes(), &ev) != nil || ev.Test == "" || strings.Contains(ev.Test, "/") {
				continue
			}
			if ev.Action == "pass" || ev.Action == "fail" || ev.Action == "skip" {
				results[testRef{pkg, ev.Test}] = ev.Action
			}
		}
		for _, n := range names {
			if _, ok := results[testRef{pkg, n}]; !ok {
				results[testRef{pkg, n}] = "fail" // did not run: build error or similar
			}
		}
	}
	return results, nil
}

// judge reports, for one id, whether its tests passed; failures are errors only for proven ids.
func judge(id string, h Hypothesis, tests []testRef, res map[testRef]string) (errs, warns []string) {
	pass := 0
	for _, t := range tests {
		switch res[t] {
		case "pass":
			pass++
		case "fail":
			msg := fmt.Sprintf("%s: %s (%s) failed", id, t.Name, t.Pkg)
			if h.Status == "proven" {
				errs = append(errs, msg)
			} else {
				warns = append(warns, msg)
			}
		}
	}
	if h.Status == "proven" && h.Evidence == "test" && pass == 0 {
		errs = append(errs, id+": proven by test, but no tagged test passed (skips do not count)")
	}
	return errs, warns
}

func check(run bool) error {
	reg, err := load()
	if err != nil {
		return err
	}
	tags, err := scanTests()
	if err != nil {
		return err
	}
	errs := validate(reg, tags)
	var warns []string
	if run && len(errs) == 0 {
		var all []testRef
		for _, id := range reg.ids {
			all = append(all, tags[id]...)
		}
		res, err := runTests(all)
		if err != nil {
			return err
		}
		for _, id := range reg.ids {
			e, w := judge(id, reg.items[id], tags[id], res)
			errs, warns = append(errs, e...), append(warns, w...)
		}
	}
	return report(reg, errs, warns)
}

func testIDs(ids []string) error {
	reg, err := load()
	if err != nil {
		return err
	}
	tags, err := scanTests()
	if err != nil {
		return err
	}
	want := []string{}
	for _, id := range ids {
		h, ok := reg.items[id]
		if !ok {
			return fmt.Errorf("%s: not in the registry", id)
		}
		want = append(want, id)
		want = append(want, h.RejectedBy...)
	}
	want = slices.Compact(slices.Sorted(slices.Values(want)))
	var refs []testRef
	for _, id := range want {
		refs = append(refs, tags[id]...)
	}
	res, err := runTests(refs)
	if err != nil {
		return err
	}
	var errs, warns []string
	for _, id := range want {
		h := reg.items[id]
		if len(tags[id]) == 0 {
			fmt.Printf("%s: %s, %s evidence, no tests\n", id, h.Status, h.Evidence)
		}
		for _, t := range tags[id] {
			fmt.Printf("%s: %-4s %s (%s)\n", id, res[t], t.Name, t.Pkg)
		}
		e, w := judge(id, h, tags[id], res)
		errs, warns = append(errs, e...), append(warns, w...)
	}
	return report(nil, errs, warns)
}

func report(reg *registry, errs, warns []string) error {
	for _, w := range warns {
		fmt.Println("warning:", w)
	}
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println("error:", e)
		}
		return fmt.Errorf("%d error(s)", len(errs))
	}
	if reg != nil {
		count := map[string]int{}
		for _, h := range reg.items {
			count[h.Status]++
		}
		fmt.Printf("ok: %d hypotheses (%d proven, %d open, %d rejected)\n",
			len(reg.ids), count["proven"], count["open"], count["rejected"])
	}
	return nil
}

func next() error {
	reg, err := load()
	if err != nil {
		return err
	}
	last := 0
	for _, id := range reg.ids {
		if m := idPattern.FindStringSubmatch(id); m != nil {
			n, _ := strconv.Atoi(m[1])
			last = max(last, n)
		}
	}
	fmt.Printf("H%02d\n", last+1)
	return nil
}

func funcName(d ast.Decl) string {
	if fd, ok := d.(*ast.FuncDecl); ok && fd.Recv == nil {
		return fd.Name.Name
	}
	return ""
}
