package sioc

// InstanceCreationMode is a marker type for requesting different instance creation modes.
type InstanceCreationMode string

const (
	// CreateNewInstance indicates that a new instance should be created.
	CreateNewInstance InstanceCreationMode = "CREATE_NEW"
)
