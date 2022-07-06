package models

type ISnippetModel interface {
	Insert(string, string, string) (int, error)
	Get(int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

type IUserModel interface {
	Insert(string, string, string) error
	Authenticate(string, string) (int, error)
	Get(int) (*User, error)
}
