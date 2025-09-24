package analyzer

import "fmt"

type FileNotFoundError struct {
	Path string
	Err  error
}

func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf("fichier introuvable: %s (%v)", e.Path, e.Err)
}

func (e *FileNotFoundError) Unwrap() error {
	return e.Err
}

type ParseLogError struct {
	Path string
	Msg  string
}

func (e *ParseLogError) Error() string {
	return fmt.Sprintf("erreur de parsing sur %s: %s", e.Path, e.Msg)
}
