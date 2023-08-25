package oauth

type Spec interface {
	GetAccessToken() (string, error)
}
