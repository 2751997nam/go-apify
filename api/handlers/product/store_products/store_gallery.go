package storeproducts

import (
	"product-service/internal/models"

	"github.com/2751997nam/go-helpers/utils"
)

func bulkStoreGallery(data []models.ProductGalleryData) {
	db := models.GetDB()
	createItems := []models.ProductGallery{}
	for _, dataItem := range data {
		gallery := dataItem.Gallery
		existedItems := []models.ProductGallery{}
		db.Model(&models.ProductGallery{}).Where("product_id = ?", dataItem.ProductId).Where("type = ?", dataItem.Type).Find(&existedItems)
		itemExistById := map[uint64]models.ProductGallery{}
		for _, value := range existedItems {
			itemExistById[value.ID] = value
		}
		for _, imageUrl := range gallery {
			galleryItem := models.ProductGallery{
				Type:      dataItem.Type,
				ProductId: dataItem.ProductId,
				ImageUrl:  imageUrl,
			}
			createItems = append(createItems, galleryItem)
		}

		err := db.Where("product_id = ?", dataItem.ProductId).Where("type = ?", dataItem.Type).Delete(&models.ProductGallery{}).Error
		if err != nil {
			utils.LogPanic(err)
		} else {
			utils.QuickLog(map[string]any{"data": data}, dataItem.ProductId, "PRODUCT_GALLERY_"+dataItem.Type, "DELETE")
		}
	}

	err := db.Model(&models.ProductGallery{}).CreateInBatches(createItems, 100).Error
	if err != nil {
		utils.LogPanic(err)
	} else {
		for _, item := range createItems {
			utils.QuickLog(item, item.ProductId, "PRODUCT_GALLERY_"+item.Type, "CREATE")
		}
	}
}
