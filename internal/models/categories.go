package models

type Category struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (Category) RelationFields() []string {
	return []string{}
}

type CategoryRequest struct {
	ID   *uint   `json:"id"`
	Name *string `json:"name"`
}

type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (b *CategoryRequest) ToModel(m *Category) {
	if b.ID != nil {
		m.ID = *b.ID
	}
	m.Name = *b.Name
}

func (r *CategoryResponse) FromModel(m *Category) {
	r.ID = m.ID
	r.Name = m.Name
}

type CategoryFilterRequest struct {
	ID *uint `json:"id"`
	PaginationRequest
}
