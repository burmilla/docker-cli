package formatter

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
)

type byPortInfo []types.Port

func (r byPortInfo) Len() int      { return len(r) }
func (r byPortInfo) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r byPortInfo) Less(i, j int) bool {
	if r[i].PrivatePort != r[j].PrivatePort {
		return r[i].PrivatePort < r[j].PrivatePort
	}

	if r[i].IP != r[j].IP {
		return r[i].IP < r[j].IP
	}

	if r[i].PublicPort != r[j].PublicPort {
		return r[i].PublicPort < r[j].PublicPort
	}

	return r[i].Type < r[j].Type
}

// DisplayablePorts returns formatted string representing open ports of container
// e.g. "0.0.0.0:80->9090/tcp, 9988/tcp"
// it's used by command 'docker ps'
func displayablePorts(ports []types.Port) string {
	type portGroup struct {
		first uint16
		last  uint16
	}
	groupMap := make(map[string]*portGroup)
	var result []string
	var hostMappings []string
	var groupMapKeys []string
	sort.Sort(byPortInfo(ports))
	for _, port := range ports {
		current := port.PrivatePort
		portKey := port.Type
		if port.IP != "" {
			if port.PublicPort != current {
				hostMappings = append(hostMappings, fmt.Sprintf("%s:%d->%d/%s", port.IP, port.PublicPort, port.PrivatePort, port.Type))
				continue
			}
			portKey = fmt.Sprintf("%s/%s", port.IP, port.Type)
		}
		group := groupMap[portKey]

		if group == nil {
			groupMap[portKey] = &portGroup{first: current, last: current}
			// record order that groupMap keys are created
			groupMapKeys = append(groupMapKeys, portKey)
			continue
		}
		if current == (group.last + 1) {
			group.last = current
			continue
		}

		result = append(result, formGroup(portKey, group.first, group.last))
		groupMap[portKey] = &portGroup{first: current, last: current}
	}
	for _, portKey := range groupMapKeys {
		g := groupMap[portKey]
		result = append(result, formGroup(portKey, g.first, g.last))
	}
	result = append(result, hostMappings...)
	return strings.Join(result, ", ")
}

func formGroup(key string, start, last uint16) string {
	parts := strings.Split(key, "/")
	groupType := parts[0]
	var ip string
	if len(parts) > 1 {
		ip = parts[0]
		groupType = parts[1]
	}
	group := strconv.Itoa(int(start))
	if start != last {
		group = fmt.Sprintf("%s-%d", group, last)
	}
	if ip != "" {
		group = fmt.Sprintf("%s:%s->%s", ip, group, group)
	}
	return fmt.Sprintf("%s/%s", group, groupType)
}
