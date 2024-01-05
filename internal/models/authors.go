package models

type Author struct {
	ID        uint    `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type AuthorRequest struct {
	ID        *uint   `json:"id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type AuthorResponse struct {
	ID        uint    `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
}

func (b *AuthorRequest) ToModel(m *Author) {
	if b.ID != nil {
		m.ID = *b.ID
	}
	m.FirstName = *b.FirstName
	m.LastName = b.LastName
}

func (r *AuthorResponse) FromModel(m *Author) {
	r.ID = m.ID
	r.FirstName = m.FirstName
	r.LastName = m.LastName
}
