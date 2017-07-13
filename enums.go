package cloudinary

type Quality string

const (
	QualityAuto     = "q_auto"
	QualityAutoLow  = "q_auto:low"
	QualityAutoEco  = "q_auto:eco"
	QualityAutoGood = "q_auto:good"
	QualityAutoBest = "q_auto:best"
)

type CropMode string

const (
	CropModeScale  = "c_scale"
	CropModeFit    = "c_fit"
	CropModeLimit  = "c_limit"
	CropModeMfit   = "c_mfit"
	CropModeFill   = "c_fill"
	CropModeLfill  = "c_lfill"
	CropModePad    = "c_pad"
	CropModeLpad   = "c_lpad"
	CropModeCrop   = "c_crop"
	CropModeThumb  = "c_thumb"
	CropModeImagga = "c_imagga_crop"
)

type Action string

const (
	ActionExplicit = "explicit"
	ActionUpload   = "upload"
)

type SourceType string

const (
	SourceUpload = "upload"
)