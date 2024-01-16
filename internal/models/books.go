package models

type Book struct {
	ID         uint      `json:"id"`
	CategoryId *uint     `json:"category_id"`
	AuthorId   *uint     `json:"author_id"`
	Name       string    `json:"name"`
	Page       *string   `json:"page"`
	Files      *[]string `json:"files"`
	Category   *Category `json:"category"`
	Author     *Author   `json:"author"`
}

func (Book) RelationFields() []string {
	return []string{"Category", "Author"}
}

type BookRequest struct {
	ID         *uint     `json:"id" form:"id"`
	Name       *string   `json:"name" form:"name"`
	Page       *string   `json:"page" form:"page"`
	Files      *[]string ``
	CategoryId *uint     `json:"category" form:"category"`
	AuthorId   *uint     `json:"author" form:"author"`
}

func (b *BookRequest) ToModel(m *Book) error {
	if b.ID != nil {
		m.ID = *b.ID
	}
	m.Name = *b.Name
	m.Page = b.Page
	m.Files = b.Files
	m.CategoryId = b.CategoryId
	m.AuthorId = b.AuthorId
	return nil
}

type BookResponse struct {
	ID       uint              `json:"id"`
	Name     string            `json:"name"`
	Page     *string           `json:"page"`
	Files    *[]string         `json:"files"`
	Category *CategoryResponse `json:"category"`
	Author   *AuthorResponse   `json:"author"`
}

func (r *BookResponse) FromModel(m *Book) {
	r.ID = m.ID
	r.Name = m.Name
	r.Page = m.Page
	if m.Category != nil {
		r.Category = &CategoryResponse{}
		r.Category.FromModel(m.Category)
	}
	if m.Author != nil {
		r.Author = &AuthorResponse{}
		r.Author.FromModel(m.Author)
	}
	if m.Files != nil {
		r.Files = &[]string{}
		for _, f := range *m.Files {
			*r.Files = append(*r.Files, fileUrl(&f))
		}
	}
}

func fileUrl(path *string) string {
	if path != nil {
		return "http://localhost:8000/web/uploads/images" + *path
	}
	return ""
}

type BookFilterRequest struct {
	ID         *uint `form:"id"`
	CategoryId *uint `form:"category_id"`
	AuthorId   *uint `form:"author_id"`
	PaginationRequest
}
