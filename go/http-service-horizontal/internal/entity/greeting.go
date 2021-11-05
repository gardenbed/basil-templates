package entity

import "fmt"

// GreetRequest is the domain model for a Greet request.
type GreetRequest struct {
	GithubUsername string
}

// String implements the fmt.Stringer interface.
func (r *GreetRequest) String() string {
	return fmt.Sprintf("GreetRequest{github_username=%s}", r.GithubUsername)
}

// GreetResponse is the domain model for a Greet response.
type GreetResponse struct {
	Greeting string
}

// String implements the fmt.Stringer interface.
func (r *GreetResponse) String() string {
	return fmt.Sprintf("GreetResponse{greeting=%s}", r.Greeting)
}
