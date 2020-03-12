package profx

import (
	"fmt"
	"profx/domain"
	"regexp"
)

type Sourcer interface {
	Load() ([]domain.Source, error)
	Match(source, url string) (*bool, error)
}

type CompiledRule struct {
	rule       string
	ruleType   domain.RuleType
	ruleRegexp *regexp.Regexp
}

func (cr *CompiledRule) Compile() error {
	var err error
	cr.ruleRegexp, err = regexp.Compile(cr.rule)
	if err != nil {
		return fmt.Errorf("failed to compile regexp: %w", err)
	}

	return nil
}

func (cr *CompiledRule) Match(s string) (*bool, error) {
	var res bool
	var err error

	if cr.ruleRegexp == nil {
		err = cr.Compile()
		if err != nil {
			return nil, err
		}
	}

	res = cr.ruleRegexp.MatchString(s)
	return &res, nil
}

type PersistentMemorySourcer struct {
	repo    CrawlerRepository
	sources []domain.Source
	rules   []domain.Rule
	memory  map[string][]*CompiledRule
}

func NewPersistentMemorySourcer(repo CrawlerRepository) Sourcer {
	return &PersistentMemorySourcer{repo: repo}
}

func (pms *PersistentMemorySourcer) prepareMemory() error {
	if pms.sources == nil || pms.rules == nil {
		return fmt.Errorf("empty data sources for sourcer")
	}

	pms.memory = map[string][]*CompiledRule{}
	for _, rule := range pms.rules {
		cr := &CompiledRule{
			rule:       rule.Rule,
			ruleType:   rule.Type,
			ruleRegexp: nil,
		}
		pms.memory[rule.Source] = append(pms.memory[rule.Source], cr)
	}

	return nil
}

func (pms *PersistentMemorySourcer) Load() ([]domain.Source, error) {
	var err error

	if pms.sources == nil {
		pms.sources, err = pms.repo.LoadSources()
		if err != nil {
			return nil, fmt.Errorf("failed to load sources from repository: %w", err)
		}
	}

	return pms.sources, nil
}

func (pms *PersistentMemorySourcer) Match(source, url string) (*bool, error) {
	var err error
	var res = false

	if pms.sources == nil {
		_, err = pms.Load()
		return nil, err
	}

	if pms.rules == nil {
		pms.rules, err = pms.repo.LoadRules()
		if err != nil {
			return nil, fmt.Errorf("failed to load source rules from repository: %w", err)
		}
	}

	if pms.memory == nil {
		err = pms.prepareMemory()
		if err != nil {
			return nil, err
		}
	}

	rules := pms.memory[source]

	if len(rules) == 0 {
		return nil, fmt.Errorf("no rules found for source %s", source)
	}

	matchedInclude := false
	matchedExclude := false

	for _, rule := range rules {
		ruleMatch, err := rule.Match(url)
		if err != nil {
			return nil, err
		}

		if rule.ruleType == domain.IncludeRuleType && *ruleMatch == true {
			matchedInclude = true
		}

		if rule.ruleType == domain.ExcludeRuleType && *ruleMatch == true {
			matchedExclude = true
		}
	}

	if matchedExclude == false && matchedInclude == true {
		res = true
	}
	return &res, nil
}
