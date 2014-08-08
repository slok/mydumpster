package mydumpster

type Censorship struct {
	Key          string
	Suffix       string
	Prefix       string
	Blank        bool
	Null         bool
	DefaultValue string
}

// Censors the string and returns the ensored string and if is a nil value (NULL)
func (c *Censorship) censore(val string) (string, bool) {

	if len(c.DefaultValue) > 0 {
		return c.DefaultValue, false
	}

	if c.Null {
		return "", true
	}

	if c.Blank {
		return "", false
	}

	if len(c.Suffix) > 0 {
		val = val + c.Suffix
	}

	if len(c.Prefix) > 0 {
		val = c.Prefix + val
	}

	return val, false
}
