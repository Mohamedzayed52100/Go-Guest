package common

func (r *CommonRepository) GetTenantDBNameFromClient(client string) string {
	var tenantDBName string
	r.SharedDbConnection.Table("tenants").Joins("JOIN tenant_profiles ON tenant_id = tenants.id").Where("name = ?", client).Select("db_name").Scan(&tenantDBName)
	return tenantDBName
}
