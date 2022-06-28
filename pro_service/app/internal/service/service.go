package service

import (
	"context"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	proto_email "github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto/email"
	proto_file "github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto/file"
)

type Stand interface {
	GetAll(*proto.GetStandsRequest) (stands []*proto.Stand, err error)
	Create(*proto.CreateStandRequest) (st *proto.IdResponse, err error)
	Update(*proto.UpdateStandRequest) error
	Delete(*proto.DeleteStandRequest) error
}

type Flange interface {
	GetAll() ([]*proto.Flange, error)
	Create(*proto.CreateFlangeRequest) (*proto.IdResponse, error)
	Update(*proto.UpdateFlangeRequest) error
	Delete(*proto.DeleteFlangeRequest) error
}

type StFl interface {
	Get() ([]*proto.StFl, error)
	Create(*proto.CreateStFlRequest) (*proto.IdResponse, error)
	Update(*proto.UpdateStFlRequest) error
	Delete(*proto.DeleteStFlRequest) error
}

type TypeFl interface {
	Get() ([]*proto.TypeFl, error)
	GetAll() ([]*proto.TypeFl, error)
	Create(*proto.CreateTypeFlRequest) (*proto.IdResponse, error)
	Update(*proto.UpdateTypeFlRequest) error
	Delete(*proto.DeleteTypeFlRequest) error
}

type Addit interface {
	GetAll() ([]*proto.Additional, error)
	Create(*proto.CreateAddRequest) (*proto.SuccessResponse, error)
	UpdateMat(*proto.UpdateAddMatRequest) (*proto.SuccessResponse, error)
	UpdateMod(*proto.UpdateAddModRequest) (*proto.SuccessResponse, error)
	UpdateTemp(*proto.UpdateAddTemRequest) (*proto.SuccessResponse, error)
	UpdateMoun(*proto.UpdateAddMounRequest) (*proto.SuccessResponse, error)
	UpdateGrap(*proto.UpdateAddGrapRequest) (*proto.SuccessResponse, error)
	UpdateFillers(*proto.UpdateAddFillersRequest) (*proto.SuccessResponse, error)
	UpdateCoating(*proto.UpdateAddCoatingRequest) (*proto.SuccessResponse, error)
	UpdateConstruction(*proto.UpdateAddConstructionRequest) (*proto.SuccessResponse, error)
	UpdateObturator(*proto.UpdateAddObturatorRequest) (*proto.SuccessResponse, error)
	UpdateBasis(*proto.UpdateAddBasisRequest) (*proto.SuccessResponse, error)
	UpdatePObturator(*proto.UpdateAddPObturatorRequest) (*proto.SuccessResponse, error)
	UpdateSealant(*proto.UpdateAddSealantRequest) (*proto.SuccessResponse, error)
}

type Size interface {
	Get(*proto.GetSizesRequest) ([]*proto.Size, []*proto.Dn, error)
	GetAll(*proto.GetSizesRequest) ([]*proto.Size, []*proto.Dn, error)
	Create(*proto.CreateSizeRequest) (*proto.IdResponse, error)
	CreateMany(*proto.CreateSizesRequest) error
	Update(*proto.UpdateSizeRequest) error
	Delete(*proto.DeleteSizeRequest) error
	DeleteAll(*proto.DeleteAllSizeRequest) error
}

type SNP interface {
	Get(*proto.GetSNPRequest) ([]*proto.SNP, error)
	Create(*proto.CreateSNPRequest) (*proto.IdResponse, error)
	Update(*proto.UpdateSNPRequest) error
	Delete(*proto.DeleteSNPRequest) error

	AddMat(id string) error
	DeleteMat(id string, materials []*proto.AddMaterials) error
	AddMoun(id string) error
	DeleteMoun(id string) error
	AddGrap(id string) error
	DeleteGrap(id string) error
	DeleteFiller(id string) error
	DeleteTemp(id string) error
	DeleteMod(id string) error
}

type PutgImage interface {
	Get(req *proto.GetPutgImageRequest) ([]*proto.PutgImage, error)
	Create(image *proto.CreatePutgImageRequest) (*proto.IdResponse, error)
	Update(image *proto.UpdatePutgImageRequest) error
	Delete(image *proto.DeletePutgImageRequest) error
}

