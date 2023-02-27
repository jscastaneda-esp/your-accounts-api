// TODO: Pendientes tests

package directive

import (
	"api-your-accounts/shared/domain"
	"api-your-accounts/shared/infrastructure/graph"
	middleware "api-your-accounts/shared/infrastructure/middleware/auth"
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var (
	validate = validator.New()
)

func auth(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	if tokenData, ok := ctx.Value(middleware.CtxAuth).(*domain.JwtCustomClaim); tokenData != nil && ok {
		return next(ctx)
	}

	return nil, &gqlerror.Error{
		Message: "Access Denied",
	}
}

func binding(ctx context.Context, _ interface{}, next graphql.Resolver, constraint string) (interface{}, error) {
	val, err := next(ctx)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: err.Error(),
		}
	}

	err = validate.Var(val, constraint)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: err.Error(),
		}
	}

	return val, nil
}

func GetDirectives() graph.DirectiveRoot {
	return graph.DirectiveRoot{
		Auth:    auth,
		Binding: binding,
	}
}
