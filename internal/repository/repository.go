package repository

// DataBaseRepo é uma interface que permite que as funcionalidades armazenadas sejam utilizadas em outros pkgs (handlers, por exemplo)
type DataBaseRepo interface {
	AllUsers() bool
}
