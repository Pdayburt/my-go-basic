package client

import (
	"context"
	intrv1 "example.com/mod/webook/api/proto/gen/intr/v1"
	"example.com/mod/webook/interactive/domain"
	"example.com/mod/webook/interactive/service"
	"google.golang.org/grpc"
)

type InteractiveServiceAdapter struct {
	svc service.InteractiveService
}

func NewInteractiveServiceAdapter(svc service.InteractiveService) *InteractiveServiceAdapter {
	return &InteractiveServiceAdapter{svc: svc}
}

func (i *InteractiveServiceAdapter) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntReq, opts ...grpc.CallOption) (*intrv1.IncrReadCntResp, error) {
	//ctx context.Context, biz string, bizId int64
	err := i.svc.IncrReadCnt(ctx, in.GetBiz(), in.GetBizId())
	return &intrv1.IncrReadCntResp{}, err
}

func (i *InteractiveServiceAdapter) Like(ctx context.Context, in *intrv1.LikeReq, opts ...grpc.CallOption) (*intrv1.LikeReq, error) {
	err := i.svc.Like(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.LikeReq{}, err
}

func (i *InteractiveServiceAdapter) CancelLike(ctx context.Context, in *intrv1.CancelLikeReq, opts ...grpc.CallOption) (*intrv1.CancelLikeResp, error) {

	err := i.svc.CancelLike(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.CancelLikeResp{}, err
}

func (i *InteractiveServiceAdapter) Collect(ctx context.Context, in *intrv1.CollectReq, opts ...grpc.CallOption) (*intrv1.CollectResp, error) {

	err := i.svc.Collect(ctx, in.GetBiz(), in.GetBizId(), in.GetCid(), in.GetUid())
	return &intrv1.CollectResp{}, err
}

func (i *InteractiveServiceAdapter) Get(ctx context.Context, in *intrv1.GetReq, opts ...grpc.CallOption) (*intrv1.GetResp, error) {

	interactive, err := i.svc.Get(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.GetResp{Intr: i.toDTO(interactive)}, err
}

func (i *InteractiveServiceAdapter) GetByIds(ctx context.Context, in *intrv1.GetByIdsReq, opts ...grpc.CallOption) (*intrv1.GetByIdsResp, error) {
	interactives, err := i.svc.GetByIds(ctx, in.GetBiz(), in.GetIds())
	if err != nil {
		return &intrv1.GetByIdsResp{}, err
	}
	m := make(map[int64]*intrv1.Interactive, len(interactives))
	for k, v := range interactives {
		m[k] = &intrv1.Interactive{
			Biz:        v.Biz,
			BizId:      v.BizId,
			ReadCnt:    v.ReadCnt,
			LikeCnt:    v.LikeCnt,
			CollectCnt: v.CollectCnt,
			Liked:      v.Liked,
			Collected:  v.Collected,
		}
	}
	return &intrv1.GetByIdsResp{Intrs: m}, nil
}

func (i *InteractiveServiceAdapter) toDTO(intr domain.Interactive) *intrv1.Interactive {
	return &intrv1.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		ReadCnt:    intr.ReadCnt,
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		Liked:      intr.Liked,
		Collected:  intr.Collected,
	}
}
