// Package gateway provides the gateway interface for the application.
package gateway

import (
	"context"
	"errors"
	"time"

	"github.com/invopop/client.go/gateway"
	"github.com/invopop/popapp/internal/domain"
	"github.com/rs/zerolog/log"
)

const queueTime = 24 * time.Hour

// Service encapsulates the gateway functionality.
type Service struct {
	gw *gateway.Client
}

// New prepares the gateway to be able to start.
func New(conf *gateway.Config, domain *domain.Setup) *Service {
	s := new(Service)
	s.gw = gateway.New(
		gateway.WithConfig(conf),
		gateway.WithTaskHandler(s.execute),
	)
	return s
}

// Start listening for incoming messages
func (s *Service) Start() error {
	return s.gw.Start()
}

// Stop waiting for messages
func (s *Service) Stop() {
	s.gw.Stop()
}

func (s *Service) execute(ctx context.Context, in *gateway.Task) *gateway.TaskResult {
	log.Info().Str("job_id", in.JobId).Str("action", in.Action).Msg("gateway: incoming item task")
	res := s.executeAction(ctx, in)
	if res.Status != gateway.TaskStatus_OK {
		log.Warn().Str("job_id", in.JobId).
			Str("status", res.Status.String()).
			Str("msg", res.Message).
			Str("code", res.Code).
			Str("action", in.Action).
			Msg("gateway: execution issue")
	} else {
		log.Info().Str("job_id", in.JobId).Str("action", in.Action).Msg("gateway: executed")
	}
	return res
}

func (s *Service) executeAction(ctx context.Context, in *gateway.Task) *gateway.TaskResult {
	switch in.Action {
	default:
		return gateway.TaskKO(errors.New("unknown action"))
	}
}

func mapDomainError(ref string, err error) *gateway.TaskResult {
	var e *domain.Error
	if errors.As(err, &e) {
		out := new(gateway.TaskResult)
		out.Code = e.Code()
		if e.Is(domain.ErrFatal) || e.Is(domain.ErrAccessDenied) || e.Is(domain.ErrInvalid) {
			out.Status = gateway.TaskStatus_KO
			out.Message = err.Error()
		} else if e.Is(domain.ErrSkip) {
			out.Status = gateway.TaskStatus_SKIP
			out.Message = e.Message()
		} else if e.Is(domain.ErrQueue) {
			out.Status = gateway.TaskStatus_QUEUED
			out.Message = e.Message()
			out.RetryIn = int32(queueTime.Seconds())
			out.Ref = ref
		} else {
			out.Status = gateway.TaskStatus_ERR
			out.Message = err.Error()
			out.Ref = ref
		}
		return out
	}
	return gateway.TaskError(err)
}
