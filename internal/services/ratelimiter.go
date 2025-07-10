package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/httplimit"

	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"
	"github.com/sethvargo/go-redisstore"
)

type Route struct {
	Router string
	Method ports.Method
	Path   string
}

const ratelimiterCaller = packageCaller + ".Ratelimiter"

type Ratelimiter struct {
	logger *utils.Logger
	stores map[Route]limiter.Store
	ep     string
	pass   string
}

func NewRatelimiterService(logger *utils.Logger) ports.RatelimiterService {
	host := os.Getenv("DRAGONFLYDB_HOST")
	if host == "" {
		logger.Fatal(context.Background(), "DRAGONFLYDB_HOST is not set")
	}
	port := os.Getenv("DRAGONFLYDB_PORT")
	if port == "" {
		logger.Fatal(context.Background(), "DRAGONFLYDB_PORT is not set")
	}
	pass := os.Getenv("DRAGONFLYDB_PASSWORD")
	if pass == "" {
		logger.Fatal(context.Background(), "DRAGONFLYDB_PASSWORD is not set")
	}

	return &Ratelimiter{
		logger: logger,
		stores: make(map[Route]limiter.Store),
		ep:     host + ":" + port,
		pass:   pass,
	}
}

func (s *Ratelimiter) Handler(
	router string,
	method ports.Method,
	path string,
	token uint64,
	duration time.Duration,
	onIP, onID bool,
) fiber.Handler {
	route := Route{
		Router: router,
		Method: method,
		Path:   path,
	}
	if _, ok := s.stores[route]; ok {
		s.logger.Fatalf(context.Background(), "route is already registered: %v", route)
	}
	var err error
	s.stores[route], err = redisstore.New(
		&redisstore.Config{
			Tokens:   token,
			Interval: duration,
			Dial: func() (redis.Conn, error) {
				return redis.Dial(
					"tcp",
					s.ep,
					redis.DialPassword(s.pass),
					redis.DialDatabase(1),
				)
			},
		},
	)
	if err != nil {
		s.logger.Fatalf(context.Background(), "failed to create redisstore: %v", err)
	}
	var fun httplimit.KeyFunc
	switch {
	case onIP && onID:
		fun = ipAndIdFunc(router, string(method), path)
	case onIP:
		fun = ipFunc(router, string(method), path)
	case onID:
		fun = idFunc(router, string(method), path)
	default:
		s.logger.Fatalf(context.Background(), "onIP and onID are not valid for: %v", route)
	}
	middleware, err := httplimit.NewMiddleware(s.stores[route], fun)
	if err != nil {
		s.logger.Fatalf(context.Background(), "failed to create middleware: %v", err)
	}
	return s.fiberMiddleware(middleware)
}

func ipFunc(router, method, path string) func(r *http.Request) (key string, err error) {
	return func(r *http.Request) (key string, err error) {
		defer func() {
			const caller = ratelimiterCaller + ".ipFunc"
			if err != nil {
				var uErr *utils.Error
				if errors.As(err, &uErr) {
					err = uErr.WithCaller(caller)
				}
			}
		}()
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			return "", errors.New("X-Forwarded-For header is missing")
		}
		return fmt.Sprintf("%s-%s-%s-%s", router, method, path, ip), nil
	}
}

func idFunc(router, method, path string) func(r *http.Request) (key string, err error) {
	return func(r *http.Request) (key string, err error) {
		defer func() {
			const caller = ratelimiterCaller + ".idFunc"
			if err != nil {
				var uErr *utils.Error
				if errors.As(err, &uErr) {
					err = uErr.WithCaller(caller)
				}
			}
		}()
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			return "", errors.New("X-User-ID header is missing")
		}
		return fmt.Sprintf("%s-%s-%s-%s", router, method, path, userID), nil
	}
}

func ipAndIdFunc(router, method, path string) func(r *http.Request) (key string, err error) {
	return func(r *http.Request) (key string, err error) {
		defer func() {
			const caller = ratelimiterCaller + ".ipAndIdFunc"
			if err != nil {
				var uErr *utils.Error
				if errors.As(err, &uErr) {
					err = uErr.WithCaller(caller)
				}
			}
		}()
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			return "", errors.New("X-Forwarded-For header is missing")
		}
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			return "", errors.New("X-User-ID header is missing")
		}
		return fmt.Sprintf("%s-%s-%s-%s%s", router, method, path, ip, userID), nil
	}
}

func (s *Ratelimiter) fiberMiddleware(m *httplimit.Middleware) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			const caller = ratelimiterCaller + ".fiberMiddleware"
			if err != nil {
				var uErr *utils.Error
				if errors.As(err, &uErr) {
					err = uErr.WithCaller(caller)
				}
			}
		}()
		httpHeader := make(http.Header)
		c.Request().Header.VisitAll(func(key, value []byte) {
			httpHeader.Add(string(key), string(value))
		})

		fasthttpURI := c.Request().URI()
		uriString := string(fasthttpURI.RequestURI())
		parsedURL, err := url.Parse(uriString)
		if err != nil {
			return err
		}

		httpReq := &http.Request{
			Method: string(c.Method()),
			URL:    parsedURL,
			Header: httpHeader,
			Body:   &fasthttpRequestBody{ctx: c},
		}

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for k, v := range r.Header {
				for _, val := range v {
					c.Request().Header.Set(k, val)
				}
			}
			if err := c.Next(); err != nil {
				var uErr *utils.Error
				if errors.As(err, &uErr) {
					utils.HandleUtilsError(s.logger, c.Context(), w, uErr)
					return
				}
				utils.HandleUnknownError(s.logger, c.Context(), w, err)
				return
			}
			c.Response().Header.VisitAll(func(key, value []byte) {
				w.Header().Add(string(key), string(value))
			})
			w.WriteHeader(c.Response().StatusCode())
		})

		handler := m.Handle(next)

		rw := &fiberResponseWriter{
			ctx:        c,
			statusCode: http.StatusOK,
			header:     make(http.Header),
		}

		handler.ServeHTTP(rw, httpReq)

		for k, v := range rw.header {
			for _, val := range v {
				c.Response().Header.Set(k, val)
			}
		}
		c.Status(rw.statusCode)
		switch rw.statusCode {
		case http.StatusTooManyRequests:
			if len(c.Response().Body()) >= 71 {
				return utils.AuthMethodOtpGetTooEarlyResponse.Clone()
			}
			return utils.TooManyRequestsResponse.Clone()
		case http.StatusInternalServerError:
			s.logger.Error(c.Context(), err, "unknown error")
			return utils.InternalServerResponse.Clone()
		}
		return nil
	}
}

type fiberResponseWriter struct {
	ctx        *fiber.Ctx
	statusCode int
	header     http.Header
	written    bool
}

func (rw *fiberResponseWriter) Header() http.Header {
	return rw.header
}

func (rw *fiberResponseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
	}
}

func (rw *fiberResponseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ctx.Write(b)
}

type fasthttpRequestBody struct {
	ctx *fiber.Ctx
}

func (b *fasthttpRequestBody) Read(p []byte) (int, error) {
	return b.ctx.Request().BodyStream().Read(p)
}

func (b *fasthttpRequestBody) Close() error {
	b.ctx.Request().ResetBody()
	return nil
}
