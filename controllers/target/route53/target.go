package route53

import (
	"context"
	"fmt"
	"strings"

	rebalancerv1 "git.pepabo.com/akichan/rebalancer/api/v1"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Target struct {
	hostedZoneId string
	recordName   string
	recordId     string
	recordType   types.RRType
	client       *route53.Client
	rr           types.ResourceRecordSet
}

func (t *Target) NewClient(ctx context.Context, r rebalancerv1.Rebalance, c client.Client) (rebalancerv1.TargetClient, error) {
	var optFns []func(*config.LoadOptions) error

	// secret ref option
	if r.Spec.Target.Route53.Auth.SecretRef != nil {
		cred, err := credFromSecretRef(ctx, &r, c)
		if err != nil {
			return nil, err
		}
		optFns = append(optFns, config.WithCredentialsProvider(cred))
	}

	// region option
	if region := r.Spec.Target.Route53.Region; region == "" {
		return nil, fmt.Errorf("route53 target require region")
	} else {
		optFns = append(optFns, config.WithRegion(region))
	}

	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return nil, fmt.Errorf("load config error: %w", err)
	}

	return &Target{
		hostedZoneId: r.Spec.Target.Route53.HostedZoneID,
		recordName:   r.Spec.Target.Route53.Resource.Name,
		recordId:     r.Spec.Target.Route53.Resource.Identifier,
		recordType:   r.Spec.Target.Route53.Resource.Type,
		client:       route53.NewFromConfig(cfg),
	}, nil
}

func credFromSecretRef(ctx context.Context, r *rebalancerv1.Rebalance, c client.Client) (credentials.StaticCredentialsProvider, error) {
	secRef := r.Spec.Target.Route53.Auth.SecretRef

	// get access key id from secret
	var ns string
	if secRef.AccessKeyID.Namespace != nil {
		ns = *secRef.AccessKeyID.Namespace
	} else {
		ns = r.Namespace
	}
	ke := client.ObjectKey{
		Name:      secRef.AccessKeyID.Name,
		Namespace: ns,
	}
	akSecret := v1.Secret{}
	err := c.Get(ctx, ke, &akSecret)
	if err != nil {
		return credentials.StaticCredentialsProvider{}, fmt.Errorf("failed to get access key id: %w", err)
	}

	// get secret access key from secret
	if secRef.SecretAccessKey.Namespace != nil {
		ns = *secRef.SecretAccessKey.Namespace
	} else {
		ns = r.Namespace
	}
	ke = client.ObjectKey{
		Name:      secRef.SecretAccessKey.Name,
		Namespace: ns,
	}
	sakSecret := v1.Secret{}
	err = c.Get(ctx, ke, &sakSecret)
	if err != nil {
		return credentials.StaticCredentialsProvider{}, fmt.Errorf("failed to get secret access key: %w", err)
	}

	ak := string(akSecret.Data[secRef.AccessKeyID.Key])
	sak := string(sakSecret.Data[secRef.SecretAccessKey.Key])
	if ak == "" {
		return credentials.StaticCredentialsProvider{}, fmt.Errorf("missing access key id")
	}
	if sak == "" {
		return credentials.StaticCredentialsProvider{}, fmt.Errorf("missing secret access key")
	}
	return credentials.NewStaticCredentialsProvider(ak, sak, ""), nil
}

func (t *Target) GetWeight(ctx context.Context) (int64, error) {
	err := t.fetchResourceRecordSets(ctx)
	if err != nil {
		return 0, err
	}
	return *t.rr.Weight, nil
}

func (p *Target) SetWeight(ctx context.Context, value int64) error {
	err := p.fetchResourceRecordSets(ctx)
	if err != nil {
		return err
	}
	p.rr.Weight = aws.Int64(value)
	changes := route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &p.hostedZoneId,
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{Action: types.ChangeActionUpsert, ResourceRecordSet: &p.rr},
			},
		},
	}
	_, err = p.client.ChangeResourceRecordSets(ctx, &changes)
	if err != nil {
		return err
	}

	return nil
}

func (t *Target) fetchResourceRecordSets(ctx context.Context) error {
	out, err := t.client.ListResourceRecordSets(ctx,
		&route53.ListResourceRecordSetsInput{
			HostedZoneId:          aws.String(t.hostedZoneId),
			StartRecordName:       aws.String(t.recordName),
			StartRecordIdentifier: aws.String(t.recordId),
			StartRecordType:       t.recordType,
		},
	)
	if err != nil {
		return err
	}

	rname := t.recordName
	if !strings.HasSuffix(".", rname) {
		rname = rname + "."
	}
	for _, rr := range out.ResourceRecordSets {
		if *rr.Name == rname && *rr.SetIdentifier == t.recordId {
			t.rr = rr
			return nil
		}
	}
	return fmt.Errorf("resource record not found")
}

func init() {
	rebalancerv1.RegisterTarget(&Target{}, &rebalancerv1.RebalanceTarget{
		Route53: &rebalancerv1.Route53Target{},
	})
}
