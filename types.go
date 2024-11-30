package main

type ShowMapsOutput struct {
	Maps []struct {
		Name       string `json:"name,omitempty"`
		UUID       string `json:"uuid,omitempty"`
		Sysfs      string `json:"sysfs,omitempty"`
		Failback   string `json:"failback,omitempty"`
		Queueing   string `json:"queueing,omitempty"`
		Paths      int    `json:"paths,omitempty"`
		WriteProt  string `json:"write_prot,omitempty"`
		DmSt       string `json:"dm_st,omitempty"`
		PathFaults int    `json:"path_faults,omitempty"`
		Vend       string `json:"vend,omitempty"`
		Prod       string `json:"prod,omitempty"`
		PathGroups []struct {
			Selector string `json:"selector,omitempty"`
			DmSt     string `json:"dm_st,omitempty"`
			Group    int    `json:"group,omitempty"`
			Paths    []struct {
				Dev         string `json:"dev,omitempty"`
				DmSt        string `json:"dm_st,omitempty"`
				DevSt       string `json:"dev_st,omitempty"`
				ChkSt       string `json:"chk_st,omitempty"`
				HostWwnn    string `json:"host_wwnn,omitempty"`
				TargetWwnn  string `json:"target_wwnn,omitempty"`
				HostWwpn    string `json:"host_wwpn,omitempty"`
				TargetWwpn  string `json:"target_wwpn,omitempty"`
				HostAdapter string `json:"host_adapter,omitempty"`
			} `json:"paths,omitempty"`
		} `json:"path_groups,omitempty"`
	} `json:"maps,omitempty"`
}
