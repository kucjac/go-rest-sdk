package forms

// Policy is a set of rules used during the process
// of model binding
type Policy struct {
	TaggedOnly  bool
	FailOnError bool
	Tag         string
}

// New creates a copy of the Policy
func (p Policy) New() *Policy {
	policy := &Policy{
		TaggedOnly:  p.TaggedOnly,
		FailOnError: p.FailOnError,
		Tag:         p.Tag,
	}
	return policy
}

// ListPolicy is a set of rules used during the process of model
// binding, enhanced with the 'List-parameters' for the list handler function.
type ListPolicy struct {
	Policy
	DefaultLimit int
	WithCount    bool
}

// New creates a copy of the list policy.
func (p ListPolicy) New() *ListPolicy {
	policy := &ListPolicy{
		Policy:       *p.Policy.New(),
		DefaultLimit: p.DefaultLimit,
		WithCount:    p.WithCount,
	}
	return policy
}

// ParamPolicy is a set of rules used during the process of
// routing/ url params.
// Enhances the Policy with DeepSearch field. This field defines if the
// Param binding function should check every model's field.
type ParamPolicy struct {
	Policy
	DeepSearch bool
}

// New creates a copy of the ParamPolicy
func (p ParamPolicy) New() *ParamPolicy {
	policy := &ParamPolicy{
		Policy:     *p.Policy.New(),
		DeepSearch: p.DeepSearch,
	}
	return policy
}

var (
	// DefaultPolicy default Policy that matches to the 'form' tag.
	// TaggedOnly and FailOnError fields are set to false.
	DefaultPolicy = Policy{
		TaggedOnly:  false,
		FailOnError: false,
		Tag:         "form",
	}

	// DefaultJSONPolicy defines the Policy for the BindJSON function.
	// by default the policy allows returning errors if occured during decoding.
	DefaultJSONPolicy = Policy{
		FailOnError: true,
	}

	// DefaultListPolicy is a default ListPolicy.
	// It uses DefaultPolicy as a base
	// Sets the DefaultLimit to 10
	// WithCount field is set to true
	DefaultListPolicy = ListPolicy{
		Policy:       DefaultPolicy,
		DefaultLimit: 10,
		WithCount:    true,
	}

	// DefaultParamPolicy is a default ParamPolicy.
	// It sets the default param tag to 'param'.
	// Every other fields are set to false (TaggedOnly, FailOnError and DeepSearch)
	// BindParam function based on this policy would search only for
	// the main ID field.
	DefaultParamPolicy = ParamPolicy{
		Policy: Policy{
			TaggedOnly:  false,
			FailOnError: false,
			Tag:         "param",
		},
		DeepSearch: false,
	}
)
