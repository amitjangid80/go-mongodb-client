package mongodb_domain

type BaseDmlModel interface {
	SetCreatedBy(by string)
	SetCreatedOn(on string)
	SetModifiedBy(by string)
	SetModifiedOn(on string)
	GetId() string
	SetId(id string)
}

type DmlModel struct {
	Id         string `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedOn  string `json:"createdOn,omitempty" bson:"createdOn,omitempty"`
	CreatedBy  string `json:"createdBy" bson:"createdBy"`
	ModifiedOn string `json:"modifiedOn,omitempty" bson:"modifiedOn,omitempty"`
	ModifiedBy string `json:"modifiedBy" bson:"modifiedBy"`
}

// Implement Auditable interface
func (d *DmlModel) SetCreatedBy(by string)  { d.CreatedBy = by }
func (d *DmlModel) SetCreatedOn(on string)  { d.CreatedOn = on }
func (d *DmlModel) SetModifiedBy(by string) { d.ModifiedBy = by }
func (d *DmlModel) SetModifiedOn(on string) { d.ModifiedOn = on }
func (d *DmlModel) GetId() string           { return d.Id }
func (d *DmlModel) SetId(id string)         { d.Id = id }
