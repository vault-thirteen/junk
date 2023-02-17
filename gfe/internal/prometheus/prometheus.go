package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vault-thirteen/junk/gfe/internal/application/config"
)

// Prometheus -- метрики.
type Prometheus struct {
	registry         *prometheus.Registry
	processCollector prometheus.Collector
	goCollector      prometheus.Collector
	metrics          *Metrics
}

// NewPrometheus -- конструктор метрик.
func NewPrometheus() (p *Prometheus, err error) {
	p = new(Prometheus)

	p.registry = prometheus.NewPedanticRegistry()

	p.metrics = new(Metrics)

	p.addIndicators()

	err = p.registerIndicators()
	if err != nil {
		return nil, err
	}

	return p, nil
}

// addIndicators добавляет показатели.
// Поле 'metrics' должно быть создано.
func (p *Prometheus) addIndicators() {
	// Стандартные коллекторы.
	p.processCollector = prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{})
	p.goCollector = prometheus.NewGoCollector()

	// Коллекторы сервиса.

	// #1.
	p.metrics.RequestsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: config.MetricNamePrefixForService + "http_requests_count",
			Help: "Счётчик количества входящих HTTP запросов.",
		},
		[]string{config.MetricParameterPath},
	)

	// #2.
	p.metrics.KafkaConsumerErrorsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: config.MetricNamePrefixForService + "kafka_consumer_errors_count",
			Help: "Счётчик количества ошибок читателя сообщений Kafka. " +
				"Данные ошибки являются критическими и обычно говорят о " +
				"неправильных или повреждённых сообщениях в очереди Kafka.",
		},
	)

	// #3.
	p.metrics.KafkaConsumerGroupErrorsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: config.MetricNamePrefixForService + "kafka_consumer_group_errors_count",
			Help: "Счётчик количества ошибок группы читателя Kafka. " +
				"Данные ошибки говорят о проблемах с очередью сообщений: " +
				"например, об отсутствии связи с сервером Kafka или Zookeeper.",
		},
	)

	// #4.
	p.metrics.DatabaseConnectionErrorsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: config.MetricNamePrefixForService + "database_connection_errors_count",
			Help: "Счётчик количества ошибок при соединении с сервером базы " +
				"данных.",
		},
	)
}

// registerIndicators производит регистрацию показателей.
func (p *Prometheus) registerIndicators() (err error) {
	// Стандартные коллекторы.

	// ProcessCollector.
	err = p.registry.Register(p.processCollector)
	if err != nil {
		return err
	}

	// GoCollector.
	err = p.registry.Register(p.goCollector)
	if err != nil {
		return err
	}

	// Коллекторы сервиса.

	// #1.
	err = p.registry.Register(p.metrics.RequestsCount)
	if err != nil {
		return err
	}

	// #2.
	err = p.registry.Register(p.metrics.KafkaConsumerErrorsCount)
	if err != nil {
		return err
	}

	// #3.
	err = p.registry.Register(p.metrics.KafkaConsumerGroupErrorsCount)
	if err != nil {
		return err
	}

	// #4.
	err = p.registry.Register(p.metrics.DatabaseConnectionErrorsCount)
	if err != nil {
		return err
	}

	return nil
}

// GetMetrics возвращает метрики.
func (p *Prometheus) GetMetrics() *Metrics {
	return p.metrics
}

// GetRegistry возвращает Registry.
func (p *Prometheus) GetRegistry() *prometheus.Registry {
	return p.registry
}
