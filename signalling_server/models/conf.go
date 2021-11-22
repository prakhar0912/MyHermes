package models

type Conf struct {
	Production  *Deployment `toml:"production"`
	Development *Deployment `toml:"development"`
	Staging     *Deployment `toml:"staging"`
}