package reception_test

import (
	"desarrollosmoyan/lambda/src/model"
	"desarrollosmoyan/lambda/src/repository/receptionrepository"
	"desarrollosmoyan/lambda/src/tests/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestIsArticleOverflow(t *testing.T) {
	handleTest, err := utils.NewHandleTest()
	assert.NoError(t, err, "handleTest don't return error")
	tx := handleTest.Begin()
	defer tx.Rollback()

	t.Run("success", func(t *testing.T) {
		article := PreSet(t, tx)

		reception := model.ReceptionArt{
			ArticleID: article.ID,
			Count:     article.Count - 10,
		}
		err = tx.Save(&reception).Error
		assert.NoError(t, err)

		t.Run("with one reception", func(t *testing.T) {
			receptionRepo := receptionrepository.ReceptionRepo{Db: tx}
			err := receptionRepo.IsArticleOverflow(tx, reception.ArticleID)
			assert.NoError(t, err)
		})

		reception2 := model.ReceptionArt{
			ArticleID: article.ID,
			Count:     article.Count - reception.Count - 2,
		}
		err = tx.Save(&reception2).Error
		assert.NoError(t, err)

		t.Run("two receptions", func(t *testing.T) {
			receptionRepo := receptionrepository.ReceptionRepo{Db: tx}
			err := receptionRepo.IsArticleOverflow(tx, reception.ArticleID)
			assert.NoError(t, err)
		})
		reception3 := model.ReceptionArt{
			ArticleID: article.ID,
			Count:     article.Count - reception.Count - reception2.Count,
		}
		err = tx.Save(&reception3).Error
		assert.NoError(t, err)

		t.Run("completed article", func(t *testing.T) {
			receptionRepo := receptionrepository.ReceptionRepo{Db: tx}
			err := receptionRepo.IsArticleOverflow(tx, reception.ArticleID)
			assert.NoError(t, err)
		})
	})

	t.Run("fail", func(t *testing.T) {

		t.Run("with one reception", func(t *testing.T) {
			article := PreSet(t, tx)

			reception := model.ReceptionArt{
				ArticleID: article.ID,
				Count:     article.Count + 1,
			}

			err = tx.Save(&reception).Error
			assert.NoError(t, err)

			receptionRepo := receptionrepository.ReceptionRepo{Db: tx}
			err := receptionRepo.IsArticleOverflow(tx, reception.ArticleID)

			assert.EqualError(t,
				err,
				receptionrepository.ErrOverflowArticles(
					article.ID,
					article.Count,
					reception.Count,
				).Error(),
			)
		})

		t.Run("with two receptions", func(t *testing.T) {
			article := PreSet(t, tx)
			receptionRepo := receptionrepository.ReceptionRepo{Db: tx}

			t.Run("past one less than article count", func(t *testing.T) {
				reception := model.ReceptionArt{
					ArticleID: article.ID,
					Count:     article.Count - 1,
				}

				err = tx.Save(&reception).Error
				assert.NoError(t, err)

				err := receptionRepo.IsArticleOverflow(tx, reception.ArticleID)
				assert.NoError(t, err)

			})

			t.Run("fail two because bigger than article count total", func(t *testing.T) {
				reception2 := model.ReceptionArt{
					ArticleID: article.ID,
					Count:     10,
				}
				err = tx.Save(&reception2).Error
				assert.NoError(t, err)
				err := receptionRepo.IsArticleOverflow(tx, reception2.ArticleID)

				assert.EqualError(t,
					err,
					receptionrepository.ErrOverflowArticles(
						article.ID,
						article.Count,
						reception2.Count+article.Count-1,
					).Error(),
				)
			})
		})
	})
}

func PreSet(t *testing.T, tx *gorm.DB) *model.Article {
	warehouse := model.Warehouse{}
	err := tx.Save(&warehouse).Error
	assert.NoError(t, err)

	suppliers := model.Supplier{}
	err = tx.Save(&suppliers).Error
	assert.NoError(t, err)

	prod := model.Product{}
	err = tx.Save(&prod).Error
	assert.NoError(t, err)

	pur := model.Purchase{
		WarehouseID: warehouse.ID,
		SupplierID:  suppliers.ID,
	}

	err = tx.Save(&pur).Error
	assert.NoError(t, err)

	article := model.Article{
		ProductID:  prod.ID,
		PurchaseID: pur.ID,
		Count:      100,
	}

	err = tx.Save(&article).Error
	assert.NoError(t, err)

	return &article
}
