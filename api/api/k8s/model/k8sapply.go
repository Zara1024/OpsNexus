package model

type ApplyManifestRequest struct {
	YAMLContent     string `json:"yamlContent" binding:"required"`
	Namespace       string `json:"namespace,omitempty"`
	DryRun          bool   `json:"dryRun"`
	ValidateOnly    bool   `json:"validateOnly"`
	ServerSideApply *bool  `json:"serverSideApply,omitempty"`
	FieldManager    string `json:"fieldManager,omitempty"`
	ForceConflicts  bool   `json:"forceConflicts"`
}

type ApplyManifestItemResult struct {
	APIVersion      string `json:"apiVersion"`
	Kind            string `json:"kind"`
	Name            string `json:"name"`
	Namespace       string `json:"namespace,omitempty"`
	Operation       string `json:"operation"`
	ResourceVersion string `json:"resourceVersion,omitempty"`
}

type ApplyManifestResponse struct {
	Success         bool                      `json:"success"`
	Message         string                    `json:"message"`
	DryRun          bool                      `json:"dryRun"`
	ValidateOnly    bool                      `json:"validateOnly"`
	ServerSideApply bool                      `json:"serverSideApply"`
	Results         []ApplyManifestItemResult `json:"results"`
}
