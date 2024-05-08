package swagger

import (
	"fmt"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	apiClient "week07/grpc_5_gateway/swagger/sess-client/client"
	auth "week07/grpc_5_gateway/swagger/sess-client/client/auth_checker"
	models "week07/grpc_5_gateway/swagger/sess-client/models"
)

// simple demo
func MainClient() {

	transport := httptransport.New("127.0.0.1:8080", "", []string{"http"})
	client := apiClient.New(transport, strfmt.Default)
	sessManager := client.AuthChecker

	// создаем сессию
	sessId, err := sessManager.AuthCheckerCreate(auth.NewAuthCheckerCreateParams().WithBody(
		&models.Grpc5GatewaySession{
			Login:     "rvasily",
			Useragent: "chrome",
		},
	))
	fmt.Println("sessId", sessId, err)

	// проверяем сессию
	sess, err := sessManager.AuthCheckerCheck(auth.
		NewAuthCheckerCheckParams().
		WithID(sessId.Payload.ID))
	fmt.Println("after create", sess, err)

	// удаляем сессию
	_, err = sessManager.AuthCheckerDelete(auth.NewAuthCheckerDeleteParams().WithBody(
		&models.Grpc5GatewaySessionID{
			ID: sessId.Payload.ID,
		},
	))

	// проверяем еще раз
	sess, err = sessManager.AuthCheckerCheck(auth.
		NewAuthCheckerCheckParams().
		WithID(sessId.Payload.ID))
	fmt.Println("after delete", sess, err)
}
