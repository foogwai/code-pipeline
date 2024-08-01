package domain

// Repository defines the methods required for a data repository that stores POST data.
type Repository interface {
	SavePostData(postData PostData) error
}
