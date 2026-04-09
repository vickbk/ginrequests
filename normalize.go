package ginrequests

// NormalizeRequests flattens multiple slices of Request objects into a single combined slice.
//
// Example usage:
//
//	getRequests := BuildRequests("GET", "/users", listUsersHandler)
//	postRequests := BuildRequests("POST", "/users", createUserHandler)
//	deleteRequests := BuildRequests("DELETE", "/users/:id", deleteUserHandler)
//
//	allRequests := NormalizeRequests(getRequests, postRequests, deleteRequests)
//	// Result: single slice containing all 3 requests
func NormalizeRequests(rs ...[]Request) []Request {
	rqs := make([]Request, 0, getTotalRequestsLength(&rs))

	for _, rq := range rs {
		rqs = append(rqs, rq...)
	}

	return rqs
}

// getTotalRequestsLength calculates the total number of Request objects across all provided slices.
func getTotalRequestsLength(rs *[][]Request) (totalLn int) {
	for _, r := range *rs {
		totalLn += len(r)
	}
	return
}
