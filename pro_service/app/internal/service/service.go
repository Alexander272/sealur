package service

import (
	"bytes"
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"github.com/Alexander272/sealur_proto/api/file_api"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	"github.com/Alexander272/sealur_proto/api/user_api"
)

type Stand interface {
	GetAll(*pro_api.GetStandsRequest) (stands []*pro_api.Stand, err error)
	Create(*pro_api.CreateStandRequest) (st *pro_api.IdResponse, err error)
	Update(*pro_api.UpdateStandRequest) error
	Delete(*pro_api.DeleteStandRequest) error
}

type Flange interface {
	GetAll() ([]*pro_api.Flange, error)
	Create(*pro_api.CreateFlangeRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdateFlangeRequest) error
	Delete(*pro_api.DeleteFlangeRequest) error
}

type StFl interface {
	Get() ([]*pro_api.StFl, error)
	Create(*pro_api.CreateStFlRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdateStFlRequest) error
	Delete(*pro_api.DeleteStFlRequest) error
}

type TypeFl interface {
	Get() ([]*pro_api.TypeFl, error)
	GetAll() ([]*pro_api.TypeFl, error)
	Create(*pro_api.CreateTypeFlRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdateTypeFlRequest) error
	Delete(*pro_api.DeleteTypeFlRequest) error
}

type Addit interface {
	GetAll() ([]*pro_api.Additional, error)
	Create(*pro_api.CreateAddRequest) (*pro_api.SuccessResponse, error)
	UpdateMat(*pro_api.UpdateAddMatRequest) (*pro_api.SuccessResponse, error)
	UpdateMod(*pro_api.UpdateAddModRequest) (*pro_api.SuccessResponse, error)
	UpdateTemp(*pro_api.UpdateAddTemRequest) (*pro_api.SuccessResponse, error)
	UpdateMoun(*pro_api.UpdateAddMounRequest) (*pro_api.SuccessResponse, error)
	UpdateGrap(*pro_api.UpdateAddGrapRequest) (*pro_api.SuccessResponse, error)
	UpdateFillers(*pro_api.UpdateAddFillersRequest) (*pro_api.SuccessResponse, error)
	UpdateCoating(*pro_api.UpdateAddCoatingRequest) (*pro_api.SuccessResponse, error)
	UpdateConstruction(*pro_api.UpdateAddConstructionRequest) (*pro_api.SuccessResponse, error)
	UpdateObturator(*pro_api.UpdateAddObturatorRequest) (*pro_api.SuccessResponse, error)
	UpdateBasis(*pro_api.UpdateAddBasisRequest) (*pro_api.SuccessResponse, error)
	UpdatePObturator(*pro_api.UpdateAddPObturatorRequest) (*pro_api.SuccessResponse, error)
	UpdateSealant(*pro_api.UpdateAddSealantRequest) (*pro_api.SuccessResponse, error)
}

type Size interface {
	Get(*pro_api.GetSizesRequest) ([]*pro_api.Size, []*pro_api.Dn, error)
	GetAll(*pro_api.GetSizesRequest) ([]*pro_api.Size, []*pro_api.Dn, error)
	Create(*pro_api.CreateSizeRequest) (*pro_api.IdResponse, error)
	CreateMany(*pro_api.CreateSizesRequest) error
	Update(*pro_api.UpdateSizeRequest) error
	Delete(*pro_api.DeleteSizeRequest) error
	DeleteAll(*pro_api.DeleteAllSizeRequest) error
}

type SNP interface {
	Get(*pro_api.GetSNPRequest) ([]*pro_api.SNP, error)
	Create(*pro_api.CreateSNPRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdateSNPRequest) error
	Delete(*pro_api.DeleteSNPRequest) error

	AddMat(id string) error
	DeleteMat(id string, materials []*pro_api.AddMaterials) error
	AddMoun(id string) error
	DeleteMoun(id string) error
	AddGrap(id string) error
	DeleteGrap(id string) error
	DeleteFiller(id string) error
	DeleteTemp(id string) error
	DeleteMod(id string) error
}

type PutgImage interface {
	Get(req *pro_api.GetPutgImageRequest) ([]*pro_api.PutgImage, error)
	Create(image *pro_api.CreatePutgImageRequest) (*pro_api.IdResponse, error)
	Update(image *pro_api.UpdatePutgImageRequest) error
	Delete(image *pro_api.DeletePutgImageRequest) error
}

type Putg interface {
	Get(*pro_api.GetPutgRequest) ([]*pro_api.Putg, error)
	Create(*pro_api.CreatePutgRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdatePutgRequest) error
	Delete(*pro_api.DeletePutgRequest) error

	DeleteGrap(id string) error
	DeleteTemp(id string) error
	DeleteMod(id string) error
	DeleteMat(id string, materials []*pro_api.AddMaterials) error
	DeleteCon(id string) error
	DeleteObt(id string) error
	DeleteMoun(id string) error
	DeleteCoating(id string) error
}

