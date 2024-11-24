package as3

import "encoding/json"

type Destination struct {
	AddressLists []Use `json:"addressLists"`
	PortLists    []Use `json:"portLists"`
}

type Source struct {
	AddressLists []Use `json:"addressLists"`
}

type NatRule struct {
	Destination       Destination `json:"destination"`
	Name              string      `json:"name"`
	Protocol          string      `json:"protocol"`
	Source            Source      `json:"source"`
	SourceTranslation Use         `json:"sourceTranslation"`
}

type NatPolicy struct {
	Class string    `json:"class"`
	Rules []NatRule `json:"rules"`
}

func (np *NatPolicy) String() string {
	data, err := json.Marshal(np)
	if err != nil {
		return "{}"
	}

	return string(data)
}

func ParseStringToSnatPolicy(data string) (*NatPolicy, error) {
	var ret NatPolicy
	if err := json.Unmarshal([]byte(data), &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

type NatSourceTranslation struct {
	Class     string   `json:"class"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses"`
}

func NewNatSourceTranslation(ips []string) NatSourceTranslation {
	return NatSourceTranslation{
		Class:     ClassNatSourceTranslation,
		Type:      "static-nat",
		Addresses: ips,
	}
}
