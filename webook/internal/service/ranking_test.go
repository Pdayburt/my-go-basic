package service

import (
	"context"
	domain2 "example.com/mod/webook/interactive/domain"
	"example.com/mod/webook/interactive/service"
	"example.com/mod/webook/internal/domain"
	"example.com/mod/webook/internal/service/svcmocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestBatchRankingService_TopN(t *testing.T) {
	now := time.Now()
	testCase := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (ArticleService, service.InteractiveService)
		wantErr error
		wantRes []domain.Article
	}{
		{
			name: "计算成功",
			//怎么模拟数据？
			mock: func(ctrl *gomock.Controller) (ArticleService, service.InteractiveService) {

				artSvc := svcmocks.NewMockArticleService(ctrl)
				artSvc.EXPECT().ListPub(gomock.Any(), 0, 3).
					Return([]domain.Article{
						{Id: 1, Utime: now, Ctime: now},
						{Id: 2, Utime: now, Ctime: now},
						{Id: 3, Utime: now, Ctime: now},
					}, nil)

				intrSvc := svcmocks.NewMockInteractiveService(ctrl)
				intrSvc.EXPECT().GetByIds(gomock.Any(),
					"article", []int64{1, 2, 3}).
					Return([]domain2.Interactive{
						{BizId: 1, LikeCnt: 1},
						{BizId: 2, LikeCnt: 2},
						{BizId: 3, LikeCnt: 3},
					}, nil)
				return artSvc, intrSvc
			},
			wantRes: []domain.Article{
				{Id: 3, Utime: now, Ctime: now},
				{Id: 2, Utime: now, Ctime: now},
				{Id: 1, Utime: now, Ctime: now},
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			artSvc, intrSvc := tc.mock(ctrl)
			rankingSvc := NewBatchRankingService(artSvc, intrSvc)
			articles, err := rankingSvc.topN(context.Background())
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRes, articles)

		})
	}
}
