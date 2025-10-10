package websocket

import (
	"fmt"

	"github.com/ali-mahdavi-dev/bunny-go/config"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/infrastructure/logging"
	"github.com/golang-jwt/jwt/v5"
	socketio "github.com/googollee/go-socket.io"
	"github.com/spf13/cast"
)

type Websocket struct {
	socket *socketio.Server
	logger logging.Logger
	// uow    unit_of_work.PGUnitOfWork
	cfg *config.Config
}

func NewWebsocket(socket *socketio.Server, logger logging.Logger, cfg *config.Config) *Websocket {
	return &Websocket{
		socket: socket,
		logger: logger,
		cfg:    cfg,
	}
}

func (w *Websocket) AddWsRoutes() {
	w.socket.OnConnect("/", func(s socketio.Conn) error {
		r := s.RemoteHeader()
		token := r.Get("Authorization")
		if token == "" {
			u := s.URL()
			token = u.Query().Get("token")
		}
		w.logger.Info(logging.WS, logging.WSOnConnect, "New connection", map[logging.ExtraKey]interface{}{
			logging.WSSocketID: s.ID(),
			logging.WSToken:    token,
		})

		if token == "" {
			s.Emit("error", "Missing token")
			s.Close()
			return fmt.Errorf("no token provided")
		}

		s.SetContext(token)

		userID, err := w.extractUserIDFromToken(token)
		if err != nil {
			return fmt.Errorf("fail to get UserID: %w", err)
		}
		s.Join(cast.ToString(userID))

		return nil
	})
}

func (w *Websocket) extractUserIDFromToken(tokenStr string) (uint64, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// validate alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return w.cfg.JWT.Secret, nil
	})
	if err != nil {
		return 0, fmt.Errorf("Websocket.extractUserFromToken fail to pars token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid JWT claims")
	}
	return cast.ToUint64(claims["user_id"]), nil
}
