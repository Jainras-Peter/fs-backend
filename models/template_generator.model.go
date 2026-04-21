package models



type TemplateGenerationRequest struct {
	TemplateType string      `json:"template_type" binding:"required"`
	Filename     string      `json:"filename" binding:"required"`
	Data         interface{} `json:"data" binding:"required"`
}

type TemplateGeneratorResult struct {
	StatusCode  int
	ContentType string
	Body        []byte // Contains PdfGeneratorUploadResponse JSON
}
