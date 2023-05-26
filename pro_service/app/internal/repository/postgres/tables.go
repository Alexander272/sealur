package postgres

const (
	StandardTable       = "standard"
	FlangeStandardTable = "flange_standard"
	FlangeTypeTable     = "flange_type"
	MaterialTable       = "material"
	MountingTable       = "mounting"
	TemperatureTable    = "temperature"

	FlangeTypeSNPTable  = "flange_type_snp"
	SnpFillerTable      = "snp_filler"
	SnpFillerNewTable   = "snp_filler_new"
	SnpMaterialTable    = "snp_material"
	SnpMaterialTableNew = "snp_material_new"
	SnpStandardTable    = "snp_standard"
	SnpTypeTable        = "snp_type"
	SnpDataTable        = "snp_data"
	SnpSizeTable        = "snp_size"

	PutgConfTable             = "putg_configuration"
	PutgStandardTable         = "putg_standard"
	PutgConstructionTable     = "putg_construction"
	PutgConstructionTableTest = "putg_construction_test"
	PutgConstructionBaseTable = "putg_construction_base"
	PutgFillerTable           = "putg_filler"
	PutgFillerTableTest       = "putg_filler_test"
	PutgFillerBaseTable       = "putg_filler_base"
	PutgFlangeTypeTable       = "putg_flange_type"
	PutgDataTable             = "putg_data"
	PutgSizeTable             = "putg_size"
	PutgSizeTableTest         = "putg_size_new"
	PutgMaterialTable         = "putg_material"
	PutgTypeTable             = "putg_type"

	OrderTable                = "order"
	PositionTable             = "position"
	PositionMainSnpTable      = "position_snp_main"
	PositionSizeSnpTable      = "position_snp_size"
	PositionMaterialSnpTable  = "position_snp_material"
	PositionDesignSnpTable    = "position_snp_design"
	PositionMainPutgTable     = "position_putg_main"
	PositionSizePutgTable     = "position_putg_size"
	PositionMaterialPutgTable = "position_putg_material"
	PositionDesignPutgTable   = "position_putg_design"

	//? можно попробовать всю аналитику вынести в отдельный сервис или даже в представление
	UserTable = "user"
)
