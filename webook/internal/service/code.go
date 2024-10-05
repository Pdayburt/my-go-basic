package service

import (
	"context"
	"example.com/mod/webook/internal/repository"
	"example.com/mod/webook/internal/service/sms"
	"math/rand"

	"fmt"
)

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, inputCode, phone string) (bool, error)
}

type CodeServiceByRedisAndTPS struct {
	//repo   *repository.CacheCodeRepository
	repo   repository.CodeRepository
	smsSvc sms.Service
	tplId  string
}

func NewCodeService(repo repository.CodeRepository, svc sms.Service) CodeService {
	return &CodeServiceByRedisAndTPS{
		repo:   repo,
		smsSvc: svc,
		tplId:  "99999",
	}
}

func (cs *CodeServiceByRedisAndTPS) Send(ctx context.Context, biz, phone string) error {
	//phone_code:$biz:$phone
	code := cs.generateCode()
	err := cs.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = cs.smsSvc.Send(ctx, cs.tplId, []string{code}, phone)
	return err
}

func (cs *CodeServiceByRedisAndTPS) Verify(ctx context.Context, biz, inputCode, phone string) (bool, error) {
	//ctx, biz, phone, inputCode
	return cs.repo.Verify(ctx, biz, phone, inputCode)
}

func (cs *CodeServiceByRedisAndTPS) generateCode() string {
	intn := rand.Intn(1000000)
	return fmt.Sprintf("%06d", intn)
}
