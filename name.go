package ark

import "fmt"

type Name struct {
	Name     string `json:"name"`
	Instance int    `json:"instance"`
}

func (n Name) String() string {
	if n.Instance != 0 {
		return fmt.Sprintf("%s:%d", n.Name, n.Instance)
	} else {
		return n.Name
	}
}

func (n Name) IsEmpty() bool {
	return n.Name == "" && n.Instance == 0
}

func (n Name) IsNone() bool {
	return n.Name == "None" && n.Instance == 0
}
