// Package spy
package spy

// type AuthenticateSpy struct {
// 	ShouldReturnError        bool
// 	ShouldReturnUserNotFound bool
// 	CustomOutput             dto.AuthenticatedOutput
// 	CallCount                int
// 	LastCPF                  string
// 	LastPassword             string
// }
//
// func NewAuthenticateSpy() *AuthenticateSpy {
// 	return &AuthenticateSpy{
// 		CustomOutput: dto.AuthenticatedOutput{
// 			PersonID:  "1",
// 			FirstName: "John",
// 			LastName:  "Doe",
// 		},
// 	}
// }
//
// func (a *AuthenticateSpy) SignIn(ctx context.Context, cpf string, password string) (dto.AuthenticatedOutput, error) {
// 	a.CallCount++
// 	a.LastCPF = cpf
// 	a.LastPassword = password
// 	if a.ShouldReturnUserNotFound {
// 		return dto.AuthenticatedOutput{}, domain.ErrUserNotFound
// 	}
// 	if a.ShouldReturnError {
// 		return dto.AuthenticatedOutput{}, errors.New("spy database error")
// 	}
// 	return a.CustomOutput, nil
// }
