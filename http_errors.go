package main

type httpStatus struct {
	status  int
	message string
}

var (
	statusCodeItemDeleted = httpStatus{200, "Item deleted"}

	statusCodeBadRequest = httpStatus{400, "Bad request"}
	statusCodeNotFound   = httpStatus{404, "Not found"}

	statusCodeInternalError = httpStatus{500, "Internal error"}
)
