package api

// FileItem 文件项目类
type FileItem struct {
	FileName string
	Content  []byte
	MimeType string
}

// NewFileItem 创建一个新的文件项目
func NewFileItem(fileName string, content []byte, mimeType string) *FileItem {
	return &FileItem{
		FileName: fileName,
		Content:  content,
		MimeType: mimeType,
	}
}

// GetFileName 获取文件名
func (f *FileItem) GetFileName() string {
	return f.FileName
}

// GetContent 获取文件内容
func (f *FileItem) GetContent() []byte {
	return f.Content
}

// GetMimeType 获取MIME类型
func (f *FileItem) GetMimeType() string {
	return f.MimeType
}
