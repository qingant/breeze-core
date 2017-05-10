package breeze

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"time"
	"log"
)

func encodeCSV(columns []string, rows []map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.Write(columns); err != nil {
		return nil, err
	}
	r := make([]string, len(columns))
	for _, row := range rows {
		for i, column := range columns {
			v, _ := json.Marshal(row[column])
			r[i] = string(v)
		}
		if err := w.Write(r); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return buf.Bytes(), nil
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
	GetStatus().logPerf(name, elapsed)
}
