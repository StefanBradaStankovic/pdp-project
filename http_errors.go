package main

type httpStatus struct {
	status  int
	message string
}

var (
	statusCodeBadRequest = httpStatus{400, "Bad request"}
	statusCodeNotFound   = httpStatus{404, "Not found"}

	statusCodeInternalError = httpStatus{500, "Internal error"}
)
