package model

type Positionable interface {
	SetPosition(int, int)	
	SetQuantity(int)
	GetType() string
}	