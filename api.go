package gisproxy

//APIMeta api meta
type APIMeta struct {
	Count int `json:"count,omitempty"`
}

//APIResource APIResource
type APIResource struct {
	Meta APIMeta `json:"meta,omitempty"`
}

//UsersAPIResponse api response
type UsersAPIResponse struct {
	APIResource
	Objects []*User `json:"objects,omitempty"`
}
