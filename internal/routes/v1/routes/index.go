package routes

func PutHandler(c *gin.Context) {
	db := middleware.GetDB(c)

	var body struct {
		UserID  string `json:"user_id"`
		GuildID string `json:"guild_id"`
		Balance int    `json:"balance"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	gu := database.GuildMember{
		GuildID: body.GuildID,
		UserID:  body.UserID,
		Balance: body.Balance,
	}
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "guild_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"balance"}),
	}).Create(&gu)
	c.JSON(http.StatusOK, gu)
}

func Register(r *gin.RouterGroup) {
	r.GET("/routes", GetHandler)

}
