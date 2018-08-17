package spy

type TestCase struct {
	Method      string              `yaml:"method"`
	URL         string              `yaml:"url"`
	Params      map[string]string   `yaml:"params"`
	Authtoken   string              `yaml:"authtoken"`
	Headers     map[string]string   `yaml:"headers"`
	Form        map[string]string   `yaml:"form"`
	MultiForm   map[string][]string `yaml:"multivalueform"`
	MultiParams map[string][]string `yaml:"multivalueparams"`
	Files       map[string]string   `yaml:"files"`
	PathParams  map[string]string   `yaml:"pathparams"`
	Body        string              `yaml:"body"`
}

type Config struct {
	Namespace      string     `yaml:"Namespace"`
	ServiceList    []string   `yaml:"ServiceList"`
	APIServerAddr  string     `yaml:"APIServerAddr"`
	TestCaseList   []TestCase `yaml:"TestCaseList"`
	ChaosList      []string   `yaml:"ChaosList"`
	GlobalSettings TestCase   `yaml:"GlobalSettings"`
}

// TODO: APIServerAddr can be a slice in case we do different chaos strategy on multiple Pods within the same one service.
type SpyConfig struct {
	NameSpace      string          `yaml:"Namespace"`
	VictimServices []VictimService `yaml:"VictimServices"`
}

type VictimService struct {
	ServiceName    string     `yaml:"ServiceName"`
	APIServerAddr  []string   `yaml:"APIServerAddr"`
	TestCaseList   []TestCase `yaml:"TestCaseList"`
	ChaosList      []string   `yaml:"ChaosList"`
	GlobalSettings TestCase   `yaml:"GlobalSettings"`
}
