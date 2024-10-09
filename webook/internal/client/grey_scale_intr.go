package client

import (
	"context"
	intrv1 "example.com/mod/webook/api/proto/gen/intr/v1"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"google.golang.org/grpc"
	"math/rand"
)

type GreyScaleInteractiveServiceClient struct {
	remote intrv1.InteractiveServiceClient
	local  intrv1.InteractiveServiceClient
	//使用随机数控制阈值 控制是 远程还是本地
	threshold *atomicx.Value[int32]
}

func NewGreyScaleInteractiveServiceClient(remote intrv1.InteractiveServiceClient,
	local intrv1.InteractiveServiceClient) *GreyScaleInteractiveServiceClient {
	return &GreyScaleInteractiveServiceClient{
		remote:    remote,
		local:     local,
		threshold: atomicx.NewValue[int32](),
	}
}

func (g *GreyScaleInteractiveServiceClient) client() intrv1.InteractiveServiceClient {
	threshold := g.threshold.Load()
	num := rand.Int31n(100)
	if num < threshold {
		return g.remote
	}
	return g.local
}

func (g *GreyScaleInteractiveServiceClient) UpdateThreshold(newThreshold int32) {

	g.threshold.Store(newThreshold)
}

func (g *GreyScaleInteractiveServiceClient) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntReq, opts ...grpc.CallOption) (*intrv1.IncrReadCntResp, error) {
	return g.client().IncrReadCnt(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Like(ctx context.Context, in *intrv1.LikeReq, opts ...grpc.CallOption) (*intrv1.LikeReq, error) {
	return g.client().Like(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) CancelLike(ctx context.Context, in *intrv1.CancelLikeReq, opts ...grpc.CallOption) (*intrv1.CancelLikeResp, error) {
	return g.client().CancelLike(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Collect(ctx context.Context, in *intrv1.CollectReq, opts ...grpc.CallOption) (*intrv1.CollectResp, error) {
	return g.client().Collect(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Get(ctx context.Context, in *intrv1.GetReq, opts ...grpc.CallOption) (*intrv1.GetResp, error) {
	return g.client().Get(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) GetByIds(ctx context.Context, in *intrv1.GetByIdsReq, opts ...grpc.CallOption) (*intrv1.GetByIdsResp, error) {
	return g.client().GetByIds(ctx, in, opts...)
}
