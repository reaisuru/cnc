package flags

func Clientside() FlagOption {
	return func(opts *FlagOptions) error {
		opts.Clientside = true
		return nil
	}
}

func AdminOnly() FlagOption {
	return func(opts *FlagOptions) error {
		opts.Admin = true
		return nil
	}
}

func Invisible() FlagOption {
	return func(opts *FlagOptions) error {
		opts.Invisible = true
		return nil
	}
}
