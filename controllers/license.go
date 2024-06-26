package controllers

import (
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/f-license/lcs"
	"github.com/crawlab-team/f-license/storage"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type licenseController struct {
}

func (ctr *licenseController) Get(c *gin.Context) {
	r := c.Request
	id := mux.Vars(r)["id"]

	var l lcs.License
	err := storage.LicenseHandler.GetByID(id, &l)
	if err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}
	controllers.HandleSuccessWithData(c, l)
}

func (ctr *licenseController) GetList(c *gin.Context) {
	var licenses []*lcs.License
	err := storage.LicenseHandler.GetAll(&licenses)
	if err != nil {
		controllers.HandleErrorInternalServerError(c, err)
		return
	}
	controllers.HandleSuccessWithListData(c, licenses, len(licenses))
}

func (ctr *licenseController) PostGenerate(c *gin.Context) {
	var l lcs.License
	if err := c.ShouldBindJSON(&l); err != nil {
		controllers.HandleErrorBadRequest(c, err)
		return
	}

	err := l.Generate()
	if err != nil {
		logrus.WithError(err).Error("License couldn't be generated")
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	err = storage.LicenseHandler.AddIfNotExisting(&l)
	if err != nil {
		logrus.WithError(err).Error("License couldn't be stored")
		controllers.HandleErrorInternalServerError(c, err)
		return
	}

	controllers.HandleSuccessWithData(c, map[string]interface{}{
		"id":    l.ID.Hex(),
		"token": l.Token,
	})
}

var LicenseController = licenseController{}
