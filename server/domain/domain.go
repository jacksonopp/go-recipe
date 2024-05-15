package domain

type DomainObject interface {
	ToDto() Dto
}

type Dto interface {
}
