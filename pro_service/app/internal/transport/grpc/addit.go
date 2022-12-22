package grpc

import (
	"context"

	"github.com/Alexander272/sealur_proto/api/pro_api"
)

func (h *Handler) GetAllAdditional(ctx context.Context, dto *pro_api.GetAllAddRequest) (*pro_api.AdditionalResponse, error) {
	addit, err := h.service.Addit.GetAll()
	if err != nil {
		return nil, err
	}

	return &pro_api.AdditionalResponse{Additionals: addit}, nil
}

func (h *Handler) CreateAdditional(ctx context.Context, dto *pro_api.CreateAddRequest) (*pro_api.SuccessResponse, error) {
	success, err := h.service.Addit.Create(dto)
	if err != nil {
		return nil, err
	}

	return success, nil
}

func (h *Handler) UpdateMat(ctx context.Context, dto *pro_api.UpdateAddMatRequest) (*pro_api.SuccessResponse, error) {
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

		if err := h.service.SNP.DeleteMat(dto.Change, addit[0].Materials); err != nil {
			return nil, err
		}

		if err := h.service.Putg.DeleteMat(dto.Change, addit[0].Materials); err != nil {
			return nil, err
		}

		if err := h.service.Putgm.DeleteMat(dto.Change, addit[0].Materials); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdateMod(ctx context.Context, dto *pro_api.UpdateAddModRequest) (*pro_api.SuccessResponse, error) {
	success, err := h.service.Addit.UpdateMod(dto)
	if err != nil {
		return nil, err
	}

	if dto.TypeCh == "delete" {
		if err := h.service.SNP.DeleteMod(dto.Change); err != nil {
			return nil, err
		}

		if err := h.service.Putg.DeleteMod(dto.Change); err != nil {
			return nil, err
		}

		if err := h.service.Putgm.DeleteMod(dto.Change); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdateTemp(ctx context.Context, dto *pro_api.UpdateAddTemRequest) (*pro_api.SuccessResponse, error) {
	success, err := h.service.Addit.UpdateTemp(dto)
	if err != nil {
		return nil, err
	}

	if dto.TypeCh == "delete" {
		if err := h.service.SNP.DeleteTemp(dto.Change); err != nil {
			return nil, err
		}

		if err := h.service.Putg.DeleteTemp(dto.Change); err != nil {
			return nil, err
		}

		if err := h.service.Putgm.DeleteTemp(dto.Change); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdateMoun(ctx context.Context, dto *pro_api.UpdateAddMounRequest) (*pro_api.SuccessResponse, error) {
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

		if err := h.service.Putg.DeleteMoun(dto.Change); err != nil {
			return nil, err
		}

		if err := h.service.Putgm.DeleteMoun(dto.Change); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdateGrap(ctx context.Context, dto *pro_api.UpdateAddGrapRequest) (*pro_api.SuccessResponse, error) {
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

		if err := h.service.Putg.DeleteGrap(dto.Change); err != nil {
			return nil, err
		}

		if err := h.service.Putgm.DeleteGrap(dto.Change); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdateFillers(ctx context.Context, dto *pro_api.UpdateAddFillersRequest) (*pro_api.SuccessResponse, error) {
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

func (h *Handler) UpdateCoating(ctx context.Context, dto *pro_api.UpdateAddCoatingRequest) (*pro_api.SuccessResponse, error) {
	success, err := h.service.UpdateCoating(dto)
	if err != nil {
		return nil, err
	}

	if dto.TypeCh == "delete" {
		if err := h.service.Putg.DeleteCoating(dto.Change); err != nil {
			return nil, err
		}

		if err := h.service.Putgm.DeleteCoating(dto.Change); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdateConstruction(ctx context.Context, dto *pro_api.UpdateAddConstructionRequest) (*pro_api.SuccessResponse, error) {
	success, err := h.service.UpdateConstruction(dto)
	if err != nil {
		return nil, err
	}

	if dto.TypeCh == "delete" {
		if err := h.service.Putg.DeleteCon(dto.Change); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdateObturator(ctx context.Context, dto *pro_api.UpdateAddObturatorRequest) (*pro_api.SuccessResponse, error) {
	success, err := h.service.UpdateObturator(dto)
	if err != nil {
		return nil, err
	}

	if dto.TypeCh == "delete" {
		if err := h.service.Putg.DeleteObt(dto.Change); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdateBasis(ctx context.Context, dto *pro_api.UpdateAddBasisRequest) (*pro_api.SuccessResponse, error) {
	success, err := h.service.UpdateBasis(dto)
	if err != nil {
		return nil, err
	}

	if dto.TypeCh == "delete" {
		if err := h.service.Putgm.DeleteCon(dto.Change); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdatePObturator(ctx context.Context, dto *pro_api.UpdateAddPObturatorRequest) (*pro_api.SuccessResponse, error) {
	success, err := h.service.UpdatePObturator(dto)
	if err != nil {
		return nil, err
	}

	if dto.TypeCh == "delete" {
		if err := h.service.Putgm.DeleteObt(dto.Change); err != nil {
			return nil, err
		}
	}

	return success, nil
}

func (h *Handler) UpdateSealant(ctx context.Context, dto *pro_api.UpdateAddSealantRequest) (*pro_api.SuccessResponse, error) {
	success, err := h.service.UpdateSealant(dto)
	if err != nil {
		return nil, err
	}

	if dto.TypeCh == "delete" {
		if err := h.service.Putgm.DeleteSeal(dto.Change); err != nil {
			return nil, err
		}
	}

	return success, nil
}
