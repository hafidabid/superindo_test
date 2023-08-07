package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strconv"
	"superindo_diksha_test/config"
	"superindo_diksha_test/controller"
	"superindo_diksha_test/database"
	"superindo_diksha_test/utils"
	"sync"
)

type Handler struct {
	config config.Config
	ctrl   controller.ProductController
}

func NewHandler(conf config.Config, appDb database.AppDatabase) *Handler {
	// seed data only for one in a life time
	ctrl := controller.NewController(conf, appDb)
	onlyOneSeed := &sync.Once{}
	go func() {
		onlyOneSeed.Do(func() {
			err := ctrl.GenerateData(500)
			if err != nil {
				log.Fatal("Error seed data -> ", err)
				return
			} else {
				log.Infof("Success seed %d data", 500)
			}
		})
	}()

	return &Handler{
		config: conf,
		ctrl:   ctrl,
	}
}

// SourceProduct godoc
// @Summary Get source product
// @Description Get source product
// @Accept json
// @Produce json
// @Param Authorization header string false "Auth key here"
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {string} utils.Response
// @Router /products/source [get]
func (h Handler) SourceProduct(c *gin.Context) {
	page, limit := paginationGetter(c)

	data, meta, err := h.ctrl.GetSource(page, limit)
	status := 200
	metadata := utils.Metadata{}

	if err != nil {
		status = 500
		log.Errorf("Error get destination, %v", err)
	}

	if meta != nil {
		metadata = *meta
	}
	log.Infof("debug data -> ", metadata)
	c.JSON(status, utils.Response{
		Status:   status,
		Success:  err == nil,
		Data:     data,
		Metadata: metadata,
		Message:  "Success get source product data",
	})
	return
}

// DestinationProduct godoc
// @Summary Get destination product
// @Description Get destination product
// @Accept json
// @Produce json
// @Param Authorization header string false "Auth key here"
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {string} utils.Response
// @Router /products/destination [get]
func (h Handler) DestinationProduct(c *gin.Context) {
	page, limit := paginationGetter(c)

	data, meta, err := h.ctrl.GetDestination(page, limit)
	status := 200
	metadata := utils.Metadata{}

	if err != nil {
		status = 500
		log.Errorf("Error get destination, %v", err)
	}
	if meta != nil {
		metadata = *meta
	}

	c.JSON(status, utils.Response{
		Status:   0,
		Success:  err == nil,
		Data:     data,
		Metadata: metadata,
		Message:  "Success get destination product data",
	})
	return
}

// UpdateCheckProduct godoc
// @Summary Check and update product
// @Description Check and update product
// @Accept json
// @Produce json
// @Param Authorization header string false "Auth key here"
// @Success 200 {string} utils.Response
// @Router /products/ [post]
func (h Handler) UpdateCheckProduct(c *gin.Context) {
	err := h.ctrl.UpdateData()

	status := 200
	if err != nil {
		status = 500
		log.Errorf("Error check and update data, %v", err)
	}

	c.JSON(status, utils.Response{
		Status:  status,
		Success: err == nil,
		Message: "Success check and update data",
	})
	return
}

func paginationGetter(c *gin.Context) (int, int) {
	pageString, isPage := c.GetQuery("page")
	limitString, isLimit := c.GetQuery("limit")

	page := 0
	limit := 20
	if isPage {
		page, _ = strconv.Atoi(pageString)
	}
	if isLimit {
		limit, _ = strconv.Atoi(limitString)
	}

	return page, limit
}
