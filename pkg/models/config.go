package models

type Configs struct {
	MongodbURI  string              `yaml:"mongodbURI"`
	DBName      string              `yaml:"dbName"`
	KafkaURL    string              `yaml:"kafkaURL"`
	Topic       string              `yaml:"topic"`
	BackLogDB   DBConfig            `yaml:"backLogDB"`
	KafkaConfig []*KafkaGroupConfig `yaml:"kafkaConfig"`
}

type DBConfig struct {
	URI    string `yaml:"uri"`
	DBName string `yaml:"dbName"`
}

type KafkaGroupConfig struct {
	Brokers           []string `yaml:"brokers"`
	Topic             string   `yaml:"topic"`
	GroupID           string   `yaml:"groupID"`
	NumOfConsumer     int      `yaml:"numOfConsumer"`
	RetryStrategy     string   `yaml:"retryStrategy"`
	DeadLetterAddress string   `yaml:"deadLetterAddress,omitempty"`
	BackLogCollection string   `yaml:"backLogCollection,omitempty"`
}
