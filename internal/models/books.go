package models

type Book struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Page       *string   `json:"page"`
	CategoryId *uint     `json:"category_id"`
	AuthorId   *uint     `json:"author_id"`
	Category   *Category `json:"category"`
	Author     *Author   `json:"author"`
}

func (Book) RelationFields() []string {
	return []string{"Category", "Author"}
}

type BookRequest struct {
	ID         *uint   `json:"id"`
	Name       *string `json:"name"`
	Page       *string `json:"page"`
	CategoryId *uint   `json:"category_id"`
	AuthorId   *uint   `json:"author_id"`
}

func (b *BookRequest) ToModel(m *Book) {
	if b.ID != nil {
		m.ID = *b.ID
	}
	m.Name = *b.Name
	m.Page = b.Page
	m.CategoryId = b.CategoryId
	m.AuthorId = b.AuthorId
}

type BookResponse struct {
	ID       uint              `json:"id"`
	Name     string            `json:"name"`
	Page     *string           `json:"page"`
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
}
