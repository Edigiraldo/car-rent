#!/bin/zsh
mockgen -source=internal/core/ports/controllers.go -destination=internal/pkg/mocks/controllers.go &&
mockgen -source=internal/core/ports/repositories.go -destination=internal/pkg/mocks/repositories.go &&
mockgen -source=internal/core/ports/services.go -destination=internal/pkg/mocks/services.go
