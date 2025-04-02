package models

// SearchRequest represents the request payload for candidate search
type SearchRequest struct {
	Collection   string   `json:"collection"`
	Query        string   `json:"query"`
	SearchFields []string `json:"search_fields,omitempty"`
	FilterFields []string `json:"filter_fields,omitempty"`
	SortBy       []string `json:"sort_by,omitempty"`
	GroupBy      []string `json:"group_by,omitempty"`
	Page         int      `json:"page,omitempty"`
	PerPage      int      `json:"per_page,omitempty"`
}

// VectorSearchRequest represents the request payload for vector-based candidate search
type VectorSearchRequest struct {
	VectorQuery  string   `json:"vector_query"`
	SearchFields []string `json:"search_fields,omitempty"`
	FilterFields []string `json:"filter_fields,omitempty"`
	SortBy       []string `json:"sort_by,omitempty"`
	Page         int      `json:"page,omitempty"`
	PerPage      int      `json:"per_page,omitempty"`
}

// SemanticSearchRequest represents the request payload for semantic candidate search
type SemanticSearchRequest struct {
	Query          string   `json:"query"`
	SearchFields   []string `json:"search_fields,omitempty"`
	FilterFields   []string `json:"filter_fields,omitempty"`
	SortBy         []string `json:"sort_by,omitempty"`
	Page           int      `json:"page,omitempty"`
	PerPage        int      `json:"per_page,omitempty"`
	EmbeddingField string   `json:"embedding_field,omitempty"`
	EmbeddingModel string   `json:"embedding_model,omitempty"`
}

// SearchResponse represents the response from a search operation
type SearchResponse struct {
	Found   int                      `json:"found"`
	Page    int                      `json:"page"`
	PerPage int                      `json:"per_page"`
	Hits    []map[string]interface{} `json:"hits"`
}

// Candidate represents a candidate profile
type Candidate struct {
	ID               string                 `json:"id"`
	FirstName        string                 `json:"first_name"`
	LastName         string                 `json:"last_name"`
	Email            string                 `json:"email"`
	Phone            string                 `json:"phone"`
	Skills           []string               `json:"skills"`
	LatestExperience string                 `json:"latest_experience"`
	HighestEducation string                 `json:"highest_education"`
	Description      string                 `json:"description"`
	LinkedIn         string                 `json:"linkedin"`
	GitHub           string                 `json:"github"`
	Twitter          string                 `json:"twitter"`
	StackOverflow    string                 `json:"stackoverflow"`
	PersonalBlog     string                 `json:"personal_blog"`
	Dribbble         string                 `json:"dribbble"`
	Behance          string                 `json:"behance"`
	GoogleScholar    string                 `json:"google_scholar"`
	ResearchGate     string                 `json:"research_gate"`
	Pronouns         string                 `json:"pronouns"`
	Custom           map[string]interface{} `json:"custom"`
	CreatedAt        string                 `json:"created_at"`
	UpdatedAt        string                 `json:"updated_at"`
}

// Attachment represents a candidate attachment
type Attachment struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	ModelName    string                 `json:"model_name"`
	ObjectKey    string                 `json:"object_key"`
	RecordID     string                 `json:"record_id"`
	Parent       string                 `json:"parent"`
	Content      string                 `json:"content"`
	AccessibleTo string                 `json:"accessible_to"`
	CreatedAt    string                 `json:"created_at"`
	UpdatedAt    string                 `json:"updated_at"`
	Document     map[string]interface{} `json:"document"`
}

// Comment represents a candidate comment
type Comment struct {
	ID           string   `json:"id"`
	Body         string   `json:"body"`
	Commentor    string   `json:"commentor"`
	ReviewID     string   `json:"review_id"`
	Source       string   `json:"source"`
	CreatedAt    int64    `json:"created_at"`
	UpdatedAt    int64    `json:"updated_at"`
	AccessibleTo []string `json:"accessible_to_user_ids,omitempty"`
}

// Review represents a candidate review
type Review struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	CandidateID     string   `json:"candidate_id"`
	Reviewed        bool     `json:"reviewed"`
	AttachmentCount int32    `json:"attachment_count"`
	CommentCount    int32    `json:"comment_count"`
	CreatedBy       string   `json:"created_by"`
	LabelIDs        []string `json:"label_ids,omitempty"`
	Members         []string `json:"members,omitempty"`
	AccessibleTo    []string `json:"accessible_to_user_ids,omitempty"`
	CreatedAt       int64    `json:"created_at"`
	UpdatedAt       int64    `json:"updated_at"`
}

// Label represents a candidate label
type Label struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color,omitempty"`
	Entity    string `json:"entity,omitempty"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
