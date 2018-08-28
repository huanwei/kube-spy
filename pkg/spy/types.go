package spy

type TestCase struct {
	Method    string            `yaml:"method"`
	URL       string            `yaml:"url"`
	Params    map[string]string `yaml:"params"`
	AuthToken string            `yaml:"authToken"`
	BasicAuth struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"basicAuth"`
	Headers     map[string]string   `yaml:"headers"`
	Form        map[string]string   `yaml:"form"`
	MultiForm   map[string][]string `yaml:"multiValueForm"`
	MultiParams map[string][]string `yaml:"multiValueParams"`
	Files       map[string]string   `yaml:"files"`
	PathParams  map[string]string   `yaml:"pathParams"`
	Body        string              `yaml:"body"`
}

type Chaos struct {
	Replica int    `yaml:"replica"`
	Ingress string `yaml:"ingress"`
	Egress  string `yaml:"egress"`
}

type ClientSetting struct {
	RetryCount   int `yaml:"retryCount"`
	RetryWait    int `yaml:"retryWait"`
	RetryMaxWait int `yaml:"retryMaxWait"`
	Timeout      int `yaml:"timeout"`
}

type VictimService struct {
	Name      string  `yaml:"name"`
	ChaosList []Chaos `yaml:"chaosList"`
}

type Config struct {
	Namespace      string          `yaml:"Namespace"`
	VictimServices []VictimService `yaml:"VictimServices"`
	APIServerAddr  string          `yaml:"APIServerAddr"`
	APISetting     TestCase        `yaml:"APISettings"`
	TestCases      []TestCase      `yaml:"TestCases"`
	ClientSetting  ClientSetting   `yaml:"ClientSetting"`
}
