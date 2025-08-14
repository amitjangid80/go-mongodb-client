package mongodb_domain

type ResponseModel struct {
	Error   bool   `json:"error" bson:"error"`
	Message string `json:"message" bson:"message"`
	Code    int    `json:"code" bson:"code"`
	Data    any    `json:"data" bson:"data"`
}
