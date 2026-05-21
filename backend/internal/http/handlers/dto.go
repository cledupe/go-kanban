package handlers

type CreateBoardRequest struct {
	Name     string  `json:"name"`
	Template *string `json:"template,omitempty"`
}

type UpdateBoardRequest struct {
	Name string `json:"name"`
}

type CreateColumnRequest struct {
	Name string `json:"name"`
}

type UpdateColumnRequest struct {
	Name     *string `json:"name,omitempty"`
	Position *int    `json:"position,omitempty"`
}

type CreateCardRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type UpdateCardRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

type MoveCardRequest struct {
	TargetColumnID string `json:"target_column_id"`
	Position       int    `json:"position"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}