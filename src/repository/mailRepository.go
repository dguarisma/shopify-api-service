package repository

import "desarrollosmoyan/lambda/src/model"

type MailRepository interface {
	GetMsg(purchase *model.Purchase) (mail *model.EmailMsg, err error)
	GetMsg2(purchase *model.Purchase) (mail *model.EmailMsg, err error)
}
