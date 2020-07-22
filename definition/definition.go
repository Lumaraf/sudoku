package definition

import (
	"errors"
	"fmt"
	"github.com/lumaraf/sudoku-solver/generator"
	"github.com/lumaraf/sudoku-solver/rules"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
)

type Definition struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
	Rules       Rules  `yaml:"rules"`
}

type Rules []generator.Rule

func (r *Rules) UnmarshalYAML(unmarshal func(interface{}) error) error {
	raw := make([]genericRule, 0)

	if err := unmarshal(&raw); err != nil {
		return err
	}

	for _, gr := range raw {
		*r = append(*r, gr.Rule)
	}
	return nil
}

type genericRule struct {
	Rule generator.Rule
}

type baseRule struct {
	Type string `yaml:"type"`
}

func (g *genericRule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	base := baseRule{}
	if err := unmarshal(&base); err != nil {
		return err
	}

	ruleTypes := map[string]reflect.Type{
		"row":               reflect.TypeOf(rules.RowRule{}),
		"column":            reflect.TypeOf(rules.ColumnRule{}),
		"box":               reflect.TypeOf(rules.BoxRule{}),
		"cross":             reflect.TypeOf(rules.CrossRule{}),
		"magic-square":      reflect.TypeOf(rules.MagicSquareRule{}),
		"anti-knights-move": reflect.TypeOf(rules.AntiKnightsMoveRule{}),
		"anti-kings-move":   reflect.TypeOf(rules.AntiKingsMoveRule{}),
		"anti-queens-move":  reflect.TypeOf(rules.AntiQueensMoveRule{}),
	}

	if t, ok := ruleTypes[base.Type]; ok {
		rule := reflect.New(t).Interface().(generator.Rule)
		if err := unmarshal(rule); err != nil {
			return err
		}
		g.Rule = rule
	} else {
		return errors.New(fmt.Sprintf("unsupported rule type '%s'", base.Type))
	}
	return nil
}

func FromFile(file string) (def Definition, err error) {
	var content []byte
	if content, err = ioutil.ReadFile(file); err != nil {
		return def, err
	}
	if err := yaml.Unmarshal(content, &def); err != nil {
		return def, err
	}
	return
}
