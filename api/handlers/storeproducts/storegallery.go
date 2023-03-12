package storeproducts

import (
	"log"
	"product-service/internal/models"
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
			log.Panic(err)
		}
	}

	err := db.Model(&models.ProductGallery{}).CreateInBatches(createItems, 100).Error
	if err != nil {
		log.Panic(err)
	}
}
