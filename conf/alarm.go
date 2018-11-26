package conf

import (
	"fmt"
	"sort"
	"strings"
)

type BasicAlert struct {
	AlarmID string `json:"alarmid"`
	Content string `json:"content"`
}

type LabelSet map[string]string

type T8TAlert struct {
	Labels      LabelSet `json:"labels"`
	Annotations LabelSet `json:"annotations"`
}

func (l LabelSet) String() string {
	lstrs := make([]string, 0, len(l))
	for l, v := range l {
		lstrs = append(lstrs, fmt.Sprintf("%s:%s", l, v))
	}

	sort.Strings(lstrs)
	return fmt.Sprintf("%s", strings.Join(lstrs, ", "))
}
