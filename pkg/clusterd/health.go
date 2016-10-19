package clusterd

import "github.com/quantum/castle/pkg/clusterd/inventory"

const (
	unhealthyNodeSecondsThreshold = 10
)

func IsNodeUnhealthy(node *inventory.NodeConfig) (int, bool) {
	age := int(node.HeartbeatAge.Seconds())
	return age, age >= unhealthyNodeSecondsThreshold
}