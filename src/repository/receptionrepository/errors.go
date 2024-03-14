package receptionrepository

import "fmt"

func ErrOverflowArticles(articleId, articleCount, receptionCount uint) error {
	return fmt.Errorf(
		"error al intentar ingresar en el articulo(%v) la cantidad de %v items cuado deberia ser %v",
		articleId, receptionCount, articleCount,
	)
}
