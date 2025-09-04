package server

const (
	ErrInvalidIdFormat ApiError = "invalid id format"
	ErrGetVideoService ApiError = "could not get video"
	ErrVideoNotFound   ApiError = "video not found"
	ErrGetImageService ApiError = "error getting image by id from service"
	ErrImageNotFound   ApiError = "image not found"
)
