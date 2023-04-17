package models

import "github.com/lib/pq"

type Product struct {
	ID   string `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name"`
}

type ProductAttributeList struct {
	Product
	SubcategoryID string         `db:"subcategory_id"`
	BrandID       string         `db:"brand_id"`
	AttributeIDs  pq.StringArray `json:"attributeId" db:"attribute_id"`
}

type ProductWeight struct {
	Product
	Weight int `json:"weight" db:"weight"`
}
