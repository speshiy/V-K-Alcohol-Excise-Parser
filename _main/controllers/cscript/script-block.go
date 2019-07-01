package cscript

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/_main/models/mscript"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/common"
)

//GetScriptBlocks return all script blocks by CompanyID and ScriptTypeID
func GetScriptBlocks(c *gin.Context) {
	var err error
	var blocks []mscript.ScriptBlock
	var typeID uint

	typeID, err = common.ParseParam(c, "TypeID")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	err = mscript.GetBlocks(c, typeID, &blocks)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Done", "data": blocks})
}

//PostScriptBlock add new block
func PostScriptBlock(c *gin.Context) {
	var block mscript.ScriptBlock
	var err error

	if err := c.ShouldBindJSON(&block); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	err = block.Post(c, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Done", "data": block})
}

//PutScriptBlock update script block
func PutScriptBlock(c *gin.Context) {
	var block mscript.ScriptBlock
	var newBlock mscript.ScriptBlock
	var err error

	if err := c.ShouldBindJSON(&block); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	replacer := strings.NewReplacer(
		" ", "_",
		"-", "_",
		"=", "_",
	)

	if block.HREF == "" {
		uid, _ := common.GetNewUUID()
		block.HREF = replacer.Replace(uid)
	}

	//Assign data to new block
	newBlock = block
	newBlock.ID = 0
	newBlock.DeletedAt = nil

	DB := c.MustGet("DB").(*gorm.DB)
	TX := DB.Begin()

	//Delete old block for history
	err = block.Delete(c, TX)
	if err != nil {
		TX.Rollback()
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	err = newBlock.Post(c, TX)
	if err != nil {
		TX.Rollback()
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	TX.Commit()

	// err = block.Put(c, nil)
	// if err != nil {
	// 	c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Done", "data": newBlock})
}

//DeleteScriptBlock delete script block
func DeleteScriptBlock(c *gin.Context) {
	var err error
	var id uint
	var block mscript.ScriptBlock

	id, err = common.ParseParam(c, "ID")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	if id < 1 {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Block ID invalid"})
		return
	}

	block.ID = id
	err = block.Delete(c, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Done"})
}
