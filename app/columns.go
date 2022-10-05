package main

type consumer struct {
	hostname, ip, ipl, fqdn, appgroup string
}

type provider struct {
	hostname, ip, ipl, fqdn, appgroup string
}

type info struct {
	transmission, port, protocol          string
	num_flows, conn_state, first_detected string
	last_detected                         string
}

type reported struct {
	policy_decision, enforcement_boundary, by string
}

type draft struct {
	policy_decision, enforcement_boundary string
}

type record struct {
	c consumer
	p provider
	r reported
	d draft
	i info
}
