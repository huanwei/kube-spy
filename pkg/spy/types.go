package spy

type TestCase struct {
	Method    string            `yaml:"method"`
	URL       string            `yaml:"url"`
	Params    map[string]string `yaml:"params"`
	Authtoken string            `yaml:"authtoken"`
	Headers   map[string]string `yaml:"headers"`
	Form      map[string]string `yaml:"form"`
	Files     map[string]string `yaml:"files"`
	Body      string            `yaml:"body"`
}

type Config struct {
	Namespace      string     `yaml:"Namespace"`
	ServiceList    []string   `yaml:"ServiceList"`
	TestCaseList   []TestCase `yaml:"TestCaseList"`
	ChaosList      []string   `yaml:"ChaosList"`
	GlobalSettings TestCase   `yaml:"GlobalSettings"`
}