type Putg interface {
	Get(*proto.GetPutgRequest) ([]*proto.Putg, error)
	Create(*proto.CreatePutgRequest) (*proto.IdResponse, error)
	Update(*proto.UpdatePutgRequest) error
	Delete(*proto.DeletePutgRequest) error

	DeleteGrap(id string) error
	DeleteTemp(id string) error
	DeleteMod(id string) error
	DeleteMat(id string, materials []*proto.AddMaterials) error
	DeleteCon(id string) error
	DeleteObt(id string) error
	DeleteMoun(id string) error
	DeleteCoating(id string) error
}

type PutgmImage interface {
	Get(req *proto.GetPutgmImageRequest) ([]*proto.PutgmImage, error)
	Create(image *proto.CreatePutgmImageRequest) (*proto.IdResponse, error)
	Update(image *proto.UpdatePutgmImageRequest) error
	Delete(image *proto.DeletePutgmImageRequest) error
}

type Putgm interface {
	Get(*proto.GetPutgmRequest) ([]*proto.Putgm, error)
	Create(*proto.CreatePutgmRequest) (*proto.IdResponse, error)
	Update(*proto.UpdatePutgmRequest) error
	Delete(*proto.DeletePutgmRequest) error

	DeleteGrap(id string) error
	DeleteTemp(id string) error
	DeleteMod(id string) error
	DeleteMat(id string, materials []*proto.AddMaterials) error
	DeleteCon(id string) error
	DeleteObt(id string) error
	DeleteSeal(id string) error
	DeleteMoun(id string) error
	DeleteCoating(id string) error
}

type Materials interface {
	GetAll(*proto.GetMaterialsRequest) ([]*proto.Materials, error)
	Create(*proto.CreateMaterialsRequest) (*proto.IdResponse, error)
	Update(*proto.UpdateMaterialsRequest) error
	Delete(*proto.DeleteMaterialsRequest) error
}

type BoltMaterials interface {
	GetAll(*proto.GetBoltMaterialsRequest) ([]*proto.BoltMaterials, error)
	Create(*proto.CreateBoltMaterialsRequest) (*proto.IdResponse, error)
	Update(*proto.UpdateBoltMaterialsRequest) error
	Delete(*proto.DeleteBoltMaterialsRequest) error
}

type SizeInt interface {
	Get(*proto.GetSizesIntRequest) ([]*proto.SizeInt, []*proto.Dn, error)
	GetAll(*proto.GetAllSizeIntRequest) ([]*proto.SizeInt, []*proto.Dn, error)
	Create(*proto.CreateSizeIntRequest) (*proto.IdResponse, error)
	CreateMany(*proto.CreateSizesIntRequest) error
	Update(*proto.UpdateSizeIntRequest) error
	Delete(*proto.DeleteSizeIntRequest) error
	DeleteAll(*proto.DeleteAllSizeIntRequest) error
}

type Interview interface {
	SendInterview(context.Context, *proto.SendInterviewRequest) error
}

type Order interface {
	GetAll(*proto.GetAllOrdersRequest) ([]*proto.Order, error)
	Create(*proto.CreateOrderRequest) (*proto.IdResponse, error)
	Delete(*proto.DeleteOrderRequest) (*proto.IdResponse, error)
	Save(*proto.SaveOrderRequest) error
}

type OrderPosition interface {
	Get(*proto.GetPositionsRequest) ([]*proto.OrderPosition, error)
	GetCur(*proto.GetCurPositionsRequest) ([]*proto.OrderPosition, error)
	Add(*proto.AddPositionRequest) (*proto.IdResponse, error)
	Update(*proto.UpdatePositionRequest) (*proto.IdResponse, error)
	Remove(*proto.RemovePositionRequest) (*proto.IdResponse, error)
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

func NewServices(repos *repository.Repositories, email proto_email.EmailServiceClient, file proto_file.FileServiceClient) *Services {
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
		Order:         NewOrderService(repos.Order),
		OrderPosition: NewPositionService(repos.OrderPosition),
	}
}
