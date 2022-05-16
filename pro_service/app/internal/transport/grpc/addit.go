package grpc

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

func (h *Handler) GetAllAdditional(ctx context.Context, dto *proto.GetAllAddRequest) (*proto.AdditionalResponse, error) {
	addit, err := h.service.Addit.GetAll()
	if err != nil {
		return nil, err
	}

	return &proto.AdditionalResponse{Additionals: addit}, nil
}

func (h *Handler) CreateAdditional(ctx context.Context, dto *proto.CreateAddRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.Addit.Create(dto)
	if err != nil {
		return nil, err
	}

	return success, nil
}

func (h *Handler) UpdateMat(ctx context.Context, dto *proto.UpdateAddMatRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.Addit.UpdateMat(dto)
	if err != nil {
		return nil, err
	}

	if dto.TypeCh == "add" {
		if err := h.service.SNP.AddMat(dto.Change); err != nil {
			return nil, err
		}
	}

	if dto.TypeCh == "delete" {
		addit, err := h.service.Addit.GetAll()
		if err != nil {
			return nil, err
		}

		//TODO исправить удаление материалов
		if err := h.service.SNP.DeleteMat(dto.Change, addit[0].Materials); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdateMod(ctx context.Context, dto *proto.UpdateAddModRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.Addit.UpdateMod(dto)
	if err != nil {
		return nil, err
	}

	if dto.TypeCh == "delete" {
		if err := h.service.SNP.DeleteMod(dto.Change); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdateTemp(ctx context.Context, dto *proto.UpdateAddTemRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.Addit.UpdateTemp(dto)
	if err != nil {
		return nil, err
	}

	if dto.TypeCh == "delete" {
		if err := h.service.SNP.DeleteTemp(dto.Change); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdateMoun(ctx context.Context, dto *proto.UpdateAddMounRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.Addit.UpdateMoun(dto)
	if err != nil {
		return nil, err
	}

	if dto.TypeCh == "add" {
		if err := h.service.SNP.AddMoun(dto.Change); err != nil {
			return nil, err
		}
	}

	if dto.TypeCh == "delete" {
		if err := h.service.SNP.DeleteMoun(dto.Change); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdateGrap(ctx context.Context, dto *proto.UpdateAddGrapRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.UpdateGrap(dto)
	if err != nil {
		return nil, err
	}

	if dto.TypeCh == "add" {
		if err := h.service.SNP.AddGrap(dto.Change); err != nil {
			return nil, err
		}
	}

	if dto.TypeCh == "delete" {
		if err := h.service.SNP.DeleteGrap(dto.Change); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdateFillers(ctx context.Context, dto *proto.UpdateAddFillersRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.UpdateFillers(dto)
	if err != nil {
		return nil, err
	}

	if dto.TypeCh == "delete" {
		if err := h.service.SNP.DeleteFiller(dto.Change); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdateCoating(ctx context.Context, dto *proto.UpdateAddCoatingRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.UpdateCoating(dto)
	if err != nil {
		return nil, err
	}

	// TODO дописать удаление
	// if dto.TypeCh == "delete" {
	// if err := h.service.SNP.DeleteFiller(dto.Change); err != nil {
	// 	return nil, err
	// }
	// }

	return success, nil
}

func (h *Handler) UpdateConstruction(ctx context.Context, dto *proto.UpdateAddConstructionRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.UpdateConstruction(dto)
	if err != nil {
		return nil, err
	}

	// TODO дописать удаление
	// if dto.TypeCh == "delete" {
	// if err := h.service.SNP.DeleteFiller(dto.Change); err != nil {
	// 	return nil, err
	// }
	// }

	return success, nil
}

func (h *Handler) UpdateObturator(ctx context.Context, dto *proto.UpdateAddObturatorRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.UpdateObturator(dto)
	if err != nil {
		return nil, err
	}

	// TODO дописать удаление
	// if dto.TypeCh == "delete" {
	// if err := h.service.SNP.DeleteFiller(dto.Change); err != nil {
	// 	return nil, err
	// }
	// }

	return success, nil
}

func (h *Handler) UpdateBasis(ctx context.Context, dto *proto.UpdateAddBasisRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.UpdateBasis(dto)
	if err != nil {
		return nil, err
	}

	// TODO дописать удаление
	// if dto.TypeCh == "delete" {
	// if err := h.service.SNP.DeleteFiller(dto.Change); err != nil {
	// 	return nil, err
	// }
	// }

	return success, nil
}

func (h *Handler) UpdateSealant(ctx context.Context, dto *proto.UpdateAddSealantRequest) (*proto.SuccessResponse, error) {
	success, err := h.service.UpdateSealant(dto)
	if err != nil {
		return nil, err
	}

	// TODO дописать удаление
	// if dto.TypeCh == "delete" {
	// if err := h.service.SNP.DeleteFiller(dto.Change); err != nil {
	// 	return nil, err
	// }
	// }

	return success, nil
}
