package aws

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
)

var svcELB = elb.New(session.New())
var svcEC2 = ec2.New(session.New())

// Status ...
type Status struct {
	config Config

	Output *StatusOutput
}

// StatusOutput ...
type StatusOutput struct {
	LoadBalancer *LoadBalancerStatus
	VPNTunnel    *VPNStatus
}

// LoadBalancerStatus ...
type LoadBalancerStatus struct {
	Healthy   int
	Unhealthy int
}

// VPNStatus ...
type VPNStatus struct {
	Tunnel1 bool
	Tunnel2 bool
}

// Config ...
type Config struct {
	LoadBalancer string `yaml:"load_balancer"`
	VPN          string `yaml:"vpn"`
}

// NewStatus ...
func NewStatus(config Config) *Status {
	return &Status{
		config: config,
		Output: &StatusOutput{
			LoadBalancer: &LoadBalancerStatus{},
			VPNTunnel:    &VPNStatus{},
		},
	}
}

// NewFromFile ...
func NewFromFile(file string) (*Status, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var c Config
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	log.Printf("%+v", c)

	return NewStatus(c), nil
}

func (s *Status) getELBStatus() {
	params := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{
			aws.String(s.config.LoadBalancer),
		},
	}
	resp, err := svcELB.DescribeLoadBalancers(params)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	params2 := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: aws.String(s.config.LoadBalancer),
		Instances:        resp.LoadBalancerDescriptions[0].Instances,
	}
	resp2, err := svcELB.DescribeInstanceHealth(params2)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	s.Output.LoadBalancer.Healthy = 0
	s.Output.LoadBalancer.Unhealthy = 0
	for _, state := range resp2.InstanceStates {
		if *state.State == "InService" {
			s.Output.LoadBalancer.Healthy++
		} else {
			s.Output.LoadBalancer.Unhealthy++
		}
	}
}

func (s *Status) getVPNStatus() {
	params := &ec2.DescribeVpnConnectionsInput{
		VpnConnectionIds: []*string{
			aws.String(s.config.VPN),
		},
	}
	resp, err := svcEC2.DescribeVpnConnections(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	if *resp.VpnConnections[0].VgwTelemetry[0].Status == "UP" {
		s.Output.VPNTunnel.Tunnel1 = true
	} else {
		s.Output.VPNTunnel.Tunnel1 = false
	}
	if *resp.VpnConnections[0].VgwTelemetry[1].Status == "UP" {
		s.Output.VPNTunnel.Tunnel2 = true
	} else {
		s.Output.VPNTunnel.Tunnel2 = false
	}
}

// Update updates the AWS status
func (s *Status) Update() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.getELBStatus()
	}()

	go func() {
		defer wg.Done()
		s.getVPNStatus()
	}()

	wg.Wait()
}
