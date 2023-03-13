package terraform

import (
	"errors"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type ModuleUsage struct {
	Path      string
	Variables map[string]int
	Locals    map[string]int
	file      *hclwrite.File
}

func NewModuleUsage(path string) (*ModuleUsage, error) {
	m := &ModuleUsage{
		Path:      path,
		Variables: map[string]int{},
		Locals:    map[string]int{},
	}

	src, err := LoadTfModule(path)
	if err != nil {
		return nil, err
	}

	f, diags := hclwrite.ParseConfig(src, "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, errors.New(path + ":" + diags.Error())
	}

	m.file = f
	err = m.processUsage()

	return m, err
}

func ListTfModules(path string) (map[string]bool, error) {
	var directories = make(map[string]bool)

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".tf" {
			module := filepath.Dir(path)
			log.Debugf("Visited: %s\n", module)
			if _, ok := directories[module]; !ok {
				directories[module] = true
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return directories, nil
}

func LoadTfModule(path string) ([]byte, error) {
	var out []byte

	files, err := os.ReadDir(path)
	if err != nil {
		return out, err
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".tf" {
			data, err := os.ReadFile(filepath.Join(path, file.Name()))
			if err != nil {
				return out, err
			}
			out = append(out, '\n')
			out = append(out, data...)
		}
	}
	return out, nil
}

func (m ModuleUsage) processUsage() error {
	body := m.file.Body()
	bodyStr := string(m.file.Bytes())
	for _, block := range body.Blocks() {
		blockType := block.Type()
		if blockType == "variable" {
			name := block.Labels()[0]
			m.Variables[name] = countPattern(bodyStr, fmt.Sprintf(`var\.%s\W`, name))
		} else if blockType == "locals" {
			attribs := block.Body().Attributes()
			for attrib := range attribs {
				m.Locals[attrib] = countPattern(bodyStr, fmt.Sprintf(`local\.%s\W`, attrib))
			}
		}

	}

	return nil
}

func countPattern(content string, pattern string) int {
	regex := regexp.MustCompile(pattern)
	matches := regex.FindAllStringIndex(content, -1)

	return len(matches)
}

func (m ModuleUsage) DisplayLocals(unusedOnly bool) error {
	return m.Display(Locals, unusedOnly)
}

func (m ModuleUsage) DisplayVariables(unusedOnly bool) error {
	return m.Display(Variables, unusedOnly)
}

type DisplayType string

const (
	All       DisplayType = "all"
	Variables DisplayType = "variables"
	Locals    DisplayType = "locals"
)

func filterUnusedOnly(items map[string]int) map[string]int {
	for name, count := range items {
		if count > 0 {
			delete(items, name)
		}
	}
	return items
}

func (m ModuleUsage) Display(dType DisplayType, unusedOnly bool) error {
	variables := map[string]int{}
	locals := map[string]int{}

	switch dType {
	case Locals:
		locals = m.Locals
		if unusedOnly {
			locals = filterUnusedOnly(locals)
		}
	case Variables:
		variables = m.Variables
		if unusedOnly {
			variables = filterUnusedOnly(variables)
		}
	case All:
		locals = m.Locals
		variables = m.Variables
		if unusedOnly {
			locals = filterUnusedOnly(locals)
			variables = filterUnusedOnly(variables)
		}
	default:
		return errors.New(fmt.Sprintf("%s is an invalid display Type", dType))
	}

	if !unusedOnly || (unusedOnly && len(locals)+len(variables) > 0) {
		fmt.Printf("\n \U0001F680 Module: %s\n", m.Path)
	}

	if dType == All || dType == Variables {
		if !unusedOnly || (unusedOnly && len(variables) > 0) {
			fmt.Printf(" \U0001F449 %d variables found\n", len(variables))
		}
		for name, count := range variables {
			fmt.Printf("%s : %d\n", name, count)
		}
	}

	if dType == All || dType == Locals {
		if !unusedOnly || (unusedOnly && len(locals) > 0) {
			fmt.Printf("\U0001F449 %d locals found\n", len(locals))
		}
		for name, count := range locals {
			fmt.Printf("%s : %d\n", name, count)
		}

	}
	return nil
}
