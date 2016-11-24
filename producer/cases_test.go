package producer

type determineFactoryCase struct {
	producerType string
	validate     validateConfigFunc
	init         initProducerFunc
	isError      bool
}

var determineFactoryCases = []determineFactoryCase{
	{"dummy", dummyValidateConfig, dummyInitProducer, false},
	{"kafka", kafkaValidateConfig, kafkaInitProducer, false},
	{"non-existant", nil, nil, true},
	{"", nil, nil, true},
}
