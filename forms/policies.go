package forms

// BindPolicy is a set of rules used during the process of model binding.
type BindPolicy struct {
	// TaggedOnly defines if a binding function should bind only to fields
	// that has a non-empty tag - defined in 'Tag' BindPolicy field.
	TaggedOnly bool

	// FailOnError defines if a binding function should return an error if occurs.
	// If set to false allows bidning functions to continue if an error occurs.
	FailOnError bool

	// Tag defines the tag used in binding function that uses this policy
	Tag string

	// SearchDepthLevel defines how deep the binding function should search
	SearchDepthLevel int
}

// Copy creates a copy of the BindPolicy
func (p BindPolicy) Copy() *BindPolicy {
	policyCopy := p
	return &policyCopy
}

// ListPolicy is a set of rules used during the process of model
// binding, enhanced with the 'List-parameters' for the list handler function.
type ListPolicy struct {
	BindPolicy
	DefaultLimit int
	WithCount    bool
}

// New creates a copy of the list policy.
func (p ListPolicy) Copy() *ListPolicy {
	policyCopy := p
	policy := &policyCopy
	return policy
}

// ParamPolicy is a set of rules used during the process of
// routing/ url params.
// Enhances the BindPolicy with DeepSearch field. This field defines if the
// Param binding function should check every model's field.
type ParamPolicy struct {
	BindPolicy

	// IDOnly a rule that forces ParamBinding to search only for the ID field.
	IDOnly bool
}

// Copy creates a copy of the ParamPolicy
func (p ParamPolicy) Copy() (policy *ParamPolicy) {
	policyCopy := p
	policy = &policyCopy
	return policy
}

var (
	// DefaultBindPolicy default BindPolicy that matches to the 'form' tag.
	// TaggedOnly and FailOnError fields are set to false.
	DefaultBindPolicy = BindPolicy{
		TaggedOnly:       false,
		FailOnError:      false,
		Tag:              "form",
		SearchDepthLevel: 0,
	}

	// DefaultJSONPolicy defines the BindPolicy for the BindJSON function.
	// by default the policy allows returning errors if occured during decoding.
	DefaultJSONPolicy = BindPolicy{
		FailOnError: true,
	}

	// DefaultListPolicy is a default ListPolicy.
	// It uses DefaultBindPolicy as a base
	// Sets the DefaultLimit to 10
	// WithCount field is set to true
	DefaultListPolicy = ListPolicy{
		BindPolicy:   DefaultBindPolicy,
		DefaultLimit: 10,
		WithCount:    true,
	}

	// DefaultParamPolicy is a default ParamPolicy.
	// It sets the default param tag to 'param'.
	// Every other fields are set to false (TaggedOnly, FailOnError and DeepSearch)
	// BindParam function based on this policy would search only for
	// the main ID field.
	DefaultParamPolicy = ParamPolicy{
		BindPolicy: BindPolicy{
			TaggedOnly:       false,
			FailOnError:      false,
			Tag:              "param",
			SearchDepthLevel: 0,
		},
		IDOnly: false,
	}
)
