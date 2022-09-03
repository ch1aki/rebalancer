package provider

import (
	"context"
	"fmt"

	rebalancerv1 "git.pepabo.com/akichan/rebalancer/api/v1"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

type Route53Provider struct {
	hostedZoneId string
	recordName   string
	recordId     string
	client       *route53.Client
	rr           types.ResourceRecordSet
}

func NewProvider(ctx context.Context, r rebalancerv1.Rebalance) (Route53Provider, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return Route53Provider{}, err
	}

	p := Route53Provider{
		hostedZoneId: r.Spec.Target.Route53.HostedZoneID,
		recordName:   r.Spec.Target.Route53.Resource.Name,
		recordId:     r.Spec.Target.Route53.Resource.Identifier,
		client:       route53.NewFromConfig(cfg),
	}

	p.fetchResourceRecordSets(ctx)
	if err != nil {
		return Route53Provider{}, err
	}

	return p, nil
}

func (p *Route53Provider) GetWeight(ctx context.Context) (int64, error) {
	return *p.rr.Weight, nil
}

func (p *Route53Provider) SetWeight(ctx context.Context, value int64) error {
	p.rr.Weight = aws.Int64(value)
	changes := route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &p.hostedZoneId,
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{Action: types.ChangeActionUpsert, ResourceRecordSet: &p.rr},
			},
		},
	}
	_, err := p.client.ChangeResourceRecordSets(ctx, &changes)
	if err != nil {
		return err
	}

	err = p.fetchResourceRecordSets(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (p *Route53Provider) fetchResourceRecordSets(ctx context.Context) error {
	out, err := p.client.ListResourceRecordSets(ctx,
		&route53.ListResourceRecordSetsInput{
			HostedZoneId:          aws.String(p.hostedZoneId),
			StartRecordName:       aws.String(p.recordName),
			StartRecordIdentifier: aws.String(p.recordId),
		},
	)
	if err != nil {
		return err
	}
	for _, rr := range out.ResourceRecordSets {
		if *rr.Name == p.recordName && *rr.SetIdentifier == p.recordId {
			p.rr = rr
			return nil
		}
	}
	return fmt.Errorf("resource record not found")
}
