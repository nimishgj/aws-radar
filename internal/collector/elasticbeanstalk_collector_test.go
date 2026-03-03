package collector

import "testing"

func TestCollector_elasticbeanstalk_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewElasticBeanstalkCollector(), "elasticbeanstalk")
}

func TestCollector_elasticbeanstalk_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewElasticBeanstalkCollector(), true)
}
