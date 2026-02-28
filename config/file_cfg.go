package config

var (
	_config *Config
)

type Config struct {
	SourceDB   *SourceDB  `yaml:"sourceDB"`
	OutPutType string     `yaml:"outPutType"`
	Mysql      *Mysql     `yaml:"mysql"`
	Post       *Post      `yaml:"post"`
	SlowLog    *SlowLog   `yaml:"slowLog"`
	StarRocks  *StarRocks `yaml:"starRocks"`
	Jobs       []Job      `yaml:"jobs"`
}

type Job struct {
	Instance  string   `yaml:"instance"`
	Interval  string   `yaml:"interval"`
	DsnSource []string `yaml:"dsnSource"`
	DsnDest   []string `yaml:"dsnDest"`
	Queries   []Query  `yaml:"queries"`
	Token     string   `yaml:"token"`
	ApiGet    string   `yaml:"apiGet"`
	ApiPost   string   `yaml:"apiPost"`
	ApiPut    string   `yaml:"apiPut"`
}

type Query struct {
	Name        string `yaml:"name"`
	Query       string `yaml:"query"`
	DestDDL     string `yaml:"destDDL"`
	DestTblName string `yaml:"destTblName"`
}

type SourceDB struct {
	DbType        string `yaml:"dbType"`
	QueryLength   int    `yaml:"queryLength"`
	CheckInterval int    `yaml:"checkInterval"`
}

type Mysql struct {
	Host      string `yaml:"host"`
	Port      uint32 `yaml:"port"`
	User      string `yaml:"user"`
	PassWord  string `yaml:"passWord"`
	Dbname    string `yaml:"dbName"`
	TableName string `yaml:"tableName"`
	Other     string `yaml:"other"`
}

type SlowLog struct {
	Prefix        string   `yaml:"prefix"`
	Path          string   `yaml:"path"`
	Env           string   `yaml:"env"`
	Instance      string   `yaml:"instance"`
	IgnoreUser    []string `yaml:"ignoreuser"`
	LongQueryTime float64  `yaml:"longquerytime"`
	ReadFirst     bool     `yaml:"readFirst"`
}

type Post struct {
	Host string `yaml:"host"`
	Url  string `yaml:"url"`
}

type StarRocks struct {
	Host                  string `yaml:"host"`
	QueryPort             uint32 `yaml:"queryport"`
	WebserverPort         uint32 `yaml:"webserverport"`
	User                  string `yaml:"user"`
	PassWord              string `yaml:"passWord"`
	Dbname                string `yaml:"dbName"`
	TableName             string `yaml:"tableName"`
	MergeCommitIntervalMs uint32 `yaml:"mergeCommitIntervalMs"`
	MergeCommitParallel   uint32 `yaml:"mergeCommitParallel"`
	PartitionLiveNumber   uint32 `yaml:"partitionLiveNumber"`
}
