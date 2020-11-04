package msg

// Create static initializers file
//go:generate go run "$REPO_ROOT/galley/pkg/config/analysis/msg/generate.main.go" messages.yaml messages.gen.go

//go:generate goimports -w -local istio.io "$REPO_ROOT/galley/pkg/config/analysis/msg/messages.gen.go"
