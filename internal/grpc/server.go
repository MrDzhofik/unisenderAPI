package grpc

import (
	"context"
	"myAwesomeProject/internal/usecase"
	"myAwesomeProject/proto/accountpb"
)

type AccountServer struct {
	accountpb.UnimplementedAccountServiceServer
	accountUsecase usecase.AccountUsecase
}

func NewAccountServer(u usecase.AccountUsecase) *AccountServer {
	return &AccountServer{accountUsecase: u}
}

func (s *AccountServer) DeleteAccountByID(ctx context.Context, req *accountpb.DeleteAccountRequest) (*accountpb.DeleteAccountResponse, error) {
	accountID := req.GetAccountId()

	err := s.accountUsecase.DeleteAccountByID(accountID)
	if err != nil {
		return &accountpb.DeleteAccountResponse{
			Message: "Ошибка удаления аккаунта: " + err.Error(),
			Success: false,
		}, nil
	}

	return &accountpb.DeleteAccountResponse{
		Message: "Аккаунт успешно удален",
		Success: true,
	}, nil
}
