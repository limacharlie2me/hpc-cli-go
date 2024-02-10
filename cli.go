package main

var cli struct { //nolint
	Region               string `kong:"flag,name='region',short='r',required,placeholder='REGION',help='Example: ad1.us-ashburn-1'"`
	ClusterNetworkOCID   string `kong:"flag,group='Search Term',xor='cluster-network,hpc-island-id,bare-metal-cluster,rack-id,instance-id,list-all',name='cluster-network',short='c',placeholder='CLUSTER-NETWORK',help='Cluster Network OCID e.g ocid1.clusternetwork.oc1.phx.bla'"`
	HpcIslandID          string `kong:"flag,group='Search Term',xor='cluster-network,hpc-island-id,bare-metal-cluster,rack-id,instance-id,list-all',name='hpc-island-id',short='d',placeholder='HPC-ISLAND-ID',help='HPC Island block ID e.g bldg1-block11'"`
	BareMetalClusterOCID string `kong:"flag,group='Search Term',xor='cluster-network,hpc-island-id,bare-metal-cluster,rack-id,instance-id,list-all',name='bare-metal-cluster',short='m',placeholder='BARE-METAL-CLUSTER',help='Bare Metal Cluster OCID e.g ocid1.baremetalcluster.oc1.phx.bla'"`
	RackSKID             string `kong:"flag,group='Search Term',xor='cluster-network,hpc-island-id,bare-metal-cluster,rack-id,instance-id,list-all',name='rack-id',short='s',placeholder='RACK-ID',help='Storekeeper Rack ID e.g sk-blalbla'"`
	InstanceOCID         string `kong:"flag,group='Search Term',xor='cluster-network,hpc-island-id,bare-metal-cluster,rack-id,instance-id,list-all',name='instance-id',short='i',placeholder='INSTANCE-ID',help='Instance OCID e.g ocid1.instance.oc1.phx.bla'"`
	ListALL              bool   `kong:"flag,group='Search Term',xor='cluster-network,hpc-island-id,bare-metal-cluster,rack-id,instance-id,list-all',name='list-all',short='l',help='Lists all HPC instances within region'"`
	ListBroken           bool   `kong:"flag,name='list-broken-hosts',short='b',help='Returns broken hosts only'"`
	HealthCheck          bool   `kong:"flag,name='run-health-check',short='p',help='Runs Pulse health check (ilom/smartnic)'"`
	Display              bool   `kong:"flag,name='show-table',short='z',help='Displays the data in a table within the terminal window'"`
}
