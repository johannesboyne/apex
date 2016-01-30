package hooks

// Hooks supported by Apex.
type Hooks struct {
	// Build command is run before creating the zip file.
	Build string `json:"build"`

	// Clean command is run after creating the zip file.
	Clean string `json:"clean"`

	// Deploy command is run after builds and before deploys.
	Deploy string `json:"deploy"`

	// PostDeploy command is run after deploys.
	PostDeploy string `json:"deployed"`
}
