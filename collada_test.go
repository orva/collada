package collada

import (
	. "launchpad.net/gocheck"
)

type ColladaSuite struct{} 
var _ = Suite(&ColladaSuite{})

func (s *ColladaSuite) TestNothing(c *C) {
	c.Check(true, Equals, true)
}