type PutgmImage interface {
	Get(req *pro_api.GetPutgmImageRequest) ([]*pro_api.PutgmImage, error)
	Create(image *pro_api.CreatePutgmImageRequest) (*pro_api.IdResponse, error)
	Update(image *pro_api.UpdatePutgmImageRequest) error
	Delete(image *pro_api.DeletePutgmImageRequest) error
}

type Putgm interface {
	Get(*pro_api.GetPutgmRequest) ([]*pro_api.Putgm, error)
	Create(*pro_api.CreatePutgmRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdatePutgmRequest) error
	Delete(*pro_api.DeletePutgmRequest) error

	DeleteGrap(id string) error
	DeleteTemp(id string) error
	DeleteMod(id string) error
	DeleteMat(id string, materials []*pro_api.AddMaterials) error
	DeleteCon(id string) error
	DeleteObt(id string) error
	DeleteSeal(id string) error
	DeleteMoun(id string) error
	DeleteCoating(id string) error
}

type Materials interface {
	GetAll(*pro_api.GetMaterialsRequest) ([]*pro_api.Materials, error)
	Create(*pro_api.CreateMaterialsRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdateMaterialsRequest) error
	Delete(*pro_api.DeleteMaterialsRequest) error
}

type BoltMaterials interface {
	GetAll(*pro_api.GetBoltMaterialsRequest) ([]*pro_api.BoltMaterials, error)
	Create(*pro_api.CreateBoltMaterialsRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdateBoltMaterialsRequest) error
	Delete(*pro_api.DeleteBoltMaterialsRequest) error
}

type SizeInt interface {
	Get(*pro_api.GetSizesIntRequest) ([]*pro_api.SizeInt, []*pro_api.Dn, error)
	GetAll(*pro_api.GetAllSizeIntRequest) ([]*pro_api.SizeInt, []*pro_api.Dn, error)
	Create(*pro_api.CreateSizeIntRequest) (*pro_api.IdResponse, error)
	CreateMany(*pro_api.CreateSizesIntRequest) error
	Update(*pro_api.UpdateSizeIntRequest) error
	Delete(*pro_api.DeleteSizeIntRequest) error
	DeleteAll(*pro_api.DeleteAllSizeIntRequest) error
}

type Interview interface {
	SendInterview(context.Context, *pro_api.SendInterviewRequest) error
}

type Order interface {
	GetAll(*pro_api.GetAllOrdersRequest) ([]*pro_api.Order, error)
	Create(*pro_api.CreateOrderRequest) (*pro_api.IdResponse, error)
	Delete(*pro_api.DeleteOrderRequest) (*pro_api.IdResponse, error)
	Save(context.Context, *pro_api.SaveOrderRequest) (*bytes.Buffer, error)
	Send(context.Context, *pro_api.SaveOrderRequest) error
	Copy(*pro_api.CopyOrderRequest) error
}

type OrderPosition interface {
	Get(*pro_api.GetPositionsRequest) ([]*pro_api.OrderPosition, error)
	GetCur(*pro_api.GetCurPositionsRequest) ([]*pro_api.OrderPosition, error)
	Add(*pro_api.AddPositionRequest) (*pro_api.IdResponse, error)
	Copy(*pro_api.CopyPositionRequest) (*pro_api.IdResponse, error)
	Update(*pro_api.UpdatePositionRequest) (*pro_api.IdResponse, error)
	Remove(*pro_api.RemovePositionRequest) (*pro_api.IdResponse, error)
}

type Services struct {
	Stand
	Flange
	StFl
	TypeFl
	Addit
	Size
	SNP
	PutgImage
	Putg
	PutgmImage
	Putgm
	Materials
	BoltMaterials
	SizeInt
	Interview
	Order
	OrderPosition
}

func NewServices(repos *repository.Repositories, email email_api.EmailServiceClient,
	file file_api.FileServiceClient, user user_api.UserServiceClient) *Services {
	return &Services{
		Stand:         NewStandService(repos.Stand),
		Flange:        NewFlangeService(repos.Flange),
		StFl:          NewStFlService(repos.StFl),
		TypeFl:        NewTypeFlService(repos.TypeFl),
		Addit:         NewAdditService(repos.Addit),
		Size:          NewSizeService(repos.Size),
		SNP:           NewSNPService(repos.SNP),
		PutgImage:     NewPutgImageService(repos.PutgImage),
		Putg:          NewPutgService(repos.Putg),
		PutgmImage:    NewPutgmImageService(repos.PutgmImage),
		Putgm:         NewPutgmService(repos.Putgm),
		Materials:     NewMatService(repos.Materials),
		BoltMaterials: NewBoltMatRepo(repos.BoltMaterials),
		SizeInt:       NewSizeIntService(repos.SizeInt),
		Interview:     NewInterviewService(email, file),
		Order:         NewOrderService(repos.Order, email, file, user),
		OrderPosition: NewPositionService(repos.OrderPosition, file),
	}
}
