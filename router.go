package clevergo

import (
	"github.com/clevergo/router"
	"github.com/clevergo/sessions"
	"github.com/valyala/fasthttp"
)

// Router for managing request handlers.
type Router struct {
	*router.Router
	middlewares  []Middleware    // Middlewares.
	sessionStore sessions.Store  // Session store for Context.
	logger       fasthttp.Logger // Logger for Context.
}

// NewRouter returns a Router's instance.
func NewRouter() *Router {
	return &Router{
		Router:      router.New(),
		middlewares: make([]Middleware, 0),
	}
}

// SetSessionStore for setting session store.
func (r *Router) SetSessionStore(store sessions.Store) {
	r.sessionStore = store
}

// SetLogger for setting logger.
func (r *Router) SetLogger(logger fasthttp.Logger) {
	r.logger = logger
}

// SetMiddlewares for setting middlewares.
func (r *Router) SetMiddlewares(middlewares []Middleware) {
	r.middlewares = middlewares
}

// AddMiddleware for adding middleware.
func (r *Router) AddMiddleware(middleware Middleware) {
	r.middlewares = append(r.middlewares, middleware)
}

// GET for adding GET request handler.
func (r *Router) GET(path string, handler Handler) {
	r.Router.GET(path, r.getHandler(handler))
}

// HEAD for adding HEAD request handler.
func (r *Router) HEAD(path string, handler Handler) {
	r.Router.HEAD(path, r.getHandler(handler))
}

// OPTIONS for adding OPTIONS request handler.
func (r *Router) OPTIONS(path string, handler Handler) {
	r.Router.OPTIONS(path, r.getHandler(handler))
}

// POST for adding POST request handler.
func (r *Router) POST(path string, handler Handler) {
	r.Router.POST(path, r.getHandler(handler))
}

// PUT for adding PUT request handler.
func (r *Router) PUT(path string, handler Handler) {
	r.Router.PUT(path, r.getHandler(handler))
}

// PATCH for adding PATCH request handler.
func (r *Router) PATCH(path string, handler Handler) {
	r.Router.PATCH(path, r.getHandler(handler))
}

// DELETE for adding DELETE request handler.
func (r *Router) DELETE(path string, handler Handler) {
	r.Router.DELETE(path, r.getHandler(handler))
}

// Handle for adding custom METHOD request handler.
func (r *Router) Handle(method, path string, handler Handler) {
	r.Router.Handle(method, path, r.getHandler(handler))
}

func (r *Router) getHandler(handler Handler) router.Handle {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i].Handle(handler)
	}

	return func(_ctx *fasthttp.RequestCtx, ps router.Params) {
		ctx := NewContext(r, _ctx, &ps)
		defer ctx.Close()
		handler.Handle(ctx)
	}
}

// RegisterController for registering controller.
//
// The Controller should implemented the ControllerInterface.
func (r *Router) RegisterController(route string, c ControllerInterface) {
	handlers := make(map[string]Handler, 0)

	// Register GET request's handler.
	var getHandler Handler
	getHandler = c.Handle(HandlerFunc(c.GET))
	handlers["GET"] = getHandler

	// Register POST request's handler.
	var postHandler Handler
	postHandler = c.Handle(HandlerFunc(c.POST))
	handlers["POST"] = postHandler

	// Register DELETE request's handler.
	var deleteHandler Handler
	deleteHandler = c.Handle(HandlerFunc(c.DELETE))
	handlers["DELETE"] = deleteHandler

	// Register PUT request's handler.
	var putHandler Handler
	putHandler = c.Handle(HandlerFunc(c.PUT))
	handlers["PUT"] = putHandler

	// Register OPTIONS request's handler.
	var optionsHandler Handler
	optionsHandler = c.Handle(HandlerFunc(c.OPTIONS))
	handlers["OPTIONS"] = optionsHandler

	// Register PATCH request's handler.
	var patchHandler Handler
	patchHandler = c.Handle(HandlerFunc(c.PATCH))
	handlers["PATCH"] = patchHandler

	for method, handler := range handlers {
		var _handler Handler
		_handler = handler
		// Register middlewares.
		for i := len(r.middlewares) - 1; i >= 0; i-- {
			_handler = r.middlewares[i].Handle(_handler)
		}
		// Add to route.
		r.Router.Handle(method, route, func(_ctx *fasthttp.RequestCtx, ps router.Params) {
			ctx := NewContext(r, _ctx, &ps)
			defer ctx.Close()
			_handler.Handle(ctx)
		})
	}
}
