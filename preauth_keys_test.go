package headscale

import (
	"time"

	"gopkg.in/check.v1"
)

func (*Suite) TestCreatePreAuthKey(c *check.C) {
	_, err := h.CreatePreAuthKey("bogus", true, nil)

	c.Assert(err, check.NotNil)

	n, err := h.CreateNamespace("test")
	c.Assert(err, check.IsNil)

	k, err := h.CreatePreAuthKey(n.Name, true, nil)
	c.Assert(err, check.IsNil)

	// Did we get a valid key?
	c.Assert(k.Key, check.NotNil)
	c.Assert(len(k.Key), check.Equals, 48)

	// Make sure the Namespace association is populated
	c.Assert(k.Namespace.Name, check.Equals, n.Name)

	_, err = h.GetPreAuthKeys("bogus")
	c.Assert(err, check.NotNil)

	keys, err := h.GetPreAuthKeys(n.Name)
	c.Assert(err, check.IsNil)
	c.Assert(len(*keys), check.Equals, 1)

	// Make sure the Namespace association is populated
	c.Assert((*keys)[0].Namespace.Name, check.Equals, n.Name)
}

func (*Suite) TestExpiredPreAuthKey(c *check.C) {
	n, err := h.CreateNamespace("test2")
	c.Assert(err, check.IsNil)

	now := time.Now()
	pak, err := h.CreatePreAuthKey(n.Name, true, &now)
	c.Assert(err, check.IsNil)

	p, err := h.checkKeyValidity(pak.Key)
	c.Assert(err, check.Equals, errorAuthKeyExpired)
	c.Assert(p, check.IsNil)
}

func (*Suite) TestPreAuthKeyDoesNotExist(c *check.C) {
	p, err := h.checkKeyValidity("potatoKey")
	c.Assert(err, check.Equals, errorAuthKeyNotFound)
	c.Assert(p, check.IsNil)
}

func (*Suite) TestValidateKeyOk(c *check.C) {
	n, err := h.CreateNamespace("test3")
	c.Assert(err, check.IsNil)

	pak, err := h.CreatePreAuthKey(n.Name, true, nil)
	c.Assert(err, check.IsNil)

	p, err := h.checkKeyValidity(pak.Key)
	c.Assert(err, check.IsNil)
	c.Assert(p.ID, check.Equals, pak.ID)
}

func (*Suite) TestAlreadyUsedKey(c *check.C) {
	n, err := h.CreateNamespace("test4")
	c.Assert(err, check.IsNil)

	pak, err := h.CreatePreAuthKey(n.Name, false, nil)
	c.Assert(err, check.IsNil)

	db, err := h.db()
	if err != nil {
		c.Fatal(err)
	}
	defer db.Close()
	m := Machine{
		ID:             0,
		MachineKey:     "foo",
		NodeKey:        "bar",
		DiscoKey:       "faa",
		Name:           "testest",
		NamespaceID:    n.ID,
		Registered:     true,
		RegisterMethod: "authKey",
		AuthKeyID:      uint(pak.ID),
	}
	db.Save(&m)

	p, err := h.checkKeyValidity(pak.Key)
	c.Assert(err, check.Equals, errorAuthKeyNotReusableAlreadyUsed)
	c.Assert(p, check.IsNil)
}

func (*Suite) TestReusableBeingUsedKey(c *check.C) {
	n, err := h.CreateNamespace("test5")
	c.Assert(err, check.IsNil)

	pak, err := h.CreatePreAuthKey(n.Name, true, nil)
	c.Assert(err, check.IsNil)

	db, err := h.db()
	if err != nil {
		c.Fatal(err)
	}
	defer db.Close()
	m := Machine{
		ID:             1,
		MachineKey:     "foo",
		NodeKey:        "bar",
		DiscoKey:       "faa",
		Name:           "testest",
		NamespaceID:    n.ID,
		Registered:     true,
		RegisterMethod: "authKey",
		AuthKeyID:      uint(pak.ID),
	}
	db.Save(&m)

	p, err := h.checkKeyValidity(pak.Key)
	c.Assert(err, check.IsNil)
	c.Assert(p.ID, check.Equals, pak.ID)
}

func (*Suite) TestNotReusableNotBeingUsedKey(c *check.C) {
	n, err := h.CreateNamespace("test6")
	c.Assert(err, check.IsNil)

	pak, err := h.CreatePreAuthKey(n.Name, false, nil)
	c.Assert(err, check.IsNil)

	p, err := h.checkKeyValidity(pak.Key)
	c.Assert(err, check.IsNil)
	c.Assert(p.ID, check.Equals, pak.ID)
}
