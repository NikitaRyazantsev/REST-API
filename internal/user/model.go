package user

// file for user struct description

type User struct {
	ID       string   `json:"id" bson:"_id,omitempty"`
	Username string   `json:"username" bson:"username"`
	Age      string   `json:"age" bson:"age"`
	Friends  []string `json:"friends" bson:"friends"`
}
