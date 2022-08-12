package franz

import "github.com/sangoisanga/kafka-go/pkg/interfaces"

type Group struct {
	Brokers        []string
	Topic          string
	GroupID        string
	NumOfConsumer  int
	ConsumeHandler interfaces.ConsumerHandler
	RetryHandler   interfaces.RetryStrategy
}
