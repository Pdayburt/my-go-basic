package webookgrpc

import (
	"context"
	"example.com/mod/webook/api/proto/gen/intr/v1"
	"example.com/mod/webook/interactive/domain"
	"example.com/mod/webook/interactive/service"
	"google.golang.org/grpc"
)

// InteractiveServiceServer 这里只是把server包装成一个grpc而已
//和grpc相关的操作都限定在这里，

type InteractiveServiceServer struct {
	//真正的业务逻辑一定在service里
	svc service.InteractiveService
	intrv1.UnimplementedInteractiveServiceServer
}

func NewInteractiveServiceServer(svc service.InteractiveService) *InteractiveServiceServer {
	return &InteractiveServiceServer{svc: svc}
}

func (i *InteractiveServiceServer) Register(server *grpc.Server) {
	intrv1.RegisterInteractiveServiceServer(server, i)
}

func (i *InteractiveServiceServer) IncrReadCnt(ctx context.Context, req *intrv1.IncrReadCntReq) (*intrv1.IncrReadCntResp, error) {
	err := i.svc.IncrReadCnt(ctx, req.GetBiz(), req.GetBizId())
	return &intrv1.IncrReadCntResp{}, err
}

func (i *InteractiveServiceServer) Like(ctx context.Context, req *intrv1.LikeReq) (*intrv1.LikeReq, error) {
	err := i.svc.Like(ctx, req.GetBiz(), req.GetBizId(), req.GetUid())
	return &intrv1.LikeReq{}, err
}

func (i *InteractiveServiceServer) CancelLike(ctx context.Context, req *intrv1.CancelLikeReq) (*intrv1.CancelLikeResp, error) {
	err := i.svc.CancelLike(ctx, req.GetBiz(), req.GetBizId(), req.GetUid())
	return &intrv1.CancelLikeResp{}, err
}

func (i *InteractiveServiceServer) Collect(ctx context.Context, req *intrv1.CollectReq) (*intrv1.CollectResp, error) {

	err := i.svc.Collect(ctx, req.GetBiz(), req.GetBizId(), req.GetCid(), req.GetUid())
	return &intrv1.CollectResp{}, err
}

func (i *InteractiveServiceServer) Get(ctx context.Context, req *intrv1.GetReq) (*intrv1.GetResp, error) {

	interactive, err := i.svc.Get(ctx, req.GetBiz(), req.GetBizId(), req.GetUid())
	return &intrv1.GetResp{Intr: i.toDTO(interactive)}, err
}

func (i *InteractiveServiceServer) GetByIds(ctx context.Context, req *intrv1.GetByIdsReq) (*intrv1.GetByIdsResp, error) {
	interactives, err := i.svc.GetByIds(ctx, req.GetBiz(), req.GetIds())
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

func (i *InteractiveServiceServer) toDTO(intr domain.Interactive) *intrv1.Interactive {
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
