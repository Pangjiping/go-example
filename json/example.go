package json

import (
	"encoding/json"
)

type Metric struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func (m *Metric) UnmarshalJSON(data []byte) error {
	type AliasMetric Metric
	t := &struct {
		Value float64 `json:"value"`
		*AliasMetric
	}{
		Value:       float64(m.Value),
		AliasMetric: (*AliasMetric)(m),
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	m.Value = int64(t.Value)
	return nil
}

func Usage(s string) (*Metric, error) {
	metric := &Metric{}
	//err := json.Unmarshal([]byte(`{"name":"tq","value":1.1}`), &metric)
	if err := metric.UnmarshalJSON([]byte(s)); err != nil {
		return nil, err
	}
	return metric, nil
}
